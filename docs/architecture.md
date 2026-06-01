# Architecture

## Overview

goboxd is organized around a configuration-driven execution pipeline.

```text
Client
  |
HTTP API
  |
Validator
  |
Language Registry
  |
Executor
  |
Compiler (optional)
  |
Runtime
  |
Response
```

## Components

### internal/api

HTTP handlers.

Endpoints:

* /healthz
* /readyz
* /info
* /run

### internal/config

Loads language configuration from:

```text
internal/config/languages.yaml
```

Provides:

* language definitions
* compiler definitions
* runtime definitions
* default limits

### internal/validator

Validates:

* required fields
* filename safety
* request structure

### internal/executor

Responsible for:

* workspace creation
* source file generation
* compilation
* execution
* test evaluation
* response generation

### internal/models

Shared request and response structures.

## Workspace Lifecycle

```text
Create Workspace
    |
Write Source File
    |
Compile (optional)
    |
Execute Tests
    |
Collect Results
    |
Cleanup
```

## Docker Architecture

All language toolchains are installed during Docker image build.

Language installation scripts:

```text
scripts/lang_install/
```

The runtime image contains:

* goboxd
* nsjail
* Python
* GCC
* G++
* Java
* Bash
* Node.js
* Verilog toolchain
