package app

import (
	"fmt"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) CircleFindPhoneTitle(chatId string) {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	rows, err := app.DbConnection.Query(context.Background(), "select title from circle_title_unique order by title")
	errExx("CircleFindPhoneTitle circle_title_unique", err)
	defer rows.Close()
	for rows.Next() {
		var title string
		rows.Scan(&title)
		button := tgbotapi.NewInlineKeyboardButtonData(title, title)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}
	keyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
	app.sendWithKeyboard(chatId, findPhoneCircleCallText, keyboardMarkup)
}

func (app *App) CircleFindPhoneFinish(chatId string, user *User, title string) {
	messagePrefix := "Кружковод"
	if title == "Театр" {
		messagePrefix = "Кружководы"
	} 
	messagePrefix = fmt.Sprintf("%s *%s*", messagePrefix, title)
	app.sendFindPhone(chatId, messagePrefix, "post='circle' and title=$1", title)
}