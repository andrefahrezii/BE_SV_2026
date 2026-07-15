import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 10 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081/api/v1';

export function setup() {
  const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
    email: 'admin@sharingvision.id',
    password: 'admin123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  const token = loginRes.json('token');
  return { token };
}

export default function (data) {
  // public list
  const r1 = http.get(`${BASE_URL}/categories`);
  check(r1, {
    'public categories status 200': (r) => r.status === 200,
    'has items': (r) => Array.isArray(r.json('items')),
  });

  // admin list
  const r2 = http.get(`${BASE_URL}/admin/categories`, {
    headers: { 'Authorization': `Bearer ${data.token}` },
  });
  check(r2, {
    'admin categories status 200': (r) => r.status === 200,
  });

  sleep(1);
}
