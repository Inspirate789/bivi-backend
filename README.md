# Backend

## Локальный запуск
```shell
go build -o app cmd/app/main.go
./app -c env/app.env
```

## Запуск в Docker

TODO

## Локальный запуск тестов

TODO

## Документация к API

Адрес документации Swagger: http://localhost:8080/swagger/index.html

### Генерация документации
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

```shell
swag fmt
swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/app/main.go -o swagger/
```

