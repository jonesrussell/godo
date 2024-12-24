#!/bin/bash

# Start Xvfb with specific screen configuration
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
XVFB_PID=$!

# Wait for Xvfb to start
sleep 1

# Set up environment
export DISPLAY=:99.0

# Start dbus session
eval $(dbus-launch --sh-syntax)

# Run the command passed as arguments
"$@"
EXIT_CODE=$?

# Cleanup
kill $XVFB_PID
kill $DBUS_SESSION_BUS_PID

exit $EXIT_CODE 