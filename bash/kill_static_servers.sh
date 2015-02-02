#!/bin/bash

# Kill any process runing exe/new_iteraid_static_server
# NOTE, it's possible this will mess up users processes if they have any outside
# of Iteraid that match
kill $(ps aux | grep 'new_iteraid_static_server' | awk '{print $2}')
