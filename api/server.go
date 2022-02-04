package api

import (
	db "code-with-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests to our services
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new server and set up the routes
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.getAllAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
