#!/bin/sh

docker compose up -d

echo start...

sleep 5

read -p "Press enter to continue"

newman run expenses.postman_collection.json

echo down

docker compose down