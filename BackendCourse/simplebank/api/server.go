package api

import (
	"fmt"

	db "github.com/ofer-sin/Courses/BackendCourse/simplebank/db/sqlc"
	"github.com/ofer-sin/Courses/BackendCourse/simplebank/token"
	"github.com/ofer-sin/Courses/BackendCourse/simplebank/util"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
// It contains the router and the store
type Server struct {
	config util.Config
	// The store is an interface that defines methods for interacting with the database.
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
// with the provided store.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}
	server.setupROuter()
	return server, nil
}

func (server *Server) setupROuter() {
	// Create a new Gin router instance
	// The router is responsible for routing incoming HTTP requests to the appropriate handler functions
	// and managing the server's routes.
	router := gin.Default()

	// Set up routes
	// When a request is made, the indicated handler method of the server is called
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount) // the ':' indicates a uri (path) parameter
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
