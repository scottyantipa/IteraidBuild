// Mark a repo as available for use
package main

import (
	"git_objects"
	"gopkg.in/mgo.v2/bson"
	"mongo_utils"
	"os"
)

func main() {
	dirName := os.Args[1]
	_, collection, session := mongo_utils.GetRepoCollection()

	var repo git_objects.RepoInstance
	collection.Find(bson.M{"dirName": dirName}).One(&repo)
	repo.IsAvailable = true
	collection.UpsertId(repo.Id, repo)
	session.Close()
}
