const http = require('http');

const BASE_URL = 'http://localhost:8080';

function request(method, path, body, token) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, BASE_URL);
    const options = {
      hostname: url.hostname,
      port: url.port,
      path: url.pathname + url.search,
      method: method,
      headers: {
        'Content-Type': 'application/json'
      }
    };
    if (token) {
      options.headers['Authorization'] = 'Bearer ' + token;
    }

    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          resolve({ status: res.statusCode, data: JSON.parse(data) });
        } catch {
          resolve({ status: res.statusCode, data: data });
        }
      });
    });
    req.on('error', reject);
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

async function runTests() {
  console.log('=== MDM Backend API Test Report ===\n');

  // Step 1: Login
  let loginResult;
  try {
    loginResult = await request('POST', '/api/v1/auth/login', { username: 'admin', password: 'admin123' });
    console.log('1. POST /api/v1/auth/login:', loginResult.status === 200 ? 'PASS' : 'FAIL', JSON.stringify(loginResult.data));
  } catch (e) {
    console.log('1. POST /api/v1/auth/login: FAIL -', e.message);
    return;
  }

  if (loginResult.status !== 200 || !loginResult.data.data?.token) {
    console.log('Login failed, cannot continue with authenticated tests.');
    return;
  }

  const token = loginResult.data.data.token;
  console.log('Token obtained successfully.\n');

  const apis = [
    // Sprint 1-8
    { method: 'GET', path: '/api/v1/devices', name: '设备列表' },
    { method: 'GET', path: '/api/v1/members', name: '会员列表' },
    { method: 'GET', path: '/api/v1/orders', name: '订单列表' },
    { method: 'GET', path: '/api/v1/coupons', name: '优惠券列表' },
    { method: 'GET', path: '/api/v1/member/cards', name: '会员卡列表' },
    { method: 'GET', path: '/api/v1/member/levels', name: '会员等级' },
    // Sprint 9-12
    { method: 'GET', path: '/api/v1/alerts', name: '告警列表' },
    { method: 'GET', path: '/api/v1/compliance/policies', name: '合规策略' },
    { method: 'GET', path: '/api/v1/ldap/config', name: 'LDAP配置' },
    { method: 'GET', path: '/api/v1/subscriptions', name: '订阅列表' },
    // Sprint 13-16
    { method: 'GET', path: '/api/v1/health/reports', name: '健康报告' },
    { method: 'GET', path: '/api/v1/emotion/records', name: '情绪记录' },
    { method: 'GET', path: '/api/v1/digital-twin/status', name: '数字孪生状态' },
    { method: 'GET', path: '/api/v1/subscriptions/billing', name: '账单' },
    // Sprint 17-20
    { method: 'GET', path: '/api/v1/ai/embodied/perception', name: '具身感知' },
    { method: 'GET', path: '/api/v1/simulation/runs', name: '仿真运行' },
    { method: 'GET', path: '/api/v1/health/tracking', name: '健康追踪' },
    // Sprint 21-32
    { method: 'GET', path: '/api/v1/ai/model/shards', name: '模型分片' },
    { method: 'GET', path: '/api/v1/insurance/products', name: '保险产品' },
    { method: 'GET', path: '/api/v1/pet-social/feed', name: '宠物社交' },
    { method: 'GET', path: '/api/v1/research/datasets', name: '研究数据集' },
  ];

  const passed = [];
  const failed = [];
  const notFound = [];

  for (const api of apis) {
    try {
      const result = await request(api.method, api.path, null, token);
      if (result.status === 404) {
        notFound.push({ api: api.name, path: api.path, status: 404 });
      } else if (result.status >= 200 && result.status < 300) {
        passed.push({ api: api.name, path: api.path, status: result.status });
      } else {
        failed.push({ api: api.name, path: api.path, status: result.status, error: result.data });
      }
    } catch (e) {
      failed.push({ api: api.name, path: api.path, status: 'ERROR', error: e.message });
    }
  }

  console.log('\n=== Test Results ===\n');
  console.log('### PASS (' + passed.length + ')');
  passed.forEach(a => console.log('- ' + a.api + ': PASS (' + a.status + ')'));

  console.log('\n### FAIL (' + failed.length + ')');
  failed.forEach(a => console.log('- ' + a.api + ': FAIL (' + a.status + ') - ' + JSON.stringify(a.error)));

  console.log('\n### NOT FOUND (404) (' + notFound.length + ')');
  notFound.forEach(a => console.log('- ' + a.api + ': 404 - ' + a.path));
}

runTests().catch(console.error);
