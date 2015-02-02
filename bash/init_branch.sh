#!/bin/bash

# Initialize a new branch for the repo
# use like: source init_branch.sh scott-feature-1 8081 RepoName_3252134

BRANCH_NAME=$1 # e.g. scott-feature1
PORT=$2  # e.g. :8081
REPO_NAME=$3 # RepoName_23452345245
BASE_DIR=$4 # base directory of repo

# relative paths from base dir file
GIT_REPO=$BASE_DIR/repo/$REPO_NAME
BUILD_DIR=$BASE_DIR/repo/builds/$BRANCH_NAME

# Checkout the repo and update it
cd $GIT_REPO
git fetch -p
git checkout master
git clean -fd
git checkout .
git pull

# Checkout specific branch
git checkout origin/$BRANCH_NAME
git submodule update

# Build the source
# TODO pass in the directory of the repo
source $BASE_DIR/user/build_branch.sh

# Create a directory to put this build in
rm -rf $BUILD_DIR
mkdir $BUILD_DIR

# run the dist script
echo "Starting dist_copy_all"
source $BASE_DIR/bash/dist_copy_all.sh $BUILD_DIR $GIT_REPO
echo "Finished dist_copy_all"

# checkout master just to clean up
cd $GIT_REPO
git clean -fd
git checkout .
git checkout master
cd ../..

# at this point we are done with the repo so we can mark it as 'not in use'
$BASE_DIR/bin/mark_repo_open $REPO_NAME

# start server for the build
$BASE_DIR/bin/new_iteraid_static_server $PORT $BUILD_DIR $BRANCH_NAME
exit
