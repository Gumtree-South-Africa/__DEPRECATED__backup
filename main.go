 package main

import (
	"log"
	"os"

	"github.corp.ebay.com/bt-siteops/backup/work"
	"runtime"
	"sync"
)

func init() {
	log.SetOutput(os.Stdout)
}
func main() {
	CpuCnt := runtime.NumCPU() // Count down in select clause
	CpuOut := CpuCnt           // Save for print report

	log.Println("Processors: ", CpuOut)

	if len(os.Args) < 3 {
		log.Println("Missing parameter, provide file name!")
		return
	}

	fileName := os.Args[1]
	dbType := os.Args[2]
	feeds, err := work.RetrieveFeeds(fileName)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	runtime.GOMAXPROCS(feeds.WorkerCnt)
	WorkQueue := make(chan work.Job)
	done := make(chan *bool, len(feeds.Data.Mysql.DBList))
	WorkerQueue := make(chan chan work.Job, feeds.WorkerCnt)

	var wg sync.WaitGroup
	//Start the workers
	work.StartDispatcher(feeds, WorkQueue, WorkerQueue, feeds.WorkerCnt, done)
	switch dbType {
	case "mysql" :
		addJob := func(feeds *work.JSONInput,wg *sync.WaitGroup) {
				wg.Add(len(feeds.Data.Mysql.DBList))
				for _, feed := range feeds.Data.Mysql.DBList {
					log.Println("Inside the db list loop", feed)
					data, err := feed.Dump(feeds.Data.Mysql.Host, feeds.Data.Mysql.Port,
						feeds.Data.Mysql.Username,
					feeds.Data.Mysql.PoolPath)
					if err != nil {
					log.Fatal("Error while dumping the mysql database ", feed.Name)
					}
					work.Generator(WorkQueue,&data,done,wg)
				}
			}
		addJob(feeds,&wg)
	case "mongo" :
		//do nothing
	}

	go func() {
		wg.Wait()
		close(done)
		close(WorkQueue)
	}()

	<-done
}
