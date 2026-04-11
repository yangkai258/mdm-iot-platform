const fs = require('fs');
const c = fs.readFileSync('C:/Users/YKing/.openclaw/workspace/mdm-project/frontend/src/views/ai/AISandboxView.vue', 'utf8');
const l = c.split('\n');
const important = ['a-card', 'a-modal', 'a-row', 'a-col', 'a-form', 'a-spin', 'a-table'];
for (let i = 0; i < l.length; i++) {
  const t = l[i].trim();
  for (const tag of important) {
    if (t === '<' + tag + '>' || t.startsWith('<' + tag + ' ')) {
      console.log(i + 1, 'OPENS', tag);
    }
    if (t === '</' + tag + '>') {
      console.log(i + 1, 'CLOSES', tag);
    }
  }
}
