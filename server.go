package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"

	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"
)

type PushNotificationRequest struct {
	Message string `json:"message"`
	ApiKey  string `json:"apikey"`
}

type Config struct {
	TelegramApiToken string `env:"TELEGRAM_API_TOKEN,required"`
	ChatId           string `env:"CHAT_ID,required"`
	ApiKey           string `env:"API_KEY,required"`
	Port             string `env:"PORT,required"`
	AllowOrigins     string `env:"ALLOW_ORIGINS,required"`
}

func getConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg := Config{
		TelegramApiToken: os.Getenv("TELEGRAM_API_TOKEN"),
		ChatId:           os.Getenv("CHAT_ID"),
		ApiKey:           os.Getenv("API_KEY"),
		Port:             os.Getenv("PORT"),
		AllowOrigins:     os.Getenv("ALLOW_ORIGINS"),
	}

	if len(cfg.TelegramApiToken) == 0 {
		log.Fatal("TELEGRAM_API_TOKEN not provided")
	}

	if len(cfg.ChatId) == 0 {
		log.Fatal("CHAT_ID not provided")
	}

	if len(cfg.ApiKey) == 0 {
		log.Fatal("API_KEY not provided")
	}

	if len(cfg.Port) == 0 {
		log.Fatal("PORT not provided")
	}

	return cfg

}

func sendTelegramNotification(message string, config Config) error {

	client := &http.Client{}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramApiToken)
	data := url.Values{
		"chat_id": {config.ChatId},
		"text":    {message},
	}
	response, err := client.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	response.Body.Close()

	return nil
}

func main() {

	config := getConfig()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "POST",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/push-notification-to-telegram", func(c *fiber.Ctx) error {
		var err error

		notification := new(PushNotificationRequest)
		if err = c.BodyParser(notification); err != nil {
			return err
		}

		if notification.ApiKey != config.ApiKey {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, invalid apikey"})
		}

		err = sendTelegramNotification(notification.Message, config)
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, cannot reach telegram"})
		}
		return c.JSON(map[string]string{"status": "success, message sent"})

	})

	appErr := app.Listen(":" + config.Port)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}
}
