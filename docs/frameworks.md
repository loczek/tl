# Backend Frameworks

## Fiber

- doesn't support request cancellation
- isn't compatible with `net/http`
- doesn't support http 1.0 or http 2.0, only http 1.1
- maybe faster than `net/http` for large quantity of small requests
