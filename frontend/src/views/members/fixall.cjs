const fs = require('fs');
const path = require('path');

const dir = __dirname;
const files = fs.readdirSync(dir).filter(f => f.endsWith('.vue'));

let totalFixed = 0;

files.forEach(file => {
  let content = fs.readFileSync(path.join(dir, file), 'utf8');
  let original = content;
  let fixes = 0;
  
  // Fix :width="xxxpx" patterns
  content = content.replace(/:width="(\d+)px"/g, ':width="$1"');
  
  // Fix placeholder/title/label corruption patterns (missing closing quote before ?)
  content = content.replace(/placeholder="([^"]*?)�\?/g, 'placeholder="$1"');
  content = content.replace(/title="([^"]*?)�\?/g, 'title="$1"');
  content = content.replace(/label="([^"]*?)�\?/g, 'label="$1"');
  
  // Fix 'text, key:' patterns where closing quote is missing
  content = content.replace(/包含会员�\? dataIndex:/g, '包含会员数", dataIndex:');
  content = content.replace(/状数, slotName:/g, '状态", slotName:');
  content = content.replace(/阅读数, dataIndex:/g, '阅读量", dataIndex:');
  content = content.replace(/兴趣分类数, :value/g, '兴趣分类数", :value');
  content = content.replace(/已标记会数, :value/g, '已标记会员数", :value');
  content = content.replace(/覆盖会员数, :value/g, '覆盖会员数", :value');
  content = content.replace(/沉默会员数, :value/g, '沉默会员数", :value');
  
  // Fix 'text?/>' patterns (missing quote before />)
  content = content.replace(/标签名�\?\/>/g, '标签名称" />');
  content = content.replace(/状�\?>/g, '状态">');
  
  // Fix </a-option> corruption
  content = content.replace(/营销类�\?\/a-option>/g, '营销类</a-option>');
  content = content.replace(/通知类�\?\/a-option>/g, '通知类</a-option>');
  content = content.replace(/验证类�\?\/a-option>/g, '验证码</a-option>');
  content = content.replace(/文章类�\?\/a-option>/g, '文章类</a-option>');
  content = content.replace(/视频类�\?\/a-option>/g, '视频类</a-option>');
  content = content.replace(/微信公众类�\?\/a-option>/g, '微信服务号</a-option>');
  content = content.replace(/APP推类�\?\/a-option>/g, 'APP推送</a-option>');
  
  // Fix mockData corruption
  content = content.replace(/美食爱好者数, category:/g, '美食爱好者", category:');
  content = content.replace(/数码产品爱好者数, status:/g, '数码产品爱好者", status:');
  content = content.replace(/新会员专享福利来数, summary:/g, '新会员专享福利来袭", summary:');
  content = content.replace(/新品上市抢先数, summary:/g, '新品上市抢先看", summary:');
  content = content.replace(/00元数, '/g, "00元', ");
  
  // Fix '已发数, : '草稿' }} patterns
  content = content.replace(/已发数, : '草稿' }}/g, "已发布', : '草稿' }}");
  
  // Fix mockData title corruption
  content = content.replace(/'新会员专享福利来数, summary:/g, "'新会员专享福利来袭', summary:");
  content = content.replace(/'新品上市抢先数, summary:/g, "'新品上市抢先看', summary:");
  
  if (content !== original) {
    fs.writeFileSync(path.join(dir, file), content);
    fixes++;
    totalFixed++;
  }
});

console.log(`Total files fixed: ${totalFixed}`);
