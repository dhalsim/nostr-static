#!/bin/bash

# Get the absolute path of the project directory
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Create the cron job command
CRON_CMD="cd $PROJECT_DIR && ./nostr-static --trigger-action=deploy"

# Add the cron job for daily run at 00:00
(crontab -l 2>/dev/null; echo "0 0 * * * $CRON_CMD") | crontab -

# Uncomment the line below for hourly runs
# (crontab -l 2>/dev/null; echo "0 * * * * $CRON_CMD") | crontab -

echo "Cron job has been set up for daily deployment at midnight."
echo "To run hourly, uncomment the hourly line in $PROJECT_DIR/scripts/setup-cron.sh" 
