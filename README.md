# Теория языков программирования и методы трансляции

## Весенний семестр. 2021-2022 гг.

Ковалев Максим Игоревич, гр. ИВ-923.

## Лабораторная работа №3. Синтаксический анализ

* **Парсер для языка:** Rust *(вариант 14)*.
* **Язык рантайма:** Go.

## Сборка и запуск

Для проекта необходимо установить компилятор **[Go](https://go.dev/dl/)**.

1. Установка ANTLR4 Go Runtime *(пакет с поддерживаемой версией ANTLR4 закачивается автоматически)*:
```
Unix    $ ./scripts/build.sh
Windows > .\scripts\build.bat
```

2. Сборка компилятора Rust:
```sh
$ go build ./cmd/ruster/ruster.go
```

3. Тестирование компилятора Rust:
```sh
$ go test ./test/...
```

Чтобы получить инструкцию по компилятору, пропишите флаг `-h`/`--help`:
```sh
$ ruster --help
```
```sh
NAME:
   ruster - A simple Rust compiler using ANTLR

USAGE:
   ruster [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value  Path to Rust code file for parsing (default: read from terminal)
   --dump-tokens            Require lexer to dump tokens in stdout (default: false)
   --dump-ast               Require parser to dump AST in stdout (default: false)
   --help, -h               show help (default: false)
```