package transcoder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Transcoder представляет основной интерфейс для транскодирования
type Transcoder struct {
	ffmpegPath string
	tempDir    string
}

// Config содержит настройки для транскодирования
type Config struct {
	InputPath    string
	OutputPath   string
	VideoCodec   string
	AudioCodec   string
	VideoBitrate string
	AudioBitrate string
	Resolution   string
	FrameRate    string
	Quality      string
	Format       string
}

// Job представляет задачу транскодирования
type Job struct {
	ID        string
	Config    Config
	Status    JobStatus
	Progress  float64
	Error     error
	StartTime time.Time
	EndTime   time.Time
}

// JobStatus представляет статус задачи
type JobStatus int

const (
	StatusPending JobStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
)

// New создает новый экземпляр транскодера
func New(ffmpegPath string) (*Transcoder, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	// Проверяем доступность FFmpeg
	if err := checkFFmpeg(ffmpegPath); err != nil {
		return nil, fmt.Errorf("FFmpeg не найден: %w", err)
	}

	tempDir := filepath.Join(os.TempDir(), "transcoder")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать временную директорию: %w", err)
	}

	return &Transcoder{
		ffmpegPath: ffmpegPath,
		tempDir:    tempDir,
	}, nil
}

// checkFFmpeg проверяет доступность FFmpeg
func checkFFmpeg(path string) error {
	cmd := exec.Command(path, "-version")
	return cmd.Run()
}

// CreateJob создает новую задачу транскодирования
func (t *Transcoder) CreateJob(config Config) *Job {
	return &Job{
		ID:     uuid.New().String(),
		Config: config,
		Status: StatusPending,
	}
}

// Execute выполняет транскодирование
func (t *Transcoder) Execute(ctx context.Context, job *Job) error {
	job.Status = StatusRunning
	job.StartTime = time.Now()

	args := t.buildFFmpegArgs(job.Config)
	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)

	if err := cmd.Run(); err != nil {
		job.Status = StatusFailed
		job.Error = err
		job.EndTime = time.Now()
		return fmt.Errorf("ошибка транскодирования: %w", err)
	}

	job.Status = StatusCompleted
	job.Progress = 100.0
	job.EndTime = time.Now()
	return nil
}

// buildFFmpegArgs строит аргументы для FFmpeg
func (t *Transcoder) buildFFmpegArgs(config Config) []string {
	args := []string{
		"-i", config.InputPath,
		"-y", // перезаписывать выходной файл
	}

	// Видео кодек
	if config.VideoCodec != "" {
		args = append(args, "-c:v", config.VideoCodec)
	}

	// Аудио кодек
	if config.AudioCodec != "" {
		args = append(args, "-c:a", config.AudioCodec)
	}

	// Битрейт видео
	if config.VideoBitrate != "" {
		args = append(args, "-b:v", config.VideoBitrate)
	}

	// Битрейт аудио
	if config.AudioBitrate != "" {
		args = append(args, "-b:a", config.AudioBitrate)
	}

	// Разрешение
	if config.Resolution != "" {
		args = append(args, "-s", config.Resolution)
	}

	// Частота кадров
	if config.FrameRate != "" {
		args = append(args, "-r", config.FrameRate)
	}

	// Качество
	if config.Quality != "" {
		args = append(args, "-crf", config.Quality)
	}

	// Формат
	if config.Format != "" {
		args = append(args, "-f", config.Format)
	}

	args = append(args, config.OutputPath)
	return args
}

// GetInfo получает информацию о медиафайле
func (t *Transcoder) GetInfo(filePath string) (map[string]interface{}, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	// Здесь можно добавить парсинг JSON, но для простоты возвращаем строку
	result := map[string]interface{}{
		"raw": string(output),
	}

	return result, nil
}
