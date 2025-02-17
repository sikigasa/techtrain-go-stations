package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	result, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *result}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	result, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: result}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	result, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *result}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) CreateTodo(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, *model.CreateTODOResponse) {
	var request model.CreateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return w, http.StatusInternalServerError, nil
	}

	if request.Subject == "" {
		return w, http.StatusBadRequest, nil
	}

	response, err := h.Create(context.Background(), &request)
	if err != nil {
		return w, http.StatusInternalServerError, nil
	}
	return w, http.StatusOK, response
}

func (h *TODOHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, *model.UpdateTODOResponse) {
	var request model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return w, http.StatusInternalServerError, nil
	}

	if request.ID == 0 || request.Subject == "" {
		return w, http.StatusBadRequest, nil
	}

	response, err := h.Update(context.Background(), &request)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(model.ErrNotFound{}) {
			return w, http.StatusNotFound, nil
		}
		return w, http.StatusInternalServerError, nil
	}
	return w, http.StatusOK, response
}

func (h *TODOHandler) ReadTodo(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, *model.ReadTODOResponse) {
	var request model.ReadTODORequest
	prevID := r.URL.Query().Get("prev_id")
	size := r.URL.Query().Get("size")

	if prevID != "" {
		prev, err := strconv.Atoi(prevID)
		if err != nil {
			return w, http.StatusBadRequest, nil
		}
		request.PrevID = int64(prev)
	}
	if size != "" {
		s, err := strconv.Atoi(size)
		if err != nil {
			return w, http.StatusBadRequest, nil
		}
		request.Size = int64(s)
	} else {
		request.Size = 5
	}

	response, err := h.Read(context.Background(), &request)
	if err != nil {
		return w, http.StatusInternalServerError, nil
	}

	return w, http.StatusOK, response
}
