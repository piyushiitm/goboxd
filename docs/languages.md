# Languages

Language support is configured entirely through:

```text
internal/config/languages.yaml
```

Each language defines:

* Display name
* Source file extension
* Version probe
* Compile command (optional)
* Run command
* Default execution limits

## Placeholders

### {source}

Absolute path to the generated source file.

Example:

```text
workspace/1234/source.cpp
```

### {artifact}

Absolute path to the compiled output.

Example:

```text
workspace/1234/program
```

### {workspace}

Request workspace directory.

Example:

```text
workspace/1234
```

## Language Lifecycle

### Interpreted Languages

Examples:

* Python
* Bash
* Node.js

Flow:

```text
Request
→ Write Source File
→ Execute Runtime
→ Evaluate Tests
```

### Compiled Languages

Examples:

* C
* C++
* Java
* Verilog

Flow:

```text
Request
→ Write Source File
→ Compile
→ Execute Artifact
→ Evaluate Tests
```

## Adding a New Language

1. Create installation script

```text
scripts/lang_install/<language>.sh
```

2. Install toolchain

Example:

```bash
apt-get install -y ruby
```

3. Register language

```yaml
ruby:
  name: "Ruby"

  extension: ".rb"

  version_command:
    command: "ruby"
    args:
      - "--version"

  run:
    - "ruby"
    - "{source}"
```

4. Rebuild Docker image

```bash
docker build -t goboxd .
```

5. Verify registration

```bash
curl localhost:8080/readyz
```

No Go code changes should be required.
