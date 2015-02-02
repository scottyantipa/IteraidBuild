package main

// Cleans the Mongo repo of all branches and repos

import (
	"mongo_utils"
)

// remove all branch and repo docs from DB
func main() {
	_, collection, session := mongo_utils.GetBranchCollection()
	collection.RemoveAll(nil)
	session.Close()

	_, collection, session = mongo_utils.GetRepoCollection()
	collection.RemoveAll(nil)
	session.Close()

	_, collection, session = mongo_utils.GetCommentsCollection()
	collection.RemoveAll(nil)
	session.Close()
}
