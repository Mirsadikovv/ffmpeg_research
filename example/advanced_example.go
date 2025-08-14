//go:build ignore
// +build ignore

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
	// Создаем транскодер с кастомным логгером
	tc, err := transcoder.New("ffmpeg")
	if err != nil {
		log.Fatal("Ошибка создания транскодера:", err)
	}

	// Устанавливаем уровень логирования
	logger := transcoder.NewDefaultLogger(transcoder.LogLevelDebug)
	tc.SetLogger(logger)

	ctx := context.Background()

	// Пример 1: Получение подробной информации о медиафайле
	fmt.Println("=== Анализ медиафайла ===")
	mediaInfo, err := tc.GetMediaInfo("input.mp4")
	if err != nil {
		fmt.Printf("Ошибка получения информации: %v\n", err)
	} else {
		fmt.Printf("Файл: %s\n", mediaInfo.Format.Filename)
		fmt.Printf("Формат: %s\n", mediaInfo.Format.FormatLongName)
		fmt.Printf("Продолжительность: %.2f секунд\n", mediaInfo.Duration.Seconds())
		fmt.Printf("Размер: %.2f MB\n", float64(mediaInfo.Size)/(1024*1024))
		fmt.Printf("Разрешение: %s\n", mediaInfo.GetResolution())
		fmt.Printf("Частота кадров: %.1f fps\n", mediaInfo.GetFrameRate())
		fmt.Printf("Краткая сводка: %s\n", mediaInfo.Summary())

		// Анализ потоков
		videoStreams := mediaInfo.GetVideoStreams()
		audioStreams := mediaInfo.GetAudioStreams()
		fmt.Printf("Видео потоков: %d, Аудио потоков: %d\n", len(videoStreams), len(audioStreams))
	}

	// Пример 2: Транскодирование с отслеживанием прогресса
	fmt.Println("\n=== Транскодирование с прогрессом ===")

	// Создаем callback для отслеживания прогресса
	progressCallback := func(progress float64, speed string, eta time.Duration) {
		fmt.Printf("\rПрогресс: %.1f%% | Скорость: %s | ETA: %v",
			progress, speed, eta.Round(time.Second))
	}

	// Создаем трекер прогресса
	progressTracker := transcoder.NewProgressTracker(progressCallback)

	// Устанавливаем продолжительность для расчета прогресса
	if mediaInfo != nil {
		progressTracker.SetDuration(mediaInfo.Duration)
	}

	// Используем пресет для высокого качества
	preset, exists := presets.GetPreset("high-quality")
	if !exists {
		log.Fatal("Пресет не найден")
	}

	config := preset.Config
	config.InputPath = "input.mp4"
	config.OutputPath = "output_hq.mp4"

	job := tc.CreateJob(config)
	fmt.Printf("Создана задача: %s\n", job.ID)

	// Выполняем транскодирование
	if err := tc.Execute(ctx, job); err != nil {
		fmt.Printf("\nОшибка: %v\n", err)
	} else {
		fmt.Printf("\nТранскодирование завершено за %v\n", job.EndTime.Sub(job.StartTime))
	}

	// Пример 3: Работа с категориями пресетов
	fmt.Println("\n=== Пресеты по категориям ===")
	categories := presets.GetPresetsByCategory()

	for category, categoryPresets := range categories {
		fmt.Printf("\n%s:\n", category)
		for _, p := range categoryPresets {
			fmt.Printf("  - %s: %s\n", p.Name, p.Description)
		}
	}

	// Пример 4: Автоматический выбор пресета на основе входного файла
	fmt.Println("\n=== Автоматический выбор пресета ===")
	if mediaInfo != nil {
		recommendedPreset := recommendPreset(mediaInfo)
		fmt.Printf("Рекомендуемый пресет: %s (%s)\n",
			recommendedPreset.Name, recommendedPreset.Description)

		// Применяем рекомендуемый пресет
		autoConfig := recommendedPreset.Config
		autoConfig.InputPath = "input.mp4"
		autoConfig.OutputPath = "output_auto.mp4"

		autoJob := tc.CreateJob(autoConfig)
		fmt.Printf("Создана автоматическая задача: %s\n", autoJob.ID)
	}

	// Пример 5: Пакетная обработка с разными пресетами
	fmt.Println("\n=== Пакетная обработка ===")

	batchConfigs := []struct {
		preset string
		suffix string
	}{
		{"web-hd", "_web_hd"},
		{"mobile", "_mobile"},
		{"small-size", "_small"},
	}

	queue := transcoder.NewQueue(tc, 2) // 2 воркера
	queue.Start()
	defer queue.Stop()

	for _, batch := range batchConfigs {
		if preset, exists := presets.GetPreset(batch.preset); exists {
			config := preset.Config
			config.InputPath = "input.mp4"
			config.OutputPath = fmt.Sprintf("output%s.mp4", batch.suffix)

			job := tc.CreateJob(config)
			queue.AddJob(job)
			fmt.Printf("Добавлена задача %s: %s\n", batch.preset, job.ID)
		}
	}

	// Ждем завершения пакетной обработки
	fmt.Println("Ожидание завершения пакетной обработки...")
	time.Sleep(10 * time.Second)

	// Пример 6: Создание миниатюр с разными временными метками
	fmt.Println("\n=== Создание миниатюр ===")
	thumbnailTimes := []string{"00:00:05", "00:00:30", "00:01:00"}

	for i, timeOffset := range thumbnailTimes {
		outputPath := fmt.Sprintf("thumbnail_%d.jpg", i+1)
		if err := tc.CreateThumbnail("input.mp4", outputPath, timeOffset); err != nil {
			fmt.Printf("Ошибка создания миниатюры %s: %v\n", timeOffset, err)
		} else {
			fmt.Printf("Создана миниатюра: %s (время: %s)\n", outputPath, timeOffset)
		}
	}

	// Пример 7: Извлечение аудио с разными форматами
	fmt.Println("\n=== Извлечение аудио ===")
	audioFormats := []struct {
		format string
		ext    string
	}{
		{"mp3", "mp3"},
		{"m4a", "m4a"},
		{"flac", "flac"},
	}

	for _, format := range audioFormats {
		outputPath := fmt.Sprintf("audio.%s", format.ext)
		if err := tc.ConvertToFormat("input.mp4", outputPath, format.format); err != nil {
			fmt.Printf("Ошибка конвертации в %s: %v\n", format.format, err)
		} else {
			fmt.Printf("Аудио извлечено: %s\n", outputPath)
		}
	}

	// Пример 8: Работа с HLS стримами (расширенная конфигурация)
	fmt.Println("\n=== Работа с HLS стримами ===")
	hlsURL := "https://example.com/playlist.m3u8"

	// Получаем информацию о плейлисте
	if info, err := tc.GetHLSInfo(hlsURL); err == nil {
		fmt.Printf("HLS плейлист: %d потоков, Live: %v\n", len(info.Streams), info.IsLive)

		// Показываем доступные потоки
		for i, stream := range info.Streams {
			fmt.Printf("  Поток %d: %s, %d kbps\n",
				i+1, stream.Resolution, stream.Bandwidth/1000)
		}

		// Загружаем с расширенной конфигурацией
		hlsConfig := dto.HLSConfig{
			URL:        hlsURL,
			OutputPath: "stream_custom.mp4",
			Quality:    "best",
			Duration:   5 * time.Minute,
			Headers: map[string]string{
				"Referer":    "https://example.com",
				"User-Agent": "Custom Transcoder Bot 1.0",
			},
			RetryAttempts:  3,
			SegmentTimeout: 15 * time.Second,
		}

		if err := tc.DownloadHLSWithConfig(ctx, hlsConfig); err != nil {
			fmt.Printf("Ошибка загрузки HLS: %v\n", err)
		} else {
			fmt.Println("HLS стрим успешно загружен с кастомными настройками")
		}
	}

	// Пример 9: Использование фильтров
	fmt.Println("\n=== Применение фильтров ===")

	// Создаем цепочку фильтров
	filterChain := transcoder.NewFilterChain().
		AddVideoFilter("scale", map[string]string{"w": "1280", "h": "720"}).
		AddVideoFilter("eq", map[string]string{"brightness": "0.1", "contrast": "1.2"}).
		AddAudioFilter("volume", map[string]string{"volume": "1.2"})

	// Применяем фильтры
	filterConfig := dto.Config{
		InputPath:  "input.mp4",
		OutputPath: "output_filtered.mp4",
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Quality:    "22",
	}

	filterJob := tc.CreateJob(filterConfig)
	if err := tc.ExecuteWithFilters(ctx, filterJob, filterChain); err != nil {
		fmt.Printf("Ошибка применения фильтров: %v\n", err)
	} else {
		fmt.Println("Фильтры успешно применены")
	}

	// Пример 10: Использование конвейеров (pipelines)
	fmt.Println("\n=== Конвейеры обработки ===")

	// Создаем веб-оптимизационный конвейер
	webPipeline := transcoder.CreateWebOptimizationPipeline(tc)
	fmt.Printf("Конвейер: %s\n", webPipeline.GetName())
	fmt.Printf("Описание: %s\n", webPipeline.GetDescription())

	// Показываем шаги конвейера
	fmt.Println("Шаги конвейера:")
	for i, step := range webPipeline.GetSteps() {
		fmt.Printf("  %d. %s: %s\n", i+1, step.GetName(), step.GetDescription())
	}

	// Выполняем конвейер
	if err := webPipeline.Execute(ctx, "input.mp4", "output_web_optimized.mp4"); err != nil {
		fmt.Printf("Ошибка выполнения конвейера: %v\n", err)
	} else {
		fmt.Println("Веб-оптимизационный конвейер завершен успешно")
	}

	// Пример 11: Мобильный конвейер
	fmt.Println("\n=== Мобильный конвейер ===")

	mobilePipeline := transcoder.CreateMobilePipeline(tc)
	if err := mobilePipeline.Execute(ctx, "input.mp4", "output_mobile.mp4"); err != nil {
		fmt.Printf("Ошибка мобильного конвейера: %v\n", err)
	} else {
		fmt.Println("Мобильный конвейер завершен успешно")
	}

	// Пример 12: Архивный конвейер
	fmt.Println("\n=== Архивный конвейер ===")

	archivePipeline := transcoder.CreateArchivePipeline(tc)
	if err := archivePipeline.Execute(ctx, "input.mp4", "output_archive.mkv"); err != nil {
		fmt.Printf("Ошибка архивного конвейера: %v\n", err)
	} else {
		fmt.Println("Архивный конвейер завершен успешно")
	}

	fmt.Println("\n=== Демонстрация завершена ===")
}

