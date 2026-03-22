$bytes = [System.IO.File]::ReadAllBytes('C:\Users\YKing\.openclaw\workspace\mdm-project\backend\mdm-backend-new.exe')[0..3]
foreach ($b in $bytes) {
    Write-Host ("{0:X2}" -f $b)
}
