package app

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) StartFindPhone(chatId string, user *User) {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	rows, err := app.DbConnection.Query(context.Background(), "select post, name from post_type_name_order_active_unique order by \"order\"")
	errExx("StartFindPhone post_type_name_order_active_unique", err)
	defer rows.Close()
	for rows.Next() {
		var post, name string
		rows.Scan(&post, &name)
		button := tgbotapi.NewInlineKeyboardButtonData(name, post)
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(button))
	}
	keyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
	app.sendWithKeyboard(chatId, findPhonePostCallText, keyboardMarkup)
}

func (app *App) ProceedFindPhone(chatId string, user *User, post string) {
	if post == "counselor" {
		app.CounselorFindPhoneSquad(chatId)
		return
	}

	if post == "circle" {
		app.CircleFindPhoneTitle(chatId)
		return
	}

	app.sendFindPhone(chatId, "", "post=$1", post)
}