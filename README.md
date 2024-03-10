# Backend

[Правила кодирования](docs/CODING_RULES.md)

## Запуск приложения
#### Локально
```shell
# Без указания кастомного конфига
go run ./cmd/app/main.go

# С кастомным конфигом
go build -o app cmd/app/main.go
./app -c env/app.local.env
```
#### В Docker
```shell
make docker-app # ARCH=arm64 if you use arm-based Mac
docker run --name bivi-backend -d -p 8080:80 -v ./content:/content bivi/backend:local
```

## Запуск тестов
#### Локально
```shell
./scripts/integration-test.sh
```
#### В Docker
```shell
make docker-integration-test # ARCH=arm64 if you use arm-based Mac
```

## Запуск линтера
#### Локально
Установка:
```shell
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2```
```
Запуск:
```shell
golangci-lint run
```
#### В Docker
```shell
make docker-lint # ARCH=arm64 if you use arm-based Mac
```

## Документация к API

Путь к документации Swagger запущенного приложения: /swagger/index.html

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

