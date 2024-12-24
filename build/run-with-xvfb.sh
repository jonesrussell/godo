#!/bin/bash

# Start Xvfb
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0

# Wait for Xvfb to start
sleep 1

# Run the command passed as arguments
exec "$@" 