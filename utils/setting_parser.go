package utils

import (
	"strings"
	"unicode"
)

type State int

const (
	DEFAULT State = iota
	SINGLE_LINE_COMMENT
	MULTI_LINE_COMMENT
	IN_SET
	IN_KEY_VALUE
	STOP
)

type ParseResult struct {
	Settings       map[string]string
	RemainingQuery string
	Errors         []string
}

func Parse(query string) ParseResult {
	trimmed := strings.TrimSpace(query)
	if !strings.HasSuffix(trimmed, ";") {
		query += ";"
	}
	return extractSetStatements(query)
}

func extractSetStatements(s string) ParseResult {
	settings := make(map[string]string)
	errors := []string{}
	excludeRanges := [][]int{}
	currentState := DEFAULT
	lowerS := strings.ToLower(s)
	i := 0
	currentStartIndex := -1

	for i < len(s) {
		switch currentState {
		case DEFAULT:
			if i <= len(s)-2 && s[i] == '-' && s[i+1] == '-' {
				currentState = SINGLE_LINE_COMMENT
				i += 2
			} else if i <= len(s)-2 && s[i] == '/' && s[i+1] == '*' {
				currentState = MULTI_LINE_COMMENT
				i += 2
			} else if i <= len(s)-3 && lowerS[i] == 's' && lowerS[i+1] == 'e' && lowerS[i+2] == 't' {
				if i+3 < len(s) && unicode.IsSpace(rune(s[i+3])) {
					currentState = IN_SET
					currentStartIndex = i
					i += 4
				} else {
					i++
				}
			} else {
				if i < len(s) && !unicode.IsSpace(rune(s[i])) {
					currentState = STOP
				}
				i++
			}
		case SINGLE_LINE_COMMENT:
			for i < len(s) && s[i] != '\n' {
				i++
			}
			if i < len(s) {
				i++
			}
			currentState = DEFAULT
		case MULTI_LINE_COMMENT:
			for i < len(s) {
				if i+1 < len(s) && s[i] == '*' && s[i+1] == '/' {
					i += 2
					currentState = DEFAULT
					break
				} else {
					i++
				}
			}
		case IN_SET:
			for i < len(s) && unicode.IsSpace(rune(s[i])) {
				i++
			}
			if i < len(s) {
				currentState = IN_KEY_VALUE
			} else {
				errors = append(errors, "Invalid SET statement: missing key-value after 'set'")
				currentStartIndex = -1
				currentState = DEFAULT
			}
		case IN_KEY_VALUE:
			keyValueStart := i
			foundSemicolon := false
			for i < len(s) {
				if s[i] == ';' {
					if i > 0 && s[i-1] != '\\' {
						foundSemicolon = true
						i++
						break
					}
				}
				i++
			}
			if foundSemicolon {
				kv := strings.TrimSpace(s[keyValueStart : i-1])
				success := parseKeyValue(kv, settings, &errors)
				if success {
					excludeRanges = append(excludeRanges, []int{currentStartIndex, i})
				}
			} else {
				errors = append(errors, "Invalid SET statement: missing semicolon")
			}
			currentState = DEFAULT
			currentStartIndex = -1
		case STOP:
			i = len(s)
		default:
			i++
		}
	}

	var remainingQuery strings.Builder
	currentPos := 0
	for _, r := range excludeRanges {
		if currentPos < r[0] {
			remainingQuery.WriteString(s[currentPos:r[0]])
		}
		currentPos = r[1]
	}
	if currentPos < len(s) {
		remainingQuery.WriteString(s[currentPos:])
	}

	return ParseResult{
		Settings:       settings,
		RemainingQuery: remainingQuery.String(),
		Errors:         errors,
	}
}

func parseKeyValue(kv string, settings map[string]string, errors *[]string) bool {
	eqIdx := strings.Index(kv, "=")
	if eqIdx == -1 {
		*errors = append(*errors, "Invalid key-value pair '"+kv+"': missing '='")
		return false
	}
	key := strings.TrimSpace(kv[:eqIdx])
	if key == "" {
		*errors = append(*errors, "Invalid key-value pair '"+kv+"': empty key")
		return false
	}
	var value string
	if eqIdx < len(kv)-1 {
		value = strings.TrimSpace(kv[eqIdx+1:])
	} else {
		value = ""
	}
	value = strings.ReplaceAll(value, "\\;", ";")
	settings[key] = value
	return true
}
