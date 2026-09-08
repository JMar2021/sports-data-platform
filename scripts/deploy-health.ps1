$containerId = docker compose -f compose.yaml ps -q app

if (-not $containerId) {
    Write-Error "Could not find the app container."
    exit 1
}

$healthy = $false

for ($i = 0; $i -lt 12; $i++) {
    $status = docker inspect $containerId --format "{{.State.Health.Status}}"

    Write-Host "Application health: $status"

    if ($status -eq "healthy") {
        $healthy = $true
        break
    }

    if ($status -eq "unhealthy") {
        break
    }

    Start-Sleep -Seconds 5
}

if (-not $healthy) {
    Write-Error "Application did not become healthy."
    docker compose -f compose.yaml ps
    exit 1
}