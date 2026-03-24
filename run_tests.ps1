# MDM Backend API Test Suite - 2026-03-24
$BASE_URL = "http://localhost:8085"
$TOKEN = $null
$TEST_START = Get-Date
$ALL_RESULTS = @()

function Test-API {
    param([string]$Method, [string]$Endpoint, [string]$Body = $null, [string]$Desc = "", [bool]$Auth = $true)
    $url = "$BASE_URL$Endpoint"
    $headers = @{"Content-Type" = "application/json"}
    if ($Auth -and $TOKEN) { $headers["Authorization"] = "Bearer $TOKEN" }
    
    $sw = [Diagnostics.Stopwatch]::StartNew()
    $status = "FAIL"
    $result = ""
    try {
        $params = @{Uri=$url; Method=$Method; Headers=$headers}
        if ($Body) { $params["Body"] = $Body }
        $resp = Invoke-RestMethod @params -TimeoutSec 30
        $sw.Stop()
        $status = "PASS"
        $result = $resp | ConvertTo-Json -Depth 3 -Compress
        if ($result.Length -gt 150) { $result = $result.Substring(0, 150) + "..." }
    } catch {
        $sw.Stop()
        $status = "FAIL"
        $result = $_.Exception.Message
        if ($result.Length -gt 150) { $result = $result.Substring(0, 150) + "..." }
    }
    
    [PSCustomObject]@{
        Status = $status
        Method = $Method
        Endpoint = $Endpoint
        Time_ms = $sw.ElapsedMilliseconds
        Desc = $Desc
        Result = $result
    }
}

# ========== Login ==========
Write-Host "`n=== [AUTH] Login ===" -ForegroundColor Cyan
$loginBody = @{"username"="admin";"password"="admin123"} | ConvertTo-Json
try {
    $r = Invoke-RestMethod -Uri "$BASE_URL/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody -TimeoutSec 10
    $TOKEN = $r.data.token
    Write-Host "PASS - Login successful, token obtained" -ForegroundColor Green
    $ALL_RESULTS += [PSCustomObject]@{Status="PASS";Method="POST";Endpoint="/api/v1/auth/login";Time_ms=0;Desc="Login with admin credentials";Result="Token obtained"}
} catch {
    Write-Host "FAIL - Login failed: $($_.Exception.Message)" -ForegroundColor Red
    $ALL_RESULTS += [PSCustomObject]@{Status="FAIL";Method="POST";Endpoint="/api/v1/auth/login";Time_ms=0;Desc="Login";Result=$_.Exception.Message}
    $TOKEN = $null
}

