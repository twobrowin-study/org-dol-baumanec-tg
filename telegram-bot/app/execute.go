package app

import (
	"strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	doneTherapistCallText = "Терапевт"
	doneTherapistLastCallText = "Остался только Терапевт"
	doneAnalysisCallText = "Анализы"
	doneDoctorCallText = "Врачи"
	findPhonePostCallText = "Выбери кого ты ищешь"
	findPhoneCircleCallText = "Выбери кружок"
	findPhoneCounselorCallText = "Выбери номер отряда"
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
					app.testUserInacticeExecute(chatId, app.StartUserCreation)
				case "help":
					app.testUserActiceExecute(chatId, app.SendFullHelp)
				case "doctorchecklist":
					app.testUserActiceExecute(chatId, app.SendDoctorChecklist)
				case "doctordone":
					app.testUserActiceExecute(chatId, app.StartDoctorDone)
				case "find2phone":
					app.testUserActiceExecute(chatId, app.StartFindPhone)
				default:
					app.SendTinyHelp(chatId)
				}
			} else {
				app.testUserActiveAndCreationInput(chatId, update.Message.Text)
			}
		} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			chatId := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
			switch update.CallbackQuery.Message.Text {
			case doneTherapistCallText:
				fallthrough
			case doneTherapistLastCallText:
				fallthrough
			case doneAnalysisCallText:
				fallthrough
			case doneDoctorCallText:
				app.testUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.SetDoctorDone)
			case findPhonePostCallText:
				app.testUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.ProceedFindPhone)
			case findPhoneCircleCallText:
				app.testUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.CircleFindPhoneFinish)
			case findPhoneCounselorCallText:
				app.testUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.CounselorFindPhoneFinish)
			}
		}
	}
}

func (app *App) testUserActiveAndCreationInput(chatId, messageText string) {
	if user := app.ActiveUser(chatId); user != nil {
		app.SendFullHelp(chatId, user)
	} else if app.WaitingUserName(chatId) == true {
		app.UserCreationName(chatId, messageText)
	} else if app.WaitingUserSex(chatId) == true {
		app.UserCreationSex(chatId, messageText)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		app.SendTinyHelp(chatId)
	}
}

func (app *App) testUserInacticeExecute(chatId string, execInactive func(string)) {
	if user := app.ActiveUser(chatId); user != nil {
		app.SendFullHelp(chatId, user)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		execInactive(chatId)
	}
}

func (app *App) testUserActiceExecute(chatId string, execActive func(string, *User)) {
	if user := app.ActiveUser(chatId); user != nil {
		execActive(chatId, user)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		app.SendTinyHelp(chatId)
	}
}

func (app *App) testUserActiceExecuteCallback(chatId string, param string, execActive func(string, *User, string)) {
	if user := app.ActiveUser(chatId); user != nil {
		execActive(chatId, user, param)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		app.SendTinyHelp(chatId)
	}
}