@echo off

set antlr_ver=4.10.1
set antlr=antlr-%antlr_ver%-complete.jar
set antlr_path=%CD%\%antlr%
set output=%CD%\parser

if not exist "%output%\" mkdir %output%

if not exist "%antlr_path%" curl -O https://www.antlr.org/download/%antlr%

go get -d github.com/antlr/antlr4/runtime/Go/antlr@%antlr_ver%
pushd build\grammar
java -jar %antlr_path% -o %output% -Dlanguage=Go RustLexer.g4
popd