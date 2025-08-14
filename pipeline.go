package transcoder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
)

// PipelineStep представляет шаг в конвейере обработки
type PipelineStep interface {
	Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (outputPath string, err error)
	GetName() string
	GetDescription() string
}

// Pipeline конвейер обработки медиафайлов
type Pipeline struct {
	name        string
	description string
	steps       []PipelineStep
	tempDir     string
	transcoder  *Transcoder
}

// NewPipeline создает новый конвейер
func NewPipeline(name, description string, transcoder *Transcoder) *Pipeline {
	tempDir := filepath.Join(transcoder.tempDir, fmt.Sprintf("pipeline_%d", time.Now().Unix()))
	os.MkdirAll(tempDir, 0755)

	return &Pipeline{
		name:        name,
		description: description,
		steps:       make([]PipelineStep, 0),
		tempDir:     tempDir,
		transcoder:  transcoder,
	}
}

// AddStep добавляет шаг в конвейер
func (p *Pipeline) AddStep(step PipelineStep) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// Execute выполняет весь конвейер
func (p *Pipeline) Execute(ctx context.Context, inputPath, outputPath string) error {
	p.transcoder.logger.Info("Запуск конвейера '%s': %s -> %s", p.name, inputPath, outputPath)

	defer func() {
		// Очищаем временные файлы
		os.RemoveAll(p.tempDir)
	}()

	currentPath := inputPath

	for i, step := range p.steps {
		p.transcoder.logger.Info("Выполнение шага %d/%d: %s", i+1, len(p.steps), step.GetName())

		var nextPath string
		var err error

		if i == len(p.steps)-1 {
			// Последний шаг - используем финальный путь
			nextPath, err = step.Execute(ctx, currentPath, p.transcoder)
			if err != nil {
				return fmt.Errorf("ошибка на шаге '%s': %w", step.GetName(), err)
			}

			// Перемещаем результат в финальное место
			if nextPath != outputPath {
				if err := os.Rename(nextPath, outputPath); err != nil {
					return fmt.Errorf("ошибка перемещения финального файла: %w", err)
				}
			}
		} else {
			// Промежуточный шаг
			nextPath, err = step.Execute(ctx, currentPath, p.transcoder)
			if err != nil {
				return fmt.Errorf("ошибка на шаге '%s': %w", step.GetName(), err)
			}
		}

		currentPath = nextPath
	}

	p.transcoder.logger.Info("Конвейер '%s' завершен успешно", p.name)
	return nil
}

// GetSteps возвращает список шагов
func (p *Pipeline) GetSteps() []PipelineStep {
	return p.steps
}

// GetName возвращает имя конвейера
func (p *Pipeline) GetName() string {
	return p.name
}

// GetDescription возвращает описание конвейера
func (p *Pipeline) GetDescription() string {
	return p.description
}

// Предустановленные шаги конвейера

// TranscodeStep шаг транскодирования
type TranscodeStep struct {
	config dto.Config
}

// NewTranscodeStep создает шаг транскодирования
func NewTranscodeStep(config dto.Config) *TranscodeStep {
	return &TranscodeStep{config: config}
}

func (s *TranscodeStep) Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (string, error) {
	outputPath := filepath.Join(transcoder.tempDir, fmt.Sprintf("transcode_%d.mp4", time.Now().UnixNano()))

	config := s.config
	config.InputPath = inputPath
	config.OutputPath = outputPath

	job := transcoder.CreateJob(config)
	if err := transcoder.Execute(ctx, job); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *TranscodeStep) GetName() string {
	return "Транскодирование"
}

func (s *TranscodeStep) GetDescription() string {
	return fmt.Sprintf("Транскодирование с кодеками %s/%s", s.config.VideoCodec, s.config.AudioCodec)
}

// FilterStep шаг применения фильтров
type FilterStep struct {
	filterChain *FilterChain
	config      dto.Config
}

// NewFilterStep создает шаг применения фильтров
func NewFilterStep(filterChain *FilterChain, config dto.Config) *FilterStep {
	return &FilterStep{
		filterChain: filterChain,
		config:      config,
	}
}

