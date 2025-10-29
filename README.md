# SBOM API

I built this project to understand how **software supply chain security** actually works - not just in theory, but through code I can run, break, and fix.

This service starts simple: a Go API that manages component records from a software bill of materials (SBOM). Each phase adds a layer of realism.

The layering starts with introducing safe vulnerabilities to generating SBOMs and scanning images, then finally, integrating secrets management with Vault. 

## Project Roadmap
Each phase will add more depth:
- Introduce a non-harmful supply chain flaw (unpinned dependency, missing validation)
- Simulate auth and token errors
- Generate and scan SBOMs (Syft/Trivy)
- Harden CI/CD (pinned actions, artifact signing)
- Integrate Vault for dynamic secrets
- Containerize and scan the image

---

## Phase 1 - CRUD API

### **Goal:** Get a working foundation before layering in security and automation. 

### Fields
- Package Name
- Version
- Checksum
- Source repo
- License

### Endpoints
| Method | Endpoint | Description |
|--------|-----------|--------------|
| POST | `/components` | Create a new component |
| GET | `/components` | List all components |
| GET | `/components?id={id}` | Get component by ID |
| PUT | `/components?id={id}` | Replace component |
| PATCH | `/components?id={id}` | Update component fields |
| DELETE | `/components?id={id}` | Delete component |

---

## Quick Start
```
git clone https://github.com/machineash/sbom-api.git
cd sbom-api
go run main.go
```
Then, open http://localhost:8080/components

Use ```curl``` or Postman to add and test records.

---

## Phase 2 - Secure Software Supply Chain

**Purpose:** Strengthen this Go API's security posture before containerization.

**Focus areas:** dependency visibility, vulnerability scanning, and secure secret management.

### Overview
Phase 2 introduced the first DevSecOps controls into the project.

The work centered on:
1. Building a **Software Bill of Materials (SBOM)** for full dependency transparency.
2. Running **automated vulnerability scans** to verify code safety.
3. Integrating **HashiCorp Vault** for runtime secret management.
4. Designing a **mock CI/CD workflow** that enforces these checks before build or deploy.

### Implementation Summary
- Reorganized project into api/handlers, api/vault, and cmd/ for cleaner modularity.

- Established a working CI pipeline (ci.yml) for automated build and test validation on every push.

- Integrated local SBOM generation and vulnerability scanning (Syft + manual verification).

- Ensured reproducible builds by aligning module paths (go.mod) and eliminating broken imports.

- Verified Vault package structure and basic function export (GetSecret) for Phase 3 integration.

### Key Artifacts
| File | Description |
|------|---------------|
| `artifacts/sbom.json` | Generated SBOM file (Phase 2 snapshot) |
| `artifacts/vuln-report.txt` | Vulnerability scan results |
| `cmd/main.go` | Entry point connecting handlers and vault packages |
| `.github/workflows/ci.yml` | CI workflow for build/test automation |
| GitHub Release | Compiled binary representing the Phase 2 Milestone |


### Outcomes
- Verified zero known vulnerabilities at time of scan.

- Removed hardcoded secrets from the codebase.

- Established repeatable dependency and secret hygiene.

- CI pipeline passing consistently with GitHub Actions checks.

- Ready to containerize and re-scan in Phase 3 (Docker hardening).


### Lessons Learned 
- **Create a .gitignore early:** prevents compiled binaries and local artifacts from polluting the main branch.

- **Capture CLI output and failures:** keeping logs makes it easier to retrace steps and show learning progress, not just results.

- **Use GitHub Releases for binaries:** each .exe can mark a milestone without cluttering the branch.

- **Be careful when refactoring folder structure:** moving files or saving mid-refactor can break imports; once saved, Go might auto-remove imports, leading to time spent troubleshooting.

- **Indentation matters in YAML:** spacing and alignment are strict; a single misplaced space can break CI.

- **Troubleshoot YAML locally before pushing:** tools like yamllint can catch indentation or schema errors early.

- **Simulate GitHub Actions locally:** you can use lightweight runners (e.g., act) to test workflows before committing; GitHub Actions remains the final validator on push.

- **Ensure line endings are set to LF:** inconsistent endings (CRLF vs LF) can confuse YAML parsers on different systems.



---

*Next phase -> containerization, SBOM + scan on built image, and runtime security validation.*
























