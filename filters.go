package transcoder

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
)

// Filter представляет FFmpeg фильтр
type Filter struct {
	Name   string
	Params map[string]string
}

// FilterChain цепочка фильтров
type FilterChain struct {
	VideoFilters []Filter
	AudioFilters []Filter
}

// NewFilterChain создает новую цепочку фильтров
func NewFilterChain() *FilterChain {
	return &FilterChain{
		VideoFilters: make([]Filter, 0),
		AudioFilters: make([]Filter, 0),
	}
}

// AddVideoFilter добавляет видео фильтр
func (fc *FilterChain) AddVideoFilter(name string, params map[string]string) *FilterChain {
	fc.VideoFilters = append(fc.VideoFilters, Filter{
		Name:   name,
		Params: params,
	})
	return fc
}

// AddAudioFilter добавляет аудио фильтр
func (fc *FilterChain) AddAudioFilter(name string, params map[string]string) *FilterChain {
	fc.AudioFilters = append(fc.AudioFilters, Filter{
		Name:   name,
		Params: params,
	})
	return fc
}

// BuildVideoFilterString строит строку видео фильтров для FFmpeg
func (fc *FilterChain) BuildVideoFilterString() string {
	if len(fc.VideoFilters) == 0 {
		return ""
	}

	var filters []string
	for _, filter := range fc.VideoFilters {
		filterStr := filter.Name
		if len(filter.Params) > 0 {
			var params []string
			for key, value := range filter.Params {
				if value == "" {
					params = append(params, key)
				} else {
					params = append(params, fmt.Sprintf("%s=%s", key, value))
				}
			}
			filterStr += "=" + strings.Join(params, ":")
		}
		filters = append(filters, filterStr)
	}

	return strings.Join(filters, ",")
}

// BuildAudioFilterString строит строку аудио фильтров для FFmpeg
func (fc *FilterChain) BuildAudioFilterString() string {
	if len(fc.AudioFilters) == 0 {
		return ""
	}

	var filters []string
	for _, filter := range fc.AudioFilters {
		filterStr := filter.Name
		if len(filter.Params) > 0 {
			var params []string
			for key, value := range filter.Params {
				if value == "" {
					params = append(params, key)
				} else {
					params = append(params, fmt.Sprintf("%s=%s", key, value))
				}
			}
			filterStr += "=" + strings.Join(params, ":")
		}
		filters = append(filters, filterStr)
	}

	return strings.Join(filters, ",")
}

// Предустановленные фильтры

// ScaleFilter масштабирование видео
func ScaleFilter(width, height int) Filter {
	return Filter{
		Name: "scale",
		Params: map[string]string{
			"w": fmt.Sprintf("%d", width),
			"h": fmt.Sprintf("%d", height),
		},
	}
}

// CropFilter обрезка видео
func CropFilter(width, height, x, y int) Filter {
	return Filter{
		Name: "crop",
		Params: map[string]string{
			"w": fmt.Sprintf("%d", width),
			"h": fmt.Sprintf("%d", height),
			"x": fmt.Sprintf("%d", x),
			"y": fmt.Sprintf("%d", y),
		},
	}
}

// RotateFilter поворот видео
func RotateFilter(angle string) Filter {
	return Filter{
		Name: "rotate",
		Params: map[string]string{
			"angle": angle,
		},
	}
}

// BlurFilter размытие
func BlurFilter(sigma float64) Filter {
	return Filter{
		Name: "gblur",
		Params: map[string]string{
			"sigma": fmt.Sprintf("%.2f", sigma),
		},
	}
}

// SharpenFilter повышение резкости
func SharpenFilter(amount float64) Filter {
	return Filter{
		Name: "unsharp",
		Params: map[string]string{
			"luma_msize_x":   "5",
			"luma_msize_y":   "5",
			"luma_amount":    fmt.Sprintf("%.2f", amount),
			"chroma_msize_x": "5",
			"chroma_msize_y": "5",
			"chroma_amount":  fmt.Sprintf("%.2f", amount*0.5),
		},
	}
}

// BrightnessContrastFilter яркость и контрастность
func BrightnessContrastFilter(brightness, contrast float64) Filter {
	return Filter{
		Name: "eq",
		Params: map[string]string{
			"brightness": fmt.Sprintf("%.2f", brightness),
			"contrast":   fmt.Sprintf("%.2f", contrast),
		},
	}
}

// SaturationFilter насыщенность
func SaturationFilter(saturation float64) Filter {
	return Filter{
		Name: "eq",
		Params: map[string]string{
			"saturation": fmt.Sprintf("%.2f", saturation),
		},
	}
}

// FadeInFilter плавное появление
func FadeInFilter(duration float64) Filter {
	return Filter{
		Name: "fade",
		Params: map[string]string{
			"type":     "in",
			"duration": fmt.Sprintf("%.2f", duration),
		},
	}
}

