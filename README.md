# Transcoder Framework

Простой и мощный фреймворк для транскодирования медиафайлов на Go с использованием FFmpeg.

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

- `web-hd` - HD качество для веб (1280x720, H.264)
- `web-sd` - SD качество для веб (854x480, H.264)
- `mobile` - Оптимизировано для мобильных устройств
- `audio-mp3` - Конвертация в MP3
- `audio-aac` - Конвертация в AAC

## API Reference

### Transcoder

- `New(ffmpegPath string) (*Transcoder, error)` - создает новый транскодер
- `CreateJob(config Config) *Job` - создает задачу транскодирования
- `Execute(ctx context.Context, job *Job) error` - выполняет транскодирование
- `GetInfo(filePath string) (map[string]interface{}, error)` - получает информацию о файле

### Queue

- `NewQueue(transcoder *Transcoder, workers int) *Queue` - создает очередь
- `AddJob(job *Job)` - добавляет задачу в очередь
- `Start()` - запускает обработку очереди
- `Stop()` - останавливает очередь

## Лицензия

MIT License