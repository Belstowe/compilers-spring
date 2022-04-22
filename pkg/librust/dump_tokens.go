package librust

import (
	"fmt"
	"io"

	"github.com/Compiler2022/compilers-1-Belstowe/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type TokenVocabulary []string

func (g *TokenVocabulary) LLVMFormat(token *antlr.Token) string {
	return fmt.Sprintf("Loc=<%d:%d>\t%s '%s'\n",
		(*token).GetLine(),
		(*token).GetColumn()+1,
		(*g)[(*token).GetTokenType()],
		(*token).GetText())
}

func DumpTokens(in io.Reader, out io.Writer) {
	b, err := io.ReadAll(in)
	if err != nil {
		panic(err)
	}

	input := antlr.NewInputStream(string(b))
	lexer := parser.NewRustLexer(input)
	var vocabulary TokenVocabulary = lexer.GetSymbolicNames()
	for _, token := range lexer.GetAllTokens() {
		_, err := out.Write([]byte(vocabulary.LLVMFormat(&token)))
		if err != nil {
			panic(err)
		}
	}
}
