package app

import (
	"fmt"
	"strconv"
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
	messagePrefix := fmt.Sprintf("Кружковод *%s*:", title)
	app.sendFindPhone(chatId, messagePrefix, "post='circle' and title=$1", title)
}

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
	messagePrefix := fmt.Sprintf("Вожатые *%s* отряда:", squad)
	app.sendFindPhone(chatId, messagePrefix, "post='counselor' and squad=$1", squad)
}

func (app *App) sendFindPhone(chatId, messagePrefix, where, param string) {
	messageText := ""
	if messagePrefix != "" {
		messageText += messagePrefix + "\n"
	}

	var post_name string
	queryStr := fmt.Sprintf("select name, phone, post_name from find_phone_active where %s", where)
	rows, err := app.DbConnection.Query(context.Background(), queryStr, param)
	errExx("sendFindPhone find_phone_active", err)
	defer rows.Close()

	rowsNext := rows.Next()
	for rowsNext {
		var name, phone string
		rows.Scan(&name, &phone, &post_name)
		messageText += fmt.Sprintf("  %s - %s", name, phone)
		rowsNext = rows.Next()
		if rowsNext {
			messageText += "\n"
		}
	}

	if messagePrefix == "" {
		messageText = fmt.Sprintf("*%s*:\n%s", post_name, messageText)
	}

	app.send(chatId, messageText)
}