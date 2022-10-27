package dev11

import "encoding/json"

type ErrorResponse struct {
	Message string `json:"error"`
}

func NewErrorResponse(err error) []byte {
	resp := ErrorResponse{err.Error()}
	jsonResp, _ := json.Marshal(resp)
	return jsonResp
}

type EditResponse struct {
	Message string `json:"result"`
}

func NewEditResponse(msg string) []byte {
	resp := EditResponse{msg}
	jsonResp, _ := json.Marshal(resp)
	return jsonResp
}

type EventsResponse struct {
	Data []Event `json:"result"`
}

func NewEventsResponse(data []Event) []byte {
	resp := EventsResponse{data}
	jsonResp, _ := json.Marshal(resp)
	return jsonResp
}
