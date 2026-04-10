package handles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
)

type CreateShareReq struct {
	Path      string `json:"path" binding:"required"`
	Label     string `json:"label"`
	ExpiresIn *int   `json:"expires_in"` // hours, nil = never
}

func CreateShare(c *gin.Context) {
	var req CreateShareReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		t := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		expiresAt = &t
	}

	st, err := db.CreateShareToken(req.Path, req.Label, expiresAt)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}

	common.SuccessResp(c, st)
}

func ListShares(c *gin.Context) {
	tokens, err := db.ListShareTokens()
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, tokens)
}

func DeleteShare(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorStrResp(c, "invalid id", 400)
		return
	}
	if err := db.DeleteShareToken(uint(id)); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c)
}

func RedeemShare(c *gin.Context) {
	token := c.Query("share_token")
	if token == "" {
		common.ErrorStrResp(c, "missing token", 400)
		return
	}

	st, err := db.GetShareToken(token)
	if err != nil {
		common.ErrorStrResp(c, "invalid token", 401)
		return
	}

	// Set cookie valid for 24h (or until expiry)
	maxAge := 86400
	if st.ExpiresAt != nil {
		maxAge = int(time.Until(*st.ExpiresAt).Seconds())
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("share_token", token, maxAge, "/", "", false, true)
	common.SuccessResp(c, gin.H{
		"path":  st.Path,
		"label": st.Label,
	})
}
