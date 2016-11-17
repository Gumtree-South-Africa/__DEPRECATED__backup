package work

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)
type Input struct {
	FileName string
	FilePath string
}
type Feed struct {
	Name string `json:"name"`
}

type Mysql struct {
	Username string `json:"mysql.user"`
	DBList []Feed `json:"mysql.dbName"`
}
type Encryptonator struct {
	Username string `json:"encryptonator.user"`
	SSHKey   string `json:"encryptonator.ssh_key"`
	Path     string `json:"encryptonator.path"`
}

type JSONInput struct {
	WorkerCnt int `json:"workerNodes"`
	Mysql Mysql `json:"mysql"`
	Encryptonator Encryptonator `json:"encryptonator"`
}


func RetrieveFeeds(dataFile string) (*JSONInput, error) {
	fmt.Println("File name:", dataFile)
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil,err
	}
	var jsonInput *JSONInput
	err = json.Unmarshal(file,&jsonInput)
	return jsonInput,err
}
