# MDM Backend API Test Suite - 2026-03-24
$BASE_URL = "http://localhost:8085"
$TOKEN = $null
$TEST_START = Get-Date

# Helper function
function Test-API {
    param([string]$Method, [string]$Endpoint, [string]$Body = $null, [string]$Desc = "", [bool]$Auth = $true)
    $url = "$BASE_URL$Endpoint"
    $headers = @{"Content-Type" = "application/json"}
    if ($Auth -and $TOKEN) { $headers["Authorization"] = "Bearer $TOKEN" }
    
    $sw = [Diagnostics.Stopwatch]::StartNew()
    try {
        $params = @{Uri=$url; Method=$Method; Headers=$headers}
        if ($Body) { $params["Body"] = $Body }
        $resp = Invoke-RestMethod @params -TimeoutSec 30
        $sw.Stop()
        $status = "PASS"
        $result = $resp | ConvertTo-Json -Depth 3 -Compress
        if ($result.Length -gt 200) { $result = $result.Substring(0, 200) + "..." }
    } catch {
        $sw.Stop()
        $status = "FAIL"
        $result = $_.Exception.Message
        if ($result.Length -gt 200) { $result = $result.Substring(0, 200) + "..." }
    }
    
    [PSCustomObject]@{
        Status = $status
        Method = $Method
        Endpoint = $Endpoint
        Time = $sw.ElapsedMilliseconds
        Desc = $Desc
        Result = $result
    }
}

# ========== PHASE 1: Login ==========
Write-Host "`n=== [1/6] AUTH - Login ===" -ForegroundColor Cyan
$r = Test-API "POST" "/api/v1/auth/login" '{"username":"admin","password":"admin123"}' "Login with admin" $false
$r | Format-Table -AutoSize
if ($r.Status -eq "PASS" -and $r.Result -match "token") {
    $TOKEN = ($r.Result -split '"token":"')[1] -split '"')[0]
    Write-Host "Token obtained: $($TOKEN.Substring(0, [Math]::Min(30, $TOKEN.Length)))..." -ForegroundColor Green
}

