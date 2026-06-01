# Goboxd

Sandboxed multi-language code execution service.

## Supported Languages

| Language | ID |
|----------|----|
| Python 3 | py3 |
| C | c |
| C++ | cpp |
| Java | java |
| Bash | bash |
| JavaScript | js |
| Verilog | verilog |

## Quick Start

### Build

```bash
docker build -t goboxd .
```

### Run

```bash
docker run --rm --privileged -p 8080:8080 goboxd
```

### Health Check

```bash
curl localhost:8080/healthz
```

### Readiness Check

```bash
curl localhost:8080/readyz | jq
```

### Info

```bash
curl localhost:8080/info | jq
```

### Execute Python

```bash
curl -X POST localhost:8080/run \
-H "Content-Type: application/json" \
-d '{
  "language":"py3",
  "source":"print(\"hello\")",
  "tests":[
    {
      "stdin":"",
      "expected_stdout":"hello"
    }
  ]
}'
```

## Architecture

Request
→ Validation
→ Compilation (optional)
→ NSJail Sandbox
→ Test Execution
→ Result Aggregation
→ Response

## Sandbox

Goboxd executes user code through NSJail inside a Docker container.

Current protections:
- Filesystem isolation
- Namespace isolation
- Time limits
- Workspace isolation

Future protections:
- Memory limits
- Process limits
- Advanced cgroup controls