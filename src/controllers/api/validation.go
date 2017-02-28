package api

import "net/http"

func FormValidation(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
