# Теория языков программирования и методы трансляции

## Весенний семестр. 2021-2022 гг.

Ковалев Максим Игоревич, гр. ИВ-923.

## Лабораторная работа №4. Таблица символов

* **Парсер для языка:** Rust *(вариант 14)*.
* **Язык рантайма:** Go.

## Сборка и запуск

### Docker

Для постройки и входа в тестовый Docker-образ в одну команду созданы следующие скрипты:
```
Unix    $ ./scripts/docker-run.sh
Windows > .\scripts\docker-run.bat
```

Полезными командами будут далее:
```
$ go test ./test/...                               # Запуск всех тестов (включая тесты на некорректный код)
$ ./ruster --dump-tokens -i ./examples/min.rs      # Вывод лексических токенов (лабораторная работа #2)
$ ./ruster --dump-ast -i ./examples/gcd.rs         # Вывод AST-дерева (лабораторная работа #3)
$ ./ruster --verbose -i ./examples/find_substr.rs  # Вывод информационных сообщений семантического визитора (лабораторная работа #4)
```

### Стандарт

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
   --verbose, -v            Print info messages as well (default: false)
   --help, -h               show help (default: false)
```
