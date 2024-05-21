#!/bin/bash

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Missing Go installation"
    exit 1
fi

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 ritual_account_name"
    exit 1
fi

if [ -z "$RITUAL_CLI_SECRET" ]; then
    echo "RTIUAL_CLI_SECRET is not defined. What is your Ritual secret key? Get one from https://ritual-api-production.up.railway.app/get-config-token?email=$1 if you don't have one/don't know."
    read -rp "Ritual secret key: " R
    export RITUAL_CLI_SECRET=$R
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

CRON_JOB="0 $LOCAL_HOUR * * 0 /usr/local/bin/ritual_cron $RITUAL_CLI_SECRET"
echo "Add \`$CRON_JOB\` to your crontab with \`crontab -e\` to receive weekly newsletters."
