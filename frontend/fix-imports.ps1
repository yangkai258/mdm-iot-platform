$files = @(
    "src\views\devices\DeviceCertificates.vue",
    "src\views\devices\DeviceGeofence.vue",
    "src\views\devices\DeviceMonitorPanel.vue",
    "src\views\devices\DevicePairing.vue",
    "src\views\devices\DeviceRemoteControl.vue",
    "src\views\advanced\FeatureConfigView.vue"
)

foreach ($f in $files) {
    if (Test-Path $f) {
        $content = Get-Content $f -Raw
        # Replace: "}; import { Message } from '@arco-design/web-vue'"
        # With:    "} from 'vue'\nimport { Message } from '@arco-design/web-vue'"
        $fixed = $content -replace "\}; import \{ Message \} from '@arco-design/web-vue'", "} from 'vue'`nimport { Message } from '@arco-design/web-vue'"
        $fixed | Set-Content $f -NoNewline
        Write-Host "Fixed: $f"
    }
}
