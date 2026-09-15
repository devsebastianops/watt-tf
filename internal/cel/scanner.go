package cel

import "strings"

type FoundInterpolation struct {
	Start int
	End   int
	Expr  string
}

const INTERPOLATION_START = "${"

func (e *Engine) FindInterpolations(s string) []FoundInterpolation {
	result := []FoundInterpolation{}

	for i := 0; i < len(s); i++ {
		if s[i] == '$' && i+1 < len(s) && s[i+1] == '{' {

			start := i
			i += 2

			depth := 1
			exprStart := i

		innerLoop:
			for i < len(s) {
				switch s[i] {
				case '{':
					depth++
				case '}':
					depth--
					if depth == 0 {
						result = append(result, FoundInterpolation{
							Start: start,
							End:   i + 1,
							Expr:  s[exprStart:i],
						})
						break innerLoop
					}
				}

				i++
			}
		}
	}

	return result
}

func (e *Engine) HasInterpolation(s string) bool {
	return s != "" && strings.Contains(s, INTERPOLATION_START)
}
