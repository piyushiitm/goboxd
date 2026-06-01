# AI Usage Log

This file records significant AI-assisted discussions that influenced goboxd. Entries are ordered by when the underlying work happened, starting from the evening of 2026-05-31.

---

## 2026-05-31 20:00 · Getting started with goboxd and the hackathon spec

**Prompt:**
I do not understand the hackathon, the spec page, or what I should build first. Explain the minimum I need to do to get moving, and keep it short.

**Response summary:**
Explained that Stage 1 is a small Go HTTP service, suggested starting with a minimal `/healthz` endpoint and a tiny server, and helped narrow the initial goal instead of trying to build the final sandbox immediately.

**What we used / didn't use:**
Used the “start small and focus on the prototype” advice. Did not use the early temptation to overbuild the final architecture before understanding the basics.

---

## 2026-05-31 20:45 · Ubuntu / UTM setup was blocking progress

**Prompt:**
Help me recover or rebuild my Ubuntu VM in UTM, because I forgot the password and the installer seems stuck.

**Response summary:**
Walked through EFI / boot-manager / grub recovery steps, then suggested deleting the broken VM and recreating it when recovery became a time sink. Also guided the Ubuntu installer screens and ARM64 VM setup.

**What we used / didn't use:**
Used the decision to stop wasting time on the broken VM and move to a fresh setup. Did not keep chasing grub once it was clear the VM setup was costing too much time.

---

## 2026-05-31 21:15 · Switching back to macOS for development

**Prompt:**
Can I build this on my native Mac instead of Linux, and how do I install Go on macOS?

**Response summary:**
Recommended using Homebrew to install Go, along with git, g++, and Python, and then using VS Code on the Mac rather than fighting Linux keyboard and VM issues.

**What we used / didn't use:**
Used the macOS setup because it was faster and more comfortable. Did not spend more time on Linux-specific editor or clipboard issues once the Mac path was chosen.

---

## 2026-05-31 21:45 · Repo fork, scaffold, and project layout

**Prompt:**
I forked the submission repo and want to inspect what is already there. What structure should I create and how should I work from the fork?

**Response summary:**
Suggested cloning the forked repository, inspecting the scaffold, and creating a simple Go layout under `cmd/` and `internal/` rather than inventing a large architecture too early.

**What we used / didn't use:**
Used the idea of keeping the repository small and structured. Did not assume any branch naming rule until the spec was checked directly.

---

## 2026-05-31 22:15 · First Go HTTP server and API basics

**Prompt:**
How do I create `/healthz`, how do I run the server, and why does `Command+S` sometimes seem to delete text in VS Code?

**Response summary:**
Explained `net/http`, HTTP status codes, JSON parsing for `/run`, and the difference between response body and status. Also explained macOS save/run shortcuts and basic VS Code behavior.

**What we used / didn't use:**
Used the minimal server and handler pattern. Did not keep any extra browser-only or UI distractions; the server first had to compile and answer requests correctly.

---

## 2026-05-31 22:45 · Designing the execution pipeline

**Prompt:**
How can the service support new languages without changing Go code, and should I use YAML with placeholders?

**Response summary:**
Explained the difference between hardcoded language branches and a config-driven registry, introduced placeholders like `{source}` and `{artifact}`, and recommended one executor that reads language settings from config.

**What we used / didn't use:**
Used the placeholder-based design because it matched the hackathon direction. Did not use a per-language `RunPython` / `RunCpp` split once the registry approach made more sense.

---

## 2026-05-31 23:30 · The spec, the prototype, and the contract mismatch

**Prompt:**
I don’t think this prototype will even make it to the second stage. We need to move faster and definitely need to change the contract that should be received by the service and the output returned by it, with the error vocabulary given.

**Response summary:**
After reading the full spec more carefully, it became clear that the final API was much richer than the prototype. The discussion shifted from a simple stdout/stderr model to status-based build and test results.

