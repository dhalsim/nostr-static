# PowerShell script for scheduled deployment
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Change to the project directory
Set-Location $PSScriptRoot/..

# Run the nostr-static command with deploy trigger
./nostr-static --trigger-action=deploy 