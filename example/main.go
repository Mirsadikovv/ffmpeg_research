package main

import (
	"context"
	"fmt"
	"log"
	"time"

	transcoder "github.com/Mirsadikovv/ffmpeg_research"
)

func main() {
	// Создаем транскодер
	tc, err := transcoder.New("ffmpeg")
	if err != nil {
		log.Fatal("Ошибка создания транскодера:", err)
	}

	// Пример 1: Простое транскодирование
	fmt.Println("=== Простое транскодирование ===")
	config := transcoder.Config{
		InputPath:    "input.mp4",
		OutputPath:   "output.mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
	}

	job := tc.CreateJob(config)
	fmt.Printf("Создана задача: %s\n", job.ID)

	ctx := context.Background()
	if err := tc.Execute(ctx, job); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Транскодирование завершено за %v\n", job.EndTime.Sub(job.StartTime))
	}

	// Пример 2: Использование пресетов
	fmt.Println("\n=== Использование пресетов ===")
	preset, exists := transcoder.GetPreset("web-hd")
	if !exists {
		log.Fatal("Пресет не найден")
	}

	presetConfig := preset.Config
	presetConfig.InputPath = "input.mp4"
	presetConfig.OutputPath = "output_hd.mp4"

	presetJob := tc.CreateJob(presetConfig)
	fmt.Printf("Создана задача с пресетом '%s': %s\n", preset.Name, presetJob.ID)

	// Пример 3: Работа с очередью
	fmt.Println("\n=== Работа с очередью ===")
	queue := transcoder.NewQueue(tc, 2) // 2 воркера
	queue.Start()

	// Добавляем несколько задач
	for i := 0; i < 3; i++ {
		config := transcoder.Config{
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

	// Пример 4: Получение информации о файле
	fmt.Println("\n=== Информация о файле ===")
	info, err := tc.GetInfo("input.mp4")
	if err != nil {
		fmt.Printf("Ошибка получения информации: %v\n", err)
	} else {
		fmt.Printf("Информация о файле получена: %d байт\n", len(info["raw"].(string)))
	}

	// Пример 5: Список доступных пресетов
	fmt.Println("\n=== Доступные пресеты ===")
	presets := transcoder.ListPresets()
	for _, p := range presets {
		fmt.Printf("- %s: %s\n", p.Name, p.Description)
	}

	// Пример 6: Быстрая конвертация в формат
	fmt.Println("\n=== Быстрая конвертация ===")
	if err := tc.ConvertToFormat("input.mp4", "output.webm", "webm"); err != nil {
		fmt.Printf("Ошибка конвертации: %v\n", err)
	} else {
		fmt.Println("Конвертация в WebM завершена")
	}

	// Пример 7: Извлечение аудио
	fmt.Println("\n=== Извлечение аудио ===")
	if err := tc.ExtractAudio("input.mp4", "audio.mp3"); err != nil {
		fmt.Printf("Ошибка извлечения аудио: %v\n", err)
	} else {
		fmt.Println("Аудио извлечено")
	}

	// Пример 8: Создание миниатюры
	fmt.Println("\n=== Создание миниатюры ===")
	if err := tc.CreateThumbnail("input.mp4", "thumbnail.jpg", "00:00:10"); err != nil {
		fmt.Printf("Ошибка создания миниатюры: %v\n", err)
	} else {
		fmt.Println("Миниатюра создана")
	}

	// Пример 9: Получение продолжительности
	fmt.Println("\n=== Продолжительность файла ===")
	duration, err := tc.GetDuration("input.mp4")
	if err != nil {
		fmt.Printf("Ошибка получения продолжительности: %v\n", err)
	} else {
		fmt.Printf("Продолжительность: %s секунд\n", duration)
	}
}
