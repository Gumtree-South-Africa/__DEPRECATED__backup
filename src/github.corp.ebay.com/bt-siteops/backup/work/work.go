 package work

import (
	"log"
	"sync"
)

type Worker struct {
	ID          int
	Work        chan Job
	WorkerQueue chan chan Job
	QuitChan    chan bool
}

type Job struct {
	input *Input
	done chan *bool
	wg *sync.WaitGroup
}


func Generator(workQueue chan Job, input *Input, done chan *bool,wg *sync.WaitGroup) {
		workQueue <- Job{input,done,wg}
}

func NewWorker(id int, workerQueue chan chan Job) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan Job),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

func (w *Worker) Start(config *JSONInput) {

	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				// Receive a work request.
				error := make(chan error, 1)
				output := func(work Job) {
					encryptedFile := work.input.encrypt(config.Encryptonator.SSHKey,error)
					log.Println(work.input.FilePath, work.input.FileName, encryptedFile)
					encryptedFile.rsync(config.Encryptonator.Username + config.Encryptonator.Path,
						config.Encryptonator.SSHKey, error)
					work.wg.Done()
				}
				go output(work)
				err := <-error
				if err != nil {
					log.Fatal(err)
				}
			case <-w.QuitChan:
				// We have been asked to stop.
				log.Println("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func StartDispatcher(config *JSONInput, WorkQueue chan Job, WorkerQueue chan chan Job, nworkers int, done chan *bool) []*Worker {
	// First, initialize the channel we are going to but the workers' work channels into.
	workerThreads := make([]*Worker, nworkers)
	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		workerThreads[i] = &worker
		worker.Start(config)
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Println("Received work requeust")
				go func(work Job) {
					worker := <-WorkerQueue

					log.Println("Dispatching work request")
					worker <- work
				}(work)
			}
		}
	}()
	return workerThreads
}
