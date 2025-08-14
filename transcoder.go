package transcoder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
	"github.com/Mirsadikovv/ffmpeg_research/utils"
	"github.com/google/uuid"
)

// Transcoder представляет основной интерфейс для транскодирования
type Transcoder struct {
	ffmpegPath string
	tempDir    string
	hls        *HLSDownloader
	logger     Logger
}

// New создает новый экземпляр транскодера
func New(ffmpegPath string) (*Transcoder, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	// Проверяем доступность FFmpeg
	if err := utils.CheckFFmpeg(ffmpegPath); err != nil {
		return nil, fmt.Errorf("FFmpeg не найден: %w", err)
	}

	tempDir := filepath.Join(os.TempDir(), "transcoder")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать временную директорию: %w", err)
	}

	transcoder := &Transcoder{
		ffmpegPath: ffmpegPath,
		tempDir:    tempDir,
		logger:     NewDefaultLogger(LogLevelInfo), // По умолчанию INFO уровень
	}

	// Инициализируем HLS загрузчик
	transcoder.hls = NewHLSDownloader(transcoder)

	return transcoder, nil
}

// CreateJob создает новую задачу транскодирования
func (t *Transcoder) CreateJob(config dto.Config) *dto.Job {
	return &dto.Job{
		ID:     uuid.New().String(),
		Config: config,
		Status: dto.StatusPending,
	}
}

// Execute выполняет транскодирование
func (t *Transcoder) Execute(ctx context.Context, job *dto.Job) error {
	// Валидируем конфигурацию перед выполнением
	if err := job.Config.Validate(); err != nil {
		job.Status = dto.StatusFailed
		job.Error = err
		t.logger.Error("Ошибка валидации конфигурации: %v", err)
		return fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	job.Status = dto.StatusRunning
	job.StartTime = time.Now()

	t.logger.Info("Начало транскодирования: %s -> %s", job.Config.InputPath, job.Config.OutputPath)

	args := utils.BuildFFmpegArgs(job.Config)
	t.logger.Debug("FFmpeg аргументы: %v", args)

	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)

	if err := cmd.Run(); err != nil {
		job.Status = dto.StatusFailed
		job.Error = err
		job.EndTime = time.Now()
		t.logger.Error("Ошибка выполнения FFmpeg: %v", err)
		return fmt.Errorf("ошибка транскодирования: %w", err)
	}

	job.Status = dto.StatusCompleted
	job.Progress = 100.0
	job.EndTime = time.Now()

	duration := job.EndTime.Sub(job.StartTime)
	t.logger.Info("Транскодирование завершено за %v", duration)

	return nil
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

// ConvertToFormat конвертирует файл в указанный формат с базовыми настройками
func (t *Transcoder) ConvertToFormat(inputPath, outputPath, format string) error {
	config := dto.Config{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Format:     format,
	}

	// Автоматически выбираем подходящие кодеки для формата
	videoCodec, audioCodec, audioBitrate := utils.GetCodecsForFormat(format)
	config.VideoCodec = videoCodec
	config.AudioCodec = audioCodec
	if audioBitrate != "" {
		config.AudioBitrate = audioBitrate
	}

	job := t.CreateJob(config)
	return t.Execute(context.Background(), job)
}

// ExtractAudio извлекает аудиодорожку из видеофайла
func (t *Transcoder) ExtractAudio(inputPath, outputPath string) error {
	config := dto.Config{
		InputPath:  inputPath,
		OutputPath: outputPath,
		AudioCodec: "copy", // копируем без перекодирования
		VideoCodec: "",     // без видео
	}

	job := t.CreateJob(config)
	return t.Execute(context.Background(), job)
}

// CreateThumbnail создает миниатюру из видео
func (t *Transcoder) CreateThumbnail(inputPath, outputPath string, timeOffset string) error {
	args := []string{
		"-i", inputPath,
		"-ss", timeOffset,
		"-frames:v", "1",
		"-y",
		outputPath,
	}

	cmd := exec.Command(t.ffmpegPath, args...)
	return cmd.Run()
}

// GetDuration получает продолжительность медиафайла
func (t *Transcoder) GetDuration(filePath string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ошибка получения продолжительности: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// HLS методы

// DownloadHLS загружает HLS стрим или плейлист
func (t *Transcoder) DownloadHLS(ctx context.Context, hlsURL, outputPath string) error {
	config := dto.HLSConfig{
		URL:        hlsURL,
		OutputPath: outputPath,
		Quality:    "best",
	}
	return t.hls.DownloadHLS(ctx, config)
}

// DownloadHLSWithConfig загружает HLS с расширенной конфигурацией
func (t *Transcoder) DownloadHLSWithConfig(ctx context.Context, config dto.HLSConfig) error {
	return t.hls.DownloadHLS(ctx, config)
}

// GetHLSInfo получает информацию о HLS плейлисте
func (t *Transcoder) GetHLSInfo(playlistURL string) (*dto.PlaylistInfo, error) {
	return t.hls.GetPlaylistInfo(playlistURL)
}

// RecordLiveStream записывает live стрим с ограничением по времени
func (t *Transcoder) RecordLiveStream(ctx context.Context, streamURL, outputPath string, duration time.Duration) error {
	return t.hls.RecordLiveStream(ctx, streamURL, outputPath, duration)
}

// ConvertHLSToFormat загружает HLS и конвертирует в указанный формат
func (t *Transcoder) ConvertHLSToFormat(ctx context.Context, hlsURL, outputPath, format string) error {
	return t.hls.ConvertHLSToFormat(ctx, hlsURL, outputPath, format)
}

// SetLogger устанавливает кастомный логгер
func (t *Transcoder) SetLogger(logger Logger) {
	t.logger = logger
}

// GetLogger возвращает текущий логгер
func (t *Transcoder) GetLogger() Logger {
	return t.logger
}
