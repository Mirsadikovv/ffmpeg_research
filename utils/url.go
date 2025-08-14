package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
)

// ResolveURL разрешает относительные URL
func ResolveURL(baseURL, relativeURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return relativeURL
	}

	rel, err := url.Parse(relativeURL)
	if err != nil {
		return relativeURL
	}

	return base.ResolveReference(rel).String()
}

// FindPlaylistURL пытается найти плейлист по URL страницы
func FindPlaylistURL(pageURL string, client *http.Client) (string, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Ищем ссылки на .m3u8 файлы
	re := regexp.MustCompile(`(https?://[^\s"']+\.m3u8[^\s"']*)`)
	matches := re.FindAllString(string(content), -1)

	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("плейлист не найден на странице")
}

// ParsePlaylist парсит содержимое плейлиста
func ParsePlaylist(content, baseURL string) (*dto.PlaylistInfo, error) {
	info := &dto.PlaylistInfo{
		URL:     baseURL,
		Streams: make([]dto.StreamInfo, 0),
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentStream dto.StreamInfo

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			// Парсим информацию о потоке
			currentStream = parseStreamInfo(line)
		} else if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			// Это live стрим
			info.IsLive = true
		} else if !strings.HasPrefix(line, "#") && line != "" {
			// Это URL потока
			if currentStream.URL == "" {
				currentStream.URL = ResolveURL(baseURL, line)
				info.Streams = append(info.Streams, currentStream)
				currentStream = dto.StreamInfo{}
			}
		}
	}

	return info, nil
}

// parseStreamInfo парсит информацию о потоке из строки EXT-X-STREAM-INF
func parseStreamInfo(line string) dto.StreamInfo {
	stream := dto.StreamInfo{}

	// Парсим bandwidth
	if match := regexp.MustCompile(`BANDWIDTH=(\d+)`).FindStringSubmatch(line); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &stream.Bandwidth)
	}

	// Парсим разрешение
	if match := regexp.MustCompile(`RESOLUTION=(\d+x\d+)`).FindStringSubmatch(line); len(match) > 1 {
		stream.Resolution = match[1]
	}

	// Парсим кодеки
	if match := regexp.MustCompile(`CODECS="([^"]+)"`).FindStringSubmatch(line); len(match) > 1 {
		stream.Codecs = match[1]
	}

	// Парсим frame rate
	if match := regexp.MustCompile(`FRAME-RATE=([\d.]+)`).FindStringSubmatch(line); len(match) > 1 {
		fmt.Sscanf(match[1], "%f", &stream.FrameRate)
	}

	return stream
}

// SelectStream выбирает поток по качеству
func SelectStream(info *dto.PlaylistInfo, quality string) (string, error) {
	if len(info.Streams) == 0 {
		return info.URL, nil // Возвращаем исходный URL если нет вариантов
	}

	switch quality {
	case "best", "":
		// Выбираем поток с максимальным bandwidth
		var bestStream dto.StreamInfo
		for _, stream := range info.Streams {
			if stream.Bandwidth > bestStream.Bandwidth {
				bestStream = stream
			}
		}
		return bestStream.URL, nil

	case "worst":
		// Выбираем поток с минимальным bandwidth
		bestStream := info.Streams[0]
		for _, stream := range info.Streams {
			if stream.Bandwidth < bestStream.Bandwidth {
				bestStream = stream
			}
		}
		return bestStream.URL, nil

	default:
		// Ищем поток с конкретным разрешением
		for _, stream := range info.Streams {
			if stream.Resolution == quality {
				return stream.URL, nil
			}
		}
		return "", fmt.Errorf("поток с качеством %s не найден", quality)
	}
}

// CreateHTTPClient создает HTTP клиент с настройками
func CreateHTTPClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &http.Client{
		Timeout: timeout,
	}
}
