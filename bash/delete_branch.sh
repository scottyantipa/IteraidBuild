#!/bin/bash

# Remove the directory of a branch
# and kill the process of the file server
# usage: delete_branch.sh feature_branch pid

BRANCH_NAME=$1
PID=$2
BASE_DIR=$3

rm -rf $BASE_DIR/repo/builds/$BRANCH_NAME

# kill pid
echo "Killing process $PID"
kill $PID || true
