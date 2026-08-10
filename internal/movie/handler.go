package movie

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type handler struct {
	movieService Service
}

func NewHandler(service Service) *handler {
	return &handler{
		movieService: service,
	}
}

func (h *handler) CreateMovie(w http.ResponseWriter, req *http.Request) {
	var m Movie

	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}

	err := h.movieService.CreateMovie(req.Context(), &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func (h *handler) GetMovies(resp http.ResponseWriter, req *http.Request) {
	movies, err := h.movieService.GetMovies(req.Context())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(movies)
}

func (h *handler) GetMovieById(resp http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(resp, "ID phim không hợp lệ", http.StatusBadRequest)
		return
	}
	movie, err := h.movieService.GetMovieById(req.Context(), id)
	if err != nil {
		http.Error(resp, "can not find movie", http.StatusNotFound)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(movie)
}

func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID phim không hợp lệ", http.StatusBadRequest)
		return
	}

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}

	m.ID = id

	if err := h.movieService.UpdateMovie(r.Context(), &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cập nhật phim thành công"}`))
}

func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID phim không hợp lệ", http.StatusBadRequest)
		return
	}

	if err := h.movieService.DeleteMovie(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Xóa phim thành công"}`))
}
