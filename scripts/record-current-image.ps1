$stateDir = "C:\sports-data-platform-deploy"
$stateFile = Join-Path $stateDir "last-known-good-image.txt"

$containerId = docker compose -f compose.yaml ps -q app

if (-not $containerId) {
    Write-Error "Could not find the app container."
    exit 1
}

$image = docker inspect $containerId --format "{{.Config.Image}}"

if (-not $image) {
    Write-Error "Could not determine the app container image."
    exit 1
}

$imageTag = $image.Substring($image.LastIndexOf(":") + 1)

New-Item -ItemType Directory -Force -Path $stateDir | Out-Null

Set-Content -Path $stateFile -Value $imageTag

Write-Host "Recorded last-known-good image tag: $imageTag"