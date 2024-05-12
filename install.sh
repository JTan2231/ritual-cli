#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 ritual_account_name"
    exit 1
fi

echo "building the CLI"
go build -o ritual ./cli

echo "building the cron job"
go build -o ritual_cron ./cron

echo "moving binaries to /usr/local/bin"
mv ritual /usr/local/bin/ritual
mv ritual_cron /usr/local/bin/ritual_cron

echo "calculating local GMT offset"
# Define the GMT target time (13:00 in 24-hour format)
TARGET_HOUR=13
TARGET_MINUTE=0

# Get the current timezone offset in hours from GMT
TIMEZONE_OFFSET=$(date +%z)  # Outputs in format +HHMM or -HHMM
OFFSET_HOURS=$((${TIMEZONE_OFFSET:0:3}))  # Extract hours with sign
OFFSET_MINUTES=$((${TIMEZONE_OFFSET:0:1}${TIMEZONE_OFFSET:3}))  # Extract minutes with correct sign

# Calculate the local time equivalent to 1 PM GMT
LOCAL_HOUR=$((TARGET_HOUR + OFFSET_HOURS))
LOCAL_MINUTE=$((TARGET_MINUTE + OFFSET_MINUTES))

# Adjust for minute overflow/underflow
if [ $LOCAL_MINUTE -lt 0 ]; then
    LOCAL_MINUTE=$((60 + LOCAL_MINUTE))
    LOCAL_HOUR=$((LOCAL_HOUR - 1))
elif [ $LOCAL_MINUTE -ge 60 ]; then
    LOCAL_MINUTE=$((LOCAL_MINUTE - 60))
    LOCAL_HOUR=$((LOCAL_HOUR + 1))
fi

# Adjust hours for overflow/underflow
if [ $LOCAL_HOUR -lt 0 ]; then
    LOCAL_HOUR=$((24 + LOCAL_HOUR))
elif [ $LOCAL_HOUR -ge 24 ]; then
    LOCAL_HOUR=$((LOCAL_HOUR - 24))
fi

CRON_JOB="0 $LOCAL_HOUR * * 0 /usr/local/bin/ritual_cron"
( crontab -l | grep -Fv "$CRON_JOB" ; echo "$CRON_JOB" ) | crontab -
echo "Cron job scheduled for $LOCAL_HOUR:00 local time."

echo 'Usage: ritual "your entry"'
echo "Generate a CLI token at https://ritual-api-production.up.railway.app/get-config-token?email=$1"
