// touched for cleanup
package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/nigdanil/Quality_Control_System_Generation/DB"
	"github.com/nigdanil/Quality_Control_System_Generation/Generation"
)

func main() {
	// Открытие файла для логов
	logFile, err := os.OpenFile("generation_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()

	// Настраиваем вывод логов в файл и консоль
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("🚀 Запуск Quality Control System")

	// Подключение к базе данных
	db, err := DB.OpenDatabase("Quality_Control.SQLite")
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	log.Println("✅ Подключение к БД успешно")

	// Получение последней записи
	logData, err := DB.GetLastGenerationLog(db)
	if err != nil {
		log.Fatalf("❌ Ошибка чтения из БД: %v", err)
	}

	log.Println("ℹ️  Загружена последняя запись из generation_logs")

	const iterations = 3
	const interval = 5 * time.Second

	for i := 1; i <= iterations; i++ {
		startTime := time.Now()
		log.Printf("▶️  Итерация %d — старт: %s\n", i, startTime.Format("15:04:05"))

		fileName, duration, err := Generation.GenerateImage(logData)
		endTime := time.Now()

		if err != nil {
			log.Printf("❌ Ошибка генерации (%s): %v\n", endTime.Format("15:04:05"), err)
		} else {
			log.Printf("✅ Сохранено: %s (%0.2f сек.) — завершено: %s\n", fileName, duration, endTime.Format("15:04:05"))
		}

		if i < iterations {
			time.Sleep(interval)
		}
	}

	log.Println("🏁 Завершено")
}
