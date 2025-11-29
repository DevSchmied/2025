package check

import (
	"net/http"
	"strings"
	"time"
)

func CheckLink(url string) bool {
	// Проверяем, есть ли у URL префикс http/https. Это необходимо для дальнейшей работы.
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Создаём HTTP-клиент с таймаутом 2 секунды.
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// Проверяем, отвечает ли сайт на запрос.
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	// Закрываем тело ответа, чтобы избежать утечки ресурсов.
	defer resp.Body.Close()

	return resp.StatusCode < 400
}
