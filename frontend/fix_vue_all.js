const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const viewsDir = 'C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/src/views';
const feDir = 'C:/Users/YKing/.openclaw/workspace/mdm-project/frontend';

function walk(d) {
  const r = [];
  for (const e of fs.readdirSync(d, { withFileTypes: true })) {
    const f = path.join(d, e.name);
    if (e.isDirectory()) r.push(...walk(f));
    else if (e.name.endsWith('.vue')) r.push(f);
  }
  return r;
}

function fixVue(content) {
  let c = content;

  // P0: Strip UTF-8 BOM
  if (c.charCodeAt(0) === 0xFEFF) {
    c = c.substring(1);
  }

  // P1: trailing </a-card> after </style>
  const styleEnd = c.indexOf('</style>');
  if (styleEnd !== -1) {
    const after = c.substring(styleEnd + 8);
    if (after.trim().startsWith('</a-card>')) {
      const idx = after.indexOf('</a-card>');
      c = c.substring(0, styleEnd + 8) + after.substring(idx + 9);
    }
  }

  // P2: orphan </a-card> after </script> (before root </template>)
  const scriptEnd = c.indexOf('</script>');
  if (scriptEnd !== -1) {
    const after = c.substring(scriptEnd + 9);
    if (after.trim().startsWith('</a-card>')) {
      const idx = after.indexOf('</a-card>');
      c = c.substring(0, scriptEnd + 9) + after.substring(idx + 9);
    }
  }

  const lines = c.split('\n');
  const fixed = [];
  let changed = false;

  function openTablesSoFar(arr) {
    let cnt = 0;
    for (const ln of arr) {
      const lt = ln.trim();
      if (/<a-table\b[^>]*\/>/.test(lt)) {}
      else if (/<a-table\b/.test(lt) && !lt.includes('/>') && !lt.includes('</a-table>')) cnt++;
      if (lt === '</a-table>') cnt--;
    }
    return cnt;
  }

  for (let i = 0; i < lines.length; i++) {
    const t = lines[i].trim();
    const lnum = i + 1;

    // P3: orphan </a-table>
    if (t === '</a-table>') {
      const cnt = openTablesSoFar(fixed);
      if (cnt <= 0) { changed = true; continue; }
      if (cnt === 1) {
        let j = i + 1;
        while (j < lines.length && lines[j].trim() === '') j++;
        const next = j < lines.length ? lines[j].trim() : '';
        if (next.startsWith('<template') || next.startsWith('</a-')) {
          changed = true; continue;
        }
      }
    }

    // P4: missing </a-table> before </a-card>
    if (t === '</a-card>') {
      let prevIdx = i - 1;
      while (prevIdx >= 0 && lines[prevIdx].trim() === '') prevIdx--;
      const prev = prevIdx >= 0 ? lines[prevIdx].trim() : '';
      const cnt = openTablesSoFar(fixed);
      if (cnt > 0 && prev !== '</template>') {
        const indent = lines[i].match(/^(\s*)/)[1];
        fixed.push(indent + '    </a-table>');
        changed = true;
      }
    }

    fixed.push(lines[i]);
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
    if (m) return { ok: false, file: m[1], line: parseInt(m[2]), type: m[4] };
    return { ok: false, file: null, out: out.substring(0, 300) };
  }
}

// Phase 1: Fix all files once
console.log('=== Phase 1: Scanning all files ===');
const files = walk(viewsDir);
let totalFixed = 0;
const results = [];

for (const f of files) {
  const original = fs.readFileSync(f, 'utf8');
  const { content, type } = fixVue(original);
  if (type) {
    fs.writeFileSync(f, content, 'utf8');
    results.push(path.relative(viewsDir, f) + ' (' + type + ')');
    totalFixed++;
  }
}
console.log('Files fixed:', totalFixed);
for (const r of results) console.log(' -', r);

// Phase 2: Iterative build fix
console.log('\n=== Phase 2: Iterative build fix ===');
for (let r = 0; r < 30; r++) {
  const result = runBuild();
  if (result.ok) { console.log('BUILD SUCCESS!'); break; }
  if (!result.file) { console.log('Unknown error:', result.out); break; }
  console.log('Round ' + r + ': ' + result.file + ' L' + result.line + ' (' + result.type + ')');
  
  // Fix this specific file
  const fullPath = path.join(viewsDir, result.file);
  const original = fs.readFileSync(fullPath, 'utf8');
  const { content, type } = fixVue(original);
  if (type) {
    fs.writeFileSync(fullPath, content, 'utf8');
    console.log('  -> Fixed! type=' + type);
  } else {
    // Show context
    const lines = fs.readFileSync(fullPath, 'utf8').split('\n');
    const el = Math.max(0, result.line - 5);
    for (let i = el; i < Math.min(lines.length, result.line + 4); i++) console.log('  L' + (i + 1) + ': ' + lines[i].substring(0, 80));
    break;
  }
}
