#!/bin/bash

# This script is used to compile the lambda function in go

echo "Compiling lambda function"

rm bootstrap
rm bootstrap.zip

GOOS=linux go build -o bootstrap $1

zip bootstrap.zip bootstrap
