package work

import (

	"time"
	"strconv"
	"os/exec"
	"io/ioutil"
	"sync"
	"log"
	"os"
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
	log.Println("Adding job to queue",feed)
	workQueue <- Job{feed,done,wg}
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
				output := func(){
					dest := mysqldump(config.Mysql.Username, work.feed.Name, "/tmp/", error)
					log.Println(dest.FilePath, dest.FileName, dest)
					rsync(dest.FilePath + dest.FileName, config.Encryptonator.Path, error)
				}

			output()
			err := <- error
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

func mysqldump(username string,tableName string, destination string, errc chan error) (Input) {
	var input Input
	if len(tableName) > 0 {
		input.FileName = tableName + strconv.FormatInt(time.Now().UnixNano(), 16) + ".sql"
		input.FilePath = destination

		cmd := exec.Command("mysqldump", "-u", username,tableName)
		log.Println(cmd.Args)
		output, execErr := cmd.Output()
		if execErr != nil {
			log.Fatal("Execution error for mysqldump", execErr)
			errc <- execErr
		}
		writeerr := ioutil.WriteFile(input.FilePath + input.FileName, output, 0644)
		if writeerr != nil {
			log.Fatal("Write error mysqldump", writeerr)
			errc <- writeerr
		}
		errc <- nil
	}
	return input
}

func rsync(source string, destination string, errc chan error) {
	if len(source) > 0 && len(destination) > 0 {
		cmd := exec.Command("/usr/bin/rsync", "-avx", "-e", "\"ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null\"",source, destination)
		out,err := os.Create("/tmp/rsync.log")
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		if err != nil {
			log.Fatal("Error while creating error log",err)
			errc <- err
		}
		execErr := cmd.Run()
		if execErr != nil {
			log.Fatal("Error while performing rsync",execErr)
			errc <- execErr
		}
		log.Println(source, destination)
	}
}

func StartDispatcher(config *JSONInput, WorkQueue chan Job, WorkerQueue chan chan Job, nworkers int, done chan *bool)([]*Worker) {
	// First, initialize the channel we are going to but the workers' work channels into.
	workerThreads := make([]*Worker, nworkers)
	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Println("Starting worker", i + 1)
		worker := NewWorker(i + 1, WorkerQueue)
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