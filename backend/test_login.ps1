$passwords = @("admin", "password", "123456", "admin123", "password123", "Admin123!", "MDM@2024", "mdm123456")

foreach ($pwd in $passwords) {
    $body = @{
        username = "admin"
        password = $pwd
    } | ConvertTo-Json
    
    try {
        $resp = Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/login' -Method POST -ContentType 'application/json' -Body $body -TimeoutSec 5 -ErrorAction Stop
        Write-Host "SUCCESS! Password: $pwd"
        $resp | ConvertTo-Json -Depth 3
        break
    } catch {
        $err = $_.Exception.Response
        if ($err) {
            $status = $err.StatusCode
            Write-Host "Password '$pwd' -> Status: $status"
        }
    }
}
