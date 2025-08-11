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

// GetPreset возвращает пресет по имени
func GetPreset(name string) (*Preset, bool) {
	presets := map[string]Preset{
		"web-hd":    PresetWebHD,
		"web-sd":    PresetWebSD,
		"mobile":    PresetMobile,
		"audio-mp3": PresetAudioMP3,
		"audio-aac": PresetAudioAAC,
	}

	preset, exists := presets[name]
	return &preset, exists
}

// ListPresets возвращает список всех доступных пресетов
func ListPresets() []Preset {
	return []Preset{
		PresetWebHD,
		PresetWebSD,
		PresetMobile,
		PresetAudioMP3,
		PresetAudioAAC,
	}
}
