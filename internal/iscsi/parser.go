package iscsi

const DefaultParser = "generic"

type TargetNameParser interface {
	Name() string
	Parse(iqn string) (string, bool)
}

var registry = map[string]TargetNameParser{}

// RegisterParser save into the registry a TargetNameParser
func RegisterParser(p TargetNameParser) {
	registry[p.Name()] = p
}

// RegisterParser save into the registry a slice of TargetNameParser
func RegisterParsers(parsers []TargetNameParser) {
	for _, p := range parsers {
		RegisterParser(p)
	}
}

// GetParser returns the parser registered under the given name, and a boolean indicating whether
// the parser was found in the registry.
func GetParser(name string) (TargetNameParser, bool) {
	p, ok := registry[name]
	return p, ok
}

// ParserExists reports whether a parser with the given name is registered.
func ParserExists(name string) bool {
	_, ok := registry[name]
	return ok
}
