package main

import (
	"log"
	"os"

	"runtime"
	"github.corp.ebay.com/golang/backup/work"
	"sync"
)

func init() {
	log.SetOutput(os.Stdout)
}
func main() {
	CpuCnt := runtime.NumCPU() // Count down in select clause
	CpuOut := CpuCnt           // Save for print report


	log.Println("Processors: ", CpuOut)

	if len(os.Args) < 2 {
		log.Println("Missing parameter, provide file name!")
		return
	}
	fileName := os.Args[1]
	feeds, err := work.RetrieveFeeds(fileName)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	runtime.GOMAXPROCS(feeds.WorkerCnt)
        WorkQueue := make(chan work.Job)
	done := make (chan *bool,len(feeds.Mysql.DBList))
	var wg sync.WaitGroup
	WorkerQueue := make(chan chan work.Job, feeds.WorkerCnt)

	//Start the workers
	work.StartDispatcher(feeds,WorkQueue,WorkerQueue,feeds.WorkerCnt,done)

	wg.Add(len(feeds.Mysql.DBList))

	for _, feed := range feeds.Mysql.DBList {
		log.Println("Inside hte db list loop", feed)
		work.Generator(WorkQueue,feed,done,&wg)
	}

	go func(){
		wg.Wait()
		close(done)
		close(WorkQueue)
	}()

	<-done
}
