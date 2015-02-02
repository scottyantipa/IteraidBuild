package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git_objects"
	"github.com/wsxiaoys/terminal/color"
	"gopkg.in/mgo.v2/bson"
	"log"
	"mongo_utils"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func getBaseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename) + "/../../" // add the /../ because we are nested in src/handle_requests/
}

// Find out if a given port is available for use by simply iterating through
// each existing branch
func portIsEmpty(port int, branches git_objects.Branches) bool {
	portString := strconv.Itoa(port)
	for _, branch := range branches {
		if branch.Port == portString {
			return false
		}
	}
	return true
}

func serveComment(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	switch r.Method {
	case "GET":
		{
			getCommentsForBranch(w, r)
		}
	case "POST":
		{
			storeComment(w, r)
		}
	case "DELETE":
		{
			deleteComment(w, r)
		}
	case "OPTIONS":
		{
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// return all commets for a given branch
func getCommentsForBranch(w http.ResponseWriter, r *http.Request) {
	branchName := r.FormValue("branchName")
	var branch git_objects.Branch
	_, collection, session := mongo_utils.GetBranchCollection()
	collection.Find(bson.M{"Name": branchName}).One(&branch)
	session.Close()

	_, collection, session = mongo_utils.GetCommentsCollection()
	var allComments git_objects.Comments
	collection.Find(bson.M{"BranchId": branch.Id}).All(&allComments)
	session.Close()

	writeJson(w, &allComments)
}

// store new comment
func storeComment(w http.ResponseWriter, r *http.Request) {
	branchName := r.FormValue("branchName")
	text := r.FormValue("text")

	// get branch object for comment
	var branch git_objects.Branch
	_, collection, session := mongo_utils.GetBranchCollection()
	collection.Find(bson.M{"Name": branchName}).One(&branch)
	session.Close()

	_, collection, session = mongo_utils.GetCommentsCollection()
	comment := &git_objects.Comment{Text: text, Id: bson.NewObjectId(), BranchId: branch.Id}
	collection.UpsertId(comment.Id, comment)
	session.Close()

	w.Write([]byte("Success storing comment"))
}

// delete comment for a given branch
func deleteComment(w http.ResponseWriter, r *http.Request) {
}

func serveBranch(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	switch r.Method {
	case "GET":
		{
			switch r.FormValue("whichBranches") {
			case "all":
				getAllBranches(w, r)
			}
		}
	case "POST":
		{
			initBranch(w, r)
		}
	case "DELETE":
		{
			deleteBranch(w, r)
		}
	case "OPTIONS":
		{
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Creates a new branch given a branch name
func initBranch(w http.ResponseWriter, r *http.Request) {
	// connect to mongo collection
	err, collection, session := mongo_utils.GetBranchCollection()
	if err != nil {
		http.Error(w, "Couldnt connect to DB", http.StatusInternalServerError)
		log.Fatal(err)
	}

	// get all branches
	var builtBranches git_objects.Branches
	collection.Find(bson.M{"Port": bson.M{"$ne": ""}}).All(&builtBranches)

	// For now do an (expensive) loop over all branches for each port
	var availablePort int
	foundPort := false
	start, end := getPortRange()
	port := start
	for {
		if port > end {
			http.Error(w, "Cannot build more than max number of branches", http.StatusInternalServerError)
			break
		}
		if portIsEmpty(port, builtBranches) {
			foundPort = true
			availablePort = port
			break
		}
		port++
	}

	if !foundPort {
		http.Error(w, "No more available ports", http.StatusNotFound)
		fmt.Println("didnt find available port so returning")
		return
	}

	newPortStr := strconv.Itoa(availablePort)

	// get name and commit, then save to db
	branchName := r.FormValue("name") // get name property from request

	var branchDoc git_objects.Branch
	collection.Find(bson.M{"Name": branchName}).One(&branchDoc)
	nilBranch := git_objects.Branch{}
	if branchDoc == nilBranch {
		http.Error(w, "Couldnt create branch because it wasnt found in git", http.StatusConflict)
		return
	}

	commitErr, commitHash := getCommitHash(branchName)
	if commitErr != nil {
		http.Error(w, "Couldnt connect to DB", http.StatusInternalServerError)
		color.Println("@r Couldnt find commit hash, not creating branch")
		fmt.Println(branchName)
		return
	}
	color.Println("@gReceived POST for branch, hash: ", branchName, commitHash)

	branchDoc.CommitHash = commitHash
	branchDoc.Status = "Waiting"
	branchDoc.Port = newPortStr

	_, saveErr := collection.UpsertId(branchDoc.Id, &branchDoc)
	if saveErr != nil {
		http.Error(w, "Couldnt save branch", http.StatusConflict)
		panic(saveErr)
	}

	session.Close() // when does this need to be closed?  seems appropriate here...

	// send this here so that status will be "Building" on client side
	w.Write([]byte("Success posting branch"))
}

// handle request to delete a branch
func deleteBranch(w http.ResponseWriter, r *http.Request) {
	branchName := r.FormValue("name") // get name property from request
	killedErr := unbuildBranch(branchName)
	if killedErr != nil {
		http.Error(w, "Couldnt delete branch", http.StatusConflict)
		return
	}
	w.Write([]byte("Success deleting branch"))
}

// utility to remove a branch directory as well as kill the server
func unbuildBranch(branchName string) error {
	err, collection, session := mongo_utils.GetBranchCollection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// get the branch so we can kill the pid and update branch Status
	var branch git_objects.Branch
	collection.Find(bson.M{"Name": branchName}).One(&branch)

	// remove branch directory
	cmd := exec.Command(getBaseDir()+"bash/delete_branch.sh", branchName, branch.Pid, getBaseDir())
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cdmErr := cmd.Run()
	if cdmErr != nil {
		return cdmErr
	}

	branch.Status = "Unbuilt"
	branch.Port = ""
	branch.Pid = ""
	collection.UpsertId(branch.Id, branch)
	session.Close()

	return nil
}

// When a branch is removed from git repo, we need to remove it from db
func branchNoLongerInGit(branchName string) error {
	err, collection, session := mongo_utils.GetBranchCollection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	collection.Remove(bson.M{"Name": branchName})
	session.Close()
	return nil
}

// returns all branches from db
func getAllBranches(w http.ResponseWriter, r *http.Request) {
	_, collection, session := mongo_utils.GetBranchCollection()

	// get all branches
	var allBranches git_objects.Branches
	collection.Find(bson.M{}).All(&allBranches)
	session.Close()

	branchMap := make(map[string]git_objects.Branch)
	for _, branch := range allBranches {
		branchMap[branch.Name] = branch
	}
	writeJson(w, &branchMap)
}

// simply returns unformatted results of git branch -r
func getAllGitBranchesFromGit() []string {
	config := mongo_utils.ParseConfig()
	repoDir := getBaseDir() + "repo/" + config.RepoName

	cmd := exec.Command("bash", "-c", "cd "+repoDir+"; git fetch -p; git branch -r")
	cmd.Env = os.Environ()
	var outbuf2 bytes.Buffer
	var errbuf2 bytes.Buffer
	cmd.Stdout = &outbuf2
	cmd.Stderr = &errbuf2
	err := cmd.Run()
	if err != nil {
		fmt.Println("error from running git branch in getAllGit: ", err.Error())
		panic(err)
	}

	branches := outbuf2.String()
	return strings.Split(branches, "\n")
}

// Get the most recent commit hash for a branch
func getCommitHash(branchName string) (error, string) {
	config := mongo_utils.ParseConfig()
	repoDir := getBaseDir() + "repo/" + config.RepoName
	cmd := exec.Command("bash", "-c", "cd "+repoDir+"; git fetch -p; git log --pretty=format:\"%H\" -1 origin/"+"'"+branchName+"'") // we put quotes around branch name so strange characters in branch name don't interfere
	cmd.Env = os.Environ()
	var outbuf bytes.Buffer
	var errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		fmt.Println("Couldn't get commit hash with error and repoDir: ", errbuf.String(), repoDir)
		return err, ""
	}
	return nil, outbuf.String()
}

// generic util for writing json to client
// got this from http://denvergophers.com/2013-04/mgo.article
func writeJson(w http.ResponseWriter, v interface{}) {
	// avoid json vulnerabilities, always wrap v in an object literal
	doc := map[string]interface{}{"data": v}

	if data, err := json.Marshal(doc); err != nil {
		log.Printf("Error marshalling json: %v", err)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

// allow cross origin requests. http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT, GET, OPTIONS")
}

// user defines their start port in config.json
// let's only allow up to 100 open ports for now
func getPortRange() (int, int) {
	config := mongo_utils.ParseConfig()
	start, _ := strconv.Atoi(config.StartPort)
	end := start + 99
	return start, end
}

var isModelingNow bool

func waitAndModelBranches() {
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if !isModelingNow {
					modelAllBranches()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// On an interval, check if there are branches waiting to be built,
// If so (and no builds are in process) build that branch
func waitAndBuild() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if !aBranchIsBuilding() { // clear to build, build the branch
					branchToBuild := getFirstWaitingBranch()
					nilBranch := git_objects.Branch{}
					if branchToBuild != nilBranch {
						buildBranch(branchToBuild)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Loop through all branches to see if they are out of date
// if so, update the hash and mark as "Waiting" if they were already built
func waitAndUpdateBranches() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// fmt.Println("WITHIN WAIT_AND_BUILD LOOP")
				// get collection
				_, collection, session := mongo_utils.GetBranchCollection()
				var allBranches git_objects.Branches
				collection.Find(bson.M{}).All(&allBranches)

				for _, branch := range allBranches {
					storedHash := branch.CommitHash
					err, currentHash := getCommitHash(branch.Name)
					if err != nil {
						color.Println("@rCouldnt find commit hash, killing branch")
						if branch.Status == "Built" {
							unbuildBranch(branch.Name)
						}
						branchNoLongerInGit(branch.Name)
						continue
					}
					if storedHash != currentHash {
						color.Println("@yBranch was out of date, marking with new hash ", branch.Name, currentHash)
						// it's out of date
						branch.CommitHash = currentHash
						if branch.Status == "Built" {
							branch.Status = "Waiting"
						}
						collection.UpsertId(branch.Id, branch)
					}
				}
				session.Close()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// return true if there are any branches being built currently
func aBranchIsBuilding() bool {
	err, collection, session := mongo_utils.GetBranchCollection()
	if err != nil {
		fmt.Println(err)
		return true
	}
	// get all branches and see if any are Building
	var branchesBeingBuilt git_objects.Branches
	collection.Find(bson.M{"Status": "Building"}).All(&branchesBeingBuilt)
	session.Close()
	return len(branchesBeingBuilt) > 0
}

func getFirstWaitingBranch() git_objects.Branch {
	err, collection, session := mongo_utils.GetBranchCollection()
	if err != nil {
		fmt.Println(err)
		return git_objects.Branch{}
	}

	// note, eventually we should return the oldest oustanding Waiting branch
	// so that it is a proper time based priority queue.
	var waitingBranch git_objects.Branch
	collection.Find(bson.M{"Status": "Waiting"}).One(&waitingBranch)
	session.Close()
	return waitingBranch
}

func buildBranch(branchToBuild git_objects.Branch) {
	// set status of branch to "Building"
	fmt.Println("Going to build branch: ", branchToBuild)
	_, branches, branchSession := mongo_utils.GetBranchCollection()
	branchToBuild.Status = "Building"
	branches.UpsertId(branchToBuild.Id, branchToBuild)
	branchSession.Close()

	baseDir := getBaseDir()

	// now find an available repo instance, or create one
	_, repos, repoSession := mongo_utils.GetRepoCollection()
	var availableRepo git_objects.RepoInstance
	repos.Find(bson.M{"isAvailable": true}).One(&availableRepo)
	var dirName string
	nilRepo := git_objects.RepoInstance{} // hacky...
	if availableRepo == nilRepo {
		repCmd := exec.Command("bash", "-c", baseDir+"bin/init_repo")
		repCmd.Env = os.Environ()
		repCmd.Stdin = os.Stdin
		repCmd.Stdout = os.Stdout
		repCmd.Stderr = os.Stderr
		repCmd.Run()
		var newRepo git_objects.RepoInstance
		repos.Find(bson.M{"isAvailable": true}).One(&newRepo)
		dirName = newRepo.DirName
		fmt.Println("Created new repo for init_branch here: ", dirName)
	} else {
		dirName = availableRepo.DirName
		fmt.Println("Found available repo instance here: ", dirName)
	}
	repoSession.Close()

	cmd := exec.Command(baseDir+"bash/init_branch.sh", branchToBuild.Name, branchToBuild.Port, dirName, baseDir)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go cmd.Run() // run static servers in paralell
}

// This gets called on intial serve and it looks at all git branches and stores them in the db
func modelAllBranches() {
	unformatted := getAllGitBranchesFromGit()
	_, collection, session := mongo_utils.GetBranchCollection()
	for index, branchName := range unformatted {

		// Don't model branch if it is
		if index == 0 || branchName == "" {
			continue
		} else if strings.Contains(branchName, "origin/master") {
			branchName = "master"
		} else {
			branchName = strings.Replace(branchName, "  origin/", "", 1)
		}

		hashErr, hash := getCommitHash(branchName)
		if hashErr != nil {
			continue // its an improper branch name like "HEAD"
		} else {
			// its a real branch so either create new branch or do nothing if its already in db
			var existingBranch git_objects.Branch
			collection.Find(bson.M{"Name": branchName}).One(&existingBranch)
			nilBranch := git_objects.Branch{}
			if existingBranch != nilBranch {
				continue
			}
			// branch doesnt exist so create it
			branchDoc := git_objects.Branch{Name: branchName, CommitHash: hash, Id: bson.NewObjectId(), Status: "Unbuilt"}
			_, saveError := collection.UpsertId(branchDoc.Id, &branchDoc)
			if saveError != nil {
				panic(saveError)
			}
		}
	}
	isModelingNow = false
	session.Close()

}

func main() {

	// jobs to do continuously, like updating/building branches
	modelAllBranches()
	waitAndModelBranches()
	waitAndBuild()
	waitAndUpdateBranches()

	// routes
	http.HandleFunc("/branch", serveBranch)
	http.HandleFunc("/comment", serveComment)

	// start the static server for the UI
	config := mongo_utils.ParseConfig()
	fmt.Println("Starting ui server on port", config.UIPort)
	go func() {
		http.ListenAndServe(":"+config.UIPort, http.FileServer(http.Dir(getBaseDir()+"ui/public")))
	}()

	http.ListenAndServe(":"+config.MainPort, nil)
}
