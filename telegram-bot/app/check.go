package app

import (
	"context"
)

func (app *App) CheckIfFuncActive(funcStr string) bool {
	ans := false
	app.DbConnection.QueryRow(context.Background(), "select true from service_active where func=$1", funcStr).Scan(&ans)
	return ans
}

func (app *App) CheckIfNotAdminChat(chatId string) bool {
	ans := true
	app.DbConnection.QueryRow(context.Background(), "select false from user_admin where chat_id=$1", chatId).Scan(&ans)
	return ans
}

func (app *App) CheckIfUserCanNotGetResponce(chatId string) bool {
	if app.CheckIfFuncActive("service") && app.CheckIfNotAdminChat(chatId) {
		app.send(chatId, "Ведутся технические работы\nМы скоро увидимся вновь!")
		return true
	}
	return false
}

func (app *App) CheckIfUserActiveAndCreationInput(chatId, messageText string) {
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

func (app *App) CheckIfUserInacticeExecute(chatId string, execInactive func(string)) {
	if user := app.ActiveUser(chatId); user != nil {
		app.SendFullHelp(chatId, user)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		execInactive(chatId)
	}
}

func (app *App) CheckIfUserActiceExecute(chatId string, execActive func(string, *User)) {
	if user := app.ActiveUser(chatId); user != nil {
		execActive(chatId, user)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		app.SendTinyHelp(chatId)
	}
}

func (app *App) CheckIfUserActiceExecuteCallback(chatId string, param string, execActive func(string, *User, string)) {
	if user := app.ActiveUser(chatId); user != nil {
		execActive(chatId, user, param)
	} else if app.FaultUser(chatId) == true {
		app.SendErrorHelp(chatId)
	} else {
		app.SendTinyHelp(chatId)
	}
}