package dodo

import (
	"os"

	"github.com/QuantumV2/dodo/pkg/helpers"
	"github.com/QuantumV2/dodo/pkg/interpreter"
	"github.com/QuantumV2/dodo/pkg/tokenizer"
)

type Config struct {
	AllowImports bool
	Debug        bool
}

func Run(input string, cfg *Config) ([]int, error) {
	t := tokenizer.NewTokenizer()
	i := interpreter.NewInterpreter()
	var err error
	i.Tokens, err = t.SplitTokens(input, cfg.AllowImports)
	if err != nil {
		return nil, err
	}
	return i.Interpret(cfg.Debug)

}

func CreateLibDirIfMissing() error {
	libDir, err := helpers.DefaultLibDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(libDir, 0755)
}