// FadeOutFilter плавное исчезновение
func FadeOutFilter(startTime, duration float64) Filter {
	return Filter{
		Name: "fade",
		Params: map[string]string{
			"type":       "out",
			"start_time": fmt.Sprintf("%.2f", startTime),
			"duration":   fmt.Sprintf("%.2f", duration),
		},
	}
}

// WatermarkFilter водяной знак
func WatermarkFilter(overlayPath string, x, y int, opacity float64) Filter {
	return Filter{
		Name: "overlay",
		Params: map[string]string{
			"x": fmt.Sprintf("%d", x),
			"y": fmt.Sprintf("%d", y),
		},
	}
}

// Аудио фильтры

// VolumeFilter громкость
func VolumeFilter(volume float64) Filter {
	return Filter{
		Name: "volume",
		Params: map[string]string{
			"volume": fmt.Sprintf("%.2f", volume),
		},
	}
}

// AudioFadeInFilter плавное появление звука
func AudioFadeInFilter(duration float64) Filter {
	return Filter{
		Name: "afade",
		Params: map[string]string{
			"type":     "in",
			"duration": fmt.Sprintf("%.2f", duration),
		},
	}
}

// AudioFadeOutFilter плавное исчезновение звука
func AudioFadeOutFilter(startTime, duration float64) Filter {
	return Filter{
		Name: "afade",
		Params: map[string]string{
			"type":       "out",
			"start_time": fmt.Sprintf("%.2f", startTime),
			"duration":   fmt.Sprintf("%.2f", duration),
		},
	}
}

// HighpassFilter высокочастотный фильтр
func HighpassFilter(frequency int) Filter {
	return Filter{
		Name: "highpass",
		Params: map[string]string{
			"f": fmt.Sprintf("%d", frequency),
		},
	}
}

// LowpassFilter низкочастотный фильтр
func LowpassFilter(frequency int) Filter {
	return Filter{
		Name: "lowpass",
		Params: map[string]string{
			"f": fmt.Sprintf("%d", frequency),
		},
	}
}

// NoiseReductionFilter шумоподавление
func NoiseReductionFilter(strength float64) Filter {
	return Filter{
		Name: "anlmdn",
		Params: map[string]string{
			"s": fmt.Sprintf("%.2f", strength),
		},
	}
}

// ExecuteWithFilters выполняет транскодирование с фильтрами
func (t *Transcoder) ExecuteWithFilters(ctx context.Context, job *dto.Job, filterChain *FilterChain) error {
	// Валидируем конфигурацию
	if err := job.Config.Validate(); err != nil {
		job.Status = dto.StatusFailed
		job.Error = err
		t.logger.Error("Ошибка валидации конфигурации: %v", err)
		return fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	job.Status = dto.StatusRunning
	job.StartTime = time.Now()

	t.logger.Info("Начало транскодирования с фильтрами: %s -> %s", job.Config.InputPath, job.Config.OutputPath)

	// Строим аргументы с фильтрами
	args := t.buildFFmpegArgsWithFilters(job.Config, filterChain)
	t.logger.Debug("FFmpeg аргументы с фильтрами: %v", args)

	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)

	if err := cmd.Run(); err != nil {
		job.Status = dto.StatusFailed
		job.Error = err
		job.EndTime = time.Now()
		t.logger.Error("Ошибка выполнения FFmpeg с фильтрами: %v", err)
		return fmt.Errorf("ошибка транскодирования с фильтрами: %w", err)
	}

	job.Status = dto.StatusCompleted
	job.Progress = 100.0
	job.EndTime = time.Now()

	duration := job.EndTime.Sub(job.StartTime)
	t.logger.Info("Транскодирование с фильтрами завершено за %v", duration)

	return nil
}

// buildFFmpegArgsWithFilters строит аргументы FFmpeg с фильтрами
func (t *Transcoder) buildFFmpegArgsWithFilters(config dto.Config, filterChain *FilterChain) []string {
	args := []string{
		"-i", config.InputPath,
		"-y", // перезаписывать выходной файл
	}

	// Добавляем видео фильтры
	if videoFilters := filterChain.BuildVideoFilterString(); videoFilters != "" {
		args = append(args, "-vf", videoFilters)
	}

	// Добавляем аудио фильтры
	if audioFilters := filterChain.BuildAudioFilterString(); audioFilters != "" {
		args = append(args, "-af", audioFilters)
	}

	// Остальные параметры как обычно
	if config.VideoCodec != "" {
		args = append(args, "-c:v", config.VideoCodec)
	}

	if config.AudioCodec != "" {
		args = append(args, "-c:a", config.AudioCodec)
	}

	if config.VideoBitrate != "" {
		args = append(args, "-b:v", config.VideoBitrate)
	}

	if config.AudioBitrate != "" {
		args = append(args, "-b:a", config.AudioBitrate)
	}

	if config.FrameRate != "" {
		args = append(args, "-r", config.FrameRate)
	}

	if config.Quality != "" {
		args = append(args, "-crf", config.Quality)
	}

	if config.Format != "" {
		args = append(args, "-f", config.Format)
	}

	args = append(args, config.OutputPath)
	return args
}
