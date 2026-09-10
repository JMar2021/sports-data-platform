$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot\..

Write-Host "Starting daily sports data pipeline..."

Write-Host "Running batch..."
docker compose run --rm --no-deps batch

if ($LASTEXITCODE -ne 0) {
    Write-Error "Batch failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

Write-Host "Batch completed successfully."

Write-Host "Triggering GitHub Pages publication..."

$workflowUrl = gh workflow run publish-data.yaml --ref main 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to trigger GitHub Pages publication with exit code $LASTEXITCODE"
    Write-Error $workflowUrl
    exit $LASTEXITCODE
}

Write-Host "GitHub Pages publication triggered."
Write-Host "Workflow: $workflowUrl"

$runId = ($workflowUrl -split "/")[-1]

if (-not $runId -or $runId -notmatch "^\d+$") {
    Write-Error "Could not determine GitHub Actions run ID from: $workflowUrl"
    exit 1
}

Write-Host "Waiting for GitHub Pages publication to complete..."

gh run watch $runId --compact --exit-status

if ($LASTEXITCODE -ne 0) {
    Write-Error "GitHub Pages publication failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

Write-Host "GitHub Pages publication completed successfully."
Write-Host "Daily sports data pipeline completed successfully."

exit 0