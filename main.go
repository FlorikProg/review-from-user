package main

// Импортируем библиотеки
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// Обработчик HTTP-запроса
func handlers(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

// Структура данных
type RequestData struct {
	Text    string `json:"text"`
	Rating  string `json:"rating"`
	Product string `json:"product"`
}

// Обработчик POST-запроса
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data RequestData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received data: %+v\n", data)

		fmt.Printf("Тип переменной raiting: %T\n", data.Rating)

		sendToTelegram(data.Text, data.Rating, data.Product)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Функция для отправки сообщения в Telegram
func sendToTelegram(text string, raiting string, product string) {

	token := "7240748730:AAERMXYL4JUcaJUs1B2T4UIz22hquHrr05M"
	chatID := "1529997307"

	message := fmt.Sprintf("<b>Новый отзыв🔥</b>\n\n<b>✅Продукт:</b> %s\n\n<b>⭐️Отзыв:</b> %s\n\n<b>💬Сообщение:</b> %s", product, raiting, text)

	encodedMessage := url.QueryEscape(message)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=HTML", token, chatID, encodedMessage)

	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending message to Telegram:", err)
	}
}

// Основная функция
func main() {
	fs := http.FileServer(http.Dir("./file"))
	http.Handle("/file/", http.StripPrefix("/file/", fs))
	http.HandleFunc("/", handlers)
	http.HandleFunc("/submit", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
