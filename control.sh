#!/bin/bash

version="1.0.5"

if [ "$1" == "build" ]; then
    echo "Building the Docker image..."
    docker build --platform linux/amd64 -t hi:$version-linux_amd64 .
elif [ "$1" == "push" ]; then
    echo "Pushing the Docker image to the private repository..."
    docker tag hi:$version-linux_amd64 119.3.172.171:5000/hi:$version-linux_amd64
    docker push 119.3.172.171:5000/hi:$version-linux_amd64
elif [ "$1" == "status" ]; then
    process_count=$(pgrep -f "app" | wc -l)
    if [ "$process_count" -eq 0 ]; then
        echo "Process 'app' not found."
        exit 1
    else
        echo "Process 'app' is running."
        exit 0
    fi
 elif [ "$1" == "stop" ]; then
    process_count=$(pgrep -f "app" | wc -l)
    if [ "$process_count" -eq 0 ]; then
        echo "Process 'app' not found."
        exit 1
    else
        echo "Stopping process 'app'..."
        pkill -15 -f "app"
        exit 0
    fi
else
    echo "Usage: $0 [build|push|status|stop]"
    exit 1
fi