**What we used / didn't use:**
Used the realization that the spec contract matters more than the prototype shape. Did not keep assuming the original stdout/stderr-only response would be enough.

---

## 2026-06-01 00:15 · Docker runtime and dependency issues

**Prompt:**
My Docker image builds but the container cannot find Python or g++. How should I fix the runtime image?

**Response summary:**
Identified missing runtime dependencies in the container, suggested installing `python3` and `g++` in the runtime image, and verified that `/healthz` and `/run` worked inside Docker after rebuilding.

**What we used / didn't use:**
Used the runtime package install fix and rebuilt the image. Did not keep the broken image state once the service ran correctly in Docker.

---

## 2026-06-01 00:45 · Unit tests and config loading

**Prompt:**
What are unit tests, which tests should I add first, and why is `go test` failing on `config/languages.yaml`?

**Response summary:**
Explained unit tests as automatic checks for small pieces of logic, suggested tests for placeholder replacement and unsupported language handling, and later moved the YAML registry into the binary with `go:embed` to avoid path issues.

**What we used / didn't use:**
Used the small executor tests and the embed approach. Did not keep the original file-path-dependent test design because it broke under different working directories.

---

## 2026-06-01 01:15 · YAML embedding and removing file-path dependence

**Prompt:**
How do I stop `config/languages.yaml` from breaking tests and Docker? Should I move it or embed it?

**Response summary:**
Recommended moving the YAML into `internal/config` and embedding it with `go:embed`, then loading it from bytes instead of reading from disk at runtime.

**What we used / didn't use:**
Used the embed-based solution because it removed path problems in tests and Docker. Did not keep the runtime `os.ReadFile(...)` version.

---

## 2026-06-01 01:45 · Cleaning error output from the executor

**Prompt:**
My executor returns ugly errors like `exit status 1| ...`. How do I make the stderr cleaner?

**Response summary:**
Suggested using the child process output as the primary error message and falling back to `err.Error()` only when there is no captured output, which removed the noisy Go-specific prefix.

**What we used / didn't use:**
Used the cleaner stderr handling. Did not keep the `exit status 1|` prefix because it was internal noise, not useful user-facing output.

---

## 2026-06-01 02:15 · Adding timeouts to stop infinite loops

**Prompt:**
How do I stop user code from running forever, and why do I need `context`?

**Response summary:**
Recommended `context.WithTimeout(...)` plus `exec.CommandContext(...)` so Go can kill long-running processes after a fixed limit, which made infinite loops terminate after about five seconds.

**What we used / didn't use:**
Used the context-based timeout approach. Did not use a manual polling or timer loop because `CommandContext` already handles process cancellation cleanly.

---

## 2026-06-01 02:45 · Reading the full spec and re-evaluating the API

**Prompt:**
Read the spec properly and tell me whether the current prototype API is enough.

**Response summary:**
After reading the full spec, realized that Stage 1 is only the prototype, but the final API contract is much richer: status vocabulary, `build`, `tests`, `readyz`, `info`, plug-and-play YAML, per-test results, and security / concurrency requirements.

**What we used / didn't use:**
Used the clarification that the prototype is only a stepping stone and that the final API would need redesign. Did not keep assuming that `stdout` / `stderr` / `exit_code` would be the final contract.

---

## 2026-06-01 03:15 · Switching chats and preserving context

**Prompt:**
This chat is getting laggy. How do I switch to a new chat without losing the context and history?

**Response summary:**
Suggested creating a concise handoff summary that captures the current architecture, what has already been implemented, the current executor behaviour, and the next planned work.

**What we used / didn't use:**
Used the idea of maintaining a structured handoff note. Did not rely on memory alone when moving between chats.

---

## 2026-06-01 04:00 · Java support through the language registry

**Prompt:**
Lets just test it by adding another language support so that we know we only need to create java.sh and change yaml file for java.

