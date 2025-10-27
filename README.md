# SBOM API

I built this project to understand how **software supply chain security** actually works - not just in theory, but through code I can run, break, and fix.

This service starts simple: a Go API that manages component records from a software bill of materials (SBOM). Each phase adds a layer of realism.

It starts with introducing safe vulnerabilities to generating SBOMs and scanning images, then finally, integrating secrets management with Vault. 

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

## What Comes Next
Each phase will add more depth:
- Introduce a non-harmful supply chain flaw (unpinned dependency, missing validation)
- Simulate auth and token errors
- Generate and scan SBOMs (Syft/Trivy)
- Harden CI/CD (pinned actions, artifact signing)
- Integrate Vault for dynamic secrets
- Containerize and scan the image

## Why This Exists
Coming from a business system support and product owner background, I wanted a hands-on way to learn the same workflows used by modern DevSecOps and product security teams: SBOMs, CI/CD hardening, and secure delivery.

This project is my way to practice those concepts and share what I learn along the way.

Thanks for stopping by!

## License
MIT (for learning / educational purposes only)
