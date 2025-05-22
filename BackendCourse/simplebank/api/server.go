package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
// It contains the router and the store
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
// with the provided store.
func NewServer(store *db.Store) *Server {
	router := gin.Default()
	server := &Server{store: store, router: router}

	// Set up routes
	// This sets up an HTTP POST endpoint at /accounts.
	// When a POST request is made to /accounts, with the createAccount method of the server as handler.
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount) // the ':' indicates a uri (path) parameter
	router.GET("/accounts", server.listAccounts)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