# Run all test phases
$phases = @(
    @{Name="Device Management"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/devices"; Desc="List devices"},
        @{Method="GET"; Endpoint="/api/v1/devices/test-device-001"; Desc="Get device by ID"},
        @{Method="PUT"; Endpoint="/api/v1/devices/test-device-001/status"; Body='{"online":true}'; Desc="Update device status"},
        @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/profile"; Desc="Get device profile"},
        @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/desired-state"; Desc="Get desired state"},
        @{Method="PUT"; Endpoint="/api/v1/devices/test-device-001/desired-state"; Body='{"power":"on"}'; Desc="Set desired state"},
        @{Method="GET"; Endpoint="/api/v1/devices/test-device-001/commands"; Desc="Get command history"},
        @{Method="POST"; Endpoint="/api/v1/devices/test-device-001/commands"; Body='{"command":"ping","params":{}}'; Desc="Send device command"},
        @{Method="GET"; Endpoint="/api/v1/device-shadow/test-device-001"; Desc="Get device shadow"},
        @{Method="PUT"; Endpoint="/api/v1/device-shadow/test-device-001"; Body='{"desired":{"power":"on"}}'; Desc="Update device shadow"},
        @{Method="GET"; Endpoint="/api/v1/devices/export"; Desc="Export devices"}
    )},
    @{Name="Member Management"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/members"; Desc="List members"},
        @{Method="GET"; Endpoint="/api/v1/member/cards"; Desc="List member cards"},
        @{Method="GET"; Endpoint="/api/v1/member/coupons"; Desc="List coupons"},
        @{Method="GET"; Endpoint="/api/v1/member/stores"; Desc="List stores"},
        @{Method="GET"; Endpoint="/api/v1/member/tags"; Desc="List tags"},
        @{Method="GET"; Endpoint="/api/v1/member/levels"; Desc="List levels"},
        @{Method="GET"; Endpoint="/api/v1/member/orders"; Desc="List orders"},
        @{Method="GET"; Endpoint="/api/v1/member/points/rules"; Desc="List points rules"},
        @{Method="GET"; Endpoint="/api/v1/members/points"; Desc="List all member points"},
        @{Method="GET"; Endpoint="/api/v1/coupons"; Desc="List coupons v2"},
        @{Method="GET"; Endpoint="/api/v1/promotions"; Desc="List promotions"},
        @{Method="POST"; Endpoint="/api/v1/members"; Body='{"name":"Test Member","email":"test@example.com"}'; Desc="Create member"}
    )},
    @{Name="Tenant Management"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/admin/tenants"; Desc="List tenants"},
        @{Method="GET"; Endpoint="/api/v1/admin/plans"; Desc="List subscription plans"},
        @{Method="GET"; Endpoint="/api/v1/tenant-approvals"; Desc="List tenant approvals"},
        @{Method="GET"; Endpoint="/api/v1/subscriptions/plans"; Desc="List subscription plans v2"},
        @{Method="GET"; Endpoint="/api/v1/subscriptions/current"; Desc="Get current subscription"}
    )},
    @{Name="Organization"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/org/companies"; Desc="List companies"},
        @{Method="GET"; Endpoint="/api/v1/org/departments"; Desc="List departments"},
        @{Method="GET"; Endpoint="/api/v1/org/departments/tree"; Desc="Department tree"},
        @{Method="GET"; Endpoint="/api/v1/org/positions"; Desc="List positions"},
        @{Method="GET"; Endpoint="/api/v1/org/employees"; Desc="List employees"},
        @{Method="GET"; Endpoint="/api/v1/org/standard-positions"; Desc="Standard positions"},
        @{Method="GET"; Endpoint="/api/v1/position-templates"; Desc="Position templates"}
    )},
    @{Name="OTA"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/ota/packages"; Desc="List OTA packages"},
        @{Method="GET"; Endpoint="/api/v1/ota/deployments"; Desc="List OTA deployments"},
        @{Method="GET"; Endpoint="/api/v1/ota/devices/test-device-001/check"; Desc="Check OTA for device"}
    )},
    @{Name="Notifications"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/notifications"; Desc="List notifications"},
        @{Method="GET"; Endpoint="/api/v1/notification-channels"; Desc="List notification channels"},
        @{Method="GET"; Endpoint="/api/v1/notification-templates"; Desc="List notification templates"},
        @{Method="GET"; Endpoint="/api/v1/announcements"; Desc="List announcements"}
    )},
    @{Name="Dashboard & Analytics"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/dashboard/stats"; Desc="Dashboard stats"},
        @{Method="GET"; Endpoint="/api/v1/dashboard/stats/simple"; Desc="Dashboard stats simple"},
        @{Method="GET"; Endpoint="/api/v1/dashboard/activity-summary"; Desc="Activity summary"},
        @{Method="GET"; Endpoint="/api/v1/activity-logs"; Desc="Activity logs"},
        @{Method="GET"; Endpoint="/api/v1/login-logs"; Desc="Login logs"},
        @{Method="GET"; Endpoint="/api/v1/logs/operations"; Desc="Operation logs"},
        @{Method="GET"; Endpoint="/api/v1/logs/login"; Desc="Login logs v2"},
        @{Method="GET"; Endpoint="/api/v1/analytics/advanced"; Desc="Advanced analytics"}
    )},
    @{Name="Health & Alerts"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/alerts"; Desc="List alerts"},
        @{Method="GET"; Endpoint="/api/v1/alerts/rules"; Desc="Alert rules"},
        @{Method="GET"; Endpoint="/api/v1/alerts/settings"; Desc="Alert settings"},
        @{Method="GET"; Endpoint="/api/v1/alerts/history"; Desc="Alert history"},
        @{Method="GET"; Endpoint="/api/v1/geofence/rules"; Desc="Geofence rules"},
        @{Method="GET"; Endpoint="/api/v1/geofence/alerts"; Desc="Geofence alerts"}
    )},
    @{Name="AI & Knowledge"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/knowledge"; Desc="List knowledge"},
        @{Method="GET"; Endpoint="/api/v1/ai/models"; Desc="List AI models"},
        @{Method="GET"; Endpoint="/api/v1/ai/inference"; Desc="List inferences"},
        @{Method="GET"; Endpoint="/api/v1/ai/training"; Desc="List training"}
    )},
    @{Name="Roles & Permissions"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/roles"; Desc="List roles"},
        @{Method="GET"; Endpoint="/api/v1/permission-groups"; Desc="Permission groups"},
        @{Method="GET"; Endpoint="/api/v1/menus"; Desc="Menus"},
        @{Method="GET"; Endpoint="/api/v1/menus/tree"; Desc="Menu tree"},
        @{Method="GET"; Endpoint="/api/v1/api-permissions"; Desc="API permissions"},
        @{Method="GET"; Endpoint="/api/v1/dicts/health_type"; Desc="Dict: health_type"},
        @{Method="GET"; Endpoint="/api/v1/dicts/alert_level"; Desc="Dict: alert_level"}
    )},
    @{Name="Performance & Monitoring"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/performance/health"; Desc="Performance health"},
        @{Method="GET"; Endpoint="/api/v1/performance/cache/stats"; Desc="Cache stats"},
        @{Method="GET"; Endpoint="/api/v1/performance/metrics"; Desc="Performance metrics"},
        @{Method="GET"; Endpoint="/api/v1/performance/db/stats"; Desc="DB stats"}
    )},
    @{Name="Pet Health & Tracking"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/health/test-pet/early-warning"; Desc="Early warnings"},
        @{Method="GET"; Endpoint="/api/v1/health/test-pet/exercise"; Desc="Exercise records"},
        @{Method="GET"; Endpoint="/api/v1/health/test-pet/sleep"; Desc="Sleep records"},
        @{Method="GET"; Endpoint="/api/v1/health/test-pet/report"; Desc="Health report"}
    )},
    @{Name="Digital Twin"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/vitals"; Desc="Get vitals"},
        @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/timeline"; Desc="Get timeline"},
        @{Method="GET"; Endpoint="/api/v1/digital-twin/test-pet/alerts"; Desc="Get alerts"}
    )},
    @{Name="Pet Finder"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/pet-finder/reports"; Desc="List reports"},
        @{Method="GET"; Endpoint="/api/v1/pet-finder/nearby"; Desc="Nearby reports"}
    )},
    @{Name="Pet Social"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/pet-social/feed"; Desc="Social feed"}
    )},
    @{Name="Insurance"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/insurance/products"; Desc="Insurance products"},
        @{Method="GET"; Endpoint="/api/v1/insurance/claims"; Desc="Insurance claims"}
    )},
    @{Name="Pet Shop"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/pet-shop/products"; Desc="Shop products"},
        @{Method="GET"; Endpoint="/api/v1/pet-shop/categories"; Desc="Shop categories"},
        @{Method="GET"; Endpoint="/api/v1/pet-shop/orders"; Desc="Shop orders"}
    )},
    @{Name="Simulation"; Tests=@(
        @{Method="GET"; Endpoint="/api/v1/simulation/virtual-pets"; Desc="Virtual pets"},
        @{Method="GET"; Endpoint="/api/v1/simulation/environments"; Desc="Environments"},
        @{Method="GET"; Endpoint="/api/v1/simulation/metrics"; Desc="Simulation metrics"}
    )},
    @{Name="Additional Modules"; Tests=@(
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
    )}
)

