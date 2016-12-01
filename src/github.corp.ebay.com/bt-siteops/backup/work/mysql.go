package work

import (
	"strconv"
	"time"
	"os/exec"
	"log"
	"io/ioutil"
	"os"
	"bytes"
	"compress/gzip"
)

func (Mysql Mysql) dump(host string, port string, username string, tableName string,
			destination string, errc chan error) Input {
	var input Input
	if len(tableName) > 0 {
		var b bytes.Buffer
		input.FileName = tableName + strconv.FormatInt(time.Now().UnixNano(), 16) + ".sql.gz"
		input.FilePath = destination

		cmd := exec.Command("mysqldump","-h",host,"-P",port,"-u", username, tableName)
		log.Println(cmd.Args)
		output, execErr := cmd.Output()
		if execErr != nil {
			log.Fatal("Execution error for mysqldump", execErr)
			errc <- execErr
		}
		w := gzip.NewWriter(&b)
		w.Write(output)
		w.Close()

		writeerr := ioutil.WriteFile(input.FilePath+input.FileName, b.Bytes(), 0644)

		if writeerr != nil {
			log.Fatal("Write error mysqldump", writeerr)
			errc <- writeerr
		}

		errc <- nil
	}
	return input
}

func (Mysql Mysql) rsync(source string, destination string, errc chan error) {
	if len(source) > 0 && len(destination) > 0 {
		cmd := exec.Command("/usr/bin/rsync", "-avx", "-e", "\"ssh -o StrictHostKeyChecking=no -o " +
			"UserKnownHostsFile=/dev/null\"", source, destination)
		out, err := os.Create("/tmp/rsync.log")
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		if err != nil {
			log.Fatal("Error while creating error log", err)
			errc <- err
		}
		execErr := cmd.Run()
		if execErr != nil {
			log.Fatal("Error while performing rsync", execErr)
			errc <- execErr
		}
		log.Println(source, destination)
	}
}

func (Mysql Mysql) encrypt(filePath Input, sshKeyPath string, errc chan error) (Input) {
	var input Input
	if len(filePath.FileName) > 0 && len(sshKeyPath) > 0 {
		input.FilePath = filePath.FilePath
		input.FileName = filePath.FileName + ".aes"
		cmd := exec.Command("/usr/bin/openssl","enc","-in",filePath.FilePath + filePath.FileName ,"-out",
			input.FilePath + input.FileName,"-e","-aes256","-k",sshKeyPath)
		log.Println(cmd.Args)
		out, err := os.Create("/tmp/ssl.log")
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		if err != nil {
			log.Fatal("Error while creating error log", err)
			errc <- err
		}

		execErr := cmd.Run()

		if execErr != nil {
			log.Fatal("Execution error for openssl encrypt", execErr)
			errc <- execErr
		}
	}
	return input
}
