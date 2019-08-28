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

case "$log" in
    *#major* ) new=$(semver bump major $t);;
    *#patch* ) new=$(semver bump patch $t);;
    * ) new=$(semver bump minor $t);;
esac

echo $new