**Response summary:**
Validated that language support could be added through the registry and install scripts rather than executor changes, and then debugged Java compilation/runtime behaviour through Docker.

**What we used / didn't use:**
Used the registry-driven approach and Java installation script. Did not add Java-specific branching to the executor.

---

## 2026-06-01 05:00 · Docker server verification for the new language

**Prompt:**
But first we need to start server with docker build.

**Response summary:**
Reinforced that all meaningful verification should happen inside the Docker image, because the image is the real unit of correctness for the hackathon.

**What we used / didn't use:**
Used Docker as the validation environment. Did not trust host-machine binaries as proof that the service would work in submission.

---

## 2026-06-01 05:30 · Adding more languages through shell scripts

**Prompt:**
We need to create sh files.

**Response summary:**
Created installation scripts for the additional languages and used them to install toolchains in the Docker image.

**What we used / didn't use:**
Used script-per-language installation. Did not keep growing the Dockerfile directly for each toolchain.

---

## 2026-06-01 06:00 · Verifying that all in-scope languages should work

**Prompt:**
Make sure again it will work for all languages that are mentioned.

**Response summary:**
Checked that the architecture could support C, C++, Java, Python 3, Bash, JavaScript (Node), and Verilog, with Java and Verilog requiring slightly different runtime configuration.

**What we used / didn't use:**
Used the idea that the registry can support the full in-scope set. Did not assume every language behaves like a simple single-binary compiled language.

---

## 2026-06-01 06:30 · Stage 1 tests and repo quality

**Prompt:**
Ok lets do that.

**Response summary:**
Focused on adding unit tests for deterministic parts of the codebase: config loading, placeholder replacement, validation, and limit resolution.

**What we used / didn't use:**
Used unit tests that do not depend on sandboxing. Did not turn unit tests into full integration tests.

---

## 2026-06-01 07:00 · README rewrite

**Prompt:**
Now refactor readme.

**Response summary:**
Reviewed the README against the spec and trimmed it down so it stays short and operational, with the longer explanations moved into docs/.

**What we used / didn't use:**
Used a short README that explains what the project is and how to run it. Did not keep feature-heavy prose or long architecture details in the README.

---

## 2026-06-01 07:30 · Documentation planning

**Prompt:**
Now lets do documentation till now.

**Response summary:**
Created separate documentation files for the API, languages, architecture, security, and benchmarks instead of keeping everything in the README.

**What we used / didn't use:**
Used a docs/ structure that matches the spec’s request for dedicated documents. Did not keep all documentation in one file.

---

## 2026-06-01 09:00 · Reviewing Stage 1 versus Stage 2

**Prompt:**
Think properly again, take time to read specs.

**Response summary:**
Re-checked the stage requirements and realized that Stage 1 is smaller than the full target architecture, while Stage 2 introduces YAML language registration, `/readyz`, `/info`, and flag allow-lists.

**What we used / didn't use:**
Used the stage separation to prioritize work correctly. Did not treat later-stage functionality as required Stage 1 scope.

---

## 2026-06-01 10:00 · Spec reread and current gaps

**Prompt:**
Also read the spec again, let me paste it for you again and kindly refactor according to it.

**Response summary:**
Re-read the hackathon specification and identified the remaining gaps: exact status vocabulary, validated build/run flags, `/readyz` and `/info` completeness, nsjail integration, concurrency, security docs, and Stage 3 work.

**What we used / didn't use:**
Used the spec as the source of truth for prioritization. Did not keep optimizing areas that were already good enough for the current stage.

---

## 2026-06-01 10:30 · Documentation completeness check

**Prompt:**
So all the document files are complete?

**Response summary:**
Reviewed the docs and concluded that the README and the docs skeleton were in place, but the AI logs still needed to be tied to real conversations and the later-stage implementation gaps still remained.

**What we used / didn't use:**
Used the docs structure as a planning tool. Did not treat the docs as finished until they matched the actual implementation.


