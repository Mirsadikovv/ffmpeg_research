package transcoder

import (
	"regexp"
	"strconv"
	"time"
)

// ProgressCallback функция для отслеживания прогресса
type ProgressCallback func(progress float64, speed string, eta time.Duration)

// ProgressTracker отслеживает прогресс транскодирования
type ProgressTracker struct {
	callback    ProgressCallback
	totalFrames int64
	duration    time.Duration
	startTime   time.Time
	frameRegex  *regexp.Regexp
	timeRegex   *regexp.Regexp
	speedRegex  *regexp.Regexp
}

// NewProgressTracker создает новый трекер прогресса
func NewProgressTracker(callback ProgressCallback) *ProgressTracker {
	return &ProgressTracker{
		callback:   callback,
		frameRegex: regexp.MustCompile(`frame=\s*(\d+)`),
		timeRegex:  regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`),
		speedRegex: regexp.MustCompile(`speed=\s*([\d.]+)x`),
		startTime:  time.Now(),
	}
}

// ParseFFmpegOutput парсит вывод FFmpeg для извлечения прогресса
func (pt *ProgressTracker) ParseFFmpegOutput(line string) {
	if pt.callback == nil {
		return
	}

	var progress float64
	var speed string
	var eta time.Duration

	// Извлекаем текущее время обработки
	if matches := pt.timeRegex.FindStringSubmatch(line); len(matches) == 5 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		seconds, _ := strconv.Atoi(matches[3])
		centiseconds, _ := strconv.Atoi(matches[4])

		currentTime := time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second +
			time.Duration(centiseconds)*time.Millisecond*10

		if pt.duration > 0 {
			progress = float64(currentTime) / float64(pt.duration) * 100
			if progress > 100 {
				progress = 100
			}
		}
	}

	// Извлекаем скорость
	if matches := pt.speedRegex.FindStringSubmatch(line); len(matches) == 2 {
		speed = matches[1] + "x"

		// Вычисляем ETA
		if speedFloat, err := strconv.ParseFloat(matches[1], 64); err == nil && speedFloat > 0 {
			elapsed := time.Since(pt.startTime)
			if progress > 0 {
				totalEstimated := time.Duration(float64(elapsed) * 100 / progress)
				eta = totalEstimated - elapsed
			}
		}
	}

	if progress > 0 || speed != "" {
		pt.callback(progress, speed, eta)
	}
}

// SetDuration устанавливает общую продолжительность для расчета прогресса
func (pt *ProgressTracker) SetDuration(duration time.Duration) {
	pt.duration = duration
}
