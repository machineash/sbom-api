package handlers

// import the essentials.

import (
	"encoding/json" /// parse + return JSON
	// lightweight validation messages
	"net/http" // handle requests and responses
	"os"
	"sbom-api/api/models"
	"strconv" // convert strings IDs from query params to integers
)

// sets up a struct holding a pointer to the in-memory store (from model.go)
// like a toolbox
type Handlers struct {
	St *models.Store
}

// NewHandlers() acts like a constructor. It builds the data store once and keeps it alive across requests
// Every handler method works with the same state
func NewHandlers() *Handlers {
	return &Handlers{St: models.NewStore()}
}

// POST /components

// Reads JSON input. Validates it.
// Locks the store (safe for concurrency)
// Auto-increments an ID and saves the record
// Returns the created object in JSON with 201 created

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Component
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := c.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.St.Mu.Lock()
	c.ID = h.St.NextID
	h.St.NextID++
	h.St.Components[c.ID] = c
	h.St.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// below change
	// _ = json.NewEncoder(w).Encode(c)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(c)

	// second pretified artifact
	b, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		_ = os.WriteFile("artifact_pretty.json", b, 0644)
	}

}

// GET /components
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	h.St.Mu.Lock()
	defer h.St.Mu.Unlock()

	out := make([]models.Component, 0, len(h.St.Components))
	for _, v := range h.St.Components {
		out = append(out, v)
	}
	w.Header().Set("Content-Type", "application/json")
	// change below
	// _ = json.NewEncoder(w).Encode(out)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)

	// second pretified artifact
	b, err := json.MarshalIndent(out, "", "  ")
	if err == nil {
		_ = os.WriteFile("artifact_pretty.json", b, 0644)
	}
}

// GET /components?id=123 (simple id query for Phase 1)
// Converts query param ?id=123 to an integer
// Fetches from map safely
// Returns 404 if not found

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	h.St.Mu.Lock()
	c, ok := h.St.Components[id]
	h.St.Mu.Unlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// below change
	// _ = json.NewEncoder(w).Encode(c)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(c)

	// second pretified artifact
	b, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		_ = os.WriteFile("artifact_pretty.json", b, 0644)
	}
}

// PUT /components?id=123 (replace)

// Replace the full record with PUT
// Validate incoming body and replace only if record exists
// test before phase 2
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var c models.Component
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := c.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	h.St.Mu.Lock()
	if _, exists := h.St.Components[id]; !exists {
		h.St.Mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.St.Components[id] = c
	h.St.Mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(c)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(c)

}

// PATCH /components?id=123 (partial)
// test before phase 2
func (h *Handlers) Patch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var in map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	h.St.Mu.Lock()
	c, ok := h.St.Components[id]
	if !ok {
		h.St.Mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if v, ok := in["name"].(string); ok && v != "" {
		c.Name = v
	}
	if v, ok := in["version"].(string); ok && v != "" {
		c.Version = v
	}
	if v, ok := in["checksum"].(string); ok && v != "" {
		c.Checksum = v
	}
	if v, ok := in["source"].(string); ok && v != "" {
		c.Source = v
	}
	if v, ok := in["license"].(string); ok {
		c.License = v
	}
	h.St.Components[id] = c
	h.St.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(c)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(c)

}

// DELETE /components?id=123
// Removes a record safely and returns 204 - No Content
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.St.Mu.Lock()
	defer h.St.Mu.Unlock() // ensures it always unlocks even if something fails

	// check existence first
	if _, ok := h.St.Components[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
  return
	}

	// delete and respond
	delete(h.St.Components, id)
	w.WriteHeader(http.StatusNoContent)
}
