package handlers

import "net/http"

type Router struct {
	Handler *HandlerApp
}

func InitRoutes(h *HandlerApp) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/archive/information", h.AnalyzeArchive)

	return mux
}
