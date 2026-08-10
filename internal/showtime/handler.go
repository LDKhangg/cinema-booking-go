package showtime

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/response"
)

type handler struct {
	showtimeService Service
}

func NewHandler(service Service) *handler {
	return &handler{showtimeService: service}
}

// CreateShowtime godoc
//
// @Summary Tạo suất chiếu mới
// @Tags showtimes
// @Accept json
// @Produce json
// @Param showtime body Showtime true "Thông tin suất chiếu"
// @Success 201 {object} Showtime
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /showtimes [post]
func (h *handler) CreateShowtime(w http.ResponseWriter, r *http.Request) {
	var st Showtime
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}
	if err := h.showtimeService.CreateShowtime(r.Context(), &st); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, st)
}

// GetShowtimeByID godoc
//
// @Summary Lấy suất chiếu theo ID
// @Tags showtimes
// @Produce json
// @Param id path int true "ID suất chiếu"
// @Success 200 {object} Showtime
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /showtimes/{id} [get]
func (h *handler) GetShowtimeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID suất chiếu không hợp lệ")
		return
	}
	st, err := h.showtimeService.GetShowtimeByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, st)
}

// GetShowtimes godoc
//
// @Summary Lấy danh sách suất chiếu theo phim hoặc rạp
// @Tags showtimes
// @Produce json
// @Param movie_id query int false "ID phim"
// @Param theater_id query int false "ID rạp"
// @Param date query string false "Ngày chiếu (YYYY-MM-DD)"
// @Success 200 {array} Showtime
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /showtimes [get]
func (h *handler) GetShowtimes(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movie_id")
	theaterIDStr := r.URL.Query().Get("theater_id")
	dateStr := r.URL.Query().Get("date")

	var targetDate time.Time
	var err error
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Định dạng ngày không hợp lệ (YYYY-MM-DD)")
			return
		}
	} else {
		targetDate = time.Now()
	}

	var showtimes []Showtime
	if movieIDStr != "" {
		movieID, err := strconv.Atoi(movieIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "ID phim không hợp lệ")
			return
		}
		showtimes, err = h.showtimeService.GetShowtimesByMovie(r.Context(), movieID, targetDate)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if theaterIDStr != "" {
		theaterID, err := strconv.Atoi(theaterIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "ID rạp không hợp lệ")
			return
		}
		showtimes, err = h.showtimeService.GetShowtimesByTheater(r.Context(), theaterID, targetDate)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		response.Error(w, http.StatusBadRequest, "Vui lòng cung cấp movie_id hoặc theater_id")
		return
	}

	response.JSON(w, http.StatusOK, showtimes)
}
