package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	// Define non-wildcard routes first
	//router.HandlerFunc(http.MethodGet, "/v1/books/search", a.searchBooksHandler) // Search books by title/author/genre

	// Routes for Books
	router.HandlerFunc(http.MethodPost, "/v1/books", a.createBookHandler)        
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", a.displayBookHandler)   
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", a.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", a.deleteBookHandler) 
	router.HandlerFunc(http.MethodGet, "/v1/books", a.listBooksHandler)        

	// Routes for Reviews
	router.HandlerFunc(http.MethodPost, "/v1/books/:id/reviews", a.createReviewHandler)   
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews/:review_id", a.displayReviewHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id/reviews/:review_id", a.updateReviewHandler)  
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id/reviews/:review_id", a.deleteReviewHandler) 
	router.HandlerFunc(http.MethodGet, "/v1/reviews", a.listReviewsHandler)  
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews", a.listBookReviewsHandler) 

	// Routes for Reading Lists
	router.HandlerFunc(http.MethodPost, "/v1/readinglists", a.createReadingListHandler)        
	router.HandlerFunc(http.MethodGet, "/v1/readinglists/:id", a.displayReadingListHandler)   
	router.HandlerFunc(http.MethodPatch, "/v1/readinglists/:id", a.updateReadingListHandler)  
	router.HandlerFunc(http.MethodDelete, "/v1/readinglists/:id", a.deleteReadingListHandler) 
	router.HandlerFunc(http.MethodGet, "/v1/readinglists", a.listReadingListsHandler)        

	// Routes for Users
	router.HandlerFunc(http.MethodPost, "/v1/users", a.createUserHandler)  
	router.HandlerFunc(http.MethodPost, "/v1/login", a.loginUserHandler)  
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.getUserProfileHandler)        
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/lists", a.getUserReadingListsHandler) 
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/reviews", a.getUserReviewsHandler)  

	return a.recoverPanic(a.rateLimit(router))
}



