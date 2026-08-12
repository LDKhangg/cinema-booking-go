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

// CreateMovie godoc
//
// @Summary Tạo phim mới
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body Movie true "Thông tin phim"
// @Success 201 {object} Movie
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies [post]
func (h *handler) CreateMovie(w http.ResponseWriter, req *http.Request) {
	var m Movie

	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}

	err := h.movieService.CreateMovie(req.Context(), &m)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, m)
}

// GetMovies godoc
//
// @Summary Lấy danh sách phim đang công chiếu
// @Tags movies
// @Produce json
// @Success 200 {array} Movie
// @Failure 500 {object} map[string]string
// @Router /movies [get]
func (h *handler) GetMovies(resp http.ResponseWriter, req *http.Request) {
	movies, err := h.movieService.GetMovies(req.Context())
	if err != nil {
		response.FromError(resp, err)
		return
	}
	response.JSON(resp, http.StatusOK, movies)
}

// GetMovieById godoc
//
// @Summary Lấy phim theo ID
// @Tags movies
// @Produce json
// @Param id path int true "ID phim"
// @Success 200 {object} Movie
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [get]
func (h *handler) GetMovieById(resp http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(resp, http.StatusBadRequest, "ID phim không hợp lệ")
		return
	}
	movie, err := h.movieService.GetMovieById(req.Context(), id)
	if err != nil {
		response.FromError(resp, err)
		return
	}
	response.JSON(resp, http.StatusOK, movie)
}

// UpdateMovie godoc
//
// @Summary Cập nhật thông tin phim
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "ID phim"
// @Param movie body Movie true "Thông tin phim cập nhật"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [put]
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
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Cập nhật phim thành công", nil)
}

// DeleteMovie godoc
//
// @Summary Xóa phim
// @Tags movies
// @Produce json
// @Param id path int true "ID phim"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [delete]
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
		response.FromError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Xóa phim thành công", nil)
}
