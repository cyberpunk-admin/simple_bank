package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	UserName         string    `json:"user_name"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	HashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		UserName:    req.UserName,
		HashPassword:  HashPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok :=  err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsq := createUserResponse{
		UserName: user.UserName,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsq)
}

