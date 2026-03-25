$body = '{"username":"admin","password":"admin123"}'
try {
    $r = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -Body $body -ContentType 'application/json'
    $r | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $_"
    $_.Exception.Response.StatusCode
}