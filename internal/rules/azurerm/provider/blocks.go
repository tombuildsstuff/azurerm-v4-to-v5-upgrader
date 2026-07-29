package provider

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// azureRMProviderBlocks returns every `provider "azurerm"` block in body
// (including aliased ones - the first label is still "azurerm").
func azureRMProviderBlocks(body *hclwrite.Body) []*hclwrite.Block {
	var out []*hclwrite.Block
	for _, b := range body.Blocks() {
		if b.Type() != "provider" {
			continue
		}
		labels := b.Labels()
		if len(labels) >= 1 && labels[0] == "azurerm" {
			out = append(out, b)
		}
	}
	return out
}

// literalBool returns the literal boolean value of attr's expression; ok is
// false for anything that isn't a bare true/false literal (var/local ref,
// interpolation, function call, non-bool).
func literalBool(attr *hclwrite.Attribute) (value bool, ok bool) {
	src := attr.Expr().BuildTokens(nil).Bytes()
	expr, diags := hclsyntax.ParseExpression(src, "", hcl.InitialPos)
	if diags.HasErrors() {
		return false, false
	}
	if _, isLit := expr.(*hclsyntax.LiteralValueExpr); !isLit {
		return false, false
	}
	if len(expr.Variables()) > 0 {
		return false, false
	}
	v, vdiags := expr.Value(nil)
	if vdiags.HasErrors() || v.IsNull() || v.Type() != cty.Bool {
		return false, false
	}
	return v.True(), true
}
