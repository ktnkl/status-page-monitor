package check

import (
	"fmt"
	"log"
	"status-page-monitor/internal/database"
	"status-page-monitor/internal/database/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

type WorkerResult struct {
	Server   models.Server
	Response *ServerStatus
}

func producer(in chan<- models.Server, tx *gorm.DB) {
	var servers []models.Server
	threshold := time.Now().Add(30 * time.Second).UTC()
	log.Println("Producer threshold", threshold)
	var total int64
	sql := tx.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(`
    SELECT * FROM servers
    WHERE nextcheckat::timestamp < ?
    FOR UPDATE SKIP LOCKED
  `, threshold)
	})
	fmt.Println("Row select * query", sql)
	err := tx.Raw(`
    SELECT * FROM servers
    WHERE nextcheckat::timestamp < ?
    FOR UPDATE SKIP LOCKED
  `, threshold).Scan(&servers).Count(&total).Error

	if err != nil {
		tx.Rollback()
	}

	for _, server := range servers {
		log.Println("Server, ", server)
		in <- server
	}
	close(in)
}

func consumer(out <-chan WorkerResult, tx *gorm.DB, done chan<- interface{}) {
	defer close(done)
	for ss := range out {
		log.Println("consumer принял ", ss.Response)
		status := "inactive"
		if ss.Response.IsActive {
			status = "active"
		}
		query := tx.Exec(`
			UPDATE servers
			SET status = ?,
			checkedat = ?,
			nextcheckat = ?
			WHERE url = ?
		`,
			status,
			time.Now().UTC().Format(time.RFC3339),
			time.Now().Add(time.Duration(ss.Server.Interval)*time.Second).UTC().Format(time.RFC3339),
			ss.Response.URL)
		if query.Error != nil {
			log.Printf("Ошибка при обновлении сервера воркером: %v", query.Error)
		}
		log.Printf("Результат обновления сервера %s, rowsAffected: %d :", ss.Server.Url, query.RowsAffected)
	}
}

func worker(sc *ServerChecker, in <-chan models.Server, out chan<- WorkerResult) {
	for server := range in {
		status, err := sc.CheckServer(server.Url)
		if err != nil {
			log.Printf("При проверке url: %s произошла ошибка %v", server.Url, err)
		}
		log.Println("worker: ", status)
		out <- WorkerResult{server, status}
	}
}

func InitCheckers() {
	ticker := time.NewTicker(30 * time.Second)

	for ; ; <-ticker.C {
		tx := database.DB.Begin()
		if tx.Error != nil {
			log.Println("Ошибка при старте транзакции")
		}

		sc := NewServerChecker(10 * time.Second)

		in := make(chan models.Server)
		out := make(chan WorkerResult)
		done := make(chan interface{})
		var wg sync.WaitGroup

		go consumer(out, tx, done)

		go producer(in, tx)

		for range 3 {
			wg.Go(func() {
				worker(sc, in, out)
			})
		}

		wg.Wait()
		close(out)
		<-done
		log.Println("Завершение цикла проверки, коммит транзакции")
		commit := tx.Commit()
		if commit.Error != nil {
			log.Println("Ошибка при коммите транзакции")
		}
	}
}
