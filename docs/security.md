# Security

This document tracks the security issues identified in the hackathon specification and the implementation status of each.

## 1. Path Traversal via Filename

Status: Closed

Mitigation:

* Validate source_filename
* Validate artifact_filename
* Reject path separators
* Reject traversal sequences

Location:

```text
internal/validator/
```

---

## 2. Shell-Based Directory Operations

Status: Closed

Mitigation:

* Use os.MkdirAll
* Use os.WriteFile
* Avoid shell execution for filesystem operations

Location:

```text
internal/executor/
```

---

## 3. Compiler Flag Injection

Status: Pending

Required:

* Per-language allow-lists
* Reject unsafe compiler flags

---

## 4. Request Size Limits

Status: Pending

Required:

* Maximum source size
* Maximum test count
* Maximum stdin size
* Maximum expected output size

---

## 5. UID / Workspace Collisions

Status: Closed

Mitigation:

* UUID-based workspace generation

Location:

```text
internal/executor/
```

---

## 6. Unbounded Child Output

Status: Pending

Required:

* Output caps
* Truncation markers

---

## 7. Stale Workspace Directories

Status: Pending

Required:

* Guaranteed cleanup
* Startup sweep of orphaned workspaces

---

## NSJail

Current Status: In Progress

Target:

* PID namespace isolation
* Mount namespace isolation
* Network isolation
* Resource limits
