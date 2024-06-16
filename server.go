package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	TelegramAPIToken string `json:"TELEGRAM_API_TOKEN"`
	ChatID           int    `json:"CHAT_ID"`
	ApiKey           string `json:"API_KEY"`
	Port             int    `json:"PORT"`
	BotName          string `json:"BOT_NAME"`
	BotShortName     string `json:"BOT_SHORT_NAME"`
	BotUrl           string `json:"BOT_URL"`
}

type PushNotificationRequest struct {
	Message string `json:"message"`
	ApiKey  string `json:"apikey"`
}

func getConfig() (Config, error) {

	file, err := os.Open("telegram.config.json")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func sendTelegramNotification(message string, config Config) error {

	client := &http.Client{}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramAPIToken)
	data := url.Values{
		"chat_id":    {strconv.Itoa(config.ChatID)},
		"text":       {message},
		"parse_mode": {"Markdown"},
	}
	response, err := client.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	fmt.Println(apiURL, response.Status)
	response.Body.Close()

	return nil
}

func handleTelegramMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var request PushNotificationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := getConfig()
	if err != nil {
		fmt.Fprintf(w, "Failed to read config: %s", err)
		return
	}

	if request.ApiKey == config.ApiKey {
		sendTelegramNotification(request.Message, config)
		fmt.Fprintf(w, "Push notification sent: %s", request.Message)
		return
	} else {
		fmt.Fprintf(w, "ApiKey is incorrect: %s", request.ApiKey)
		return
	}

}

func main() {

	config, err := getConfig()
	if err != nil {
		log.Panicf("Failed to read config: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/send-push-notification-to-telegram", handleTelegramMessage)

	fmt.Println("Server is started...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
