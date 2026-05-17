package dodo

import (
	"fmt"
	"os"
	"time"

	"codeberg.org/QuantumV/dodo/pkg/helpers"
	"codeberg.org/QuantumV/dodo/pkg/interpreter"
	"codeberg.org/QuantumV/dodo/pkg/tokenizer"
)

type Config struct {
	AllowImports bool
	Debug        bool
	DebugDelayMS time.Duration
}

func Run(input string, cfg *Config) ([]int, error) {
	t := tokenizer.NewTokenizer()
	i := interpreter.NewInterpreter()
	var err error
	i.Tokens, err = t.SplitTokens(input, cfg.AllowImports)
	if cfg.Debug {
		fmt.Printf("%v", i.Tokens)
	}
	if err != nil {
		return nil, err
	}
	return i.Interpret(cfg.Debug, cfg.DebugDelayMS)

}

func CreateLibDirIfMissing() error {
	libDir, err := helpers.DefaultLibDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(libDir, 0755)
}
