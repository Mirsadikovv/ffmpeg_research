package transcoder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
	"github.com/Mirsadikovv/ffmpeg_research/utils"
)

// HLSDownloader управляет загрузкой HLS стримов и плейлистов
type HLSDownloader struct {
	transcoder *Transcoder
	client     *http.Client
	userAgent  string
}

// NewHLSDownloader создает новый загрузчик HLS
func NewHLSDownloader(transcoder *Transcoder) *HLSDownloader {
	return &HLSDownloader{
		transcoder: transcoder,
		client:     utils.CreateHTTPClient(30 * time.Second),
		userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}

// DownloadHLS загружает HLS стрим или плейлист
func (h *HLSDownloader) DownloadHLS(ctx context.Context, config dto.HLSConfig) error {
	// Валидируем конфигурацию HLS
	if err := config.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации HLS конфигурации: %w", err)
	}

	// Проверяем, является ли URL плейлистом или прямой ссылкой
	if strings.HasSuffix(config.URL, ".m3u8") || strings.Contains(config.URL, "m3u8") {
		return h.downloadPlaylist(ctx, config)
	}

	// Пытаемся найти плейлист по URL
	playlistURL, err := utils.FindPlaylistURL(config.URL, h.client)
	if err != nil {
		return fmt.Errorf("не удалось найти плейлист: %w", err)
	}

	config.URL = playlistURL
	return h.downloadPlaylist(ctx, config)
}

// downloadPlaylist загружает плейлист
func (h *HLSDownloader) downloadPlaylist(ctx context.Context, config dto.HLSConfig) error {
	// Получаем информацию о плейлисте
	info, err := h.GetPlaylistInfo(config.URL)
	if err != nil {
		return fmt.Errorf("ошибка получения информации о плейлисте: %w", err)
	}

	// Выбираем поток по качеству
	streamURL, err := utils.SelectStream(info, config.Quality)
	if err != nil {
		return fmt.Errorf("ошибка выбора потока: %w", err)
	}

	// Используем FFmpeg для загрузки
	return h.downloadWithFFmpeg(ctx, streamURL, config)
}

// downloadWithFFmpeg загружает стрим с помощью FFmpeg
func (h *HLSDownloader) downloadWithFFmpeg(ctx context.Context, streamURL string, config dto.HLSConfig) error {
	args := utils.BuildHLSArgs(h.transcoder.ffmpegPath, streamURL, config)
	cmd := exec.CommandContext(ctx, h.transcoder.ffmpegPath, args...)

	// Настраиваем вывод для отслеживания прогресса
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetPlaylistInfo получает информацию о плейлисте
func (h *HLSDownloader) GetPlaylistInfo(playlistURL string) (*dto.PlaylistInfo, error) {
	resp, err := h.client.Get(playlistURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки плейлиста: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка HTTP: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения плейлиста: %w", err)
	}

	return utils.ParsePlaylist(string(content), playlistURL)
}

// RecordLiveStream записывает live стрим с ограничением по времени
func (h *HLSDownloader) RecordLiveStream(ctx context.Context, streamURL, outputPath string, duration time.Duration) error {
	config := dto.HLSConfig{
		URL:        streamURL,
		OutputPath: outputPath,
		Duration:   duration,
		Quality:    "best",
	}

	return h.DownloadHLS(ctx, config)
}

// ConvertHLSToFormat конвертирует загруженный HLS в другой формат
func (h *HLSDownloader) ConvertHLSToFormat(ctx context.Context, hlsURL, outputPath, format string) error {
	// Создаем временный файл
	tempFile := filepath.Join(h.transcoder.tempDir, fmt.Sprintf("temp_hls_%d.ts", time.Now().Unix()))
	defer os.Remove(tempFile)

	// Сначала загружаем HLS
	config := dto.HLSConfig{
		URL:        hlsURL,
		OutputPath: tempFile,
		Quality:    "best",
	}

	if err := h.DownloadHLS(ctx, config); err != nil {
		return fmt.Errorf("ошибка загрузки HLS: %w", err)
	}

	// Затем конвертируем в нужный формат
	return h.transcoder.ConvertToFormat(tempFile, outputPath, format)
}