# ========== PHASE 2: Device Management ==========
Write-Host "`n=== [2/6] Device Management ===" -ForegroundColor Cyan
$deviceTests = @(
    @{Method="GET"; Endpoint="/api/v1/devices"; Desc="List devices"},
    @{Method="GET"; Endpoint="/api/v1/devices/test-device-001"; Desc="Get device by ID"},
    @{Method="PUT"; Endpoint="/api/v1/devices/test-device-001/status"; Body='{"online":true}'; Desc="Update device status"},
    @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/profile"; Desc="Get device profile"},
    @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/desired-state"; Desc="Get desired state"},
    @{Method="PUT"; Endpoint="/api/v1/devices/test-device-001/desired-state"; Body='{"power":"on","mode":"normal"}'; Desc="Set desired state"},
    @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/commands"; Desc="Get command history"},
    @{Method="POST"; Endpoint="/api/v1/devices/test-device-001/commands"; Body='{"command":"ping","params":{}}'; Desc="Send device command"},
    @{Method="GET"; Endpoint="/api/v1/device-shadow/test-device-001"; Desc="Get device shadow"},
    @{Method="PUT"; Endpoint="/api/v1/device-shadow/test-device-001"; Body='{"desired":{"power":"on"}}'; Desc="Update device shadow"},
    @{Method="GET"; Endpoint="/api/v1/devices/export"; Desc="Export devices"}
)
$results = @()
foreach ($t in $deviceTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 3: Member Management ==========
Write-Host "`n=== [3/6] Member Management ===" -ForegroundColor Cyan
$memberTests = @(
    @{Method="GET"; Endpoint="/api/v1/members"; Desc="List members"},
    @{Method="POST"; Endpoint="/api/v1/members"; Body='{"name":"Test Member","email":"test@example.com"}'; Desc="Create member"},
    @{Method="GET"; Endpoint="/api/v1/member/cards"; Desc="List member cards"},
    @{Method="GET"; Endpoint="/api/v1/member/coupons"; Desc="List coupons"},
    @{Method="GET"; Endpoint="/api/v1/member/stores"; Desc="List stores"},
    @{Method="GET"; Endpoint="/api/v1/member/tags"; Desc="List tags"},
    @{Method="GET"; Endpoint="/api/v1/member/levels"; Desc="List levels"},
    @{Method="GET"; Endpoint="/api/v1/member/orders"; Desc="List orders"},
    @{Method="GET"; Endpoint="/api/v1/member/points/rules"; Desc="List points rules"},
    @{Method="GET"; Endpoint="/api/v1/members/points"; Desc="List all member points"},
    @{Method="GET"; Endpoint="/api/v1/coupons"; Desc="List coupons v2"},
    @{Method="GET"; Endpoint="/api/v1/promotions"; Desc="List promotions"}
)
$results = @()
foreach ($t in $memberTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 4: Tenant Management ==========
Write-Host "`n=== [4/6] Tenant Management ===" -ForegroundColor Cyan
$tenantTests = @(
    @{Method="GET"; Endpoint="/api/v1/admin/tenants"; Desc="List tenants"},
    @{Method="GET"; Endpoint="/api/v1/admin/plans"; Desc="List subscription plans"},
    @{Method="GET"; Endpoint="/api/v1/tenant-approvals"; Desc="List tenant approvals"},
    @{Method="GET"; Endpoint="/api/v1/subscriptions/plans"; Desc="List subscription plans v2"},
    @{Method="GET"; Endpoint="/api/v1/subscriptions/current"; Desc="Get current subscription"}
)
$results = @()
foreach ($t in $tenantTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 5: Organization ==========
Write-Host "`n=== [5/6] Organization ===" -ForegroundColor Cyan
$orgTests = @(
    @{Method="GET"; Endpoint="/api/v1/org/companies"; Desc="List companies"},
    @{Method="GET"; Endpoint="/api/v1/org/departments"; Desc="List departments"},
    @{Method="GET"; Endpoint="/api/v1/org/departments/tree"; Desc="Department tree"},
    @{Method="GET"; Endpoint="/api/v1/org/positions"; Desc="List positions"},
    @{Method="GET"; Endpoint="/api/v1/org/employees"; Desc="List employees"},
    @{Method="GET"; Endpoint="/api/v1/org/standard-positions"; Desc="Standard positions"},
    @{Method="GET"; Endpoint="/api/v1/position-templates"; Desc="Position templates"}
)
$results = @()
foreach ($t in $orgTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 6: OTA ==========
Write-Host "`n=== [6/6] OTA ===" -ForegroundColor Cyan
$otaTests = @(
    @{Method="GET"; Endpoint="/api/v1/ota/packages"; Desc="List OTA packages"},
    @{Method="GET"; Endpoint="/api/v1/ota/deployments"; Desc="List OTA deployments"},
    @{Method="GET"; Endpoint="/api/v1/ota/devices/test-device-001/check"; Desc="Check OTA for device"}
)
$results = @()
foreach ($t in $otaTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 7: Notifications ==========
Write-Host "`n=== [7/6] Notifications ===" -ForegroundColor Cyan
$notifTests = @(
    @{Method="GET"; Endpoint="/api/v1/notifications"; Desc="List notifications"},
    @{Method="GET"; Endpoint="/api/v1/notification-channels"; Desc="List notification channels"},
    @{Method="GET"; Endpoint="/api/v1/notification-templates"; Desc="List notification templates"},
    @{Method="GET"; Endpoint="/api/v1/announcements"; Desc="List announcements"}
)
$results = @()
foreach ($t in $notifTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 8: Dashboard & Analytics ==========
Write-Host "`n=== [8/6] Dashboard & Analytics ===" -ForegroundColor Cyan
$dashTests = @(
    @{Method="GET"; Endpoint="/api/v1/dashboard/stats"; Desc="Dashboard stats"},
    @{Method="GET"; Endpoint="/api/v1/dashboard/stats/simple"; Desc="Dashboard stats simple"},
    @{Method="GET"; Endpoint="/api/v1/dashboard/activity-summary"; Desc="Activity summary"},
    @{Method="GET"; Endpoint="/api/v1/activity-logs"; Desc="Activity logs"},
    @{Method="GET"; Endpoint="/api/v1/login-logs"; Desc="Login logs"},
    @{Method="GET"; Endpoint="/api/v1/logs/operations"; Desc="Operation logs"},
    @{Method="GET"; Endpoint="/api/v1/logs/login"; Desc="Login logs v2"},
    @{Method="GET"; Endpoint="/api/v1/analytics/advanced"; Desc="Advanced analytics"}
)
$results = @()
foreach ($t in $dashTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 9: Health & Alerts ==========
Write-Host "`n=== [9/6] Health & Alerts ===" -ForegroundColor Cyan
$healthTests = @(
    @{Method="GET"; Endpoint="/api/v1/alerts"; Desc="List alerts"},
    @{Method="GET"; Endpoint="/api/v1/alerts/rules"; Desc="Alert rules"},
    @{Method="GET"; Endpoint="/api/v1/alerts/settings"; Desc="Alert settings"},
    @{Method="GET"; Endpoint="/api/v1/alerts/history"; Desc="Alert history"},
    @{Method="GET"; Endpoint="/api/v1/geofence/rules"; Desc="Geofence rules"},
    @{Method="GET"; Endpoint="/api/v1/geofence/alerts"; Desc="Geofence alerts"},
    @{Method="GET"; Endpoint="/api/v1/notification-channels"; Desc="Notification channels"}
)
$results = @()
foreach ($t in $healthTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 10: AI & Knowledge ==========
Write-Host "`n=== [10/6] AI & Knowledge ===" -ForegroundColor Cyan
$aiTests = @(
    @{Method="GET"; Endpoint="/api/v1/knowledge"; Desc="List knowledge"},
    @{Method="GET"; Endpoint="/api/v1/ai/models"; Desc="List AI models"},
    @{Method="GET"; Endpoint="/api/v1/ai/inference"; Desc="List inferences"},
    @{Method="GET"; Endpoint="/api/v1/ai/training"; Desc="List training"}
)
$results = @()
foreach ($t in $aiTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 11: Role/Permission ==========
Write-Host "`n=== [11/6] Roles & Permissions ===" -ForegroundColor Cyan
$roleTests = @(
    @{Method="GET"; Endpoint="/api/v1/roles"; Desc="List roles"},
    @{Method="GET"; Endpoint="/api/v1/permission-groups"; Desc="Permission groups"},
    @{Method="GET"; Endpoint="/api/v1/menus"; Desc="Menus"},
    @{Method="GET"; Endpoint="/api/v1/menus/tree"; Desc="Menu tree"},
    @{Method="GET"; Endpoint="/api/v1/api-permissions"; Desc="API permissions"},
    @{Method="GET"; Endpoint="/api/v1/dicts/health_type"; Desc="Dict: health_type"},
    @{Method="GET"; Endpoint="/api/v1/dicts/alert_level"; Desc="Dict: alert_level"}
)
$results = @()
foreach ($t in $roleTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 12: Performance ==========
Write-Host "`n=== [12/6] Performance & Monitoring ===" -ForegroundColor Cyan
$perfTests = @(
    @{Method="GET"; Endpoint="/api/v1/performance/health"; Desc="Performance health"},
    @{Method="GET"; Endpoint="/api/v1/performance/cache/stats"; Desc="Cache stats"},
    @{Method="GET"; Endpoint="/api/v1/performance/metrics"; Desc="Performance metrics"},
    @{Method="GET"; Endpoint="/api/v1/performance/db/stats"; Desc="DB stats"}
)
$results = @()
foreach ($t in $perfTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 13: Pet Health & Tracking ==========
Write-Host "`n=== [13/6] Pet Health & Tracking ===" -ForegroundColor Cyan
$petTests = @(
    @{Method="GET"; Endpoint="/api/v1/health/test-pet/early-warning"; Desc="Early warnings"},
    @{Method="GET"; Endpoint="/api/v1/health/test-pet/exercise"; Desc="Exercise records"},
    @{Method="GET"; Endpoint="/api/v1/health/test-pet/sleep"; Desc="Sleep records"},
    @{Method="GET"; Endpoint="/api/v1/health/test-pet/report"; Desc="Health report"}
)
$results = @()
foreach ($t in $petTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 14: Digital Twin ==========
Write-Host "`n=== [14/6] Digital Twin ===" -ForegroundColor Cyan
$twinTests = @(
    @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/vitals"; Desc="Get vitals"},
    @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/timeline"; Desc="Get timeline"},
    @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/alerts"; Desc="Get alerts"}
)
$results = @()
foreach ($t in $twinTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 15: Pet Finder ==========
Write-Host "`n=== [15/6] Pet Finder ===" -ForegroundColor Cyan
$finderTests = @(
    @{Method="GET"; Endpoint="/api/v1/pet-finder/reports"; Desc="List reports"},
    @{Method="GET"; Endpoint="/api/v1/pet-finder/nearby"; Desc="Nearby reports"}
)
$results = @()
foreach ($t in $finderTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 16: Pet Social ==========
Write-Host "`n=== [16/6] Pet Social ===" -ForegroundColor Cyan
$socialTests = @(
    @{Method="GET"; Endpoint="/api/v1/pet-social/feed"; Desc="Social feed"}
)
$results = @()
foreach ($t in $socialTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 17: Insurance ==========
Write-Host "`n=== [17/6] Insurance ===" -ForegroundColor Cyan
$insTests = @(
    @{Method="GET"; Endpoint="/api/v1/insurance/products"; Desc="Insurance products"},
    @{Method="GET"; Endpoint="/api/v1/insurance/claims"; Desc="Insurance claims"}
)
$results = @()
foreach ($t in $insTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 18: Pet Shop ==========
Write-Host "`n=== [18/6] Pet Shop ===" -ForegroundColor Cyan
$shopTests = @(
    @{Method="GET"; Endpoint="/api/v1/pet-shop/products"; Desc="Shop products"},
    @{Method="GET"; Endpoint="/api/v1/pet-shop/categories"; Desc="Shop categories"},
    @{Method="GET"; Endpoint="/api/v1/pet-shop/orders"; Desc="Shop orders"}
)
$results = @()
foreach ($t in $shopTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 19: Simulation ==========
Write-Host "`n=== [19/6] Simulation ===" -ForegroundColor Cyan
$simTests = @(
    @{Method="GET"; Endpoint="/api/v1/simulation/virtual-pets"; Desc="Virtual pets"},
    @{Method="GET"; Endpoint="/api/v1/simulation/environments"; Desc="Environments"},
    @{Method="GET"; Endpoint="/api/v1/simulation/metrics"; Desc="Simulation metrics"}
)
$results = @()
foreach ($t in $simTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

# ========== PHASE 20: More Modules ==========
Write-Host "`n=== [20/6] Additional Modules ===" -ForegroundColor Cyan
$moreTests = @(
    @{Method="GET"; Endpoint="/api/v1/emotion/logs"; Desc="Emotion logs"},
    @{Method="GET"; Endpoint="/api/v1/emotion/reports"; Desc="Emotion reports"},
    @{Method="GET"; Endpoint="/api/v1/emotion/actions"; Desc="Emotion actions"},
    @{Method="GET"; Endpoint="/api/v1/emotion/family-map"; Desc="Family emotion map"},
    @{Method="GET"; Endpoint="/api/v1/miniclaw/firmwares"; Desc="MiniClaw firmwares"},
    @{Method="GET"; Endpoint="/api/v1/mesh/networks"; Desc="Mesh networks"},
    @{Method="GET"; Endpoint="/api/v1/mesh/devices"; Desc="Mesh devices"},
    @{Method="GET"; Endpoint="/api/v1/compliance/policies"; Desc="Compliance policies"},
    @{Method="GET"; Endpoint="/api/v1/developer/apps"; Desc="Developer apps"},
    @{Method="GET"; Endpoint="/api/v1/store/apps"; Desc="Store apps"},
    @{Method="GET"; Endpoint="/api/v1/i18n/translations"; Desc="i18n translations"},
    @{Method="GET"; Endpoint="/api/v1/i18n/locales"; Desc="i18n locales"},
    @{Method="GET"; Endpoint="/api/v1/regions"; Desc="Regions"},
    @{Method="GET"; Endpoint="/api/v1/timezones"; Desc="Timezones"},
    @{Method="GET"; Endpoint="/api/v1/daas/contracts"; Desc="DaaS contracts"},
    @{Method="GET"; Endpoint="/api/v1/daas/devices"; Desc="DaaS devices"},
    @{Method="GET"; Endpoint="/api/v1/research/datasets"; Desc="Research datasets"},
    @{Method="GET"; Endpoint="/api/v1/research/experiments"; Desc="Research experiments"},
    @{Method="GET"; Endpoint="/api/v1/ai/fairness/tests"; Desc="AI fairness tests"},
    @{Method="GET"; Endpoint="/api/v1/ai/audit/logs"; Desc="AI audit logs"}
)
$results = @()
foreach ($t in $moreTests) {
    $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
    $results += $r
}
$results | Format-Table -AutoSize

Write-Host "`n=== Test Complete ===" -ForegroundColor Green
