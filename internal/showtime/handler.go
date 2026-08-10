package showtime

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	showtimeService Service
}

func NewHandler(service Service) *handler {
	return &handler{showtimeService: service}
}

func (h *handler) CreateShowtime(w http.ResponseWriter, r *http.Request) {
	var st Showtime
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		http.Error(w, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := h.showtimeService.CreateShowtime(r.Context(), &st); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(st)
}

func (h *handler) GetShowtimeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID suất chiếu không hợp lệ", http.StatusBadRequest)
		return
	}
	st, err := h.showtimeService.GetShowtimeByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(st)
}

func (h *handler) GetShowtimes(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.URL.Query().Get("movie_id")
	theaterIDStr := r.URL.Query().Get("theater_id")
	dateStr := r.URL.Query().Get("date")

	var targetDate time.Time
	var err error
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Định dạng ngày không hợp lệ (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		targetDate = time.Now()
	}

	var showtimes []Showtime
	if movieIDStr != "" {
		movieID, err := strconv.Atoi(movieIDStr)
		if err != nil {
			http.Error(w, "ID phim không hợp lệ", http.StatusBadRequest)
			return
		}
		showtimes, err = h.showtimeService.GetShowtimesByMovie(r.Context(), movieID, targetDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if theaterIDStr != "" {
		theaterID, err := strconv.Atoi(theaterIDStr)
		if err != nil {
			http.Error(w, "ID rạp không hợp lệ", http.StatusBadRequest)
			return
		}
		showtimes, err = h.showtimeService.GetShowtimesByTheater(r.Context(), theaterID, targetDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Vui lòng cung cấp movie_id hoặc theater_id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(showtimes)
}
