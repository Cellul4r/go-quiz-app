#!/bin/bash

APP_ENV=$(grep -w "APP_ENV" .env | cut -d '=' -f2 | tr -d '[:space:]')
APP_DEBUG=$(grep -w "APP_DEBUG" .env | cut -d '=' -f2 | tr -d '[:space:]')

echo "APP_ENV: $APP_ENV"
echo "APP_DEBUG: $APP_DEBUG"

if [ "$APP_ENV" = "production" ] || [ "$APP_DEBUG" = "false" ]; then
  echo "Building production image..."
  docker compose up --build -d
else
  echo "Building development image..."
  docker compose up -d
fi