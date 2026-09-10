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

Write-Host "Running exporter..."
docker compose run --rm --no-deps export

if ($LASTEXITCODE -ne 0) {
    Write-Error "Export failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

Write-Host "Export completed successfully."
Write-Host "Daily sports data pipeline completed successfully."

exit 0