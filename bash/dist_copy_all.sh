#!/bin/bash

# Creates a directory with everything you neeed to run the app
# For now, we just copy the entire repo over
# Maybe at some point allow users to include their own script

DEST_DIR=$1 # directory where we will put the dist
REPO_DIR=$2 # the git repo to make dist from

rm -rf $DEST_DIR
cp -R $REPO_DIR $DEST_DIR
