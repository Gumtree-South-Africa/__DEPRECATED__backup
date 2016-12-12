package work

import (
	"strconv"
	"log"
	"github.corp.ebay.com/bt-siteops/backup/utils"
	"time"
)

func (Mysql MySqlDb) Dump(host string, port string, username string, destination string) (Input,error) {
	tableName := Mysql.Name
	var input Input
	if len(tableName) > 0 {
		input.FileName = tableName + strconv.FormatInt(time.Now().UnixNano(), 16) + ".sql"
		input.FilePath = destination
		err := utils.OsExecStdOut("mysqldump",input.FilePath+input.FileName,"-h",host,"-P",
			port,"-u",username,tableName)
		if err != nil {
			return input,err
		}
		gzipErr := utils.OsExecStdOut("gzip","",input.FilePath+input.FileName)
		if gzipErr != nil {
			log.Fatal("Error while compressing ", gzipErr)
			return input,err
		}
		input.FileName = input.FileName + ".gz"
	}
	return input,nil
}

