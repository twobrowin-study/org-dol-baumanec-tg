package app

import (
	"fmt"
	"strconv"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) CounselorFindPhoneSquad(chatId string) {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var keyrow     []tgbotapi.InlineKeyboardButton
	
	rows, err := app.DbConnection.Query(context.Background(), "select squad from counselor_squad_unique order by squad")
	errExx("CounselorFindPhoneSquad counselor_squad_unique", err)
	defer rows.Close()
	
	keyrowIdx := 0
	for rows.Next() {
		var squad int
		rows.Scan(&squad)
		squadStr := strconv.Itoa(squad)
		
		button := tgbotapi.NewInlineKeyboardButtonData(squadStr, squadStr)
		keyrow = append(keyrow, button)
		if (keyrowIdx+1) % 4 == 0 {
			keyboard = append(keyboard, keyrow[keyrowIdx-3:keyrowIdx+1])
		}
		keyrowIdx += 1
	}
	keyboard = append(keyboard, keyrow[(keyrowIdx/4)*4:])

	keyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
	app.sendWithKeyboard(chatId, findPhoneCounselorCallText, keyboardMarkup)
}

func (app *App) CounselorFindPhoneFinish(chatId string, user *User, squad string) {
	messagePrefix := fmt.Sprintf("Вожатые *%s отряда*", squad)
	app.sendFindPhone(chatId, messagePrefix, "post='counselor' and squad=$1", squad)
}