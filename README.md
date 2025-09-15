# Doc Analyzer

Система на Go для загрузки, хранения и анализа документов.  
Архитектура построена на микросервисах:

- **Gateway** — принимает клиентские запросы и маршрутизирует их.
- **Storage** — сервис для сохранения файлов.
- **Analyzer** — сервис для анализа текста.

Сервисы общаются через **gRPC**. Для запуска используется **Docker Compose**, для сборки и тестирования — `Makefile`.

## Запуск

### Вручную
```bash
go build -o bin/gateway ./cmd/gateway
go build -o bin/storage ./cmd/storage
go build -o bin/analyzer ./cmd/analyzer

./bin/gateway &
./bin/storage &
./bin/analyzer &
```

### Docker Compose
```bash
docker-compose up --build
```

## Пример

1. Загрузка файла:
   ```
   curl -F "file=@sample.txt" http://localhost:8080/upload
   ```
   → `{"file_id": "12345"}`

2. Анализ файла:
   ```
   curl http://localhost:8080/analyze/12345
   ```
   → `{"words": 350, "unique_words": 120}`

## Технологии

- Go 1.21+
- gRPC, net/http
- PostgreSQL
- Docker, Docker Compose
# doc-analyzer
