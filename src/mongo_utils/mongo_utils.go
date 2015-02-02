package mongo_utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"path"
	"runtime"
)

type Config struct {
	RepoName  string `json:"repoName"`
	StartPort string `json:"startPort"`
	MongoPort string `json:"mongoPort"`
	HostName  string `json:"hostName"`
	MainPort  string `json:"mainPort"`
	RepoSSH   string `json:"repoSSH"`
	UIPort    string `json:"UIPort"`
}

func getCollectionForDoc(docName string) (error, *mgo.Collection, *mgo.Session) {
	config := ParseConfig()
	session, err := mgo.Dial(":" + config.MongoPort)
	if err != nil {
		fmt.Println("ERROR TRYING TO DIAL: ", config.HostName+":"+config.MongoPort)
		fmt.Println(err.Error())
		return err, nil, nil
	}
	collection := session.DB("iteraid").C(docName)
	return nil, collection, session
}

func GetBranchCollection() (error, *mgo.Collection, *mgo.Session) {
	return getCollectionForDoc("branches")
}

func GetRepoCollection() (error, *mgo.Collection, *mgo.Session) {
	return getCollectionForDoc("repos")
}

func getBaseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func GetCommentsCollection() (error, *mgo.Collection, *mgo.Session) {
	return getCollectionForDoc("comments")
}

func ParseConfig() Config {
	configFile, err := os.Open(getBaseDir() + "/../../user/config.json")
	if err != nil {
		fmt.Println("PARSE CONFIG ERR:", err.Error())
		return Config{}
	}
	var config Config
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	configFile.Close()
	return config
}
