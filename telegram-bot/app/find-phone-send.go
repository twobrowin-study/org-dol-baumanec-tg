package app

import (
	"fmt"
	"context"
)

func (app *App) sendFindPhone(chatId, messagePrefix, where, param string) {
	queryStr := fmt.Sprintf("select name, phone, post_name, address, is_address_equal from find_phone_active_address where %s", where)
	rows, err := app.DbConnection.Query(context.Background(), queryStr, param)
	errExx("sendFindPhone find_phone_active_address", err)
	defer rows.Close()
	
	listText := ""
	var post_name, address string
	var isAddressEqual bool

	rowsNext := rows.Next()
	for rowsNext {
		var name, phone string
		err := rows.Scan(&name, &phone, &post_name, &address, &isAddressEqual)
		errExx("sendFindPhone find_phone_active_address scan", err)

		if isAddressEqual == true || address == "" {
			listText += fmt.Sprintf("  %s %s", name, phone)
		} else {
			listText += fmt.Sprintf("  %s %s\n      %s", name, phone, address)
		}
		
		rowsNext = rows.Next()
		if rowsNext {
			listText += "\n"
		}
	}

	messageText := ""
	if messagePrefix == "" && isAddressEqual == true {
		messageText = fmt.Sprintf("*%s*: %s\n%s", post_name, address, listText)
	} else if messagePrefix == "" && isAddressEqual == false {
		messageText = fmt.Sprintf("*%s*:\n%s", post_name, listText)
	} else if messagePrefix != "" && isAddressEqual == true {
		messageText = fmt.Sprintf("%s: %s\n%s", messagePrefix, address, listText)
	} else if messagePrefix != "" && isAddressEqual == false {
		messageText = fmt.Sprintf("%s:\n%s", messagePrefix, listText)
	}

	app.send(chatId, messageText)
}