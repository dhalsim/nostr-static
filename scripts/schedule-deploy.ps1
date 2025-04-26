# PowerShell script for scheduled deployment
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Change to the project directory
Set-Location $PSScriptRoot/..

# Run the nostr-static command to trigger the deploy action
./nostr-static generate

# Check if there are any changes
$changes = git status --porcelain
if ($changes) {
    git add .
    git commit -m "Update Nostr static site"
    git push origin main
    Write-Host "Changes detected and pushed successfully."
} else {
    Write-Host "No changes detected, skipping commit and push."
}
