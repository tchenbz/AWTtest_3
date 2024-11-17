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

	// Define wildcard routes after
	router.HandlerFunc(http.MethodPost, "/v1/books", a.createBookHandler)        // Create a new book
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", a.displayBookHandler)   // Get a specific book by ID
	//router.HandlerFunc(http.MethodPatch, "/v1/books/:id", a.updateBookHandler)  // Update a book by ID
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", a.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", a.deleteBookHandler) // Delete a book by ID
	router.HandlerFunc(http.MethodGet, "/v1/books", a.listBooksHandler)         // List all books

	// Routes for Reviews
	router.HandlerFunc(http.MethodPost, "/v1/books/:id/reviews", a.createReviewHandler)   // Create a review for a book
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews/:review_id", a.displayReviewHandler) // Get a specific review
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id/reviews/:review_id", a.updateReviewHandler) // Update a review
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id/reviews/:review_id", a.deleteReviewHandler) // Delete a review
	router.HandlerFunc(http.MethodGet, "/v1/reviews", a.listReviewsHandler)   // List all reviews
	router.HandlerFunc(http.MethodGet, "/v1/books/:id/reviews", a.listBookReviewsHandler) // List reviews for a specific book

	// Routes for Reading Lists
	router.HandlerFunc(http.MethodPost, "/v1/readinglists", a.createReadingListHandler)        // Create a new reading list
	router.HandlerFunc(http.MethodGet, "/v1/readinglists/:id", a.displayReadingListHandler)   // Get a specific reading list by ID
	router.HandlerFunc(http.MethodPatch, "/v1/readinglists/:id", a.updateReadingListHandler)  // Update a reading list by ID
	router.HandlerFunc(http.MethodDelete, "/v1/readinglists/:id", a.deleteReadingListHandler) // Delete a reading list by ID
	router.HandlerFunc(http.MethodGet, "/v1/readinglists", a.listReadingListsHandler)         // List all reading lists

	// Routes for Users
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.getUserProfileHandler)        // Get user profile
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/lists", a.getUserReadingListsHandler) // Get user's reading lists
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/reviews", a.getUserReviewsHandler)  // Get user's reviews


	// Return the router with rate-limiting and panic recovery
	return a.recoverPanic(a.rateLimit(router))
}



