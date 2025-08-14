package main

import (
	"context"
	"fmt"
	"log"
	"time"

	transcoder "github.com/Mirsadikovv/ffmpeg_research"
	"github.com/Mirsadikovv/ffmpeg_research/dto"
	"github.com/Mirsadikovv/ffmpeg_research/presets"
)

func main() {
	// Создаем транскодер
	tc, err := transcoder.New("ffmpeg")
	if err != nil {
		log.Fatal("Ошибка создания транскодера:", err)
	}

	ctx := context.Background()

	// Пример 1: Простое транскодирование
	fmt.Println("=== Простое транскодирование ===")
	config := dto.Config{
		InputPath:    "input.mp4",
		OutputPath:   "output.mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
	}

	job := tc.CreateJob(config)
	fmt.Printf("Создана задача: %s\n", job.ID)

	if err := tc.Execute(ctx, job); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Транскодирование завершено за %v\n", job.EndTime.Sub(job.StartTime))
	}

	// Пример 2: Загрузка HLS плейлиста
	fmt.Println("\n=== Загрузка HLS плейлиста ===")
	hlsURL := "https://example.com/playlist.m3u8"
	outputPath := "downloaded_stream.mp4"

	if err := tc.DownloadHLS(ctx, hlsURL, outputPath); err != nil {
		fmt.Printf("Ошибка загрузки HLS: %v\n", err)
	} else {
		fmt.Println("HLS стрим успешно загружен")
	}

	// Пример 3: Получение информации о HLS плейлисте
	fmt.Println("\n=== Информация о HLS плейлисте ===")
	info, err := tc.GetHLSInfo(hlsURL)
	if err != nil {
		fmt.Printf("Ошибка получения информации: %v\n", err)
	} else {
		fmt.Printf("Найдено %d потоков\n", len(info.Streams))
		fmt.Printf("Live стрим: %v\n", info.IsLive)
		for i, stream := range info.Streams {
			fmt.Printf("  Поток %d: %s, %d kbps\n", i+1, stream.Resolution, stream.Bandwidth/1000)
		}
	}

	// Пример 4: Загрузка HLS с расширенной конфигурацией
	fmt.Println("\n=== HLS с расширенной конфигурацией ===")
	hlsConfig := dto.HLSConfig{
		URL:        hlsURL,
		OutputPath: "stream_hd.mp4",
		Quality:    "1920x1080",      // Конкретное разрешение
		Duration:   30 * time.Minute, // Ограничение по времени
		Headers: map[string]string{
			"Referer":    "https://example.com",
			"User-Agent": "Custom User Agent",
		},
		RetryAttempts:  3,
		SegmentTimeout: 10 * time.Second,
	}

	if err := tc.DownloadHLSWithConfig(ctx, hlsConfig); err != nil {
		fmt.Printf("Ошибка загрузки с конфигурацией: %v\n", err)
	} else {
		fmt.Println("HLS загружен с расширенной конфигурацией")
	}

	// Пример 5: Запись live стрима с ограничением по времени
	fmt.Println("\n=== Запись live стрима ===")
	liveStreamURL := "https://example.com/live/stream.m3u8"
	recordDuration := 5 * time.Minute

	if err := tc.RecordLiveStream(ctx, liveStreamURL, "live_record.mp4", recordDuration); err != nil {
		fmt.Printf("Ошибка записи live стрима: %v\n", err)
	} else {
		fmt.Printf("Live стрим записан в течение %v\n", recordDuration)
	}

	// Пример 6: Конвертация HLS в другой формат
	fmt.Println("\n=== Конвертация HLS в WebM ===")
	if err := tc.ConvertHLSToFormat(ctx, hlsURL, "stream.webm", "webm"); err != nil {
		fmt.Printf("Ошибка конвертации HLS: %v\n", err)
	} else {
		fmt.Println("HLS успешно конвертирован в WebM")
	}

	// Пример 7: Использование пресетов
	fmt.Println("\n=== Использование пресетов ===")
	preset, exists := presets.GetPreset("web-hd")
	if !exists {
		log.Fatal("Пресет не найден")
	}

	presetConfig := preset.Config
	presetConfig.InputPath = "input.mp4"
	presetConfig.OutputPath = "output_hd.mp4"

	presetJob := tc.CreateJob(presetConfig)
	fmt.Printf("Создана задача с пресетом '%s': %s\n", preset.Name, presetJob.ID)

	// Пример 8: Работа с очередью
	fmt.Println("\n=== Работа с очередью ===")
	queue := transcoder.NewQueue(tc, 2) // 2 воркера
	queue.Start()

	// Добавляем несколько задач
	for i := 0; i < 3; i++ {
		config := dto.Config{
			InputPath:  fmt.Sprintf("input%d.mp4", i+1),
			OutputPath: fmt.Sprintf("output%d.mp4", i+1),
			VideoCodec: "libx264",
			AudioCodec: "aac",
		}
		job := tc.CreateJob(config)
		queue.AddJob(job)
		fmt.Printf("Добавлена задача в очередь: %s\n", job.ID)
	}

	// Ждем завершения
	time.Sleep(5 * time.Second)
	queue.Stop()

	// Пример 9: Список доступных пресетов
	fmt.Println("\n=== Доступные пресеты ===")
	presetsList := presets.ListPresets()
	for _, p := range presetsList {
		fmt.Printf("- %s: %s\n", p.Name, p.Description)
	}

	// Пример 10: Быстрые операции
	fmt.Println("\n=== Быстрые операции ===")

	// Конвертация в формат
	if err := tc.ConvertToFormat("input.mp4", "output.webm", "webm"); err != nil {
		fmt.Printf("Ошибка конвертации: %v\n", err)
	} else {
		fmt.Println("Конвертация в WebM завершена")
	}

	// Извлечение аудио
	if err := tc.ExtractAudio("input.mp4", "audio.mp3"); err != nil {
		fmt.Printf("Ошибка извлечения аудио: %v\n", err)
	} else {
		fmt.Println("Аудио извлечено")
	}

	// Создание миниатюры
	if err := tc.CreateThumbnail("input.mp4", "thumbnail.jpg", "00:00:10"); err != nil {
		fmt.Printf("Ошибка создания миниатюры: %v\n", err)
	} else {
		fmt.Println("Миниатюра создана")
	}

	// Получение продолжительности
	duration, err := tc.GetDuration("input.mp4")
	if err != nil {
		fmt.Printf("Ошибка получения продолжительности: %v\n", err)
	} else {
		fmt.Printf("Продолжительность: %s секунд\n", duration)
	}
}
