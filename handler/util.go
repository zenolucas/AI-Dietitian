package handler

import (
	"AI-Dietitian/types"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	// do we have a user key? and is that key an AuthenticatedUser?
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		// no? return an empty AuthenticatedUser
		return types.AuthenticatedUser{}
	}
	// yes, return an AuthenticatedUser
	return user
}

// as a way to make a centralized error handling system
func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}

