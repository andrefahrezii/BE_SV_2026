# Performance Tests

Performance tests using [k6](https://k6.io/).

## Scripts

- `public_articles.js` - load test public article listing
- `admin_articles.js` - load test admin article listing (requires login)
- `login.js` - load test login endpoint

## Run

```bash
k6 run tests/performance/k6_scripts/public_articles.js
k6 run tests/performance/k6_scripts/admin_articles.js
k6 run tests/performance/k6_scripts/login.js
```

With custom base URL:

```bash
BASE_URL=http://localhost:8080/api/v1 k6 run tests/performance/k6_scripts/admin_articles.js
```
