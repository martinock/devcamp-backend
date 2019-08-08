package internal

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var jobs = make(chan []string)

func (h *Handler) RunWorkerPool(data []string) {
	var wg sync.WaitGroup
	start := time.Now()

	// create worker
	go h.CreateWorker(1, &wg)

	// create jobs
	h.CreateJobs(data)

	log.Println("[Worker] Finished with time", time.Since(start).Seconds(), "s")

}

// CreateWorker create worker as much as {total}
func (h *Handler) CreateWorker(total int, wg *sync.WaitGroup) {
	for i := 0; i < total; i++ {
		go h.Worker(total, wg)
		time.Sleep(1 * time.Second)
	}
}

// CreateJobs create jobs to send data to jobs channel
func (h *Handler) CreateJobs(data []string) {
	batchSize := 3

	for start := 0; start < len(data); start += batchSize {
		end := start + batchSize

		if end > len(data) {
			end = len(data) - 1
		}
		jobs <- data[start:end]
	}
	close(jobs)
}

// Worker used to run process
func (h *Handler) Worker(workerIndex int, wg *sync.WaitGroup) {
	wg.Add(1)

	for job := range jobs {
		h.InsertToDB(job)
	}

	wg.Done()
}

// InsertToDB insert data to DB with multiple queries
func (h *Handler) InsertToDB(rows []string) {
	queries := ""

	for _, row := range rows {

		// Split rows to column
		columns := strings.Split(row, ",")
		query := fmt.Sprintf("INSERT INTO books (id, title, author, isbn, stock) VALUES (%s, '%s', '%s', '%s', %s);", columns[0], columns[1], columns[2], columns[3], columns[4])

		queries += query
	}

	_, err := h.DB.Exec(queries)
	if err != nil {
		log.Println(err)
	}
}
