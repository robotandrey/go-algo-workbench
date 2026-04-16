# go-algo-workbench

Локальный Go-workbench для решения задач LeetCode:
- автоскелет для новой задачи
- генерация `problem.md`, `solution.go`, `solution_test.go`
- запуск тестов и бенчмарков через `make`

## Требования

- Go `1.21+`
- Доступ в интернет к `https://leetcode.com/graphql` для `fetch`

## Быстрый старт

```bash
# 1) Подтянуть задачу (ВАЖНО: через переменную, не позиционным аргументом)
make fetch slug=two-sum
# или
make fetch url=https://leetcode.com/problems/two-sum/

# 2) Решить задачу в problems/two_sum/solution.go

# 3) Запустить тесты только этой задачи
make one name=two_sum

# 4) Прогнать все тесты проекта
make test
```

## Команды

```bash
make fetch slug=<leetcode-slug>
make fetch url=<leetcode-url>
make one name=<problem_dir>
make test
make bench
make fmt
```

Что важно помнить:
- `make fetch https://...` не сработает. Нужен формат `make fetch url=https://...`.
- `name` в `make one` — это имя директории в `problems/` (например `longest_common_prefix`).

## Что делает fetch

Команда `make fetch ...`:
1. Собирает бинарник `bin/fetch`.
2. Запрашивает задачу из LeetCode GraphQL API.
3. Создает директорию `problems/<slug_with_underscores>/`.
4. Генерирует файлы:
- `problem.md` — условие, ограничения, примеры.
- `solution.go` — сигнатура функции из LeetCode.
- `solution_test.go` — table-driven тесты на основе примеров.

Генератор тестов пытается:
- подставить входы из `exampleTestcaseList`
- вытащить `Output` из HTML-примеров задачи
- преобразовать литералы в Go-вид (`[1,2]` -> `[]int{1,2}` и т.д.)

Если формат конкретной задачи нестандартный, часть полей может остаться с fallback-значениями — это нормально, их нужно дописать вручную.

## Рекомендуемый workflow

1. `make fetch slug=<slug>`
2. Открыть `problems/<name>/solution.go` и реализовать решение.
3. Проверить `problems/<name>/solution_test.go`, при необходимости уточнить кейсы.
4. `make one name=<name>`
5. После стабилизации — `make test` и `make bench`.

## Структура проекта

```text
go-algo-workbench/
├── cmd/
│   └── fetch/                  # CLI для получения и генерации задач
├── internal/
│   └── leetcode/
│       ├── client.go           # GraphQL клиент LeetCode
│       ├── parse.go            # парсинг примеров/outputs
│       ├── codegen.go          # генерация solution + tests
│       ├── markdown.go         # генерация problem.md
│       └── scaffold.go         # создание директории и файлов
├── problems/
│   ├── _template/
│   └── two_sum/                # пример готовой задачи
├── bin/
├── Makefile
└── go.mod
```

## Troubleshooting

### 1) `usage: make fetch ...`
Причина: неверный синтаксис вызова `make fetch`.

Правильно:
```bash
make fetch url=https://leetcode.com/problems/longest-common-prefix/
# или
make fetch slug=longest-common-prefix
```

### 2) `context deadline exceeded` при fetch
Причина: сетевой таймаут до LeetCode API.

Проверь:
- доступность `https://leetcode.com/graphql`
- VPN/прокси/фаервол
- попробуй повторить команду через 10-30 секунд

### 3) `problem directory already exists`
Причина: задача уже была сгенерирована ранее.

Варианты:
- работать в существующей папке
- удалить папку задачи и выполнить `fetch` повторно

### 4) `make test` падает на незавершенной задаче
`make test` запускает все пакеты `./...`. Если в одной задаче недописан `return` или некомпилируемый код, упадет весь прогон.

Для локальной работы по одной задаче используй:
```bash
make one name=<problem_dir>
```
