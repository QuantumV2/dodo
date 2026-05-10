package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/QuantumV2/dodo"
	"github.com/QuantumV2/dodo/pkg/interpreter"
	"github.com/QuantumV2/dodo/pkg/tokenizer"
)

var default_config dodo.Config = dodo.Config{AllowImports: true}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runREPL()
		return
	}

	switch args[0] {
	case "repl":
		runREPL()
		return

	case "init-lib":
		if err := dodo.CreateLibDirIfMissing(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating lib dir: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Library directory created.")
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
	fmt.Println("  init-lib   Create the global Dodo library directory")
	fmt.Println("  help       Show this help")
	fmt.Println()
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
		_, err = i.Interpret()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

	}

}
