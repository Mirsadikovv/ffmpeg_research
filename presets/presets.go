package presets

import "github.com/Mirsadikovv/ffmpeg_research/dto"

// Предустановленные конфигурации
var (
	// Веб-оптимизированные пресеты
	WebHD = dto.Preset{
		Name:        "web-hd",
		Description: "HD качество для веб (1280x720, H.264)",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1280x720",
			VideoBitrate: "2500k",
			AudioBitrate: "128k",
			Quality:      "23",
		},
	}

	WebSD = dto.Preset{
		Name:        "web-sd",
		Description: "SD качество для веб (854x480, H.264)",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "854x480",
			VideoBitrate: "1000k",
			AudioBitrate: "96k",
			Quality:      "25",
		},
	}

	// Мобильные пресеты
	Mobile = dto.Preset{
		Name:        "mobile",
		Description: "Оптимизировано для мобильных устройств",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "640x360",
			VideoBitrate: "500k",
			AudioBitrate: "64k",
			Quality:      "28",
		},
	}

	// Аудио пресеты
	AudioMP3 = dto.Preset{
		Name:        "audio-mp3",
		Description: "Конвертация в MP3",
		Config: dto.Config{
			AudioCodec:   "mp3",
			AudioBitrate: "192k",
			Format:       "mp3",
		},
	}

	AudioAAC = dto.Preset{
		Name:        "audio-aac",
		Description: "Конвертация в AAC",
		Config: dto.Config{
			AudioCodec:   "aac",
			AudioBitrate: "128k",
			Format:       "m4a",
		},
	}

	// 4K пресеты
	FourK = dto.Preset{
		Name:        "4k",
		Description: "4K качество (3840x2160, H.264)",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "3840x2160",
			VideoBitrate: "15000k",
			AudioBitrate: "192k",
			Quality:      "20",
		},
	}

	// Streaming пресеты
	Twitch = dto.Preset{
		Name:        "twitch",
		Description: "Оптимизировано для Twitch стрима",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "6000k",
			AudioBitrate: "160k",
			FrameRate:    "60",
			Quality:      "23",
		},
	}

	YouTube = dto.Preset{
		Name:        "youtube",
		Description: "Оптимизировано для YouTube",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "8000k",
			AudioBitrate: "192k",
			FrameRate:    "30",
			Quality:      "21",
		},
	}

	// Архивные пресеты
	Archive = dto.Preset{
		Name:        "archive",
		Description: "Высокое качество для архивирования",
		Config: dto.Config{
			VideoCodec: "libx265",
			AudioCodec: "flac",
			Quality:    "18",
		},
	}

	// Новые продвинутые пресеты
	WebOptimized = dto.Preset{
		Name:        "web-optimized",
		Description: "Оптимизировано для веб с быстрой загрузкой",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "3000k",
			AudioBitrate: "128k",
			Quality:      "22",
			Format:       "mp4",
		},
	}

	HighQuality = dto.Preset{
		Name:        "high-quality",
		Description: "Высокое качество с H.265",
		Config: dto.Config{
			VideoCodec:   "libx265",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "4000k",
			AudioBitrate: "192k",
			Quality:      "20",
		},
	}

	FastEncode = dto.Preset{
		Name:        "fast-encode",
		Description: "Быстрое кодирование с приемлемым качеством",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1280x720",
			VideoBitrate: "1500k",
			AudioBitrate: "96k",
			Quality:      "26",
		},
	}

	SmallSize = dto.Preset{
		Name:        "small-size",
		Description: "Минимальный размер файла",
		Config: dto.Config{
			VideoCodec:   "libx265",
			AudioCodec:   "aac",
			Resolution:   "854x480",
			VideoBitrate: "800k",
			AudioBitrate: "64k",
			Quality:      "28",
		},
	}

	// Специализированные пресеты
	Animation = dto.Preset{
		Name:        "animation",
		Description: "Оптимизировано для анимации",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "5000k",
			AudioBitrate: "192k",
			Quality:      "18",
			FrameRate:    "24",
		},
	}

	Gaming = dto.Preset{
		Name:        "gaming",
		Description: "Оптимизировано для игрового контента",
		Config: dto.Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "6000k",
			AudioBitrate: "160k",
			Quality:      "21",
			FrameRate:    "60",
		},
	}

	Podcast = dto.Preset{
		Name:        "podcast",
		Description: "Оптимизировано для подкастов",
		Config: dto.Config{
			AudioCodec:   "libmp3lame",
			AudioBitrate: "128k",
			Format:       "mp3",
		},
	}

	AudiobookMP3 = dto.Preset{
		Name:        "audiobook-mp3",
		Description: "Оптимизировано для аудиокниг (MP3)",
		Config: dto.Config{
			AudioCodec:   "libmp3lame",
			AudioBitrate: "64k",
			Format:       "mp3",
		},
	}

	AudiobookM4A = dto.Preset{
		Name:        "audiobook-m4a",
		Description: "Оптимизировано для аудиокниг (M4A)",
		Config: dto.Config{
			AudioCodec:   "aac",
			AudioBitrate: "64k",
			Format:       "m4a",
		},
	}
)

// GetPreset возвращает пресет по имени
func GetPreset(name string) (*dto.Preset, bool) {
	presets := map[string]dto.Preset{
		"web-hd":        WebHD,
		"web-sd":        WebSD,
		"mobile":        Mobile,
		"audio-mp3":     AudioMP3,
		"audio-aac":     AudioAAC,
		"4k":            FourK,
		"twitch":        Twitch,
		"youtube":       YouTube,
		"archive":       Archive,
		"web-optimized": WebOptimized,
		"high-quality":  HighQuality,
		"fast-encode":   FastEncode,
		"small-size":    SmallSize,
		"animation":     Animation,
		"gaming":        Gaming,
		"podcast":       Podcast,
		"audiobook-mp3": AudiobookMP3,
		"audiobook-m4a": AudiobookM4A,
	}

	preset, exists := presets[name]
	return &preset, exists
}

// ListPresets возвращает список всех пресетов
func ListPresets() []dto.Preset {
	return []dto.Preset{
		WebHD,
		WebSD,
		Mobile,
		AudioMP3,
		AudioAAC,
		FourK,
		Twitch,
		YouTube,
		Archive,
		WebOptimized,
		HighQuality,
		FastEncode,
		SmallSize,
		Animation,
		Gaming,
		Podcast,
		AudiobookMP3,
		AudiobookM4A,
	}
}

// GetPresetsByCategory возвращает пресеты по категориям
func GetPresetsByCategory() map[string][]dto.Preset {
	return map[string][]dto.Preset{
		"Веб и мобильные": {
			WebHD,
			WebSD,
			WebOptimized,
			Mobile,
		},
		"Высокое качество": {
			FourK,
			HighQuality,
			Archive,
		},
		"Стриминг": {
			Twitch,
			YouTube,
			Gaming,
		},
		"Аудио": {
			AudioMP3,
			AudioAAC,
			Podcast,
			AudiobookMP3,
			AudiobookM4A,
		},
		"Специализированные": {
			Animation,
			FastEncode,
			SmallSize,
		},
	}
}
