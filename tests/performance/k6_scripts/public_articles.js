import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081/api/v1';

export default function () {
  const res = http.get(`${BASE_URL}/articles?limit=10&offset=0`);
  check(res, {
    'status is 200': (r) => r.status === 200,
    'has items array': (r) => Array.isArray(r.json('items')),
  });

  sleep(1);
}
