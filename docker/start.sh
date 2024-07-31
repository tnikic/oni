#!/bin/sh

# Start the backend process
cd backend
./oni &

# Start the fronend process
cd ../frontend
npm run start
