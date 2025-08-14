package dto

import (
	"time"
)

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

// JobStatus представляет статус задачи
type JobStatus int

const (
	StatusPending JobStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
)

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

// HLSConfig конфигурация для работы с HLS
type HLSConfig struct {
	URL            string            // URL плейлиста или стрима
	OutputPath     string            // Путь для сохранения
	Quality        string            // Качество (best, worst, или конкретное разрешение)
	Duration       time.Duration     // Максимальная длительность записи (0 = без ограничений)
	Headers        map[string]string // Дополнительные HTTP заголовки
	Cookies        string            // Cookies для авторизации
	UserAgent      string            // User-Agent
	RetryAttempts  int               // Количество попыток при ошибках
	SegmentTimeout time.Duration     // Таймаут для загрузки сегментов
}

// PlaylistInfo информация о плейлисте
type PlaylistInfo struct {
	URL       string
	Streams   []StreamInfo
	Duration  string
	IsLive    bool
	Title     string
	Bandwidth int64
}

// StreamInfo информация о потоке
type StreamInfo struct {
	URL        string
	Resolution string
	Bandwidth  int64
	Codecs     string
	FrameRate  float64
}

// Preset представляет предустановленную конфигурацию
type Preset struct {
	Name        string
	Description string
	Config      Config
}
