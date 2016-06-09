#!/bin/sh

echo "Restore vendors"

gvt restore


mkdir -p /go/src/github.com/postgres-ci/notifier/

cp -r . /go/src/github.com/postgres-ci/notifier/