package work

import (
	"github.corp.ebay.com/bt-siteops/backup/utils"
	"log"
	"os"
)

func (Input Input) rsync(destination string, sshKey string, errc chan error) {
	source := Input.FilePath + Input.FileName
	if len(source) > 0 && len(destination) > 0 {
		rsyncErr := utils.OsExecStdOut("rsync","","--progress","--remove-source-files","--bwlimit=8750",
			"-azh","-e ssh",source,destination)
		if rsyncErr != nil {
			errc <- rsyncErr
		}
		log.Println(source, destination)
	}
}

func (filePath Input) encrypt(sshKeyPath string, errc chan error) (Input) {
	var input Input
	if len(filePath.FileName) > 0 && len(sshKeyPath) > 0 {
		input.FilePath = filePath.FilePath
		input.FileName = filePath.FileName + ".aes"
		execErr := utils.OsExecStdOut("openssl","","enc","-in",filePath.FilePath + filePath.FileName ,"-out",
			input.FilePath + input.FileName,"-e","-aes256","-k",sshKeyPath)
		if execErr != nil {
			log.Fatal("Execution error for openssl encrypt", execErr)
			errc <- execErr
		}
		execErr = os.Remove(filePath.FilePath + filePath.FileName)
		if execErr != nil {
			log.Fatal("Execution error for removing the original file", execErr)
			errc <- execErr
		}
	}
	return input
}