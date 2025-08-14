package transcoder

import (
	"context"
	"sync"

	"github.com/Mirsadikovv/ffmpeg_research/dto"
)

// Queue управляет очередью задач транскодирования
type Queue struct {
	transcoder *Transcoder
	jobs       []*dto.Job
	workers    int
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewQueue создает новую очередь
func NewQueue(transcoder *Transcoder, workers int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		transcoder: transcoder,
		jobs:       make([]*dto.Job, 0),
		workers:    workers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// AddJob добавляет задачу в очередь
func (q *Queue) AddJob(job *dto.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = append(q.jobs, job)
}

// Start запускает обработку очереди
func (q *Queue) Start() {
	for i := 0; i < q.workers; i++ {
		go q.worker()
	}
}

// Stop останавливает обработку очереди
func (q *Queue) Stop() {
	q.cancel()
}

// worker обрабатывает задачи из очереди
func (q *Queue) worker() {
	for {
		select {
		case <-q.ctx.Done():
			return
		default:
			job := q.getNextJob()
			if job != nil {
				q.transcoder.Execute(q.ctx, job)
			}
		}
	}
}

// getNextJob получает следующую задачу из очереди
func (q *Queue) getNextJob() *dto.Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, job := range q.jobs {
		if job.Status == dto.StatusPending {
			q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
			return job
		}
	}

	return nil
}

// GetJobs возвращает все задачи
func (q *Queue) GetJobs() []*dto.Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	jobs := make([]*dto.Job, len(q.jobs))
	copy(jobs, q.jobs)
	return jobs
}
