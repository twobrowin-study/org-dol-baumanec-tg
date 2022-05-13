package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"org.dol.baumanec/tgbot/app"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to database")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal("Unable to authorize on telegram bot account:", err)
	}
	log.Println("Authorized on telegram bot account", bot.Self.UserName)

	tgAdmin := os.Getenv("TELEGRAM_ADMIN")

	appObj, err := app.NewApp(conn, bot, tgAdmin)
	if err != nil {
		log.Fatal("Unable to start an application:", err)
	}
	log.Println("Started application")
	appObj.Execute()
}