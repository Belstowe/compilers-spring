parser grammar RustParser;

options
{
    tokenVocab = RustLexer;
}

crate : item* EOF ;

item
    : useDeclaration
    | function
    | struct
    | typeAlias
    | constantItem ;

useDeclaration : 'use' useTree ';' ;
useTree
    : (simplePath? '::')? IncludeAll='*'
    | simplePath ('as' (identifier | '_'))? ;

function
    : 'fn' identifier '(' functionParameters? ')' functionReturnType? (blockExpression | ';') ;
functionParameters: functionParam (',' functionParam)* ','? ;
functionParam : ((identifier | '_') ':')? type ;
functionReturnType: '->' type ;

struct : 'struct' identifier ('{' structFields? '}' | ';') ;
structFields : structField (',' structField)* ','? ;
structField : identifier ':' type ;

typeAlias : 'type' identifier ('=' type)? ';' ;

constantItem : 'const' (identifier | '_') ':' type ('=' expression)? ';' ;

statement
	: ';'
	| item
	| letStatement
	| expressionStatement ;

letStatement : 'let' pattern (':' type)? ('=' expression)? ';' ;

expressionStatement
	: expression ';'
	| expressionWithBlock ';'? ;

expression
	: literalExpression									                        # LiteralExpression_
	| pathExpression										                    # PathExpression_
	| expression '.' simplePathSegment '(' callParams? ')'                      # MethodCallExpression
	| expression '.' identifier								                    # FieldExpression
	| expression '(' callParams? ')'						                    # CallExpression
	| expression '.' tupleIndex								                    # TupleIndexingExpression
	| expression '[' expression ']'							                    # IndexExpression
	| expression '?'                                                            # ErrorPropagationExpression
	| RefToken=('&' | '&&') MutToken='mut'? expression                          # BorrowExpression
	| '*' expression										                    # DereferenceExpression
	| Op=('-' | '!') Val=expression								                # UnaryOpExpression
	| expression 'as' type									                    # TypeCastExpression
	| LHS=expression Op=('*' | '/' | '%' | '+' | '-') RHS=expression            # BinaryOpExpression
	| LHS=expression Op=('<<' | '>>') RHS=expression					        # BinaryOpExpression
	| LHS=expression Op=('&' | '^' | '|') RHS=expression				        # BinaryOpExpression
	| LHS=expression Op=('==' | '!=' | '>' | '<' | '>=' | '<=') RHS=expression	# BinaryOpExpression
	| LHS=expression Op=('&&' | '||') RHS=expression					        # BinaryOpExpression
	| LHS=expression Op=('..' | '..=') RHS=expression?					        # RangeExpression
	| Op=('..' | '..=') Val=expression?								            # RHSRangeExpression
	| LHS=expression '=' RHS=expression							                # AssignmentExpression
	| LHS=expression compoundAssignOperator RHS=expression			            # CompoundAssignmentExpression
	| 'continue' expression?								                    # ContinueExpression
	| 'break' expression?									                    # BreakExpression
	| 'return' expression?									                    # ReturnExpression
	| '(' expression ')'									                    # NestedExpression
	| '[' arrayElements? ']'								                    # ArrayExpression
	| '(' tupleElements? ')'								                    # TupleExpression
	| structExpression										                    # StructExpression_
	| expressionWithBlock									                    # ExpressionWithBlock_
	;

compoundAssignOperator
	: '+=' | '-=' | '*=' | '/=' | '%=' | '&=' | '|=' | '^=' | '<<=' | '>>=' ;

literalExpression
	: Literal=(CHAR_LITERAL
	| STRING_LITERAL
	| INTEGER_LITERAL
	| KW_TRUE
	| KW_FALSE) ;

blockExpression : '{' statements? '}' ;
statements
	: statement+ expression?
	| expression ;

arrayElements
	: expression (',' expression)* ','?
	| expression ';' expression ;

tupleElements : (expression ',')+ expression? ;
tupleIndex : INTEGER_LITERAL ;

structExpression : pathExpression '{' structExprFields? '}' ;
structExprFields : structExprField (',' structExprField)* ','? ;
structExprField : (identifier | (identifier | tupleIndex) ':' expression) ;

