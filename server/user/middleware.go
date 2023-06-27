package user

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

type CacheMiddleware struct {
	SessionManager SessionManager
}

// SessionGuard checks for a session cookie on an incoming request. If found, it checks for a valid session under that ID in the cache.
// If session is valid, it forwards the request to the next handler in the chain. If not, it blocks with a 401 unauthorized
func (middleware *CacheMiddleware) SessionGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("user-session")

		// No cookie found, return immediately and block
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, nil)
			return
		}

		sessionId := cookie.Value

		// Attempt to pull the session data from our session cache (redis)
		sessionData, err := middleware.SessionManager.GetSession(r.Context(), sessionId)

		if err != nil {
			// User session is not present or has expired
			ClearAuthCookie(w)
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, nil)
			return
		}

		// Assume there is now a valid session. Append the userId to the context to make it available in upstream handlers
		ctx := context.WithValue(r.Context(), "userId", sessionData.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
