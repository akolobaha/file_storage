
### Установка плагина protobuf
`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

`export PATH="$PATH:$(go env GOPATH)/bin"`

### Установка taskfile 
https://taskfile.dev/installation/

### Генерация grpc кода
`task proto`

### Docker
БД в докере
`docker-cmpose up -d`

### Миграции
- Установить [goose](https://github.com/pressly/goose)
- Выполнить `task migrateUp`
- Создать `task migrateCreate name={{name}}`
- Откатить `task migrateDown`
- Статус `task migrateStatus`

### Запуск приложения
`cp .env.sample .env`
`go run cmd/main.go`

### Нагрузочные тесы
- Устнвить [ghz](https://ghz.sh/docs/install)
- `task testList`
- `task testUpload`
