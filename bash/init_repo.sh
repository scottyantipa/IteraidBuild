#!/bin/bash
# Initialize the repo
# usage $ ./init_repo.sh url_to_git_repo

GIT_URL=$1
DIR_NAME=$2
BASE_DIR=$3

mkdir -p $BASE_DIR/repo
mkdir -p $BASE_DIR/repo/builds # create a dir for the builds
echo "Cloning git url: "
echo $GIT_URL
git clone $GIT_URL $BASE_DIR/repo/$DIR_NAME # clone repo

cd $BASE_DIR/repo/$DIR_NAME
git submodule init
cd $BASE_DIR

# Should probably do an initial build here on master just to get all dependencies


