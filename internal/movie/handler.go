package movie

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LDKhangg/cinema-booking-go/pkg/response"
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
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}

	err := h.movieService.CreateMovie(req.Context(), &m)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, m)
}

func (h *handler) GetMovies(resp http.ResponseWriter, req *http.Request) {
	movies, err := h.movieService.GetMovies(req.Context())
	if err != nil {
		response.Error(resp, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(resp, http.StatusOK, movies)
}

func (h *handler) GetMovieById(resp http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(resp, http.StatusBadRequest, "ID phim không hợp lệ")
		return
	}
	movie, err := h.movieService.GetMovieById(req.Context(), id)
	if err != nil {
		response.Error(resp, http.StatusNotFound, "can not find movie")
		return
	}
	response.JSON(resp, http.StatusOK, movie)
}

func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.URL.Path[len("/movies/"):]
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID phim không hợp lệ")
		return
	}

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}

	m.ID = id

	if err := h.movieService.UpdateMovie(r.Context(), &m); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Cập nhật phim thành công", nil)
}

func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.URL.Path[len("/movies/"):]
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID phim không hợp lệ")
		return
	}

	if err := h.movieService.DeleteMovie(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Xóa phim thành công", nil)
}
