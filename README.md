# Теория языков программирования и методы трансляции

## Весенний семестр. 2021-2022 гг.

Ковалев Максим Игоревич, гр. ИВ-923.

## Лабораторная работа №2. Лексический анализ

* **Парсер для языка:** Rust *(вариант 14)*.
* **Язык рантайма:** Go.

## Сборка и запуск

Для проекта необходимо установить компилятор **[Go](https://go.dev/dl/)**.

1. Установка ANTLR4 Go Runtime *(пакет с ANTLR4 закачивается автоматически)*:
```
Unix    $ ./scripts/build.sh
Windows > .\scripts\build.bat
```

2. Сборка лексера Rust:
```sh
$ go build ./cmd/rust-parser/rust-parser.go
```

3. Тестирование лексера Rust:
```sh
$ go test ./test/...
```

Чтобы получить инструкцию по компилятору, пропишите флаг `-h`/`--help`:
```sh
$ rust-parser --help
```
```sh
NAME:
   rust-parser - A simple Rust compiler using ANTLR

USAGE:
   rust-parser [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value  Path to Rust code file for parsing (default: read from terminal)
   --dump-tokens            Require lexer to dump tokens in stdout (default: false)
   --help, -h               show help (default: false)
```