$phaseNum = 1
$totalPhases = $phases.Count
foreach ($phase in $phases) {
    Write-Host "`n=== [$phaseNum/$totalPhases] $($phase.Name) ===" -ForegroundColor Cyan
    $passCount = 0
    $failCount = 0
    foreach ($t in $phase.Tests) {
        $r = Test-API $t.Method $t.Endpoint $t.Body $t.Desc
        $ALL_RESULTS += $r
        if ($r.Status -eq "PASS") {
            $passCount++
            Write-Host "  [PASS] $($r.Method) $($r.Endpoint) - $($r.Time_ms)ms" -ForegroundColor Green
        } else {
            $failCount++
            Write-Host "  [FAIL] $($r.Method) $($r.Endpoint) - $($r.Time_ms)ms - $($r.Result)" -ForegroundColor Red
        }
    }
    Write-Host "  Phase result: $passCount passed, $failCount failed" -ForegroundColor Yellow
    $phaseNum++
}

# Summary
Write-Host "`n`n========================================" -ForegroundColor Magenta
Write-Host "         TEST SUMMARY" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta
$total = $ALL_RESULTS.Count
$passTotal = ($ALL_RESULTS | Where-Object { $_.Status -eq "PASS" }).Count
$failTotal = $total - $passTotal
Write-Host "Total Tests: $total" -ForegroundColor White
Write-Host "Passed:      $passTotal ($([Math]::Round($passTotal/$total*100))%)" -ForegroundColor Green
Write-Host "Failed:      $failTotal ($([Math]::Round($failTotal/$total*100))%)" -ForegroundColor Red

