package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"

	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"
)

type BotToken struct {
	TelegramApiToken string `json:"telegramApiToken"`
}

type PushNotificationContent struct {
	MessageType string `json:"messageType" db:"MessageType"`
	Email       string `json:"email" db:"Email"`
	Message     string `json:"message" db:"Message"`
	ApiKey      string `json:"apikey"`
}

type MsgTableRow struct {
	MessageType string `json:"messageType" db:"MessageType"`
	Email       string `json:"email" db:"Email"`
	Message     string `json:"message" db:"Message"`
	Timestamp   string `json:"timestamp" db:"Timestamp"`
}

var schema = `
CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    MessageType TEXT,
    Email TEXT,
    Message TEXT,
	Timestamp TEXT
);`

type Config struct {
	TelegramApiToken   string `env:"TELEGRAM_API_TOKEN"`
	ChatId             string `env:"CHAT_ID"`
	ApiKey             string `env:"API_KEY"`
	Port               string `env:"PORT"`
	MaxRequestsPerHour int    `env:"MAX_REQUESTS_PER_HOUR"`
}

func getConfig() Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	maxRequestsPerHour := 1
	if os.Getenv("MAX_REQUESTS_PER_HOUR") != "" {
		maxRequestsPerHour, err = strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_HOUR"))
		if err != nil {
			log.Fatal("MAX_REQUESTS_PER_HOUR must be a number")
		}
	}

	cfg := Config{
		TelegramApiToken:   os.Getenv("TELEGRAM_API_TOKEN"),
		ChatId:             os.Getenv("CHAT_ID"),
		ApiKey:             os.Getenv("API_KEY"),
		Port:               os.Getenv("PORT"),
		MaxRequestsPerHour: maxRequestsPerHour,
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

	dbPath := "./database.sqlite"

	db := sqlx.MustConnect("sqlite3", dbPath)
	db.MustExec(schema)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "POST, DELETE",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               config.MaxRequestsPerHour,
		Expiration:        1 * time.Hour,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Delete("/clear-db", func(c *fiber.Ctx) error {
		var err error

		msg := new(BotToken)
		if err = c.BodyParser(msg); err != nil {
			return err
		}

		if msg.TelegramApiToken != config.TelegramApiToken {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, invalid apikey"})
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			c.Status(fiber.StatusNotFound)
			return c.JSON(map[string]string{"status": "failed, file not found"})
		}

		_, err = db.Exec("DELETE FROM messages")
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, unable to delete messages"})
		}

		return c.JSON(map[string]string{"status": "success, all messages deleted"})

	})

	app.Post("/download-db", func(c *fiber.Ctx) error {
		var err error

		msg := new(BotToken)
		if err = c.BodyParser(msg); err != nil {
			return err
		}

		if msg.TelegramApiToken != config.TelegramApiToken {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, invalid apikey"})
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			c.Status(404)
			return c.JSON(map[string]string{"status": "failed, file not found"})
		}

		return c.Download(dbPath, filepath.Base(dbPath))

	})

	app.Post("/push-notification-to-telegram", func(c *fiber.Ctx) error {
		var err error

		msg := new(PushNotificationContent)
		if err = c.BodyParser(msg); err != nil {
			return err
		}

		if msg.ApiKey != config.ApiKey {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, invalid apikey"})
		}

		if len(msg.Email) == 0 || len(msg.Message) == 0 || len(msg.MessageType) == 0 {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, body"})
		}

		const tmpl = `
MESSAGE_TYPE: 
{{.MessageType}}
EMAIL: 
{{.Email}}
MESSAGE:
{{.Message}}
`

		t := template.Must(template.New("contact").Parse(tmpl))
		buf := new(bytes.Buffer)
		err = t.Execute(buf, msg)
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, invalid text"})
		}

		err = sendTelegramNotification(buf.String(), config)
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, cannot reach telegram"})
		}

		tx := db.MustBegin()
		_, dberr := tx.NamedExec("INSERT INTO messages (MessageType, Email, Message, Timestamp) VALUES (:MessageType, :Email, :Message, :Timestamp)", &MsgTableRow{
			MessageType: msg.MessageType,
			Email:       msg.Email,
			Message:     msg.Message,
			Timestamp:   time.Now().Format(time.RFC3339),
		})
		if dberr != nil {
			tx.Rollback()
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, cannot save message"})
		}
		dberr = tx.Commit()
		if dberr != nil {
			c.Status(500)
			return c.JSON(map[string]string{"status": "failed, cannot save message"})
		}

		return c.JSON(map[string]string{"status": "success, message sent"})

	})

	appErr := app.Listen(":" + config.Port)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}
}
