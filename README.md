# Backend

[Правила кодирования](docs/CODING_RULES.md)

## Локальный запуск
```shell
go build -o app cmd/app/main.go
./app -c env/app.env
```

## Запуск в Docker

TODO

## Локальный запуск тестов

TODO

## Запуск линтера
Установка:
```shell
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2```
```
Запуск:
```shell
golangci-lint run
```

## Документация к API

Адрес документации Swagger запущенного приложения: http://localhost:8080/swagger/index.html

### Генерация документации
Установка:
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```
Запуск:
```shell
swag fmt
swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/app/main.go -o swagger/
```

