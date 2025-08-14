package utils

import (
	"fmt"
	"os/exec"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
)

// CheckFFmpeg проверяет доступность FFmpeg
func CheckFFmpeg(path string) error {
	cmd := exec.Command(path, "-version")
	return cmd.Run()
}

// BuildFFmpegArgs строит аргументы для FFmpeg
func BuildFFmpegArgs(config dto.Config) []string {
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

// GetCodecsForFormat возвращает подходящие кодеки для формата
func GetCodecsForFormat(format string) (videoCodec, audioCodec, audioBitrate string) {
	switch format {
	case "mp4":
		return "libx264", "aac", "128k"
	case "webm":
		return "libvpx-vp9", "libopus", "128k"
	case "avi":
		return "libx264", "mp3", "192k"
	case "mp3":
		return "", "libmp3lame", "192k"
	case "m4a":
		return "", "aac", "128k"
	case "flac":
		return "", "flac", ""
	case "wav":
		return "", "pcm_s16le", ""
	default:
		return "libx264", "aac", "128k"
	}
}

// BuildHLSArgs строит аргументы для загрузки HLS
func BuildHLSArgs(ffmpegPath, streamURL string, config dto.HLSConfig) []string {
	args := []string{
		"-i", streamURL,
		"-c", "copy", // копируем без перекодирования для скорости
		"-y", // перезаписывать файл
	}

	// Добавляем заголовки если есть
	if len(config.Headers) > 0 {
		headers := make([]string, 0)
		for key, value := range config.Headers {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
		args = append(args, "-headers", fmt.Sprintf("%s", headers))
	}

	// User-Agent
	userAgent := config.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}
	args = append(args, "-user_agent", userAgent)

	// Ограничение по времени
	if config.Duration > 0 {
		args = append(args, "-t", fmt.Sprintf("%.0f", config.Duration.Seconds()))
	}

	// Выходной файл
	args = append(args, config.OutputPath)
	return args
}
