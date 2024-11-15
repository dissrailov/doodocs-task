package handlers

import "net/http"

func InformationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
