# SBOM API

I built this project to understand how **software supply chain security** actually works - not just in theory, but through code I can run, break, and fix.

This service starts simple: a Go API that manages component records from a software bill of materials (SBOM). Each phase adds a layer of complexity and realism.

It starts with building the foundation (CRUD, dummy data) then moves onto scanning images, Vault integration, and fixing non-harmful yet realistic vulnerabilities.

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

**Purpose:** Strengthen the API's security posture before containerization.

**Focus areas:** dependency visibility, vulnerability scanning, and secure secret management.

### Overview
Phase 2 introduced the first DevSecOps control into the picture.

This phase centered on:
1. Building a **Software Bill of Materials (SBOM)** for full dependency transparency.
2. Running **automated vulnerability scans** to verify code safety.
3. Integrating **HashiCorp Vault** for runtime secret management.
4. Designing a **mock CI/CD workflow** that enforces these checks before build or deploy.
5. Adding 100 records of randomized dummy data to fill the database (in progress)

### Implementation Summary
- Reorganized project into api/handlers, api/vault, and cmd/ for clarity.

- Established a working CI pipeline (ci.yml) for automated build and test validation on every push.

- Integrated local SBOM generation and vulnerability scanning (Syft + manual verification).

- Ensured reproducible builds by aligning module paths (go.mod) and eliminating broken imports.

- Verified Vault package structure and basic function export (GetSecret) for Phase 3 integration. Lots of troubleshooting and fun!

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

- Established repeatable dependency and secret hygiene.

- CI pipeline passing consistently with GitHub Actions checks.

- Ready to containerize and scan in Phase 3 (Docker hardening).


### Lessons Learned 
- **Create a .gitignore early:** prevents compiled binaries and local artifacts from polluting the main branch.

- **Capture CLI output and failures:** keeping logs makes it easier to retrace steps and document learning progress, not just the results.

- **Use GitHub Releases for binaries:** each .exe can mark a milestone without cluttering the main branch. This goes hand-in-hand with .gitignore.

- **Be careful when editing folder structure:** moving files to different folders can break imports. One moment, everything is fine, and the next, the page is flooded with red squiggly lines. Go might even auto-remove these underlined imports when saving, leading to time spent troubleshooting.

- **Indentation matters in YAML:** spacing and alignment are strict; a single misplaced space can break CI and lead to failures. Make sure the right option is selected in the bottom right (LF vs CRLF). Spacing is important and can be the difference between getting back 0 errors or 5.

- **Troubleshoot YAML locally before pushing:** tools like yamllint can catch indentation or schema errors early and before pushing to GitHub. 

- **Simulate GitHub Actions locally:** you can use lightweight runners (e.g., act) to test workflows before committing. 

---

*Next phase -> prettify JSON file, containerization via Docker, SBOM + scan on built image, and runtime security validation.*
























