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
			if app.CheckIfUserCanNotGetResponce(chatId) {
				continue
			}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					app.CheckIfUserInacticeExecute(chatId, app.StartUserCreation)
				case "help":
					app.CheckIfUserActiceExecute(chatId, app.SendFullHelp)
				case "doctorchecklist":
					if app.CheckIfFuncActive("doctorchecklist") {
						app.CheckIfUserActiceExecute(chatId, app.SendDoctorChecklist)
					}
				case "doctordone":
					if app.CheckIfFuncActive("doctordone") {
						app.CheckIfUserActiceExecute(chatId, app.StartDoctorDone)
					}
				case "find2phone":
					if app.CheckIfFuncActive("find2phone") {
						app.CheckIfUserActiceExecute(chatId, app.StartFindPhone)
					}
				default:
					app.SendTinyHelp(chatId)
				}
			} else {
				app.CheckIfUserActiveAndCreationInput(chatId, update.Message.Text)
			}
		} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {

			chatId := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
			if app.CheckIfUserCanNotGetResponce(chatId) {
				continue
			}

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := app.TelegramApi.Request(callback); err != nil {
				panic(err)
			}

			switch update.CallbackQuery.Message.Text {
			case doneTherapistCallText:
				fallthrough
			case doneTherapistLastCallText:
				fallthrough
			case doneAnalysisCallText:
				fallthrough
			case doneDoctorCallText:
				app.CheckIfUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.SetDoctorDone)
			case findPhonePostCallText:
				app.CheckIfUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.ProceedFindPhone)
			case findPhoneCircleCallText:
				app.CheckIfUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.CircleFindPhoneFinish)
			case findPhoneCounselorCallText:
				app.CheckIfUserActiceExecuteCallback(chatId, update.CallbackQuery.Data, app.CounselorFindPhoneFinish)
			}
		}
	}
}