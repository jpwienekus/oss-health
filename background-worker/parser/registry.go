package parser

type ParserRegistry struct {
	parsers []Parser
}

func NewParserRegistry() *ParserRegistry {
	return &ParserRegistry{
		parsers: []Parser{
			RequirementsParser{},
			NpmParser{},
			PoetryPyProjectParser{},
			Pep621PyProjectParser{},
		},
	}
}

func (r *ParserRegistry) GetParsers(filename string) ([]Parser, bool) {
	parsers := []Parser{}
	hasParser := false

	for _, parser := range r.parsers {
		if parser.CanParse(filename) {
			parsers = append(parsers, parser)
			hasParser = true
		}
	}

	return parsers, hasParser
}
