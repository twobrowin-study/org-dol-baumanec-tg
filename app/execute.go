package app

import (
	"strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) Execute() {
	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 60

	updates := app.TelegramApi.GetUpdatesChan(updater)

	for update := range updates {
		if update.Message != nil {
			chatId := strconv.FormatInt(update.Message.Chat.ID, 10)
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					if user := app.ActiveUser(chatId); user != nil {
						app.SendFullHelp(chatId, user)
					} else if app.FaultUser(chatId) == true {
						app.SendErrorHelp(chatId)
					} else {
						app.StartUserCreation(chatId)
					}
				case "help":
					if user := app.ActiveUser(chatId); user != nil {
						app.SendFullHelp(chatId, user)
					} else if app.FaultUser(chatId) == true {
						app.SendErrorHelp(chatId)
					} else {
						app.SendTinyHelp(chatId)
					}
				case "doctorchecklist":
					if user := app.ActiveUser(chatId); user != nil {
						app.SendDoctorChecklist(chatId, user)
					} else if app.FaultUser(chatId) == true {
						app.SendErrorHelp(chatId)
					} else {
						app.SendTinyHelp(chatId)
					}
				case "doctordone":
					if user := app.ActiveUser(chatId); user != nil {
						app.StartDoctorDone(chatId, user)
					} else if app.FaultUser(chatId) == true {
						app.SendErrorHelp(chatId)
					} else {
						app.SendTinyHelp(chatId)
					}
				default:
					app.SendTinyHelp(chatId)
				}
			} else {
				if user := app.ActiveUser(chatId); user != nil {
					app.SendFullHelp(chatId, user)
				} else if app.WaitingUserName(chatId) == true {
					app.UserCreationName(chatId, update.Message.Text)
				} else if app.WaitingUserSex(chatId) == true {
					app.UserCreationSex(chatId, update.Message.Text)
				} else if app.FaultUser(chatId) == true {
					app.SendErrorHelp(chatId)
				} else {
					app.SendTinyHelp(chatId)
				}
			}
		} else if update.CallbackQuery != nil {
			chatId := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
			if user := app.ActiveUser(chatId); user != nil {
				app.SetDoctorDone(chatId, user, update.CallbackQuery.Data)
			} else if app.FaultUser(chatId) == true {
				app.SendErrorHelp(chatId)
			} else {
				app.SendTinyHelp(chatId)
			}
		}
	}
}