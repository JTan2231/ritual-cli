#!/bin/bash

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Missing Go installation"
    exit 1
fi

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 ritual_account_name"
    exit 1
fi

echo "Building the CLI..."
go build -o ritual ./cli

echo "Building the cron job..."
go build -o ritual_cron ./cron

echo "Moving binaries to /usr/local/bin..."
mv ritual /usr/local/bin/ritual
mv ritual_cron /usr/local/bin/ritual_cron

echo "Calculating local GMT offset..."
TARGET_HOUR=13

TIMEZONE_OFFSET=$(date +%z)
OFFSET_HOURS=$((${TIMEZONE_OFFSET:0:3}))

LOCAL_HOUR=$((TARGET_HOUR + OFFSET_HOURS))

if [ $LOCAL_HOUR -lt 0 ]; then
    LOCAL_HOUR=$((24 + LOCAL_HOUR))
elif [ $LOCAL_HOUR -ge 24 ]; then
    LOCAL_HOUR=$((LOCAL_HOUR - 24))
fi

CRON_JOB="0 $LOCAL_HOUR * * 0 /usr/local/bin/ritual_cron"
( crontab -l | grep -Fv "$CRON_JOB" ; echo "$CRON_JOB" ) | crontab -
echo "Cron job scheduled for $LOCAL_HOUR:00 local time."
echo
echo 'Usage: ritual "your entry"'
echo
echo "Generate a CLI secret key at https://ritual-api-production.up.railway.app/get-config-token?email=$1"
echo "Make sure to set your RITUAL_CLI_SECRET environment variable with your generated secret key."
