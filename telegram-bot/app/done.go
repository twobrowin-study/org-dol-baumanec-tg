package app

import (
	"strconv"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) StartDoctorDone(chatId string, user *User) {
	app.send(chatId, "Выбери среди доступных вариантов:")

	var id int64
	var title string
	var show bool = false
	err := app.DbConnection.QueryRow(context.Background(), "select id, title, show from user_therapist_to_go where chat_id=$1 and ready = false", chatId).Scan(&id, &title, &show)
	errExx("StartDoctorDone user_therapist_to_go", err)
	if show {
		idStr := strconv.FormatInt(id, 10)
		therapistKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(title, idStr),
			),
		)
		app.sendWithKeyboard(chatId, "Терапевт", therapistKeyboard)
		return
	}

	var analysisKeyboard [][]tgbotapi.InlineKeyboardButton
	rowsAnalysis, err := app.DbConnection.Query(context.Background(), "select id, title from user_analysis_to_go where chat_id=$1 and ready = false order by id", chatId)
	errExx("StartDoctorDone user_analysis_to_go", err)
	defer rowsAnalysis.Close()
	for rowsAnalysis.Next() {
		var id int64
		var title string
		rowsAnalysis.Scan(&id, &title)
		idStr := strconv.FormatInt(id, 10)
		
		button := tgbotapi.NewInlineKeyboardButtonData(title, idStr)
		analysisKeyboard = append(analysisKeyboard, tgbotapi.NewInlineKeyboardRow(button))
	}
	analysisKeyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: analysisKeyboard,
	}
	app.sendWithKeyboard(chatId, "Анализы", analysisKeyboardMarkup)

	var doctorKeyboard [][]tgbotapi.InlineKeyboardButton
	rowsDoctor, err := app.DbConnection.Query(context.Background(), "select id, title from user_doctor_to_go where chat_id=$1 and ready = false order by id", chatId)
	errExx("StartDoctorDone user_doctor_to_go", err)
	defer rowsDoctor.Close()
	for rowsDoctor.Next() {
		var id int64
		var title string
		rowsDoctor.Scan(&id, &title)
		idStr := strconv.FormatInt(id, 10)
		
		button := tgbotapi.NewInlineKeyboardButtonData(title, idStr)
		doctorKeyboard = append(doctorKeyboard, tgbotapi.NewInlineKeyboardRow(button))
	}
	doctorKeyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: doctorKeyboard,
	}
	app.sendWithKeyboard(chatId, "Врачи", doctorKeyboardMarkup)
}

func (app *App) SetDoctorDone(chatId string, user *User, doctorId string) {
	rows, err := app.DbConnection.Query(context.Background(), "insert into user_doctor_done (user_id, doctor_id) values ($1, $2)", user.Id, doctorId)
	rows.Close()
	errExx("SetDoctorDone user_doctor_done", err)

	var doctorTitle string
	var doctorType string
	var doctorFun bool
	err = app.DbConnection.QueryRow(context.Background(), "select title, doctor_type, is_fun from doctor_active_fun where id=$1 and user_id=$2", doctorId, user.Id).Scan(&doctorTitle, &doctorType, &doctorFun)
	errExx("SetDoctorDone doctor", err)
	app.send(chatId, "Отмечено ✅ " + doctorTitle)

	if doctorFun == true {
		app.SendSticker(chatId, app.GetSticker("omg"))
	}

	if doctorType == "therapist" {
		app.send(chatId, "Супер! Всё готово, ты можешь быть свободен!")
		app.SendSticker(chatId, app.GetSticker("super"))
	}

	var id int64
	var title string
	var show bool = false
	err = app.DbConnection.QueryRow(context.Background(), "select id, title, show from user_therapist_to_go where chat_id=$1 and ready = false", chatId).Scan(&id, &title, &show)
	errExx("SetDoctorDone user_therapist_to_go", err)
	if show {
		idStr := strconv.FormatInt(id, 10)
		therapistKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(title, idStr),
			),
		)
		app.sendWithKeyboard(chatId, "Остался только *" + title + "*", therapistKeyboard)
		return
	}
}