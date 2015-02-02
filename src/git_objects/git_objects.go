package git_objects

import (
	"gopkg.in/mgo.v2/bson"
)

// Struct describing a branch to be served
// Name = name of git branch
// Port = port to serve it on
type Branch struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name       string        `json:"Name" bson:"Name"`
	Port       string        `json:"Port" bson:"Port"`
	Pid        string        `json:"Pid" bson:"Pid"`
	Status     string        `json:"Status" bson:"Status"`
	CommitHash string        `json:"CommitHash" bson:"CommitHash"`
}

type Branches []Branch

type RepoInstance struct {
	Id          bson.ObjectId `bson:"id" bson:"_id,omitempty"`
	DirName     string        `bson:"dirName"`
	IsAvailable bool          `bson:"isAvailable"`
	IsMaster    bool          `bson:"isMaster"` // we always want one repo that just stays on master and is never used other than for printing branches, etc.
}

type Repos []RepoInstance

type Comment struct {
	Id       bson.ObjectId `json:"Id" bson:"_Id,omitempty"`
	BranchId bson.ObjectId `json:"BranchId" bson:"BranchId"`
	Text     string        `json:"Text" bson:"Text"`
}

type Comments []Comment
