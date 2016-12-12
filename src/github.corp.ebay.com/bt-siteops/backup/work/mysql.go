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

//func (Mysql Mysql) rsync(source string, destination string, sshKey string, errc chan error) {
//	if len(source) > 0 && len(destination) > 0 {
//		rsyncErr := utils.OsShellExec("rsync","rsync","--progress","--remove-source-files","--bwlimit=50",
//			"-azh","-e ssh",source,destination)
//		if rsyncErr != nil {
//			errc <- rsyncErr
//		}
//		log.Println(source, destination)
//	}
//}
//
//func (Mysql Mysql) encrypt(filePath Input, sshKeyPath string, errc chan error) (Input) {
//	var input Input
//	if len(filePath.FileName) > 0 && len(sshKeyPath) > 0 {
//		input.FilePath = filePath.FilePath
//		input.FileName = filePath.FileName + ".aes"
//		execErr := utils.OsExecStdOut("openssl","","enc","-in",filePath.FilePath + filePath.FileName ,"-out",
//			input.FilePath + input.FileName,"-e","-aes256","-k",sshKeyPath)
//		if execErr != nil {
//			log.Fatal("Execution error for openssl encrypt", execErr)
//			errc <- execErr
//		}
//		execErr = os.Remove(filePath.FilePath + filePath.FileName)
//		if execErr != nil {
//			log.Fatal("Execution error for removing the original file", execErr)
//			errc <- execErr
//		}
//	}
//	return input
//}
