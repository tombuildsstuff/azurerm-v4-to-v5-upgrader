package rules

var registry []Rule

// Register adds a rule to the global registry (call from init()).
func Register(r Rule) { registry = append(registry, r) }

// All returns a copy of the registered rules in registration order.
func All() []Rule { return append([]Rule(nil), registry...) }

// resetRegistry clears the registry (tests only).
func resetRegistry() { registry = nil }
