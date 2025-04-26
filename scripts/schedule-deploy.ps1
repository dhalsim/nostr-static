# PowerShell script for scheduled deployment
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Change to the project directory
Set-Location $PSScriptRoot/..

# Set up logging
$LOG_FILE = Join-Path $PSScriptRoot/.. "cron.log"

# Run the nostr-static command to trigger the deploy action
Add-Content -Path $LOG_FILE -Value "$(Get-Date): Starting deployment"
./nostr-static generate

# Check if there are any changes
$changes = git status --porcelain
if ($changes) {
    git add .
    git commit -m "Update Nostr static site"
    git push origin main
    Add-Content -Path $LOG_FILE -Value "$(Get-Date): Changes detected and pushed successfully"
} else {
    Add-Content -Path $LOG_FILE -Value "$(Get-Date): No changes detected, skipping commit and push"
}
