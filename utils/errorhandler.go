package utils

import "log"

func HandlePanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}
