package logger

import (
	"log"
	"os"
)

// Определяем глобальные логгеры
var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

// InitLogger инициализирует логгер, записывая логи в файл и в консоль.
func InitLogger(logFilePath string) error {
	var err error

	// Открываем или создаем файл для логов
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Создаем логгеры
	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Дублируем вывод логов в консоль
	multiWriter := os.Stdout
	log.SetOutput(multiWriter)

	return nil
}

// LogInfo записывает информационные сообщения.
func LogInfo(message string) {
	infoLogger.Println(message)
}

// LogError записывает сообщения об ошибках.
func LogError(message string) {
	errorLogger.Println(message)
}

// CloseLogger закрывает файл логов (нужно вызывать при завершении программы).
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
