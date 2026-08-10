package theather

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type handler struct {
	theaterService Service
}

func NewHandler(service Service) *handler {
	return &handler{theaterService: service}
}

func (h *handler) GetTheaters(w http.ResponseWriter, r *http.Request) {
	theaters, err := h.theaterService.GetTheaters(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theaters)
}

func (h *handler) GetTheaterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID rạp không hợp lệ", http.StatusBadRequest)
		return
	}
	theater, err := h.theaterService.GetTheaterByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theater)
}

func (h *handler) CreateTheater(w http.ResponseWriter, r *http.Request) {
	var t Theater
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := h.theaterService.CreateTheater(r.Context(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *handler) GetRoomsByTheaterID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID rạp không hợp lệ", http.StatusBadRequest)
		return
	}
	rooms, err := h.theaterService.GetRoomsByTheaterID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h *handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var rm Room
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := h.theaterService.CreateRoom(r.Context(), &rm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rm)
}

type CreateSeatsRequest struct {
	Rows        int `json:"rows"`
	SeatsPerRow int `json:"seats_per_row"`
}

func (h *handler) CreateSeats(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID phòng không hợp lệ", http.StatusBadRequest)
		return
	}
	var req CreateSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := h.theaterService.CreateSeats(r.Context(), roomID, req.Rows, req.SeatsPerRow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Tạo ghế thành công"}`))
}

func (h *handler) GetSeatsByRoomID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID phòng không hợp lệ", http.StatusBadRequest)
		return
	}
	seats, err := h.theaterService.GetSeatsByRoomID(r.Context(), roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(seats)
}
