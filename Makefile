.PHONY: build test clean run example deps fmt vet lint

# Переменные
BINARY_NAME=transcoder
GO_FILES=$(shell find . -name "*.go" -type f)

# Сборка
build:
	go build -o bin/$(BINARY_NAME) ./example

# Запуск примера
run: build
	./bin/$(BINARY_NAME)

# Запуск примера напрямую
example:
	go run example/main.go

# Запуск расширенного примера
example-advanced:
	go run example/advanced_example.go

# Тесты
test:
	go test -v ./...

# Тесты с покрытием
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Бенчмарки
bench:
	go test -bench=. -benchmem ./...

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Форматирование кода
fmt:
	go fmt ./...

# Проверка кода
vet:
	go vet ./...

# Линтер (требует golangci-lint)
lint:
	golangci-lint run

# Очистка
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Проверка FFmpeg
check-ffmpeg:
	@which ffmpeg > /dev/null || (echo "FFmpeg не найден. Установите FFmpeg для работы с фреймворком." && exit 1)
	@which ffprobe > /dev/null || (echo "FFprobe не найден. Установите FFmpeg для работы с фреймворком." && exit 1)
	@echo "FFmpeg найден: $$(ffmpeg -version | head -n1)"

# Полная проверка
check: check-ffmpeg fmt vet test

# Установка инструментов разработки
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Создание релиза
release: clean check build
	@echo "Релиз готов в bin/$(BINARY_NAME)"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build         - Собрать проект"
	@echo "  run           - Собрать и запустить пример"
	@echo "  example       - Запустить пример напрямую"
	@echo "  example-advanced - Запустить расширенный пример"
	@echo "  test          - Запустить тесты"
	@echo "  test-coverage - Запустить тесты с покрытием"
	@echo "  bench         - Запустить бенчмарки"
	@echo "  deps          - Установить зависимости"
	@echo "  fmt           - Форматировать код"
	@echo "  vet           - Проверить код"
	@echo "  lint          - Запустить линтер"
	@echo "  clean         - Очистить сборочные файлы"
	@echo "  check-ffmpeg  - Проверить наличие FFmpeg"
	@echo "  check         - Полная проверка проекта"
	@echo "  install-tools - Установить инструменты разработки"
	@echo "  release       - Создать релиз"
	@echo "  help          - Показать эту справку"