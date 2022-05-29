package app

import (
	"strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sexKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(maleText),
		tgbotapi.NewKeyboardButton(femaleText),
	),
)

func (app *App) avaliableFuncsDesc() string {
	mesg := ""
	if app.isFuncActive("doctordone") {
		mesg += "\n\n" +
				"Введи /doctordone для того чтобы отметить прохождение врача в чеклисте"
	}
	if app.isFuncActive("doctorchecklist") {
		mesg += "\n\n" +
				"Введи /doctorchecklist для просмотра твоего чеклиста врачей"
	}
	if app.isFuncActive("find2phone") {
		mesg += "\n\n" +
				"Eсли ты хочешь найти телефон кого-то из нашей команды - введи /find2phone"
	}
	return mesg
}

func (app *App) SendFullHelp(chatId string, user *User) {
	mesg := "Привет, " + user.FirstName + "!" +
			"\n\n" +
			"Ты уже зарегистрирован!"
	mesg += app.avaliableFuncsDesc()
	app.send(chatId, mesg)
}

func (app *App) SendErrorHelp(chatId string) {
	mesg := "Хмм... произошла ошибка пользователя" +
			"\n\n" +
			"Обратись, пожалуйста, к " + app.TelegramAdmin
	app.send(chatId, mesg)
}

func (app *App) SendTinyHelp(chatId string) {
	mesg := "Это ещё не релизованная функция или произошла ошибка!" +
			"\n\n" +
			"Баг был найден и уже обрабатывается!"
	app.send(chatId, mesg)
}

func (app *App) SendNewUserWaitForName(chatId string) {
	mesg := "Привет! Добро пожаловать!" +
			"\n\n" +
			"Введи, пожалуйста, свои *Фамилию* и *Имя* через пробел"
	app.send(chatId, mesg)
}

func (app *App) SendNotName(chatId string) {
	mesg := "Это не похоже на *Фамилию* и *Имя* через пробел" +
			"\n\n" +
			"Попробуй ещё раз"
	app.send(chatId, mesg)
}

func (app *App) SendUserNameWaitForSex(chatId, firstName string) {
	mesg := "Привет, " + firstName + "!" +
			"\n\n" +
			"Теперь выбери свой пол"
	app.sendWithKeyboard(chatId, mesg, sexKeyboard)
}

func (app *App) SendUserDone(chatId string) {
	mesg := "Ура! Регистрация завершена!"
	app.send(chatId, mesg)
}

func (app *App) SendNextCommands(chatId string) {
	mesg := "Теперь сделай вот что:"
	mesg += app.avaliableFuncsDesc()
	app.send(chatId, mesg)
}

func (app *App) send(chatId, textMarkdown string) {
	chatIdInt, _ := strconv.ParseInt(chatId, 10, 64)
	msg := tgbotapi.NewMessage(chatIdInt, textMarkdown)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	app.TelegramApi.Send(msg)
}

func (app *App) sendWithKeyboard(chatId, textMarkdown string, keyboard interface{}) {
	chatIdInt, _ := strconv.ParseInt(chatId, 10, 64)
	msg := tgbotapi.NewMessage(chatIdInt, textMarkdown)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = keyboard
	app.TelegramApi.Send(msg)
}

func (app *App) SendSticker(chatId, sticker string) {
	chatIdInt, _ := strconv.ParseInt(chatId, 10, 64)
	msg := tgbotapi.NewSticker(chatIdInt, tgbotapi.FileID(sticker))
	app.TelegramApi.Send(msg)
}