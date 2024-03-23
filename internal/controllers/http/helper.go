package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, data any) {
	payload, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(payload))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(payload); err != nil {
		slog.Warn(err.Error())
	}
}
