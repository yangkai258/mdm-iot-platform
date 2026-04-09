Get-ChildItem -Path "C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views" -Recurse -Filter "*.vue" | ForEach-Object {
    $file = $_.FullName
    $lines = Get-Content $file
    $n = $lines.Count
    $changed = $false
    for ($i = 0; $i -lt $n; $i++) {
        if ($lines[$i] -match '</a-table>' -and $i + 1 -lt $n -and $lines[$i + 1] -match '<template\s') {
            $j = $i + 1
            while ($j -lt $n -and $lines[$j] -notmatch '</a-table>') { $j++ }
            if ($j -lt $n -and $lines[$j] -match '</a-table>') {
                $inside = @()
                for ($k = $i + 1; $k -lt $j; $k++) { $inside += $lines[$k] }
                $newLines = @()
                for ($k = 0; $k -lt $n; $k++) {
                    if ($k -eq $i) {
                        $newLines += $lines[$k]
                        $newLines += $inside
                        $changed = $true
                    } elseif ($k -ge $i + 1 -and $k -le $j - 1) {
                        continue
                    } else {
                        $newLines += $lines[$k]
                    }
                }
                $lines = $newLines
                $n = $lines.Count
                $i = -1
            }
        }
    }
    if ($changed) {
        Set-Content -Path $file -Value $lines -NoNewline -Encoding UTF8
        Write-Host "Fixed $($_.Name)"
    }
}
