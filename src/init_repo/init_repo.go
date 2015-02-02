// Creates a new instance of the Repo.  There must always be one instance that just
// stays on master and is not used for creating branches -- this "master" repo instance is
// used for printing out available branches and other general tasks

package main

import (
	"fmt"
	"git_objects"
	"gopkg.in/mgo.v2/bson"
	"mongo_utils"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func getBaseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename) + "/../../"
}

// Create a new repo instance and store in db
func main() {
	config := mongo_utils.ParseConfig()
	_, repos, repoSession := mongo_utils.GetRepoCollection()
	var allRepos git_objects.Repos
	repos.Find(bson.M{}).All(&allRepos)

	id := bson.NewObjectId()
	var repoDirName string
	var newRepo git_objects.RepoInstance
	if len(allRepos) == 0 { // create the initial repo as isMaster=true
		repoDirName = config.RepoName // use the base repo name for the master repo
		newRepo = git_objects.RepoInstance{IsAvailable: false, DirName: repoDirName, Id: id, IsMaster: true}
	} else {
		repoDirName = config.RepoName + "_" + id.Hex()
		newRepo = git_objects.RepoInstance{IsAvailable: true, DirName: repoDirName, Id: id, IsMaster: false}
	}

	repos.UpsertId(newRepo.Id, &newRepo)
	repoSession.Close()

	baseDir := getBaseDir()
	cmd := exec.Command(baseDir+"bash/init_repo.sh", config.RepoSSH, repoDirName, baseDir)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR IN CMD INIT_REPO: ", err.Error())
	}
}
