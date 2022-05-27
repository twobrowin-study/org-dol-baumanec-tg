package app

import (
	"context"
)

func (app *App) SendDoctorChecklist(chatId string, user *User) {
	mesg := user.FirstName + ", твой чеклист врачей\n\n"

	mesg += "*Анализы*:\n"
	rowsAnalysis, err := app.DbConnection.Query(context.Background(), "select title, ready from user_analysis_to_go where chat_id=$1 order by id", chatId)
	errExx("SendDoctorChecklist user_analysis_to_go", err)
	defer rowsAnalysis.Close()
	for rowsAnalysis.Next() {
		var title string
		var ready bool
		rowsAnalysis.Scan(&title, &ready)
		mesg += "    " + compileWithReady(title, ready) + "\n"
	}

	mesg += "\n*Врачи*:\n"
	rowsDoctor, err := app.DbConnection.Query(context.Background(), "select title, ready from user_doctor_to_go where chat_id=$1 order by id", chatId)
	errExx("SendDoctorChecklist user_doctor_to_go", err)
	defer rowsDoctor.Close()
	for rowsDoctor.Next() {
		var title string
		var ready bool
		rowsDoctor.Scan(&title, &ready)
		mesg += "    " + compileWithReady(title, ready) + "\n"
	}

	var title string
	var ready bool
	err = app.DbConnection.QueryRow(context.Background(), "select title, ready from user_therapist_to_go where chat_id=$1", chatId).Scan(&title, &ready)
	errExx("SendDoctorChecklist user_therapist_to_go", err)
	mesg += "\nИ *" + compileWithReady(title, ready) + "* в конце!"
	mesg += "\n\nВведи /doctordone для того чтобы отметить прохождение врача в чеклисте"

	app.send(chatId, mesg)
}

func compileWithReady(title string, ready bool) string {
	if ready {
		return "✅ " + title
	}
	return "🔹 " + title
}