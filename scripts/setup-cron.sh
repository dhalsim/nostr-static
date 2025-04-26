#!/bin/bash

# Get the absolute path of the project directory
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_FILE="$PROJECT_DIR/cron.log"

# Create the cron job command with git operations and logging
CRON_CMD="cd $PROJECT_DIR && echo \"\$(date): Starting cron job\" >> $LOG_FILE && ./nostr-static && if [ -n \"\$(git status --porcelain)\" ]; then git add . && git commit -m 'Update Nostr static site' && git push origin main && echo \"\$(date): Changes detected and pushed\" >> $LOG_FILE; else echo \"\$(date): No changes detected\" >> $LOG_FILE; fi"

# Add the cron job for testing (runs every minute)
#(crontab -l 2>/dev/null; echo "* * * * * $CRON_CMD") | crontab -
#echo "Cron job has been set up for testing (runs every minute)."

# Add the cron job for daily runs at 00:00
#(crontab -l 2>/dev/null; echo "0 0 * * * $CRON_CMD") | crontab -
#echo "Cron job has been set up for daily runs at 00:00."

# Add the cron job for hourly runs at 00:00
(crontab -l 2>/dev/null; echo "0 * * * * $CRON_CMD") | crontab -
echo "Cron job has been set up for hourly runs at 00:00."

echo "Logs will be written to: $LOG_FILE"
