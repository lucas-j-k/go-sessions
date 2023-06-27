package user

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"github.com/lucas-j-k/go-sessions/httpTools"
)

// Login checks incoming user details against the stored Email + PasswordHash and sets a session if valid.
// SessionID is written into a http-only cookie
func Login(service *UserService, sessionManager SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// Unmarshall and validate incoming body
		reqBody := &LoginRequestBody{}
		if err := render.Bind(r, reqBody); err != nil {
			render.Render(w, r, httpTools.ErrInvalidRequest())
			return
		}

		userRows, err := service.FindUserByEmail(FindUserByEmailPayload{
			Email: reqBody.Email,
		})

		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		// User doesn't exist
		if len(userRows) < 1 {
			render.Render(w, r, httpTools.ErrUnauthorized())
			return
		}

		user := userRows[0]

		// Verify the password hash
		validPassword := PasswordsMatch(user.PasswordHash, reqBody.Password)
		if !validPassword {
			render.Render(w, r, httpTools.ErrUnauthorized())
			return
		}

		// Generate a new session ID
		sessionId, err := uuid.NewV4()
		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		// Generate new session contents
		session := UserSession{
			UserID: user.Id,
		}

		// Persist the session to cache
		err = sessionManager.SetSession(ctx, sessionId.String(), session)
		if err != nil {
			render.Render(w, r, httpTools.ErrInternalServer())
			return
		}

		// Set the session ID in the server-side cookie
		SetAuthCookie(sessionId.String(), w)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
	}
}
