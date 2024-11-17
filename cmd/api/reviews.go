package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tchenbz/AWTtest_3/internal/data"
	"github.com/tchenbz/AWTtest_3/internal/validator"
)

func (a *applicationDependencies) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := a.readIDParam(r)  
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	var input struct {
		Content   string `json:"content"`
		Author    string `json:"author"`
		Rating    int    `json:"rating"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Create a review object
	review := &data.Review{
		BookID:  bookID, 
		Content: input.Content,
		Author:  input.Author,
		Rating:  input.Rating,
	}

	log.Printf("Inserting review: %+v", review)

	// Insert the review into the database
	err = a.reviewModel.Insert(review)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the created review
	data := envelope{"review": review}
	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}


func (a *applicationDependencies) displayReviewHandler(w http.ResponseWriter, r *http.Request) {
    // Extract both the book_id and review_id from the URL path
    params := httprouter.ParamsFromContext(r.Context())
    bookID := params.ByName("id")  
    reviewID := params.ByName("review_id")  

    // Check if either ID is empty or invalid
    if bookID == "" || reviewID == "" {
        a.notFoundResponse(w, r)
        return
    }

    // Convert the extracted parameters into integers
    bookIDInt, err := strconv.ParseInt(bookID, 10, 64)
    if err != nil {
        a.badRequestResponse(w, r, fmt.Errorf("invalid book ID"))
        return
    }

    reviewIDInt, err := strconv.ParseInt(reviewID, 10, 64)
    if err != nil {
        a.badRequestResponse(w, r, fmt.Errorf("invalid review ID"))
        return
    }

    // Fetch the review from the database using the bookID and reviewID
    review, err := a.reviewModel.Get(bookIDInt, reviewIDInt)
    if err != nil {
        log.Printf("Error fetching review with book_id %d and review_id %d: %v", bookIDInt, reviewIDInt, err)
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    // Return the review in the response
    data := envelope{"review": review}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}


func (a *applicationDependencies) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
    // Extract both the book_id and review_id from the URL path
    params := httprouter.ParamsFromContext(r.Context())
    bookID := params.ByName("id")           
    reviewID := params.ByName("review_id")  

    // Ensure both bookID and reviewID are integers (as they are passed as strings)
    bookIDInt, err := strconv.ParseInt(bookID, 10, 64)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    reviewIDInt, err := strconv.ParseInt(reviewID, 10, 64)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    log.Printf("Book ID: %d, Review ID: %d", bookIDInt, reviewIDInt)  

    // Use both bookID and reviewID for fetching the review
    review, err := a.reviewModel.Get(bookIDInt, reviewIDInt) 
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)  
        default:
            a.serverErrorResponse(w, r, err)  
        }
        return
    }

    // Parse the input JSON for updates
    var input struct {
        Content      *string `json:"content"`
        Author       *string `json:"author"`
        Rating       *int    `json:"rating"`
        HelpfulCount *int    `json:"helpful_count"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)  
        return
    }

    // Update the review fields if provided
    if input.Content != nil {
        review.Content = *input.Content
    }
    if input.Author != nil {
        review.Author = *input.Author
    }
    if input.Rating != nil {
        review.Rating = *input.Rating
    }
    if input.HelpfulCount != nil {
        review.HelpfulCount = *input.HelpfulCount
    }

    // Validate the updated review 
    v := validator.New()
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Update the review in the database
    err = a.reviewModel.Update(review)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Return the updated review as a response
    data := envelope{"review": review}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}


func (a *applicationDependencies) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
    // Extract both the book_id and review_id from the URL path
    params := httprouter.ParamsFromContext(r.Context())
    bookID := params.ByName("id")          
    reviewID := params.ByName("review_id")  

    // Ensure both bookID and reviewID are integers (as they are passed as strings)
    bookIDInt, err := strconv.ParseInt(bookID, 10, 64)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    reviewIDInt, err := strconv.ParseInt(reviewID, 10, 64)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    log.Printf("Book ID: %d, Review ID: %d", bookIDInt, reviewIDInt)  

    // Use both bookID and reviewID for deleting the review
    err = a.reviewModel.Delete(bookIDInt, reviewIDInt)  
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)  
        default:
            a.serverErrorResponse(w, r, err)  
        }
        return
    }

    // Respond with a success message
    data := envelope{"message": "review successfully deleted"}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

func (a *applicationDependencies) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string
		Author  string
		Rating  int
		data.Filters
	}

	query := r.URL.Query()
	input.Content = a.getSingleQueryParameter(query, "content", "")
	input.Author = a.getSingleQueryParameter(query, "author", "")
	input.Rating = a.getSingleIntegerParameter(query, "rating", 0, validator.New())
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "rating", "helpful_count", "-id", "-rating", "-helpful_count"}

	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, metadata, err := a.reviewModel.GetAll(input.Content, input.Author, input.Rating, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"reviews":  reviews,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listBookReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL
	bookID, err := a.readIDParam(r) // This will read the book ID from the URL
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Set up the query parameters for filtering and pagination
	var input struct {
		Content string
		Author  string
		Rating  int
		data.Filters
	}

	query := r.URL.Query()
	input.Content = a.getSingleQueryParameter(query, "content", "")
	input.Author = a.getSingleQueryParameter(query, "author", "")
	input.Rating = a.getSingleIntegerParameter(query, "rating", 0, validator.New())
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "rating", "helpful_count", "-id", "-rating", "-helpful_count"}

	// Validate the query parameters (pagination, filters)
	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Fetch reviews for the book from the database using the provided filters and pagination
	reviews, metadata, err := a.reviewModel.GetAllForBook(bookID, input.Content, input.Author, input.Rating, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the reviews in JSON format, along with metadata (pagination info)
	data := envelope{
		"reviews":  reviews,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