# Performance summary
$avgTime = ($ALL_RESULTS | Where-Object { $_.Time_ms -gt 0 } | Measure-Object -Property Time_ms -Average).Average
Write-Host "Avg Response Time: $([Math]::Round($avgTime))ms" -ForegroundColor White

# Export results to CSV
$csvPath = "C:\Users\YKing\.openclaw\workspace\mdm-project\test_results_2026-03-24.csv"
$ALL_RESULTS | Export-Csv -Path $csvPath -NoTypeInformation -Encoding UTF8
Write-Host "`nResults exported to: $csvPath" -ForegroundColor Cyan

# Save report
$reportPath = "C:\Users\YKing\.openclaw\workspace\TEST_REPORT_2026-03-24.md"
$TEST_END = Get-Date
$duration = $TEST_END - $TEST_START
$failed = $ALL_RESULTS | Where-Object { $_.Status -eq "FAIL" }

$report = @"
# MDM Backend Test Report - 2026-03-24

## Test Summary
- **Test Date**: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
- **Backend URL**: http://localhost:8085
- **Total Tests**: $total
- **Passed**: $passTotal ($([Math]::Round($passTotal/$total*100))%)
- **Failed**: $failTotal ($([Math]::Round($failTotal/$total*100))%)
- **Avg Response Time**: $([Math]::Round($avgTime))ms

## Detailed Results

"@

foreach ($phase in $phases) {
    $report += "`n### $($phase.Name)`n`n"
    $report += "| Status | Method | Endpoint | Time(ms) | Description |`n"
    $report += "|--------|--------|----------|----------|-------------|`n"
    foreach ($t in $phase.Tests) {
        $r = $ALL_RESULTS | Where-Object { $_.Endpoint -eq $t.Endpoint -and $_.Method -eq $t.Method }
        $emoji = if ($r.Status -eq "PASS") { ":white_check_mark:" } else { ":x:" }
        $report += "| $emoji | $($r.Method) | $($r.Endpoint) | $($r.Time_ms) | $($t.Desc) |`n"
    }
}

if ($failed.Count -gt 0) {
    $report += "`n## Failed Tests`n`n"
    foreach ($f in $failed) {
        $report += "- **$($f.Method) $($f.Endpoint)**: $($f.Result)`n"
    }
}

$report += "`n## Notes`n"
$report += "- Backend was started on port 8085 (port 8080 is occupied by OpenClaw)`n"
$report += "- All tests were run with admin credentials`n"
$report += "- Some 404 responses may indicate routes that exist in the API spec but not implemented yet`n"

$report | Out-File -FilePath $reportPath -Encoding UTF8
Write-Host "Report saved to: $reportPath" -ForegroundColor Cyan
