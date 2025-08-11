# ПОЛНЫЙ СПРАВОЧНИК FFMPEG

## 1. ОБЩАЯ СТРУКТУРА КОМАНДЫ
```bash
ffmpeg [общие_опции] [опции_входа] -i входной_файл [опции_выхода] выходной_файл
```

## 2. ОБЩИЕ ОПЦИИ
```bash
-version          # Показать версию ffmpeg
-formats          # Список всех форматов ввода/вывода
-codecs           # Список всех кодеков
-encoders         # Список энкодеров
-decoders         # Список декодеров
-filters          # Список фильтров
-pix_fmts         # Список форматов пикселей
-sample_fmts      # Список форматов аудио-сэмплов
-layouts          # Список конфигураций каналов
-devices          # Список устройств захвата
-protocols        # Список протоколов
-h, -help         # Справка
```

## 3. ОПЦИИ ВВОДА
```bash
-i file           # Входной файл или поток
-f fmt            # Принудительный формат ввода
-ss time          # Начало обрезки (до -i — быстрее, после - точнее)
-t duration       # Продолжительность захвата
-to time          # Конечная точка
-stream_loop n    # Повторить файл n раз
-itsoffset time   # Сдвиг времени для входа
-analyzeduration  # Длительность анализа
-probesize size   # Размер анализа
```

## 4. ОПЦИИ ВЫВОДА
```bash
-f fmt            # Формат вывода
-c codec          # Кодек для всех потоков
-c:v codec        # Кодек видео
-c:a codec        # Кодек аудио
-b:v bitrate      # Битрейт видео
-b:a bitrate      # Битрейт аудио
-r fps            # Частота кадров
-s WxH            # Разрешение
-aspect ratio     # Соотношение сторон
-map              # Выбор дорожек
-y                # Перезаписать файл
-n                # Не перезаписывать
-an               # Без аудио
-vn               # Без видео
-sn               # Без субтитров
-preset           # Скорость кодирования (ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow)
-crf value        # Качество (0 — без потерь, 23 — по умолчанию, выше — хуже)
-g frames         # GOP размер
```

## 5. ОСНОВНЫЕ ВИДЕОКОДЕКИ
```bash
libx264           # H.264
libx265           # H.265/HEVC
mpeg4             # MPEG-4
vp8 / libvpx      # VP8
vp9 / libvpx-vp9  # VP9
av1 / libaom-av1  # AV1
libtheora         # Theora
```

## 6. ОСНОВНЫЕ АУДИОКОДЕКИ
```bash
aac               # AAC
libmp3lame        # MP3
libopus           # Opus
libvorbis         # Vorbis
pcm_s16le         # PCM 16-bit little-endian
flac              # FLAC
ac3               # Dolby AC3
```

## 7. ПОПУЛЯРНЫЕ ФИЛЬТРЫ ВИДЕО
```bash
scale=WxH         # Изменить размер
crop=WxH:x:y      # Обрезать
fps=N             # FPS
hflip             # Горизонтальный переворот
vflip             # Вертикальный переворот
rotate=PI/2       # Поворот
transpose=N       # Поворот (0 — 90° по часовой, 1 — 90° против часовой)
drawtext=text     # Текст
overlay           # Наложение видео/картинки
eq=...            # Яркость/контраст/гамма
fade              # Появление/затухание
format=pix_fmt    # Формат пикселей
setsar, setdar    # Pixel Aspect Ratio / Display Aspect Ratio
tblend            # Смешивание кадров
zoompan           # Зум и панорамирование
```

## 8. ПОПУЛЯРНЫЕ ФИЛЬТРЫ АУДИО
```bash
volume=val        # Громкость (1.0 — без изменений)
atrim             # Обрезка
aecho             # Эхо
bass, treble      # НЧ/ВЧ коррекция
silenceremove     # Удаление тишины
pan               # Микширование каналов
asetrate          # Частота дискретизации
aresample         # Ресемплинг
compand           # Компрессор/экспандер
```

## 9. ПОЛЕЗНЫЕ ПРИМЕРЫ

### Конвертировать видео:
```bash
ffmpeg -i in.mov out.mp4
```

### Сжать видео:
```bash
ffmpeg -i in.mp4 -vcodec libx264 -crf 23 out.mp4
```

### Вырезать фрагмент:
```bash
ffmpeg -ss 00:01:00 -t 30 -i in.mp4 -c copy cut.mp4
```

### Извлечь аудио:
```bash
ffmpeg -i in.mp4 -q:a 0 -map a out.mp3
```

### Извлечь кадр:
```bash
ffmpeg -i in.mp4 -ss 00:00:05 -frames:v 1 frame.png
```

### Сделать видео из картинок:
```bash
ffmpeg -framerate 30 -i frame_%03d.png -c:v libx264 out.mp4
```

### Сделать GIF:
```bash
ffmpeg -i in.mp4 -vf "fps=10,scale=320:-1" out.gif
```

### Объединить видео (одинаковые кодеки):
```bash
echo "file 'video1.mp4'" > list.txt
echo "file 'video2.mp4'" >> list.txt
ffmpeg -f concat -safe 0 -i list.txt -c copy out.mp4
```

### Разделить видео на части:
```bash
ffmpeg -i in.mp4 -c copy -map 0 -segment_time 00:05:00 -f segment out%03d.mp4
```

### Изменить разрешение:
```bash
ffmpeg -i in.mp4 -vf scale=1280:720 out.mp4
```

### Добавить водяной знак:
```bash
ffmpeg -i video.mp4 -i watermark.png -filter_complex "overlay=10:10" out.mp4
```

### Конвертировать в разные форматы:
```bash
# В WebM
ffmpeg -i in.mp4 -c:v libvpx-vp9 -c:a libopus out.webm

# В AVI
ffmpeg -i in.mp4 -c:v libx264 -c:a mp3 out.avi

# В MOV
ffmpeg -i in.mp4 -c:v libx264 -c:a aac out.mov
```

### Работа с аудио:
```bash
# Изменить громкость
ffmpeg -i in.mp3 -af "volume=0.5" out.mp3

# Нормализация аудио
ffmpeg -i in.mp3 -af loudnorm out.mp3

# Конвертировать в разные аудиоформаты
ffmpeg -i in.wav -c:a libmp3lame -b:a 192k out.mp3
ffmpeg -i in.mp3 -c:a aac -b:a 128k out.m4a
```

### Потоковое вещание:
```bash
# RTMP стрим
ffmpeg -i in.mp4 -c:v libx264 -c:a aac -f flv rtmp://server/live/stream

# HLS стрим
ffmpeg -i in.mp4 -c:v libx264 -c:a aac -f hls -hls_time 10 -hls_list_size 0 stream.m3u8
```