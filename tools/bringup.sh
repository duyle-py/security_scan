#!/bin/bash
set -e

cd src
go build -o main
./main