# Backend

[Правила кодирования](docs/CODING_RULES.md)

## Запуск приложения
#### Локально
```shell
# Без указания кастомного конфига
go run ./cmd/app/main.go

# С кастомным конфигом
go build -o app cmd/app/main.go
./app -c env/app.local.yaml
```
#### В Docker
```shell
make docker-app # ARCH=arm64 if you use arm-based PC
docker run --name bivi-backend -d -p 8080:80 -v ./content:/content -v ./logs:/logs bivi/backend:local
```

## Запуск тестов
#### Локально
```shell
./scripts/e2e-test.sh # Использует конфиг app.test.yaml
```
Посмотреть отчёты allure:
```shell
allure serve test-reports/allure-results
```
#### В Docker
```shell
make docker-e2e-test # ARCH=arm64 if you use arm-based PC
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
make docker-lint # ARCH=arm64 if you use arm-based PC
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

