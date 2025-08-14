package dto

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("поле '%s': %s", e.Field, e.Message)
}

// ValidationErrors коллекция ошибок валидации
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors проверяет наличие ошибок
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Validate валидирует конфигурацию транскодирования
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Проверка входного файла
	if c.InputPath == "" {
		errors = append(errors, ValidationError{
			Field:   "InputPath",
			Message: "путь к входному файлу не может быть пустым",
		})
	} else if !fileExists(c.InputPath) {
		errors = append(errors, ValidationError{
			Field:   "InputPath",
			Message: fmt.Sprintf("файл '%s' не существует", c.InputPath),
		})
	}

	// Проверка выходного файла
	if c.OutputPath == "" {
		errors = append(errors, ValidationError{
			Field:   "OutputPath",
			Message: "путь к выходному файлу не может быть пустым",
		})
	} else {
		// Проверяем, что директория для выходного файла существует
		outputDir := filepath.Dir(c.OutputPath)
		if !dirExists(outputDir) {
			errors = append(errors, ValidationError{
				Field:   "OutputPath",
				Message: fmt.Sprintf("директория '%s' не существует", outputDir),
			})
		}
	}

	// Валидация видео кодека
	if c.VideoCodec != "" {
		if err := validateVideoCodec(c.VideoCodec); err != nil {
			errors = append(errors, ValidationError{
				Field:   "VideoCodec",
				Message: err.Error(),
			})
		}
	}

	// Валидация аудио кодека
	if c.AudioCodec != "" {
		if err := validateAudioCodec(c.AudioCodec); err != nil {
			errors = append(errors, ValidationError{
				Field:   "AudioCodec",
				Message: err.Error(),
			})
		}
	}

	// Валидация битрейта видео
	if c.VideoBitrate != "" {
		if err := validateBitrate(c.VideoBitrate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "VideoBitrate",
				Message: err.Error(),
			})
		}
	}

	// Валидация битрейта аудио
	if c.AudioBitrate != "" {
		if err := validateBitrate(c.AudioBitrate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "AudioBitrate",
				Message: err.Error(),
			})
		}
	}

	// Валидация разрешения
	if c.Resolution != "" {
		if err := validateResolution(c.Resolution); err != nil {
			errors = append(errors, ValidationError{
				Field:   "Resolution",
				Message: err.Error(),
			})
		}
	}

	// Валидация частоты кадров
	if c.FrameRate != "" {
		if err := validateFrameRate(c.FrameRate); err != nil {
			errors = append(errors, ValidationError{
				Field:   "FrameRate",
				Message: err.Error(),
			})
		}
	}

	// Валидация качества (CRF)
	if c.Quality != "" {
		if err := validateQuality(c.Quality); err != nil {
			errors = append(errors, ValidationError{
				Field:   "Quality",
				Message: err.Error(),
			})
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// Validate валидирует конфигурацию HLS
func (h *HLSConfig) Validate() error {
	var errors ValidationErrors

	// Проверка URL
	if h.URL == "" {
		errors = append(errors, ValidationError{
			Field:   "URL",
			Message: "URL не может быть пустым",
		})
	} else if !isValidURL(h.URL) {
		errors = append(errors, ValidationError{
			Field:   "URL",
			Message: "некорректный URL",
		})
	}

	// Проверка выходного пути
	if h.OutputPath == "" {
		errors = append(errors, ValidationError{
			Field:   "OutputPath",
			Message: "путь к выходному файлу не может быть пустым",
		})
	}

	// Проверка качества
	if h.Quality != "" && !isValidQuality(h.Quality) {
		errors = append(errors, ValidationError{
			Field:   "Quality",
			Message: "некорректное значение качества (используйте 'best', 'worst' или разрешение вида '1920x1080')",
		})
	}

	// Проверка продолжительности
	if h.Duration < 0 {
		errors = append(errors, ValidationError{
			Field:   "Duration",
			Message: "продолжительность не может быть отрицательной",
		})
	}

	// Проверка количества попыток
	if h.RetryAttempts < 0 {
		errors = append(errors, ValidationError{
			Field:   "RetryAttempts",
			Message: "количество попыток не может быть отрицательным",
		})
	}

	// Проверка таймаута сегментов
	if h.SegmentTimeout < 0 {
		errors = append(errors, ValidationError{
			Field:   "SegmentTimeout",
			Message: "таймаут сегментов не может быть отрицательным",
		})
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// Вспомогательные функции валидации

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

func validateVideoCodec(codec string) error {
	validCodecs := []string{
		"libx264", "libx265", "libvpx", "libvpx-vp9", "libav1",
		"h264", "hevc", "vp8", "vp9", "av1", "copy",
	}

	for _, valid := range validCodecs {
		if codec == valid {
			return nil
		}
	}

	return fmt.Errorf("неподдерживаемый видео кодек '%s'", codec)
}

func validateAudioCodec(codec string) error {
	validCodecs := []string{
		"aac", "libmp3lame", "libopus", "libvorbis", "flac",
		"mp3", "opus", "vorbis", "pcm_s16le", "copy",
	}

	for _, valid := range validCodecs {
		if codec == valid {
			return nil
		}
	}

	return fmt.Errorf("неподдерживаемый аудио кодек '%s'", codec)
}

func validateBitrate(bitrate string) error {
	// Проверяем формат битрейта (например: "1000k", "2M", "128")
	re := regexp.MustCompile(`^(\d+)([kKmM]?)$`)
	matches := re.FindStringSubmatch(bitrate)

	if len(matches) != 3 {
		return fmt.Errorf("некорректный формат битрейта '%s' (используйте формат: 1000k, 2M, 128)", bitrate)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil || value <= 0 {
		return fmt.Errorf("некорректное значение битрейта '%s'", bitrate)
	}

	return nil
}

func validateResolution(resolution string) error {
	// Проверяем формат разрешения (например: "1920x1080", "1280x720")
	re := regexp.MustCompile(`^(\d+)x(\d+)$`)
	matches := re.FindStringSubmatch(resolution)

	if len(matches) != 3 {
		return fmt.Errorf("некорректный формат разрешения '%s' (используйте формат: 1920x1080)", resolution)
	}

	width, err1 := strconv.Atoi(matches[1])
	height, err2 := strconv.Atoi(matches[2])

	if err1 != nil || err2 != nil || width <= 0 || height <= 0 {
		return fmt.Errorf("некорректные значения разрешения '%s'", resolution)
	}

	// Проверяем разумные пределы
	if width > 7680 || height > 4320 { // 8K максимум
		return fmt.Errorf("разрешение '%s' слишком большое (максимум 7680x4320)", resolution)
	}

	if width < 64 || height < 64 { // Минимальное разрешение
		return fmt.Errorf("разрешение '%s' слишком маленькое (минимум 64x64)", resolution)
	}

	return nil
}

func validateFrameRate(frameRate string) error {
	rate, err := strconv.ParseFloat(frameRate, 64)
	if err != nil {
		return fmt.Errorf("некорректное значение частоты кадров '%s'", frameRate)
	}

	if rate <= 0 {
		return fmt.Errorf("частота кадров должна быть положительной")
	}

	if rate > 120 {
		return fmt.Errorf("частота кадров '%s' слишком высокая (максимум 120)", frameRate)
	}

	return nil
}

func validateQuality(quality string) error {
	crf, err := strconv.Atoi(quality)
	if err != nil {
		return fmt.Errorf("некорректное значение качества '%s' (используйте число от 0 до 51)", quality)
	}

	if crf < 0 || crf > 51 {
		return fmt.Errorf("значение качества '%s' вне допустимого диапазона (0-51)", quality)
	}

	return nil
}

func isValidURL(url string) bool {
	// Простая проверка URL
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidQuality(quality string) bool {
	if quality == "best" || quality == "worst" {
		return true
	}

	// Проверяем, является ли это разрешением
	re := regexp.MustCompile(`^\d+x\d+$`)
	return re.MatchString(quality)
}
