package librust

import (
	"fmt"
	"os"

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

func DumpTokens(in string, out *os.File) {
	input := antlr.NewInputStream(in)
	lexer := parser.NewRustLexer(input)
	var vocabulary TokenVocabulary = lexer.GetSymbolicNames()
	for _, token := range lexer.GetAllTokens() {
		_, err := out.WriteString(vocabulary.LLVMFormat(&token))
		if err != nil {
			panic(err)
		}
	}
}
