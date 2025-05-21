package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

// This is one API handler function that handles the creation of a new account.
// It is called when a POST request is made to the /accounts endpoint.
// The handler was set by the router in the NewServer function by calling:
// router.POST("/accounts", server.createAccount)
func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	// Bind JSON request to the CreateAccountRequest struct
	// and validate the input
	// If the binding fails, return a 400 Bad Request response
	// with the error message
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// translate the api request to the db request
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// Call the store to create the account in the database
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// This is one API handler function that handles the retrieval of an account.
// It is called when a GET request is made to the /accounts/:id endpoint.
// The :id part of the URL is a path parameter that is passed to the handler.
// The handler was set by the router in the NewServer function by calling:
// router.GET("/accounts/:id", server.getAccount)
func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	// Bind URI request to the GetAccountRequest struct
	// and validate the input
	// If the binding fails, return a 400 Bad Request response
	// with the error message
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Call the store to get the account from the database
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the account is not found, return a 404 Not Found response
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// If there is any other error, return a 500 Internal Server Error response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
