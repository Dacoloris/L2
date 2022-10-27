package dev11

import (
	"io/ioutil"
	"net/http"
)

type handler struct {
	service   *Service
	validator *Validator
}

func NewHandler(s *Service, v *Validator) *handler {
	return &handler{
		service:   s,
		validator: v,
	}
}

func (h *handler) Run() {
	http.HandleFunc("/create_event", MakeLoggingHandler(h.createHandler))
	http.HandleFunc("/update_event", MakeLoggingHandler(h.updateHandler))
	http.HandleFunc("/delete_event", MakeLoggingHandler(h.deleteHandler))
	http.HandleFunc("/events_for_day", MakeLoggingHandler(h.dayHandler))
	http.HandleFunc("/events_for_week", MakeLoggingHandler(h.weekHandler))
	http.HandleFunc("/events_for_month", MakeLoggingHandler(h.monthHandler))
}

func (h *handler) createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.validator.ValidateEventCreate(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.service.CreateEvent(body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(NewEditResponse("event created"))
}

func (h *handler) updateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.validator.ValidateEventID(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.service.UpdateEvent(body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.Write(NewEditResponse("event updated"))
}

func (h *handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.validator.ValidateEventID(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	err = h.service.DeleteEvent(body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.Write(NewEditResponse("event deleted"))
}

func (h *handler) dayHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	err := h.validator.ValidateDate(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	jsonResp, err := h.service.GetEventsForDay(params)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.Write(jsonResp)
}

func (h *handler) weekHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	err := h.validator.ValidateDate(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	jsonResp, err := h.service.GetEventsForWeek(params)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.Write(jsonResp)
}

func (h *handler) monthHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	err := h.validator.ValidateDate(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse(err))
		return
	}

	jsonResp, err := h.service.GetEventsForMonth(params)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(NewErrorResponse(err))
		return
	}

	w.Write(jsonResp)
}
