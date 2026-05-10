package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/QuantumV2/dodo"
	"github.com/QuantumV2/dodo/pkg/helpers"
	"github.com/QuantumV2/dodo/pkg/interpreter"
	"github.com/QuantumV2/dodo/pkg/tokenizer"
)

var default_config dodo.Config = dodo.Config{AllowImports: true, Debug: false}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runREPL()
		return
	}
	if len(args) > 1 {
		switch args[1] {
		case "debug":
			default_config.Debug = true
		}
	}

	switch args[0] {
	case "repl":
		runREPL()
		return

	case "init-lib":
		if err := dodo.CreateLibDirIfMissing(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating lib dir: %s\n", err.Error())
			os.Exit(1)
		}
		a, _ := helpers.DefaultLibDir()
		fmt.Printf("Library directory at path %s created.\n", a)
		return

	case "help", "--help", "-h":
		printHelp()
		return
	}

	runFile(args[0])
}

func runFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	result, err := dodo.Run(string(content), &default_config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime Error: %v\n", err)
		os.Exit(1)
	}

	if result != nil {
		fmt.Println(result)
	}
}

func printHelp() {
	fmt.Println("dodo [command] [file]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  repl       Start the interactive REPL")
	fmt.Println("  init-lib   Create the global DODO library directory")
	fmt.Println("  help       Show this help")
	fmt.Println()
	fmt.Println("Add a \"debug\" after any command to enter into debug mode and see things like the stack while the program is running.")
	fmt.Println("If the first argument is not a command, it is treated as a script file.")
}

func runREPL() {
	t := tokenizer.NewTokenizer()
	i := interpreter.NewInterpreter()
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("DODO Interactive REPL (Press Ctrl+C to exit)")
	for {
		fmt.Printf("\n> ")

		if !scanner.Scan() {
			continue
		}
		input = scanner.Text()

		toks, err := t.SplitTokens(input, default_config.AllowImports)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		i.Tokens = toks
		_, err = i.Interpret(default_config.Debug)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

	}

}
