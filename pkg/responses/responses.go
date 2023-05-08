package responses

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) Build(s int, m string, d interface{}) {
	r.Status = s
	r.Message = m
	r.Data = d
}
func (r *Response) Respond(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	json.NewEncoder(w).Encode(r)
}
