package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(t.Run())
}

