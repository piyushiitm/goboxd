# Benchmarks

Benchmarking is scheduled for Stage 3.

This document will contain:

* Requests per second
* p50 latency
* p95 latency
* p99 latency

Measured under:

* 1 concurrent client
* 10 concurrent clients
* 50 concurrent clients
* 100 concurrent clients

## Test Case

```json
{
  "language": "py3",
  "source": "print(\"hello\")",
  "tests": [
    {
      "stdin": "",
      "expected_stdout": "hello"
    }
  ]
}
```

## Results

Not yet measured.