func (s *FilterStep) Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (string, error) {
	outputPath := filepath.Join(transcoder.tempDir, fmt.Sprintf("filter_%d.mp4", time.Now().UnixNano()))

	config := s.config
	config.InputPath = inputPath
	config.OutputPath = outputPath

	job := transcoder.CreateJob(config)
	if err := transcoder.ExecuteWithFilters(ctx, job, s.filterChain); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *FilterStep) GetName() string {
	return "Применение фильтров"
}

func (s *FilterStep) GetDescription() string {
	videoFilters := len(s.filterChain.VideoFilters)
	audioFilters := len(s.filterChain.AudioFilters)
	return fmt.Sprintf("Применение %d видео и %d аудио фильтров", videoFilters, audioFilters)
}

// ThumbnailStep шаг создания миниатюр
type ThumbnailStep struct {
	timeOffsets []string
	outputDir   string
}

// NewThumbnailStep создает шаг создания миниатюр
func NewThumbnailStep(timeOffsets []string, outputDir string) *ThumbnailStep {
	return &ThumbnailStep{
		timeOffsets: timeOffsets,
		outputDir:   outputDir,
	}
}

func (s *ThumbnailStep) Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (string, error) {
	os.MkdirAll(s.outputDir, 0755)

	for i, timeOffset := range s.timeOffsets {
		thumbnailPath := filepath.Join(s.outputDir, fmt.Sprintf("thumbnail_%d.jpg", i+1))
		if err := transcoder.CreateThumbnail(inputPath, thumbnailPath, timeOffset); err != nil {
			return "", fmt.Errorf("ошибка создания миниатюры %s: %w", timeOffset, err)
		}
	}

	// Возвращаем исходный путь, так как это не изменяет основной файл
	return inputPath, nil
}

func (s *ThumbnailStep) GetName() string {
	return "Создание миниатюр"
}

func (s *ThumbnailStep) GetDescription() string {
	return fmt.Sprintf("Создание %d миниатюр", len(s.timeOffsets))
}

// ExtractAudioStep шаг извлечения аудио
type ExtractAudioStep struct {
	outputPath string
	format     string
}

// NewExtractAudioStep создает шаг извлечения аудио
func NewExtractAudioStep(outputPath, format string) *ExtractAudioStep {
	return &ExtractAudioStep{
		outputPath: outputPath,
		format:     format,
	}
}

func (s *ExtractAudioStep) Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (string, error) {
	if err := transcoder.ConvertToFormat(inputPath, s.outputPath, s.format); err != nil {
		return "", err
	}

	// Возвращаем исходный путь, так как это побочная операция
	return inputPath, nil
}

func (s *ExtractAudioStep) GetName() string {
	return "Извлечение аудио"
}

func (s *ExtractAudioStep) GetDescription() string {
	return fmt.Sprintf("Извлечение аудио в формат %s", s.format)
}

// AnalyzeStep шаг анализа медиафайла
type AnalyzeStep struct {
	outputPath string
}

// NewAnalyzeStep создает шаг анализа
func NewAnalyzeStep(outputPath string) *AnalyzeStep {
	return &AnalyzeStep{outputPath: outputPath}
}

func (s *AnalyzeStep) Execute(ctx context.Context, inputPath string, transcoder *Transcoder) (string, error) {
	mediaInfo, err := transcoder.GetMediaInfo(inputPath)
	if err != nil {
		return "", err
	}

	// Сохраняем анализ в файл
	analysisText := fmt.Sprintf("Анализ файла: %s\n", inputPath)
	analysisText += fmt.Sprintf("Разрешение: %s\n", mediaInfo.GetResolution())
	analysisText += fmt.Sprintf("Частота кадров: %.1f fps\n", mediaInfo.GetFrameRate())
	analysisText += fmt.Sprintf("Продолжительность: %.1f секунд\n", mediaInfo.Duration.Seconds())
	analysisText += fmt.Sprintf("Размер: %.2f MB\n", float64(mediaInfo.Size)/(1024*1024))
	analysisText += fmt.Sprintf("Краткая сводка: %s\n", mediaInfo.Summary())

	if err := os.WriteFile(s.outputPath, []byte(analysisText), 0644); err != nil {
		return "", fmt.Errorf("ошибка сохранения анализа: %w", err)
	}

	transcoder.logger.Info("Анализ сохранен в %s", s.outputPath)

	// Возвращаем исходный путь
	return inputPath, nil
}

