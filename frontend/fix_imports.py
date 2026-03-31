import re
import os

files = [
    r"src\views\devices\DeviceCertificates.vue",
    r"src\views\devices\DeviceGeofence.vue",
    r"src\views\devices\DeviceMonitorPanel.vue",
    r"src\views\devices\DevicePairing.vue",
    r"src\views\devices\DeviceRemoteControl.vue",
    r"src\views\advanced\FeatureConfigView.vue"
]

for f in files:
    if os.path.exists(f):
        with open(f, 'r', encoding='utf-8') as file:
            content = file.read()
        
        original = content
        
        # Pattern: "}; import { X } from '@arco-design/web-vue'"
        # Should become: "} from 'vue'\nimport { X } from '@arco-design/web-vue'"
        # This handles multiple imports on the same line like: "}; import { X } from '...'; import draggable from '...'"
        pattern = r"import (\{[^{}]+\}); import (\{[^}]+\}) from '([^']+)' import (\w+) from '([^']+)'"
        replacement = r"import \1 from 'vue'\nimport \2 from '\3'\nimport \4 from '\5'"
        content = re.sub(pattern, replacement, content)
        
        # Pattern: "}; import { X } from '@arco-design/web-vue'  import { Y } from '@arco-design/web-vue'"
        pattern2 = r"import (\{[^{}]+\});  import (\{[^}]+\}) from '([^']+)'"
        replacement2 = r"import \1 from 'vue'\nimport \2 from '\3'"
        content = re.sub(pattern2, replacement2, content)
        
        # Pattern: "}; import { X } from '@arco-design/web-vue'"
        pattern3 = r"import (\{[^{}]+\}); import (\{[^}]+\}) from '([^']+)'"
        replacement3 = r"import \1 from 'vue'\nimport \2 from '\3'"
        content = re.sub(pattern3, replacement3, content)
        
        if content != original:
            with open(f, 'w', encoding='utf-8') as file:
                file.write(content)
            print(f"Fixed: {f}")
        else:
            print(f"No change: {f}")
    else:
        print(f"Not found: {f}")
