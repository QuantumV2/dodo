package tokenizer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/QuantumV2/dodo/pkg/helpers"
)

type Tokenizer struct {
	filesImported map[string]struct{}
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{filesImported: map[string]struct{}{}}
}

func (t *Tokenizer) SplitTokens(code string, allowImports bool) ([]string, error) {
	codeRunes := []rune(code)
	var tokens []string
	i := 0
	n := len(codeRunes)

	for i < n {

		for i < n && unicode.IsSpace(codeRunes[i]) {
			i++
		}
		if i >= n {
			break
		}

		if codeRunes[i] == ';' {
			for i < n && codeRunes[i] != '\n' {
				i++
			}
			continue
		}

		if codeRunes[i] == '"' {
			start := i
			i++
			for i < n && codeRunes[i] != '"' {
				i++
			}
			if i < n {
				i++
			}
			tokens = append(tokens, string(codeRunes[start:i]))
			continue
		}

		start := i
		for i < n && !unicode.IsSpace(codeRunes[i]) && codeRunes[i] != ';' && codeRunes[i] != '"' {
			i++
		}

		tokens = append(tokens, strings.ToLower(string(codeRunes[start:i])))

		if len(tokens) >= 2 && tokens[len(tokens)-1] == "import" {
			if !allowImports {
				return tokens, fmt.Errorf("Cannot use \"IMPORT\" because imports are disabled")
			}
			path := tokens[len(tokens)-2]
			path = path[1 : len(path)-1]
			if _, ok := t.filesImported[path]; ok {
				return tokens, fmt.Errorf("File %s already imported. check for circular references", path)
			}
			content, err := helpers.GetLibrary(path)
			if err != nil {
				return tokens, err
			}
			t.filesImported[path] = struct{}{}
			insertTokens, err := t.SplitTokens(string(content), allowImports)
			if err != nil {
				return tokens, err
			}
			tokens = tokens[:len(tokens)-2]
			tokens = append(tokens, insertTokens...)

		}
	}
	return tokens, nil
}
