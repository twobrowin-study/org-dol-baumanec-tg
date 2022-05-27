package app

import (
	"log"
)

func errExx(mod string, err error) bool {
	if err != nil {
		log.Println("Error in", mod, "::", err)
		return true
	}
	return false
}