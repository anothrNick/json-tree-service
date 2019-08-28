#!/bin/bash

# get latest tag
t=$(git describe --tags `git rev-list --tags --max-count=1`) > /dev/null 2>&1

# if there are none, start tags at 0.0.0
if [ -z "$t" ]
then
    log=$(git log --pretty=oneline)
    t=0.0.0
else
    log=$(git log $t..HEAD --pretty=oneline)
fi

echo $t

case "$log" in
    *#major* ) echo "major bump";;
    *#patch* ) echo "patch bump";;
    * ) echo "minor bump";;
esac
