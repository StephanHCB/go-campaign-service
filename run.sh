#! /bin/bash

go build main.go
./main --config-path=. --secrets-path=.

