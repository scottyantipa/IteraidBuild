// Sets up a static file server for a given port and directory
package main

import (
	"fmt"
	"git_objects"
	"gopkg.in/mgo.v2/bson"
	"mongo_utils"
	"net/http"
	"os"
	"strconv"
)

// create a static file server at the port (Args[1]) for directory (Args[2])
// store the process in the mongo doc so we can kill it later
func main() {
	_, collection, session := mongo_utils.GetBranchCollection()
	port := os.Args[1]
	dir := os.Args[2]
	branchName := os.Args[3]

	// retrieve the mongo branch and update it with the proces id
	var branch git_objects.Branch
	collection.Find(bson.M{"Name": branchName}).One(&branch)
	pid := os.Getpid() // process id of this static file server
	branch.Pid = strconv.Itoa(pid)

	branch.Status = "Built"
	collection.UpsertId(branch.Id, branch)
	session.Close()

	fmt.Println("Starting branch server: ", branch.Name+" ", port)
	http.ListenAndServe(":"+port, http.FileServer(http.Dir(dir)))
}
