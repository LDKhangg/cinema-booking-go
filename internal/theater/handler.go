package theater

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LDKhangg/cinema-booking-go/pkg/response"
)

type handler struct {
	theaterService Service
}

func NewHandler(service Service) *handler {
	return &handler{theaterService: service}
}

// GetTheaters godoc
//
// @Summary Lấy danh sách rạp chiếu
// @Tags theaters
// @Produce json
// @Success 200 {array} Theater
// @Failure 500 {object} map[string]string
// @Router /theaters [get]
func (h *handler) GetTheaters(w http.ResponseWriter, r *http.Request) {
	theaters, err := h.theaterService.GetTheaters(r.Context())
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, theaters)
}

// GetTheaterByID godoc
//
// @Summary Lấy rạp chiếu theo ID
// @Tags theaters
// @Produce json
// @Param id path int true "ID rạp"
// @Success 200 {object} Theater
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /theaters/{id} [get]
func (h *handler) GetTheaterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID rạp không hợp lệ")
		return
	}
	theater, err := h.theaterService.GetTheaterByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, theater)
}

// CreateTheater godoc
//
// @Summary Tạo rạp chiếu mới
// @Tags theaters
// @Accept json
// @Produce json
// @Param theater body Theater true "Thông tin rạp"
// @Success 201 {object} Theater
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /theaters [post]
func (h *handler) CreateTheater(w http.ResponseWriter, r *http.Request) {
	var t Theater
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}
	if err := h.theaterService.CreateTheater(r.Context(), &t); err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, t)
}

// GetRoomsByTheaterID godoc
//
// @Summary Lấy danh sách phòng chiếu của rạp
// @Tags rooms
// @Produce json
// @Param id path int true "ID rạp"
// @Success 200 {array} Room
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /theaters/{id}/rooms [get]
func (h *handler) GetRoomsByTheaterID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID rạp không hợp lệ")
		return
	}
	rooms, err := h.theaterService.GetRoomsByTheaterID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, rooms)
}

// CreateRoom godoc
//
// @Summary Tạo phòng chiếu mới
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body Room true "Thông tin phòng chiếu"
// @Success 201 {object} Room
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /theaters/rooms [post]
func (h *handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var rm Room
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}
	if err := h.theaterService.CreateRoom(r.Context(), &rm); err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, rm)
}

type CreateSeatsRequest struct {
	Rows        int `json:"rows"`
	SeatsPerRow int `json:"seats_per_row"`
}

// CreateSeats godoc
//
// @Summary Tạo danh sách ghế cho phòng chiếu
// @Tags seats
// @Accept json
// @Produce json
// @Param id path int true "ID phòng chiếu"
// @Param request body CreateSeatsRequest true "Số hàng và số ghế mỗi hàng"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rooms/{id}/seats [post]
func (h *handler) CreateSeats(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID phòng không hợp lệ")
		return
	}
	var req CreateSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ")
		return
	}
	if err := h.theaterService.CreateSeats(r.Context(), roomID, req.Rows, req.SeatsPerRow); err != nil {
		response.FromError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "Tạo ghế thành công", nil)
}

// GetSeatsByRoomID godoc
//
// @Summary Lấy danh sách ghế của phòng chiếu
// @Tags seats
// @Produce json
// @Param id path int true "ID phòng chiếu"
// @Success 200 {array} Seat
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rooms/{id}/seats [get]
func (h *handler) GetSeatsByRoomID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID phòng không hợp lệ")
		return
	}
	seats, err := h.theaterService.GetSeatsByRoomID(r.Context(), roomID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, seats)
}
