package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/henningstorck/monkey-interpreter/lexer"
	"github.com/henningstorck/monkey-interpreter/parser"
)

const prompt = ">> "

const monkeyFace = `           __,__
  .--.  .-"     "-.  .--.
 / .. \/  .-. .-.  \/ .. \
| |  '|  /   Y   \  |'  | |
| \   \  \ 0 | 0 /  /   / |
 \ '- ,\.-"""""""-./, -' /
  ''-' /_   ^ ^   _\ '-''
      |  \._   _./  |
      \   \ '~' /   /
       '._ '-=-' _.'
          '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.NewLexer(line)
		par := parser.NewParser(lex)
		program := par.ParseProgram()

		if len(par.Errors()) != 0 {
			printParseErrors(out, par.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, monkeyFace)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, "Parser errors:\n")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
