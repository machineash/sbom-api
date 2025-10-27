package api

// import the essentials.

import (
	"encoding/json" /// parse + return JSON
	"errors"        // lightweight validation messages
	"net/http"      // handle requests and responses
	"strconv"       // convert strings IDs from query params to integers
	"strings"       // basic text cleanup for validation
)

// Handler struct + constructor

// sets up a struct holding a pointer to the in-memory store (from model.go)
type Handlers struct {
	st *store
}

// NewHandlers() acts like a constructor. It builds the data store once and keps it alive across requests
// Every handler method works with the same state
func NewHandlers() *Handlers {
	return &Handlers{st: newStore()}
}

// simple validation for phase 1 (update later as needed)
func validate(c *Component) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(c.Version) == "" {
		return errors.New("version is required")
	}
	if strings.TrimSpace(c.Checksum) == "" { // CHECK
		return errors.New("checksum is required")
	}
	if strings.TrimSpace(c.Source) == "" {
		return errors.New("source is required")
	}
	return nil
}

// POST /components

// Reads JSON input. Validates it.
// Locks the store (safe for concurrency)
// Auto-increments an ID and saves the record
// Returns the created object in JSON with 201 created

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var c Component
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := validate(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.st.mu.Lock()
	c.ID = h.st.nextID
	h.st.nextID++
	h.st.components[c.ID] = c
	h.st.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// GET /components
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	h.st.mu.RLock()
	defer h.st.mu.RUnlock()

	out := make([]Component, 0, len(h.st.components))
	for _, v := range h.st.components {
		out = append(out, v)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GET /components?id=123 (simple id query for Phase 1)
// Converts query param ?id=123419 to an integer
// Fetches from map safely
// Returns 404 if not found

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	h.st.mu.RLock()
	c, ok := h.st.components[id]
	h.st.mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
}

// PUT /components?id=123 (replace)

// Replace the full record with PUT
// Validate incoming body and replace only if record exists
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var c Component
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := validate(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	h.st.mu.Lock()
	if _, exists := h.st.components[id]; !exists {
		h.st.mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.st.components[id] = c
	h.st.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
}

// PATCH /components?id=123 (partial)
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
	h.st.mu.Lock()
	c, ok := h.st.components[id]
	if !ok {
		h.st.mu.Unlock()
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
	h.st.components[id] = c
	h.st.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
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

	h.st.mu.Lock()
	defer h.st.mu.Unlock() // ensures it always unlocks even if something fails

	// check existence first
	if _, ok := h.st.components[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
	}

	// delete and respond
	delete(h.st.components, id)
	w.WriteHeader(http.StatusNoContent)
}
