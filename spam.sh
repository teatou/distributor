#!/bin/bash

count=10

for ((i=1; i<=$count; i++))
do
    curl http://localhost:8000
    echo "request made"
done