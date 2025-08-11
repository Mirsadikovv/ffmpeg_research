package transcoder

// Preset представляет предустановленную конфигурацию
type Preset struct {
	Name        string
	Description string
	Config      Config
}

// Предустановленные конфигурации
var (
	// Веб-оптимизированные пресеты
	PresetWebHD = Preset{
		Name:        "web-hd",
		Description: "HD качество для веб (1280x720, H.264)",
		Config: Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1280x720",
			VideoBitrate: "2500k",
			AudioBitrate: "128k",
			Quality:      "23",
		},
	}

	PresetWebSD = Preset{
		Name:        "web-sd",
		Description: "SD качество для веб (854x480, H.264)",
		Config: Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "854x480",
			VideoBitrate: "1000k",
			AudioBitrate: "96k",
			Quality:      "25",
		},
	}

	// Мобильные пресеты
	PresetMobile = Preset{
		Name:        "mobile",
		Description: "Оптимизировано для мобильных устройств",
		Config: Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "640x360",
			VideoBitrate: "500k",
			AudioBitrate: "64k",
			Quality:      "28",
		},
	}

	// Аудио пресеты
	PresetAudioMP3 = Preset{
		Name:        "audio-mp3",
		Description: "Конвертация в MP3",
		Config: Config{
			AudioCodec:   "mp3",
			AudioBitrate: "192k",
			Format:       "mp3",
		},
	}

	PresetAudioAAC = Preset{
		Name:        "audio-aac",
		Description: "Конвертация в AAC",
		Config: Config{
			AudioCodec:   "aac",
			AudioBitrate: "128k",
			Format:       "m4a",
		},
	}
)

// Дополнительные пресеты
var (
	// 4K пресеты
	Preset4K = Preset{
		Name:        "4k",
		Description: "4K качество (3840x2160, H.264)",
		Config: Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "3840x2160",
			VideoBitrate: "15000k",
			AudioBitrate: "192k",
			Quality:      "20",
		},
	}

	// Streaming пресеты
	PresetTwitch = Preset{
		Name:        "twitch",
		Description: "Оптимизировано для Twitch стрима",
		Config: Config{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			Resolution:   "1920x1080",
			VideoBitrate: "6000k",
			AudioBitrate: "160k",
			FrameRate:    "60",
			Quality:      "23",
		},
	}

	PresetYouTube = Preset{
		Name:        "youtube",
		Description: "Оптимизировано для YouTube",
		Config: Config{
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
	PresetArchive = Preset{
		Name:        "archive",
		Description: "Высокое качество для архивирования",
		Config: Config{
			VideoCodec: "libx265",
			AudioCodec: "flac",
			Quality:    "18",
		},
	}
)

// GetPreset обновленная версия с новыми пресетами
func GetPreset(name string) (*Preset, bool) {
	presets := map[string]Preset{
		"web-hd":    PresetWebHD,
		"web-sd":    PresetWebSD,
		"mobile":    PresetMobile,
		"audio-mp3": PresetAudioMP3,
		"audio-aac": PresetAudioAAC,
		"4k":        Preset4K,
		"twitch":    PresetTwitch,
		"youtube":   PresetYouTube,
		"archive":   PresetArchive,
	}

	preset, exists := presets[name]
	return &preset, exists
}

// ListPresets обновленная версия
func ListPresets() []Preset {
	return []Preset{
		PresetWebHD,
		PresetWebSD,
		PresetMobile,
		PresetAudioMP3,
		PresetAudioAAC,
		Preset4K,
		PresetTwitch,
		PresetYouTube,
		PresetArchive,
	}
}
