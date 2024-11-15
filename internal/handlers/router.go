package handlers

import "net/http"

func InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/archive/information", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			InformationHandler(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
