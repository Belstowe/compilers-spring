lexer grammar RustLexer;

KW_AS:          'as';
KW_BREAK:       'break';
KW_CONST:       'const';
KW_CONTINUE:    'continue';
KW_CRATE:       'crate';
KW_ELSE:        'else';
KW_FALSE:       'false';
KW_FN:          'fn';
KW_FOR:         'for';
KW_IF:          'if';
KW_IN:          'in';
KW_MUT:         'mut';
KW_RETURN:      'return';
KW_STATIC:      'static';
KW_STRUCT:      'struct';
KW_SUPER:       'super';
KW_TRUE:        'true';
KW_TYPE:        'type';
KW_USE:         'use';
KW_WHERE:       'where';
KW_WHILE:       'while';

ID: [a-zA-Z][a-zA-Z0-9_]* | '_' [a-zA-Z0-9_]+;

LINE_COMMENT: ('//' ~[\r\n]* | '//') -> skip;
WHITESPACE: [\p{Zs}] -> skip;
NEWLINE: ('\r\n' | [\r\n]) -> skip;

CHAR_LITERAL : '\'' ( ~['\\\n\r\t] | '\\' ['nrt\\0] ) '\'' ;

STRING_LITERAL : '"' ( ~["] | '\\' ["nrt\\0] )* '"' ;

INTEGER_LITERAL : DEC_LITERAL INTEGER_SUFFIX?;

DEC_LITERAL: DEC_DIGIT (DEC_DIGIT | '_')*;

fragment INTEGER_SUFFIX
   : 'u8' | 'u16' | 'u32' | 'u64' | 'u128' | 'usize'
   | 'i8' | 'i16' | 'i32' | 'i64' | 'i128' | 'isize' ;

fragment DEC_DIGIT: [0-9];

PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
PERCENT: '%';

CARET: '^';
NOT: '!';
AND: '&';
OR: '|';
ANDAND: '&&';
OROR: '||';
SHL: '<<';
SHR: '>>';
EQ: '=';
EQEQ: '==';
NE: '!=';
GT: '>';
LT: '<';
GE: '>=';
LE: '<=';

PLUSEQ: '+=';
MINUSEQ: '-=';
STAREQ: '*=';
SLASHEQ: '/=';
PERCENTEQ: '%=';

CARETEQ: '^=';
ANDEQ: '&=';
OREQ: '|=';
SHLEQ: '<<=';
SHREQ: '>>=';

UNDERSCORE: '_';
DOT: '.';
AT: '@';
COMMA: ',';
SEMI: ';';
COLON: ':';
PATHSEP: '::';
RARROW: '->';
POUND: '#';
DOLLAR: '$';
QUESTION: '?';

LCURLYBRACE: '{';
RCURLYBRACE: '}';
LSQUAREBRACKET: '[';
RSQUAREBRACKET: ']';
LPAREN: '(';
RPAREN: ')';