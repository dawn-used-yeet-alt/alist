package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/alist-org/alist/v3/internal/db"
	"github.com/gin-gonic/gin"
)

const ShareTokenHeader = "X-Share-Token"
const ShareTokenCookie = "share_token"

// ShareTokenMiddleware checks if the request carries a valid share token,
// and if so, restricts path access to the token's root path.
func ShareTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from header, then query param, then cookie
		token := c.GetHeader(ShareTokenHeader)
		if token == "" {
			token = c.Query("share_token")
		}
		if token == "" {
			token, _ = c.Cookie(ShareTokenCookie)
		}

		if token == "" {
			c.Next()
			return
		}

		st, err := db.GetShareToken(token)
		if err != nil {
			// Invalid token — treat as unauthenticated
			c.Next()
			return
		}

		// Check expiry
		if st.ExpiresAt != nil && time.Now().After(*st.ExpiresAt) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Share link has expired",
			})
			return
		}

		// Store the allowed root path in context for later checks
		c.Set("share_token_path", st.Path)
		c.Set("is_share_token_request", true)

		c.Next()
	}
}

// EnforceShareTokenPath checks that the requested path is within the
// share token's allowed root. Call this in fs list/get handlers.
func EnforceShareTokenPath(c *gin.Context, requestedPath string) bool {
	val, exists := c.Get("share_token_path")
	if !exists {
		return true // not a share token request, allow normally
	}

	allowedRoot := val.(string)
	// Normalize: make sure requestedPath starts with allowedRoot
	if !strings.HasPrefix(requestedPath, allowedRoot) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Access denied: path outside share scope",
		})
		return false
	}
	return true
}
