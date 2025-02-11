package logging

import (
	"log"
)

func LogInfo(message string) {
	log.Println("INFO:", message)
}

func LogError(message string) {
	log.Println("ERROR:", message)
}
