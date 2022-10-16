package librust

import (
	"fmt"
	"io"

	"github.com/Compiler2022/compilers-1-Belstowe/parser"
	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/ast"
	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/llvmir"
	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/semantics"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/llir/llvm/ir"
	"gopkg.in/yaml.v3"
)

type Error struct {
	Line    int
	Column  int
	Message string
}

func (err Error) Format() string {
	return fmt.Sprintf("<%d:%d>\t%s\n", err.Line, err.Column, err.Message)
}

type StreamErrorListener struct {
	*antlr.DefaultErrorListener
	errors []Error
}

func NewStreamErrorListener() StreamErrorListener {
	var sel StreamErrorListener
	sel.errors = make([]Error, 0)
	return sel
}

func (sel StreamErrorListener) Errors() []Error {
	return sel.errors
}

func (sel StreamErrorListener) HasErrors() bool {
	return len(sel.errors) != 0
}

func (sel *StreamErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	sel.errors = append(sel.errors, Error{line, column, msg})
}

type TokenVocabulary []string

func (g *TokenVocabulary) LLVMFormat(token *antlr.Token) string {
	return fmt.Sprintf("Loc=<%d:%d>\t%s '%s'\n",
		(*token).GetLine(),
		(*token).GetColumn()+1,
		(*g)[(*token).GetTokenType()],
		(*token).GetText())
}

func Parse(in io.Reader, out io.Writer, to_dump_tokens bool, to_dump_ast bool, verbose bool) {
	b, err := io.ReadAll(in)
	if err != nil {
		panic(err)
	}

	input := antlr.NewInputStream(string(b))
	lexer := parser.NewRustLexer(input)

	if to_dump_tokens {
		var vocabulary TokenVocabulary = lexer.GetSymbolicNames()
		for _, token := range lexer.GetAllTokens() {
			_, err := out.Write([]byte(vocabulary.LLVMFormat(&token)))
			if err != nil {
				panic(err)
			}
		}
	}

	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewRustParser(stream)

	sel := NewStreamErrorListener()
	p.RemoveErrorListeners()
	p.AddErrorListener(&sel)

	p.BuildParseTrees = true
	parseTree := p.Crate()

	if sel.HasErrors() {
		DumpErrors(sel.Errors(), out)
		return
	}

	builder := ast.NewANTLRRusterVisitor()
	ast := builder.Visit(parseTree).(ast.Crate)

	if to_dump_ast {
		DumpAST(ast, out)
	}

	symtabBuilder := semantics.NewANTLRSemVisitor()
	ast.Accept(symtabBuilder)

	numOfErrors := 0
	for _, log := range symtabBuilder.DumpLogs().([]semantics.Message) {
		if log.Type == semantics.INFO && !verbose {
			continue
		}
		if log.Type == semantics.ERROR {
			numOfErrors += 1
		}
		_, err := out.Write([]byte(log.String() + "\n"))
		if err != nil {
			panic(err)
		}
	}

	if numOfErrors != 0 {
		out.Write([]byte(fmt.Sprintf("Semantics analyzer found %d errors, can't continue.", numOfErrors)))
		return
	}

	llvmctx := llvmir.NewLLVMContext()
	fmt.Println(llvmctx.Visit(ast).(*ir.Module))
}

func DumpErrors(errors []Error, out io.Writer) {
	for _, e := range errors {
		_, err := out.Write([]byte(e.Format()))
		if err != nil {
			panic(err)
		}
	}
}

func DumpAST(tree ast.Crate, out io.Writer) {
	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)
	err := enc.Encode(tree)
	if err != nil {
		panic(err)
	}
}