callParams : expression (',' expression)* ','? ;

pathExpression : '::'? simplePathSegment ('::' simplePathSegment)* ;

expressionWithBlock
	: loopExpression
	| ifExpression
	| matchExpression ;

ifExpression
   : 'if' Cond=expression IfBody=blockExpression
   (
      'else' (ElseBody=blockExpression | ElseIf=ifExpression)
   )? ;

matchExpression : 'match' expression '{' matchArms? '}' ;
matchArms : (matchArm '=>' matchArmExpression)* matchArm '=>' expression ','? ;
matchArmExpression
    : expression ','
    | expressionWithBlock ','? ;
matchArm :
    '|'? pattern ('|' pattern)* ;

loopExpression
	: infiniteLoopExpression
	| predicateLoopExpression
	| iteratorLoopExpression ;

infiniteLoopExpression : 'loop' blockExpression ;
predicateLoopExpression : 'while' expression blockExpression ;
iteratorLoopExpression : 'for' pattern 'in' expression blockExpression;

pattern
	: nonRangePattern
	| rangePattern ;

nonRangePattern
	: literalPattern
	| identifierPattern
	| wildcardPattern
	| restPattern
	| referencePattern
	| structPattern
	| tuplePattern
	| groupedPattern
	| slicePattern
	| pathPattern
    | rangePattern ;

literalPattern
	: Literal=(KW_TRUE
	| KW_FALSE
	| CHAR_LITERAL
	| STRING_LITERAL
	| INTEGER_LITERAL) ;

identifierPattern
	: 'ref'? 'mut'? identifier ;

wildcardPattern : '_' ;
restPattern : '..' ;

referencePattern : ('&' | '&&') 'mut'? nonRangePattern ;

structPattern : pathExpression '{' structPatternElements? '}' ;
structPatternElements
    : structPatternFields (',' '..'?)?
    | '..' ;
structPatternFields : structPatternField (',' structPatternField)* ;
structPatternField
    : tupleIndex ':' pattern
    | identifier ':' pattern
    | 'ref'? 'mut'? identifier ;

tuplePattern : '(' tuplePatternItems? ')' ;
tuplePatternItems
   : pattern ','
   | restPattern
   | pattern (',' pattern)+ ','? ;

groupedPattern : '(' pattern ')' ;

slicePattern : '[' slicePatternItems? ']' ;
slicePatternItems : pattern (',' pattern)* ','? ;

pathPattern : pathExpression ;

rangePattern : rangePatternBound '..=' rangePatternBound ;
rangePatternBound
   : CHAR_LITERAL
   | '-'? INTEGER_LITERAL
   | pathExpression ;

type
	: nestedType
	| typePath
	| tupleType
	| neverType
	| pointerType
	| referenceType
	| arrayType
	| sliceType
	| inferredType
	| functionType ;

nestedType : '(' type ')' ;

neverType : '!' ;
inferredType : '_' ;

tupleType : '(' ((type ',')+ type?)? ')' ;
arrayType : '[' type ';' expression ']' ;
sliceType : '[' type ']' ;

referenceType : '&' 'mut'? type ;
pointerType : '*' ('mut' | 'const') type ;

functionType: 'fn' '(' functionParameters? ')' functionReturnType? ;

typePath : '::'? typePathSegment ('::' typePathSegment)* ;
typePathSegment : simplePathSegment '::'? typePathFn? ;
typePathFn : '(' typePathInputs? ')' ('->' ReturnType=type)? ;
typePathInputs : type (',' type)* ','? ;

simplePath : '::'? simplePathSegment ('::' simplePathSegment)* ;
simplePathSegment : identifier | 'super' | 'self' ;

identifier : ID ;

keyword
    : KW_AS
    | KW_BREAK
    | KW_CONST
    | KW_CONTINUE
    | KW_CRATE
    | KW_ELSE
    | KW_FALSE
    | KW_FN
    | KW_FOR
    | KW_IF
    | KW_IN
    | KW_MUT
    | KW_RETURN
    | KW_STATIC
    | KW_STRUCT
    | KW_SUPER
    | KW_TRUE
    | KW_TYPE
    | KW_USE
    | KW_WHERE
    | KW_WHILE ;