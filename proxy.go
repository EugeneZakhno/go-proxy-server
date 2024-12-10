package main

import (
	"io"
	"log"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Создаем новый HTTP клиент
	client := &http.Client{}

	// Задаем целевой URL (например, http://example.com)
	targetURL := "https://openapi.mypay.kz" + r.URL.Path // Прокси на другой сервер

	// Создаем новый запрос с теми же методами и заголовками
	req, err := http.NewRequest(r.Method, targetURL, nil) // nil для GET-запроса
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из оригинального запроса
	for key, value := range r.Header {
		req.Header[key] = value
	}

	// Подмена заголовка X-Forwarded-For
	if ip := r.RemoteAddr; ip != "" {
		var myIP = "185.209.30.57"
		req.Header.Set("X-Forwarded-For", myIP)
	}

	// Ваш IP-адрес
	//myIP := "192.168.1.1"

	// Подмена заголовка X-Forwarded-For
	//req.Header.Set("X-Forwarded-For", myIP)

	// Отправляем запрос на целевой сервер
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем статус и заголовки ответа
	w.WriteHeader(resp.StatusCode)
	for key, value := range resp.Header {
		w.Header()[key] = value
	}

	// Копируем тело ответа
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Println("Proxy server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
