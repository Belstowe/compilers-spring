# Теория языков программирования и методы трансляции

## Осенний семестр. 2022-2023 гг.

Ковалев Максим Игоревич, гр. ИВ-923.

## Лабораторная работа №1. Введение в LLVM IR

Изучить основы работы с LLVM IR. Разработать приложение для построения графа вызовов и def-use графа по IR.

## Сборка и запуск

Тестировалось с LLVM/Clang 14, CMake 3.21. 

Графы сериализуются в формат `dot`, который можно вывести, например, через Graphviz.

## Процесс тестирования

Ожидается использование Linux.

Компиляция:
```sh
cmake -S . -B build
cmake --build build
```

Выведение графов call и def-use для LLVM IR:
```sh
clang -S -O0 -emit-llvm assets/example_factorial_recursive.c
./build/bin/llir-graph --dot-callgraph --input example_factorial_recursive.ll | dot -Tsvg > f_r_c.svg
./build/bin/llir-graph --dot-def-use --input example_factorial_recursive.ll | dot -Tsvg > f_r_d.svg
```

То же, но через вложенные скрипты:
```sh
./to_ir.sh assets/example_factorial_recursive.c
./dot-call.sh example_factorial_recursive.ll f_r_c.svg
./dot-def.sh example_factorial_recursive.ll f_r_d.svg
```
