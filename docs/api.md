# API

## GET /healthz

Liveness probe.

Returns HTTP 200 when the service process is running.

Response:

```json
{
  "status": "ok"
}
```

---

## GET /readyz

Readiness probe.

Checks:

* nsjail availability
* configured language toolchains
* version commands

Returns HTTP 200 when all checks pass.

Returns HTTP 503 if any dependency fails.

Example:

```json
{
  "status": "ok",
  "languages": {
    "py3": {
      "ok": true,
      "version": "Python 3.11.2"
    },
    "cpp": {
      "ok": true,
      "version": "g++ 12.2.0"
    }
  }
}
```

---

## GET /info

Returns service metadata.

Example:

```json
{
  "build_info": {
    "version": "0.1.0",
    "commit": "abc123",
    "go_version": "go1.23"
  },
  "languages": []
}
```

---

## POST /run

Compiles and/or executes source code.

### Request

```json
{
  "language": "cpp",
  "source": "#include <iostream>\nint main(){std::cout<<\"hi\";}",
  "source_filename": "solution.cpp",
  "artifact_filename": "solution",
  "build": {
    "limits": {
      "wall_time_s": 5,
      "memory_kb": 1048576,
      "max_processes": 100
    },
    "flags": ["-O2"]
  },
  "run": {
    "limits": {
      "wall_time_s": 3,
      "memory_kb": 524288,
      "max_processes": 64
    },
    "flags": []
  },
  "tests": [
    {
      "stdin": "",
      "expected_stdout": "hi"
    }
  ]
}
```

### Response

```json
{
  "status": "accepted",
  "build": {
    "status": "ok",
    "stdout": "",
    "stderr": "",
    "duration_ms": 120
  },
  "tests": [
    {
      "status": "accepted",
      "stdout": "hi",
      "stderr": "",
      "duration_ms": 4,
      "memory_peak_kb": 0
    }
  ]
}
```

### Error Response

```json
{
  "error": {
    "code": "unknown_language",
    "message": "unknown language"
  }
}
```

### Supported Languages

* py3
* c
* cpp
* java
* bash
* js
* verilog
