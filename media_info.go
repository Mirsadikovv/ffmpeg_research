package transcoder

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// MediaInfo структурированная информация о медиафайле
type MediaInfo struct {
	Format   FormatInfo   `json:"format"`
	Streams  []StreamInfo `json:"streams"`
	Duration time.Duration
	Size     int64
	Bitrate  int64
	HasVideo bool
	HasAudio bool
}

// FormatInfo информация о формате файла
type FormatInfo struct {
	Filename       string            `json:"filename"`
	FormatName     string            `json:"format_name"`
	FormatLongName string            `json:"format_long_name"`
	Duration       string            `json:"duration"`
	Size           string            `json:"size"`
	Bitrate        string            `json:"bit_rate"`
	Tags           map[string]string `json:"tags"`
}

// StreamInfo информация о потоке
type StreamInfo struct {
	Index         int               `json:"index"`
	CodecName     string            `json:"codec_name"`
	CodecLongName string            `json:"codec_long_name"`
	CodecType     string            `json:"codec_type"`
	Width         int               `json:"width,omitempty"`
	Height        int               `json:"height,omitempty"`
	PixelFormat   string            `json:"pix_fmt,omitempty"`
	FrameRate     string            `json:"r_frame_rate,omitempty"`
	Duration      string            `json:"duration,omitempty"`
	Bitrate       string            `json:"bit_rate,omitempty"`
	SampleRate    string            `json:"sample_rate,omitempty"`
	Channels      int               `json:"channels,omitempty"`
	Tags          map[string]string `json:"tags"`
}

// GetMediaInfo получает подробную информацию о медиафайле
func (t *Transcoder) GetMediaInfo(filePath string) (*MediaInfo, error) {
	t.logger.Debug("Получение информации о файле: %s", filePath)

	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		t.logger.Error("Ошибка выполнения ffprobe: %v", err)
		return nil, fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	var rawInfo struct {
		Format  FormatInfo   `json:"format"`
		Streams []StreamInfo `json:"streams"`
	}

	if err := json.Unmarshal(output, &rawInfo); err != nil {
		t.logger.Error("Ошибка парсинга JSON: %v", err)
		return nil, fmt.Errorf("ошибка парсинга информации о файле: %w", err)
	}

	info := &MediaInfo{
		Format:  rawInfo.Format,
		Streams: rawInfo.Streams,
	}

	// Парсим продолжительность
	if duration, err := strconv.ParseFloat(rawInfo.Format.Duration, 64); err == nil {
		info.Duration = time.Duration(duration * float64(time.Second))
	}

	// Парсим размер файла
	if size, err := strconv.ParseInt(rawInfo.Format.Size, 10, 64); err == nil {
		info.Size = size
	}

	// Парсим битрейт
	if bitrate, err := strconv.ParseInt(rawInfo.Format.Bitrate, 10, 64); err == nil {
		info.Bitrate = bitrate
	}

	// Определяем типы потоков
	for _, stream := range rawInfo.Streams {
		switch stream.CodecType {
		case "video":
			info.HasVideo = true
		case "audio":
			info.HasAudio = true
		}
	}

	t.logger.Info("Информация о файле получена: %s (%.2fs, %d байт)",
		filePath, info.Duration.Seconds(), info.Size)

	return info, nil
}

// GetVideoStreams возвращает только видео потоки
func (info *MediaInfo) GetVideoStreams() []StreamInfo {
	var videoStreams []StreamInfo
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			videoStreams = append(videoStreams, stream)
		}
	}
	return videoStreams
}

// GetAudioStreams возвращает только аудио потоки
func (info *MediaInfo) GetAudioStreams() []StreamInfo {
	var audioStreams []StreamInfo
	for _, stream := range info.Streams {
		if stream.CodecType == "audio" {
			audioStreams = append(audioStreams, stream)
		}
	}
	return audioStreams
}

// GetResolution возвращает разрешение первого видео потока
func (info *MediaInfo) GetResolution() string {
	videoStreams := info.GetVideoStreams()
	if len(videoStreams) > 0 {
		stream := videoStreams[0]
		if stream.Width > 0 && stream.Height > 0 {
			return fmt.Sprintf("%dx%d", stream.Width, stream.Height)
		}
	}
	return ""
}

// GetFrameRate возвращает частоту кадров первого видео потока
func (info *MediaInfo) GetFrameRate() float64 {
	videoStreams := info.GetVideoStreams()
	if len(videoStreams) > 0 {
		frameRate := videoStreams[0].FrameRate
		if frameRate != "" {
			// Парсим дробь вида "30/1" или "25000/1001"
			parts := strings.Split(frameRate, "/")
			if len(parts) == 2 {
				num, err1 := strconv.ParseFloat(parts[0], 64)
				den, err2 := strconv.ParseFloat(parts[1], 64)
				if err1 == nil && err2 == nil && den != 0 {
					return num / den
				}
			}
		}
	}
	return 0
}

// IsVideo проверяет, является ли файл видео
func (info *MediaInfo) IsVideo() bool {
	return info.HasVideo
}

// IsAudio проверяет, является ли файл только аудио
func (info *MediaInfo) IsAudio() bool {
	return info.HasAudio && !info.HasVideo
}

// GetTitle возвращает название из метаданных
func (info *MediaInfo) GetTitle() string {
	if title, exists := info.Format.Tags["title"]; exists {
		return title
	}
	return ""
}

// GetArtist возвращает исполнителя из метаданных
func (info *MediaInfo) GetArtist() string {
	if artist, exists := info.Format.Tags["artist"]; exists {
		return artist
	}
	return ""
}

// Summary возвращает краткую сводку о файле
func (info *MediaInfo) Summary() string {
	var parts []string

	if info.IsVideo() {
		resolution := info.GetResolution()
		frameRate := info.GetFrameRate()
		if resolution != "" {
			if frameRate > 0 {
				parts = append(parts, fmt.Sprintf("Видео: %s @ %.1f fps", resolution, frameRate))
			} else {
				parts = append(parts, fmt.Sprintf("Видео: %s", resolution))
			}
		} else {
			parts = append(parts, "Видео")
		}
	}

	if info.HasAudio {
		audioStreams := info.GetAudioStreams()
		if len(audioStreams) > 0 {
			stream := audioStreams[0]
			if stream.SampleRate != "" && stream.Channels > 0 {
				parts = append(parts, fmt.Sprintf("Аудио: %s Hz, %d каналов", stream.SampleRate, stream.Channels))
			} else {
				parts = append(parts, "Аудио")
			}
		}
	}

	if info.Duration > 0 {
		parts = append(parts, fmt.Sprintf("Длительность: %.1fs", info.Duration.Seconds()))
	}

	if info.Size > 0 {
		parts = append(parts, fmt.Sprintf("Размер: %.1f MB", float64(info.Size)/(1024*1024)))
	}

	return strings.Join(parts, ", ")
}
