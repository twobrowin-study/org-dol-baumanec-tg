package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
)

const (
	maleText = "муж."
	femaleText = "жен."

	doneTherapistCallText = "Терапевт"
	doneTherapistLastCallText = "Остался только Терапевт"
	doneAnalysisCallText = "Анализы"
	doneDoctorCallText = "Врачи"
	findPhonePostCallText = "Выбери кого ты ищешь"
	findPhoneCircleCallText = "Выбери кружок"
	findPhoneCounselorCallText = "Выбери номер отряда"
)

type App struct {
	DbConnection	*pgx.Conn
	TelegramApi		*tgbotapi.BotAPI
	TelegramAdmin	string
}

func NewApp(DbConnection *pgx.Conn, TelegramApi *tgbotapi.BotAPI, TelegramAdmin string) (*App, error) {
	return &App{DbConnection, TelegramApi, TelegramAdmin}, nil
}