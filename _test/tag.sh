#!/bin/bash

# get latest tag
t=git describe --tags `git rev-list --tags --max-count=1` > /dev/null 2>&1

if [ -z "$t" ]
then
      t=0.0.0
fi

echo $t
