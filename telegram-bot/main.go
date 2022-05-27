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

	var tgToken, tgAdmin string
	err = conn.QueryRow(context.Background(), "select token, admin from telegram_active_single").Scan(&tgToken, &tgAdmin)
	if err != nil {
		log.Fatal("Unable to get telegram token and admin data:", err)
	}
	log.Println("Got telegram token and admin data")

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal("Unable to authorize on telegram bot account:", err)
	}
	log.Println("Authorized on telegram bot account", bot.Self.UserName)

	appObj, err := app.NewApp(conn, bot, tgAdmin)
	if err != nil {
		log.Fatal("Unable to start an application:", err)
	}
	log.Println("Started application")
	appObj.Execute()
}