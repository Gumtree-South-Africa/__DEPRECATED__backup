package work

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Input struct {
	FileName string
	FilePath string
}
type MySqlDb struct {
	Name string `json:"name"`
}

type MongoDb struct {
	Name string `json:"name"`
}

type Mysql struct {
	Username string `json:"mysql.user"`
	DBList   []MySqlDb `json:"mysql.dbName"`
	Host     string `json:"mysql.host"`
	Port     string `json:"mysql.port"`
	PoolPath string `json:"mysql.zfspool"`
}

type Mongo struct {
	Host    string `json:"mongo.host"`
	DBList []MongoDb  `json:"mongo.dbName"`
	Port    string `json:"mongo.port"`
}

type Encryptonator struct {
	Username string `json:"encryptonator.user"`
	SSHKey   string `json:"encryptonator.ssh_key"`
	Path     string `json:"encryptonator.path"`
}

type Data struct {
	Mysql Mysql `json:"mysql"`
	Mongo Mongo `json:"mongo"`
}

type JSONInput struct {
	WorkerCnt     int           `json:"workerNodes"`
	Data          Data          `json:"data"`
	Encryptonator Encryptonator `json:"encryptonator"`
}

type Transfer interface {
	rsync(destination string, errc chan error)
	encrypt(sshKeyPath string, errc chan error) Input
}

type DBDump interface {
	dump(host string, port string, username string, tableName string, destination string) (Input,error)
}

func RetrieveFeeds(dataFile string) (*JSONInput, error) {
	fmt.Println("File name:", dataFile)
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var jsonInput *JSONInput
	err = json.Unmarshal(file, &jsonInput)
	return jsonInput, err
}
