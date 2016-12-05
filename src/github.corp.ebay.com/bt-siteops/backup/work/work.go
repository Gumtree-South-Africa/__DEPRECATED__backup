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
	feed Feed
	done chan *bool
	wg   *sync.WaitGroup
}


func Generator(workQueue chan Job, feed Feed, done chan *bool, wg *sync.WaitGroup) {
	log.Println("Adding job to queue", feed)
	workQueue <- Job{feed, done, wg}
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
				output := func() {
					dest := config.Data.Mysql.dump(config.Data.Mysql.Host, config.Data.Mysql.Port,
						config.Data.Mysql.Username, work.feed.Name, config.Data.Mysql.PoolPath, error)
					encryptedFile := config.Data.Mysql.encrypt(dest,config.Encryptonator.SSHKey,error)
					log.Println(dest.FilePath, dest.FileName, encryptedFile)
					config.Data.Mysql.rsync(encryptedFile.FilePath + encryptedFile.FileName,
						config.Encryptonator.Username + config.Encryptonator.Path,
						config.Encryptonator.SSHKey, error)
				}

				output()
				err := <-error
				if err != nil {
					log.Fatal(err)
				}
				work.wg.Done()
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
