# Правила кодирования

## Требования к оформлению кода
Общепринятые соглашения: https://go.dev/doc/effective_go

## Антипаттерны
- Изобретать велосипед.
- Смешивать бизнес-логику и адаптеры для её реализации.
- Использовать рефлексию там, где её можно было бы избежать.
- Использовать `any` там, где можно использовать дженерик.
- Не преаллоцировать слайсы и мапы, если есть такая возможность. Пример:
    ```go
    // Неэффективно: 
    slice := make([]int, 0)
    for i := 0; i < 10; i++ {
        slice = append(slice, i)
    }
  
    // Эффективно:
    slice := make([]int, 0, 10)
    for i := 0; i < 10; i++ {
        slice = append(slice, i)
    }
    ```

## Структура каталогов проекта
Макет со стандартными названиями каталогов, на который нужно ориентироваться: https://github.com/golang-standards/project-layout

## Правила структуризации кода
- Все зависимости объединяются в пакете main.
- Пакеты группируются по зависимостям.
- Все доменные типы должны быть в отдельном (корневом) пакете.
- Работать с внешними зависимостями через интерфейсы.
- Код приложения в `internal` (за исключением `internal/pkg`) делится на горизонтальные слои
  - адаптеры для внешних зависимостей;
  - прикладные бизнес-правила (use-case);
  - слой доступа к данным (как правило, по паттерну repository);
  
  и вертикальные срезы по группам use-case (как правило, связанным с одной доменной сущностью).

## Пример структуры проекта:
```
.
├── assets
│   ├── influx.png
│   └── swagger.png
├── cmd
│   └── app
│       └── main.go
├── db
│   └── init.sql
├── docker-compose.yaml
├── Dockerfile
├── env
│   └── app.env
├── go.mod
├── go.sum
├── internal
│   ├── models
│   │   └── segment_event.go
│   ├── pkg
│   │   └── app
│   │       ├── fiber_app.go
│   │       └── interface.go
│   ├── segment
│   │   ├── delivery
│   │   │   ├── delivery.go
│   │   │   └── interfaces.go
│   │   ├── repository
│   │   │   ├── queries.go
│   │   │   └── repository.go
│   │   └── usecase
│   │       ├── dto
│   │       │   └── dto.go
│   │       ├── errors
│   │       │   └── errors.go
│   │       ├── interfaces.go
│   │       └── usecase.go
│   └── user
│       ├── delivery
│       │   ├── delivery.go
│       │   └── interfaces.go
│       ├── repository
│       │   ├── fs
│       │   │   └── repository.go
│       │   └── sql
│       │       ├── queries.go
│       │       └── repository.go
│       └── usecase
│           ├── dto
│           │   └── dto.go
│           ├── errors
│           │   └── errors.go
│           ├── interfaces.go
│           └── usecase.go
├── makefile
├── pkg
│   ├── influx
│   │   └── influx.go
│   └── sqlx_utils
│       └── utils.go
├── readme.md
├── swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── task.md
```
