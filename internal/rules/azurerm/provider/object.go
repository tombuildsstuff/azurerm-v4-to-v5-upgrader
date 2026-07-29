package provider

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// parseObjectExpr parses attr's expression and, if it is an object
// constructor (e.g. { source = "x" }), returns it as an
// *hclsyntax.ObjectConsExpr along with the exact source bytes it was parsed
// from. Those bytes are attr.Expr().BuildTokens(nil).Bytes() - the same
// token stream a caller would walk to make surgical edits - so byte offsets
// in the returned expression's ranges line up 1:1 with offsets computed via
// tokenOffsets over that same token stream.
//
// ok is false when the expression fails to parse or isn't an
// *hclsyntax.ObjectConsExpr at all (e.g. a bare traversal) - i.e. not a
// shape callers can safely edit in place.
func parseObjectExpr(attr *hclwrite.Attribute) (oce *hclsyntax.ObjectConsExpr, tokens hclwrite.Tokens, ok bool) {
	tokens = attr.Expr().BuildTokens(nil)
	expr, diags := hclsyntax.ParseExpression(tokens.Bytes(), "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, nil, false
	}
	oce, ok = expr.(*hclsyntax.ObjectConsExpr)
	if !ok {
		return nil, nil, false
	}
	return oce, tokens, true
}

// findObjectItem returns the item in oce whose key is the literal string
// key, e.g. {version = ...} for key "version".
func findObjectItem(oce *hclsyntax.ObjectConsExpr, key string) (item *hclsyntax.ObjectConsItem, ok bool) {
	for i, it := range oce.Items {
		kv, diags := it.KeyExpr.Value(nil)
		if diags.HasErrors() || kv.Type().FriendlyName() != "string" {
			continue
		}
		if kv.AsString() == key {
			return &oce.Items[i], true
		}
	}
	return nil, false
}

// tokenOffsets returns, for each token in tokens, the byte offset (within
// tokens.Bytes()) at which that token's own bytes begin, i.e. immediately
// after its SpacesBefore padding. hclwrite.Tokens.Bytes() writes each
// token's SpacesBefore run of spaces followed by its Bytes, in order, so
// this mirrors that layout exactly - letting hcl.Range byte offsets
// (computed by parsing tokens.Bytes()) be mapped back onto token
// boundaries.
func tokenOffsets(tokens hclwrite.Tokens) []int {
	offsets := make([]int, len(tokens))
	pos := 0
	for i, t := range tokens {
		pos += t.SpacesBefore
		offsets[i] = pos
		pos += len(t.Bytes)
	}
	return offsets
}

// replaceExprRange returns a NEW token slice with the tokens spanning rng
// replaced by replacement; tokens itself is left untouched (its underlying
// Token pointers are shared with the live document tree, so callers must
// never mutate them in place - see TraversalRenamesNameToID for the same
// discipline). tokens must be the exact stream rng's byte offsets were
// computed against, i.e. parsed from tokens.Bytes().
//
// The first replacement token inherits the SpacesBefore of the first
// replaced token, so the spacing before the value (e.g. the space after
// "=") is preserved.
func replaceExprRange(tokens hclwrite.Tokens, rng hcl.Range, replacement hclwrite.Tokens) hclwrite.Tokens {
	offsets := tokenOffsets(tokens)

	start := len(tokens)
	for i, off := range offsets {
		if off >= rng.Start.Byte {
			start = i
			break
		}
	}
	end := len(tokens)
	for i, off := range offsets {
		if off >= rng.End.Byte {
			end = i
			break
		}
	}

	repl := make(hclwrite.Tokens, len(replacement))
	copy(repl, replacement)
	if len(repl) > 0 && start < len(tokens) {
		first := *repl[0]
		first.SpacesBefore = tokens[start].SpacesBefore
		repl[0] = &first
	}

	out := make(hclwrite.Tokens, 0, len(tokens)-(end-start)+len(repl))
	out = append(out, tokens[:start]...)
	out = append(out, repl...)
	out = append(out, tokens[end:]...)
	return out
}

// insertObjectItem returns a NEW token slice with a `key = <value>` item
// inserted into the object described by tokens/oce, just before its closing
// brace. Indentation matches the last existing item's key when the object
// has one; otherwise it falls back to two spaces.
func insertObjectItem(tokens hclwrite.Tokens, oce *hclsyntax.ObjectConsExpr, key string, value hclwrite.Tokens) hclwrite.Tokens {
	indent := 2
	if n := len(oce.Items); n > 0 {
		keyStart := oce.Items[n-1].KeyExpr.Range().Start.Byte
		for i, off := range tokenOffsets(tokens) {
			if off == keyStart {
				indent = tokens[i].SpacesBefore
				break
			}
		}
	}

	item := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(key), SpacesBefore: indent},
		{Type: hclsyntax.TokenEqual, Bytes: []byte("="), SpacesBefore: 1},
	}
	val := make(hclwrite.Tokens, len(value))
	copy(val, value)
	if len(val) > 0 {
		first := *val[0]
		first.SpacesBefore = 1
		val[0] = &first
	}
	item = append(item, val...)
	item = append(item, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})

	insertAt := len(tokens) - 1 // just before the closing brace
	if insertAt < 0 {
		insertAt = 0
	}
	out := make(hclwrite.Tokens, 0, len(tokens)+len(item))
	out = append(out, tokens[:insertAt]...)
	out = append(out, item...)
	out = append(out, tokens[insertAt:]...)
	return out
}
