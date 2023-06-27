package user

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lucas-j-k/go-sessions/httpTools"
)

// Signup inserts a new user record into the DB
func Signup(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody := &SignupJSONBody{}
		if err := render.Bind(r, reqBody); err != nil {
			render.Render(w, r, httpTools.ErrInvalidRequest())
			return
		}

		// Search for existing user with same email, and reject if found
		existingRows, err := service.FindUserByEmail(FindUserByEmailPayload{
			Email: reqBody.Email,
		})

		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		if len(existingRows) > 0 {
			render.Render(w, r, httpTools.ErrInvalidRequest())
			return
		}

		hashedPassword, err := HashPassword(reqBody.Password)
		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		insertedId, err := service.InsertUser(InsertUserPayload{
			FirstName:    reqBody.FirstName,
			LastName:     reqBody.LastName,
			Email:        reqBody.Email,
			PasswordHash: hashedPassword,
		})

		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, insertedId)
	}
}
