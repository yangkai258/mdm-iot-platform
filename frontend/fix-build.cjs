const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const iconv = require('iconv-lite');

const viewsDir = 'C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/src/views';
const feDir = 'C:/Users/YKing/.openclaw/workspace/mdm-project/frontend';

// FIXER: remove orphan </a-card> and </a-table> after </style></script>
function fixVue(content) {
  let c = content;
  let changed = false;
  if (c.charCodeAt(0) === 0xFEFF) { c = c.substring(1); changed = true; }
  
  // Check for GBK encoding issues (replacement chars)
  if (c.includes('\ufffd')) {
    try {
      const buf = Buffer.from(content, 'utf8');
      c = iconv.decode(buf, 'gbk');
      changed = true;
    } catch (e) {}
  }
  
  const lines = c.split('\n');
  const fixed = [];
  let tableDepth = 0;
  let styleDepth = 0;
  let scriptDepth = 0;
  for (let i = 0; i < lines.length; i++) {
    const lt = lines[i].trim();
    const raw = lines[i];

    if (lt.startsWith('<style')) styleDepth++;
    if (lt === '</style>') styleDepth--;
    if (lt.startsWith('<script')) scriptDepth++;
    if (lt === '</script>') scriptDepth--;

    if (styleDepth === 0 && scriptDepth === 0) {
      // Remove orphan </a-card> after </style> or </script> (when prev was </style> or </script>)
      if ((lt === '</a-card>' || lt === '</a-table>') && fixed.length > 0) {
        const prev = fixed[fixed.length - 1].trim();
        if (prev === '</style>' || prev === '</script>' || prev === '</a-table>' || prev === '</a-card>') {
          changed = true;
          continue;
        }
      }

      if (lt === '</a-table>') {
        if (tableDepth === 0) { changed = true; continue; }
        tableDepth--;
      }
      if (/<a-table\b/.test(lt) && !lt.includes('/>') && !lt.includes('</a-table>')) {
        tableDepth++;
      }
    }
    fixed.push(raw);
  }
  return changed ? { content: fixed.join('\n'), type: 'fixed' } : { content, type: null };
}

function runBuild() {
  try {
    execSync('npm run build', { cwd: feDir, stdio: 'pipe', encoding: 'utf8', timeout: 60000 });
    return { ok: true };
  } catch (e) {
    const out = (e.stderr || e.stdout || '').toString();
    const m = out.match(/src\/views\/([^:]+\.vue) \((\d+):(\d+)\): (Invalid end tag|Element is missing end tag)/);
    if (m) return { ok: false, file: m[1], line: parseInt(m[2]) };
    return { ok: false, file: null, out: out.substring(0, 500) };
  }
}

for (let r = 0; r < 50; r++) {
  const result = runBuild();
  if (result.ok) { console.log('BUILD SUCCESS at round ' + r); break; }
  if (!result.file) { console.log('Unknown:', result.out); break; }
  console.log('Round ' + r + ': ' + result.file + ' L' + result.line);
  const fp = path.join(viewsDir, result.file);
  const orig = fs.readFileSync(fp, 'utf8');
  const { content } = fixVue(orig);
  fs.writeFileSync(fp, content, 'utf8');
}