// recommendPreset рекомендует пресет на основе характеристик входного файла
func recommendPreset(info *transcoder.MediaInfo) *dto.Preset {
	resolution := info.GetResolution()
	duration := info.Duration
	isVideo := info.IsVideo()

	// Если это только аудио
	if !isVideo {
		if duration > 30*time.Minute {
			// Длинное аудио - вероятно подкаст или аудиокнига
			preset, _ := presets.GetPreset("audiobook-mp3")
			return preset
		} else {
			// Короткое аудио - обычная музыка
			preset, _ := presets.GetPreset("audio-mp3")
			return preset
		}
	}

	// Для видео анализируем разрешение
	switch resolution {
	case "3840x2160", "4096x2160": // 4K
		preset, _ := presets.GetPreset("high-quality")
		return preset
	case "1920x1080": // Full HD
		if duration > 2*time.Hour {
			// Длинное видео - используем эффективное сжатие
			preset, _ := presets.GetPreset("small-size")
			return preset
		} else {
			// Обычное видео - веб качество
			preset, _ := presets.GetPreset("web-optimized")
			return preset
		}
	case "1280x720": // HD
		preset, _ := presets.GetPreset("web-hd")
		return preset
	default:
		// Низкое разрешение или неизвестное - мобильный пресет
		preset, _ := presets.GetPreset("mobile")
		return preset
	}
}
