# 🎬 Transcoder Framework

Профессиональный фреймворк для транскодирования медиафайлов на Go с использованием FFmpeg.

**Версия 2.0** - Полностью переработанная архитектура с расширенными возможностями

📖 **[Полный справочник FFmpeg](FFMPEG_REFERENCE.md)** - подробное руководство по всем возможностям FFmpeg

## ✨ Возможности

### 🎯 Основные функции
- 🎥 **Транскодирование** видео и аудио файлов с валидацией
- 📺 **HLS стримы** - загрузка и обработка .m3u8 плейлистов
- 🔴 **Live стримы** - запись с ограничением по времени
- 🔧 **Гибкая настройка** параметров кодирования
- 📋 **19 пресетов** для популярных форматов и сценариев
- 🔄 **Система очередей** с многопоточностью

### 🚀 Новые возможности v2.0
- 📊 **Расширенный анализ** медиафайлов с структурированными данными
- 🎛️ **Система фильтров** FFmpeg (масштабирование, яркость, звук)
- 🔗 **Конвейеры обработки** для сложных операций
- 📈 **Отслеживание прогресса** в реальном времени
- 📝 **Профессиональное логирование** с уровнями
- ✅ **Автоматическая валидация** всех параметров
- 🏗️ **Модульная архитектура** DTO/Utils/Presets

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

## 🚀 Быстрый старт

### Базовое использование
```go
package main

import (
    "context"
    "log"
    
    transcoder "github.com/Mirsadikovv/ffmpeg_research"
    "github.com/Mirsadikovv/ffmpeg_research/dto"
)

func main() {
    // Создаем транскодер с автоматическим логированием
    tc, err := transcoder.New("ffmpeg")
    if err != nil {
        log.Fatal(err)
    }

    // Настраиваем параметры с автоматической валидацией
    config := dto.Config{
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

### Расширенный пример с прогрессом
```go
// Отслеживание прогресса в реальном времени
progressCallback := func(progress float64, speed string, eta time.Duration) {
    fmt.Printf("\rПрогресс: %.1f%% | Скорость: %s | ETA: %v", 
        progress, speed, eta)
}

tracker := transcoder.NewProgressTracker(progressCallback)
// Используйте tracker для мониторинга выполнения
```

## 🎛️ Новые возможности v2.0

### Система фильтров FFmpeg
```go
// Создаем цепочку фильтров
filterChain := transcoder.NewFilterChain().
    AddVideoFilter("scale", map[string]string{"w": "1920", "h": "1080"}).
    AddVideoFilter("eq", map[string]string{"brightness": "0.1", "contrast": "1.2"}).
    AddAudioFilter("volume", map[string]string{"volume": "1.2"})

// Применяем фильтры
config := dto.Config{
    InputPath:  "input.mp4",
    OutputPath: "output_filtered.mp4",
    VideoCodec: "libx264",
    AudioCodec: "aac",
}

job := tc.CreateJob(config)
err := tc.ExecuteWithFilters(ctx, job, filterChain)
```

### Конвейеры обработки (Pipelines)
```go
// Создаем веб-оптимизационный конвейер
pipeline := transcoder.CreateWebOptimizationPipeline(tc)

// Выполняем полный конвейер обработки
err := pipeline.Execute(ctx, "input.mp4", "output_web.mp4")

// Доступные конвейеры:
// - CreateWebOptimizationPipeline() - веб-оптимизация
// - CreateMobilePipeline() - мобильная оптимизация  
// - CreateArchivePipeline() - архивирование
```

### Расширенный анализ медиафайлов
```go
// Получаем подробную информацию
mediaInfo, err := tc.GetMediaInfo("video.mp4")
if err == nil {
    fmt.Printf("Разрешение: %s\n", mediaInfo.GetResolution())
    fmt.Printf("Частота кадров: %.1f fps\n", mediaInfo.GetFrameRate())
    fmt.Printf("Продолжительность: %.1f сек\n", mediaInfo.Duration.Seconds())
    fmt.Printf("Краткая сводка: %s\n", mediaInfo.Summary())
    
    // Анализ потоков
    videoStreams := mediaInfo.GetVideoStreams()
    audioStreams := mediaInfo.GetAudioStreams()
}
```

### Профессиональное логирование
```go
// Настройка уровня логирования
logger := transcoder.NewDefaultLogger(transcoder.LogLevelDebug)
tc.SetLogger(logger)

// Или отключение логов в production
tc.SetLogger(&transcoder.NoOpLogger{})
```

## 📋 Использование пресетов (19 штук!)

### Получение пресета
```go
// Получаем пресет из новой системы
preset, exists := presets.GetPreset("web-hd")
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

### Пресеты по категориям
```go
// Получаем пресеты по категориям
categories := presets.GetPresetsByCategory()

for category, categoryPresets := range categories {
    fmt.Printf("%s:\n", category)
    for _, preset := range categoryPresets {
        fmt.Printf("  - %s: %s\n", preset.Name, preset.Description)
    }
}
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

## Работа с HLS стримами

### Загрузка HLS плейлиста
```go
// Простая загрузка HLS стрима
ctx := context.Background()
err := tc.DownloadHLS(ctx, "https://example.com/playlist.m3u8", "stream.mp4")
```

### Расширенная конфигурация HLS
```go
config := transcoder.HLSConfig{
    URL:        "https://example.com/playlist.m3u8",
    OutputPath: "stream_hd.mp4",
    Quality:    "1920x1080", // или "best", "worst"
    Duration:   30 * time.Minute, // ограничение по времени
    Headers: map[string]string{
        "Referer":    "https://example.com",
        "User-Agent": "Custom User Agent",
    },
    RetryAttempts:  3,
    SegmentTimeout: 10 * time.Second,
}

err := tc.DownloadHLSWithConfig(ctx, config)
```

### Получение информации о плейлисте
```go
info, err := tc.GetHLSInfo("https://example.com/playlist.m3u8")
if err == nil {
    fmt.Printf("Найдено %d потоков\n", len(info.Streams))
    fmt.Printf("Live стрим: %v\n", info.IsLive)
    
    for _, stream := range info.Streams {
        fmt.Printf("Поток: %s, %d kbps\n", stream.Resolution, stream.Bandwidth/1000)
    }
}
```

### Запись live стрима
```go
// Записать live стрим в течение 10 минут
duration := 10 * time.Minute
err := tc.RecordLiveStream(ctx, "https://example.com/live.m3u8", "live_record.mp4", duration)
```

### Конвертация HLS в другой формат
```go
// Загрузить HLS и сконвертировать в WebM
err := tc.ConvertHLSToFormat(ctx, "https://example.com/playlist.m3u8", "stream.webm", "webm")
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