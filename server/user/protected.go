package user

import (
	"net/http"

	"github.com/go-chi/render"
)

// This is a demo protected endpoint, to check the userId is available after being retrieved from the session
func Protected(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, _ := ctx.Value("userId").(int64)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, userId)
	}
}
