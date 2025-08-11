# Transcoder Framework

Простой и мощный фреймворк для транскодирования медиафайлов на Go с использованием FFmpeg.

📖 **[Полный справочник FFmpeg](FFMPEG_REFERENCE.md)** - подробное руководство по всем возможностям FFmpeg

## Возможности

- 🎥 Транскодирование видео и аудио файлов
- 🔧 Гибкая настройка параметров кодирования
- 📋 Предустановленные пресеты для популярных форматов
- 🔄 Система очередей с поддержкой многопоточности
- 📊 Получение информации о медиафайлах
- 🚀 Простой и понятный API

## Установка

```bash
go get github.com/your-username/transcoder
```

Убедитесь, что FFmpeg установлен в вашей системе:

```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# macOS
brew install ffmpeg

# Windows
# Скачайте с https://ffmpeg.org/download.html
```

## Быстрый старт

```go
package main

import (
    "context"
    "log"
    "transcoder"
)

func main() {
    // Создаем транскодер
    tc, err := transcoder.New("ffmpeg")
    if err != nil {
        log.Fatal(err)
    }

    // Настраиваем параметры
    config := transcoder.Config{
        InputPath:    "input.mp4",
        OutputPath:   "output.mp4",
        VideoCodec:   "libx264",
        AudioCodec:   "aac",
        VideoBitrate: "1000k",
        AudioBitrate: "128k",
    }

    // Создаем и выполняем задачу
    job := tc.CreateJob(config)
    ctx := context.Background()
    
    if err := tc.Execute(ctx, job); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Транскодирование завершено!")
}
```

## Использование пресетов

```go
// Получаем пресет
preset, exists := transcoder.GetPreset("web-hd")
if !exists {
    log.Fatal("Пресет не найден")
}

// Настраиваем пути
config := preset.Config
config.InputPath = "input.mp4"
config.OutputPath = "output_hd.mp4"

// Выполняем транскодирование
job := tc.CreateJob(config)
tc.Execute(ctx, job)
```

## Работа с очередью

```go
// Создаем очередь с 2 воркерами
queue := transcoder.NewQueue(tc, 2)
queue.Start()

// Добавляем задачи
for i := 0; i < 5; i++ {
    config := transcoder.Config{
        InputPath:  fmt.Sprintf("input%d.mp4", i),
        OutputPath: fmt.Sprintf("output%d.mp4", i),
        VideoCodec: "libx264",
        AudioCodec: "aac",
    }
    job := tc.CreateJob(config)
    queue.AddJob(job)
}

// Останавливаем очередь
defer queue.Stop()
```

## Доступные пресеты

### Веб и мобильные
- `web-hd` - HD качество для веб (1280x720, H.264)
- `web-sd` - SD качество для веб (854x480, H.264)
- `mobile` - Оптимизировано для мобильных устройств

### Высокое качество
- `4k` - 4K качество (3840x2160, H.264)
- `archive` - Высокое качество для архивирования (H.265, FLAC)

### Стриминг
- `twitch` - Оптимизировано для Twitch стрима (1080p60)
- `youtube` - Оптимизировано для YouTube (1080p30)

### Аудио
- `audio-mp3` - Конвертация в MP3
- `audio-aac` - Конвертация в AAC

## Дополнительные возможности

### Быстрая конвертация
```go
// Конвертация в популярные форматы одной командой
tc.ConvertToFormat("input.mp4", "output.webm", "webm")
tc.ConvertToFormat("video.mov", "audio.mp3", "mp3")
```

### Извлечение аудио
```go
// Извлечь аудиодорожку из видео
tc.ExtractAudio("video.mp4", "audio.mp3")
```

### Создание миниатюр
```go
// Создать миниатюру на 10-й секунде
tc.CreateThumbnail("video.mp4", "thumb.jpg", "00:00:10")
```

### Получение информации
```go
// Узнать продолжительность файла
duration, err := tc.GetDuration("video.mp4")
```

## API Reference

### Transcoder

#### Основные методы
- `New(ffmpegPath string) (*Transcoder, error)` - создает новый транскодер
- `CreateJob(config Config) *Job` - создает задачу транскодирования
- `Execute(ctx context.Context, job *Job) error` - выполняет транскодирование
- `GetInfo(filePath string) (map[string]interface{}, error)` - получает информацию о файле

#### Удобные методы
- `ConvertToFormat(inputPath, outputPath, format string) error` - быстрая конвертация
- `ExtractAudio(inputPath, outputPath string) error` - извлечение аудио
- `CreateThumbnail(inputPath, outputPath, timeOffset string) error` - создание миниатюры
- `GetDuration(filePath string) (string, error)` - получение продолжительности

### Queue

- `NewQueue(transcoder *Transcoder, workers int) *Queue` - создает очередь
- `AddJob(job *Job)` - добавляет задачу в очередь
- `Start()` - запускает обработку очереди
- `Stop()` - останавливает очередь

## Разработка

### Makefile команды
```bash
make help          # Показать все доступные команды
make check-ffmpeg  # Проверить наличие FFmpeg
make example       # Запустить пример
make test          # Запустить тесты
make test-coverage # Тесты с покрытием кода
make bench         # Бенчмарки
make fmt           # Форматировать код
make check         # Полная проверка проекта
```

### Тестирование
```bash
# Запуск всех тестов
go test -v ./...

# Тесты с покрытием
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Бенчмарки
go test -bench=. -benchmem ./...
```

### Структура проекта
```
transcoder/
├── transcoder.go      # Основной модуль
├── presets.go         # Предустановленные конфигурации
├── queue.go           # Система очередей
├── transcoder_test.go # Тесты
├── example/
│   └── main.go        # Примеры использования
├── FFMPEG_REFERENCE.md # Справочник FFmpeg
├── Makefile           # Команды сборки
└── README.md          # Документация
```

## Лицензия

MIT License