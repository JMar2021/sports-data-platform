$stateFile = "C:\sports-data-platform-deploy\last-known-good-image.txt"

if (-not (Test-Path $stateFile)) {
    Write-Error "Last-known-good image state file was not found."
    exit 1
}

$imageTag = (Get-Content $stateFile -Raw).Trim()

if (-not $imageTag) {
    Write-Error "Last-known-good image tag is empty."
    exit 1
}

Write-Host "Rolling back to image tag: $imageTag"

$env:IMAGE_TAG = $imageTag

docker compose -f compose.yaml pull

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to pull rollback image."
    exit 1
}

docker compose -f compose.yaml up -d

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to start rollback deployment."
    exit 1
}

Write-Host "Rollback deployment started."