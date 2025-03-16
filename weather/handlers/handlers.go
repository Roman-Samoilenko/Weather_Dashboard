package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func HandleGetWeatherAPIOpenWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		slog.Error("API_KEY не установлен")
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}

	city := r.PathValue("city")
	if city == "" {
		http.Error(w, `{"error": "Город не указан"}`, http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?"+
			"q=%s&appid=%s&units=metric&lang=ru", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Ошибка:", err)
		http.Error(w, "Ошибка получения погоды", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("API вернуло статус:", resp.StatusCode)
		http.Error(w, "Ошибка получения данных о погоде", resp.StatusCode)
		return
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("Ошибка декодирования JSON:", err)
		http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Ошибка кодирования ответа:", err)
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}

type SimpleWeather struct {
	Temp       string `json:"temp"`
	Sunrise    string `json:"sunrise"`
	Sunset     string `json:"sunset"`
	Wind       string `json:"wind"`
	Humidity   string `json:"humidity"`
	Pressure   string `json:"pressure"`
	Cloudiness string `json:"cloudiness"`
}

func HandleGetWeatherParseMailRu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	city := r.PathValue("city")
	if city == "" {
		http.Error(w, `{"error": "Город не указан"}`, http.StatusBadRequest)
		return
	}

	searchURL := "https://pogoda.mail.ru/search/?q=" + url.QueryEscape(city)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		slog.Error("Ошибка с куки:", err)
	}

	// Создаем HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 15 * time.Second,
		Jar:     jar, // Используйте cookie-менеджер
	}

	// Создаем новый запрос
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		slog.Error("Ошибка создания запроса:", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cookie", "CONSENT=YES+cb; SOCS=CAESEwgDEgk0ODE3Nzk3MjQaAmVuIAEaBgiA_LuWBg") // Актуальные куки

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Ошибка HTTP-запроса:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		slog.Error("статус ошибки:", resp.StatusCode, resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Ошибка чтения тела ответа:", err)
	}

	// Проверяем на наличие капчи
	if strings.Contains(string(body), "CAPTCHA") {
		slog.Error("Обнаружена капча. Попробуйте позже или смените IP")
	}

	// Парсим HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		slog.Error("Ошибка парсинга HTML:", err)
	}

	var weather SimpleWeather

	doc.Find(".ab64e36fe5").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".c3132db061.e6255c6329").Text()
		value := s.Find(".c3132db061").Last().Text()
		weather.Cloudiness = "404"
		switch strings.TrimSpace(title) {
		case "Восход":
			weather.Sunrise = value
		case "Заход":
			weather.Sunset = value
		case "Ветер":
			weather.Wind = value
		case "Влажность":
			weather.Humidity = value
		case "Давление":
			weather.Pressure = value
		case "Облачность":
			weather.Cloudiness = value
		}
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		slog.Error("Ошибка кодирования ответа:", err)
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}

}
