package transcoder

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Тест создания транскодера
	tc, err := New("ffmpeg")
	if err != nil {
		t.Skip("FFmpeg не найден, пропускаем тест")
	}

	if tc == nil {
		t.Error("Транскодер не должен быть nil")
	}

	if tc.ffmpegPath != "ffmpeg" {
		t.Errorf("Ожидался путь 'ffmpeg', получен '%s'", tc.ffmpegPath)
	}
}

func TestCreateJob(t *testing.T) {
	tc, err := New("ffmpeg")
	if err != nil {
		t.Skip("FFmpeg не найден, пропускаем тест")
	}

	config := Config{
		InputPath:  "test_input.mp4",
		OutputPath: "test_output.mp4",
		VideoCodec: "libx264",
		AudioCodec: "aac",
	}

	job := tc.CreateJob(config)

	if job == nil {
		t.Error("Job не должен быть nil")
	}

	if job.ID == "" {
		t.Error("Job ID не должен быть пустым")
	}

	if job.Status != StatusPending {
		t.Errorf("Ожидался статус %d, получен %d", StatusPending, job.Status)
	}

	if job.Config.InputPath != config.InputPath {
		t.Error("Конфигурация job не соответствует переданной")
	}
}

func TestGetPreset(t *testing.T) {
	// Тест существующего пресета
	preset, exists := GetPreset("web-hd")
	if !exists {
		t.Error("Пресет 'web-hd' должен существовать")
	}

	if preset.Name != "web-hd" {
		t.Errorf("Ожидалось имя 'web-hd', получено '%s'", preset.Name)
	}

	// Тест несуществующего пресета
	_, exists = GetPreset("nonexistent")
	if exists {
		t.Error("Несуществующий пресет не должен быть найден")
	}
}

func TestListPresets(t *testing.T) {
	presets := ListPresets()

	if len(presets) == 0 {
		t.Error("Список пресетов не должен быть пустым")
	}

	// Проверяем, что все основные пресеты присутствуют
	expectedPresets := []string{"web-hd", "web-sd", "mobile", "audio-mp3", "audio-aac"}
	presetNames := make(map[string]bool)

	for _, preset := range presets {
		presetNames[preset.Name] = true
	}

	for _, expected := range expectedPresets {
		if !presetNames[expected] {
			t.Errorf("Пресет '%s' не найден в списке", expected)
		}
	}
}

func TestBuildFFmpegArgs(t *testing.T) {
	tc := &Transcoder{ffmpegPath: "ffmpeg"}

	config := Config{
		InputPath:    "input.mp4",
		OutputPath:   "output.mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
		Resolution:   "1280x720",
		FrameRate:    "30",
		Quality:      "23",
		Format:       "mp4",
	}

	args := tc.buildFFmpegArgs(config)

	// Проверяем основные аргументы
	expectedArgs := []string{
		"-i", "input.mp4",
		"-y",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-b:v", "1000k",
		"-b:a", "128k",
		"-s", "1280x720",
		"-r", "30",
		"-crf", "23",
		"-f", "mp4",
		"output.mp4",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Ожидалось %d аргументов, получено %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if i >= len(args) || args[i] != expected {
			t.Errorf("Аргумент %d: ожидался '%s', получен '%s'", i, expected, args[i])
		}
	}
}

func TestQueue(t *testing.T) {
	tc, err := New("ffmpeg")
	if err != nil {
		t.Skip("FFmpeg не найден, пропускаем тест")
	}

	queue := NewQueue(tc, 2)

	if queue == nil {
		t.Error("Queue не должна быть nil")
	}

	if queue.workers != 2 {
		t.Errorf("Ожидалось 2 воркера, получено %d", queue.workers)
	}

	// Тест добавления задач
	config := Config{
		InputPath:  "test.mp4",
		OutputPath: "output.mp4",
	}

	job := tc.CreateJob(config)
	queue.AddJob(job)

	jobs := queue.GetJobs()
	if len(jobs) != 1 {
		t.Errorf("Ожидалась 1 задача в очереди, получено %d", len(jobs))
	}

	// Останавливаем очередь
	queue.Stop()
}

// Benchmark тесты
func BenchmarkCreateJob(b *testing.B) {
	tc, err := New("ffmpeg")
	if err != nil {
		b.Skip("FFmpeg не найден, пропускаем бенчмарк")
	}

	config := Config{
		InputPath:  "input.mp4",
		OutputPath: "output.mp4",
		VideoCodec: "libx264",
		AudioCodec: "aac",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.CreateJob(config)
	}
}

func BenchmarkBuildFFmpegArgs(b *testing.B) {
	tc := &Transcoder{ffmpegPath: "ffmpeg"}
	config := Config{
		InputPath:    "input.mp4",
		OutputPath:   "output.mp4",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.buildFFmpegArgs(config)
	}
}

// Пример теста интеграции (требует реальные файлы)
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме")
	}

	// Создаем тестовый файл (можно пропустить если нет ffmpeg)
	tc, err := New("ffmpeg")
	if err != nil {
		t.Skip("FFmpeg не найден, пропускаем интеграционный тест")
	}

	// Проверяем, что можем получить информацию о файле
	// (этот тест будет работать только если есть тестовый файл)
	testFile := "test_video.mp4"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Тестовый файл не найден, пропускаем интеграционный тест")
	}

	info, err := tc.GetInfo(testFile)
	if err != nil {
		t.Errorf("Ошибка получения информации о файле: %v", err)
	}

	if info == nil {
		t.Error("Информация о файле не должна быть nil")
	}
}
