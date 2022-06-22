package api

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"
	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	seesion, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if seesion.IsBlocked {
		err := fmt.Errorf("blocked ssession")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if seesion.UserName != refreshPayload.Username {
		err := fmt.Errorf("incorrect ssesion user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.RefreshToken != seesion.RefreshToken {
		err := fmt.Errorf("mismatch RefreshToken")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if time.Now().After(seesion.ExpiresAt) {
		err := fmt.Errorf("expired token")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsq := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsq)
}
