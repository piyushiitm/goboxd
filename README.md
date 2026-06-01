# goboxd

goboxd is a Go HTTP service that compiles and executes untrusted code inside isolated sandboxes.

## Framework

net/http

Reason: the API is small and does not require additional routing features.

## Supported Languages

- Python 3
- C
- C++
- Java
- Bash
- JavaScript (Node.js)
- Verilog

## Requirements

- Docker
- Docker Compose

## Run

Build:

```bash
make build
```

Start:

```bash
make run
```

Health check:

```bash
curl localhost:8080/healthz
```

## Tests

```bash
make test
```

## Load Testing

```bash
make load
```

## Documentation

- docs/api.md
- docs/languages.md
- docs/security.md
- docs/architecture.md
- docs/benchmarks.md

## License

GPL-3.0