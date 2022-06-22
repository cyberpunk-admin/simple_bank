package gapi

import (
	"fmt"

	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/pb"
	"github.com/simplebank/token"
	"github.com/simplebank/util"
)

// Server serves GRPC requests for banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config util.Config
	store	db.Store
	tokenMaker token.Maker
}

// NewServer creates a new GRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot craete tokenMaker %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}


