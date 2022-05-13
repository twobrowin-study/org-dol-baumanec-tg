package app

import (
	"context"
	"strings"
)

type User struct {
	Id			uint
	ChatId		string
	FirstName	string
	LastName	string
	Sex			string
}

func (app *App) ActiveUser(chatId string) *User {
	var user User
	err := app.DbConnection.QueryRow(context.Background(), "select * from user_active where chat_id=$1", chatId).Scan(&user.Id, &user.ChatId, &user.FirstName, &user.LastName, &user.Sex)
	if errExx("ActiveUser", err) {
		return nil
	}
	return &user
}

func (app *App) FaultUser(chatId string) bool {
	var ans bool
	err := app.DbConnection.QueryRow(context.Background(), "select true from user_with_errors where chat_id=$1", chatId).Scan(&ans)
	if errExx("FaultUser", err) {
		return false
	}
	return ans
}

func (app *App) StartUserCreation(chatId string) {
	rows, err := app.DbConnection.Query(context.Background(), "insert into \"user\" (chat_id) values ($1)", chatId)
	rows.Close()
	errExx("StartUserCreation", err)
	app.SendNewUserWaitForName(chatId)
}

func (app *App) WaitingUserName(chatId string) bool {
	var ans bool
	err := app.DbConnection.QueryRow(context.Background(), "select true from user_waiting_for_name_and_sex where chat_id=$1", chatId).Scan(&ans)
	if errExx("WaitingUserName", err) {
		return false
	}
	return ans
}

func (app *App) UserCreationName(chatId, message string) {
	nameFields := strings.Fields(message)
	if len(nameFields) != 2 {
		app.SendNotName(chatId)
		return
	}
	rows, err := app.DbConnection.Query(context.Background(), "update \"user\" set first_name=$3, last_name=$2 where chat_id=$1", chatId, nameFields[0], nameFields[1])
	rows.Close()
	errExx("UserCreationName", err)
	app.SendUserNameWaitForSex(chatId, nameFields[1])
}

func (app *App) WaitingUserSex(chatId string) bool {
	var ans bool
	err := app.DbConnection.QueryRow(context.Background(), "select true from user_waiting_for_sex_but_not_name where chat_id=$1", chatId).Scan(&ans)
	if errExx("WaitingUserSex", err) {
		return false
	}
	return ans
}

func (app *App) UserCreationSex(chatId, message string) {
	sex := "female"
	if message == maleText {
		sex = "male"
	}
	rows, err := app.DbConnection.Query(context.Background(), "update \"user\" set sex=$2 where chat_id=$1", chatId, sex)
	rows.Close()
	errExx("UserCreationSex", err)
	app.SendUserDone(chatId)
	app.SendSticker(chatId, app.GetSticker("hello"))
	app.SendNextCommands(chatId)
}