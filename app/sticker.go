package app

import (
	"context"
)

func (app *App) GetSticker(name string) string {
	var sticker string
	err := app.DbConnection.QueryRow(context.Background(), "select file from sticker_active where name=$1", name).Scan(&sticker)
	errExx("GetSticker", err)
	return sticker
}