func (s *AnalyzeStep) GetName() string {
	return "Анализ медиафайла"
}

func (s *AnalyzeStep) GetDescription() string {
	return "Анализ характеристик медиафайла"
}

// Предустановленные конвейеры

// CreateWebOptimizationPipeline создает конвейер для веб-оптимизации
func CreateWebOptimizationPipeline(transcoder *Transcoder) *Pipeline {
	pipeline := NewPipeline(
		"Веб-оптимизация",
		"Полная оптимизация видео для веб с созданием миниатюр и извлечением аудио",
		transcoder,
	)

	// Шаг 1: Анализ исходного файла
	pipeline.AddStep(NewAnalyzeStep("analysis.txt"))

	// Шаг 2: Применение фильтров для улучшения качества
	filterChain := NewFilterChain().
		AddVideoFilter("scale", map[string]string{"w": "1920", "h": "1080"}).
		AddVideoFilter("eq", map[string]string{"brightness": "0.05", "contrast": "1.1"}).
		AddAudioFilter("volume", map[string]string{"volume": "0.9"})

	pipeline.AddStep(NewFilterStep(filterChain, dto.Config{
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Quality:    "22",
	}))

	// Шаг 3: Создание миниатюр
	pipeline.AddStep(NewThumbnailStep(
		[]string{"00:00:05", "00:00:30", "00:01:00"},
		"thumbnails",
	))

	// Шаг 4: Извлечение аудио
	pipeline.AddStep(NewExtractAudioStep("audio.mp3", "mp3"))

	// Шаг 5: Финальное транскодирование для веб
	pipeline.AddStep(NewTranscodeStep(dto.Config{
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "2000k",
		AudioBitrate: "128k",
		Quality:      "23",
		Format:       "mp4",
	}))

	return pipeline
}

// CreateArchivePipeline создает конвейер для архивирования
func CreateArchivePipeline(transcoder *Transcoder) *Pipeline {
	pipeline := NewPipeline(
		"Архивирование",
		"Высококачественное архивирование с анализом и резервными копиями",
		transcoder,
	)

	// Анализ
	pipeline.AddStep(NewAnalyzeStep("archive_analysis.txt"))

	// Создание миниатюр для каталогизации
	pipeline.AddStep(NewThumbnailStep(
		[]string{"00:00:01", "00:01:00", "00:05:00", "00:10:00"},
		"archive_thumbnails",
	))

	// Извлечение аудио в lossless формат
	pipeline.AddStep(NewExtractAudioStep("archive_audio.flac", "flac"))

	// Финальное архивное кодирование
	pipeline.AddStep(NewTranscodeStep(dto.Config{
		VideoCodec: "libx265",
		AudioCodec: "flac",
		Quality:    "18", // Очень высокое качество
		Format:     "mkv",
	}))

	return pipeline
}

// CreateMobilePipeline создает конвейер для мобильных устройств
func CreateMobilePipeline(transcoder *Transcoder) *Pipeline {
	pipeline := NewPipeline(
		"Мобильная оптимизация",
		"Оптимизация для мобильных устройств с минимальным размером файла",
		transcoder,
	)

	// Применение агрессивных фильтров для уменьшения размера
	filterChain := NewFilterChain().
		AddVideoFilter("scale", map[string]string{"w": "854", "h": "480"}).
		AddVideoFilter("eq", map[string]string{"contrast": "1.05"}).
		AddAudioFilter("volume", map[string]string{"volume": "0.95"}).
		AddAudioFilter("lowpass", map[string]string{"f": "15000"}) // Обрезаем высокие частоты

	pipeline.AddStep(NewFilterStep(filterChain, dto.Config{
		VideoCodec: "libx264",
		AudioCodec: "aac",
		Quality:    "28", // Более агрессивное сжатие
	}))

	// Создание одной миниатюры
	pipeline.AddStep(NewThumbnailStep([]string{"00:00:10"}, "mobile_thumb"))

	// Финальное кодирование для мобильных
	pipeline.AddStep(NewTranscodeStep(dto.Config{
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		VideoBitrate: "800k",
		AudioBitrate: "64k",
		Quality:      "28",
		Format:       "mp4",
	}))

	return pipeline
}
