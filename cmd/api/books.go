package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tchenbz/AWTtest_3/internal/data"
	"github.com/tchenbz/AWTtest_3/internal/validator"
	"github.com/dgrijalva/jwt-go"
)


func (a *applicationDependencies) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title          string   `json:"title"`
		Authors        []string `json:"authors"` 
		ISBN           string   `json:"isbn"`
		PublicationDate string   `json:"publication_date"`
		Genre          string   `json:"genre"`
		Description    string   `json:"description"`
		AverageRating  float64  `json:"average_rating"` 
	}

	// Read and parse the JSON request body
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Create a new Book instance
	book := &data.Book{
		Title:          input.Title,
		Authors:        input.Authors,
		ISBN:           input.ISBN,
		PublicationDate: input.PublicationDate,
		Genre:          input.Genre,
		Description:    input.Description,
		AverageRating:  input.AverageRating,
	}

	// Insert the new book into the database
	err = a.bookModel.Insert(book)  
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Set the Location header for the newly created resource
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	// Respond with the created book
	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.Get(id)  
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}


func (a *applicationDependencies) updateBookHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Fetch the book from the database
    book, err := a.bookModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    // Decode the incoming JSON request
    var input struct {
        Title           *string   `json:"title"`
        Authors         []string  `json:"authors"`
        ISBN            *string   `json:"isbn"`
        PublicationDate *string   `json:"publication_date"`
        Genre           *string   `json:"genre"`
        Description     *string   `json:"description"`
        AverageRating   *float64  `json:"average_rating"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Update only the fields provided in the request
    if input.Title != nil {
        book.Title = *input.Title
    }
    if input.Authors != nil {
        book.Authors = input.Authors
    }
    if input.ISBN != nil {
        book.ISBN = *input.ISBN
    }
    if input.PublicationDate != nil {
        book.PublicationDate = *input.PublicationDate
    }
    if input.Genre != nil {
        book.Genre = *input.Genre
    }
    if input.Description != nil {
        book.Description = *input.Description
    }
    if input.AverageRating != nil {
        book.AverageRating = *input.AverageRating
    }

    // Save the updated book
    err = a.bookModel.Update(book)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{"book": book}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}


func (a *applicationDependencies) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.bookModel.Delete(id) 
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"message": "book successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string
		Author   string 
		Genre    string
		data.Filters
	}

	query := r.URL.Query()
	input.Title = a.getSingleQueryParameter(query, "title", "")
	input.Author = a.getSingleQueryParameter(query, "author", "") 
	input.Genre = a.getSingleQueryParameter(query, "genre", "")
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "author", "genre", "-id", "-title", "-author", "-genre"}

	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := a.bookModel.GetAll(input.Title, input.Author, input.Genre, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"books":    books,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) searchBooksHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the query parameters
    query := r.URL.Query()
    title := query.Get("title")
    author := query.Get("author")
    genre := query.Get("genre")

    // Parse filters (pagination, sorting)
    filters := data.Filters{
        Page:     a.getSingleIntegerParameter(query, "page", 1, validator.New()),
        PageSize: a.getSingleIntegerParameter(query, "page_size", 10, validator.New()),
        Sort:     a.getSingleQueryParameter(query, "sort", "id"),
        SortSafeList: []string{"id", "title", "author", "genre", "-id", "-title", "-author", "-genre"}, // Add sort-safe fields here
    }

    // Call the SearchBooks method with filters
    books, metadata, err := a.bookModel.SearchBooks(title, author, genre, filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Return the search results and metadata
    data := envelope{
        "books":    books,
        "metadata": metadata,
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}



func (a *applicationDependencies) authenticateUser(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	// Get the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// Pass the error as a map
		a.failedValidationResponse(w, r, map[string]string{"error": "missing authorization token"})
		return nil, fmt.Errorf("missing authorization token")
	}

	// Call parseToken to validate the token and parse the claims
	token, err := a.parseToken(tokenString)
	if err != nil {
		a.failedValidationResponse(w, r, map[string]string{"error": err.Error()})
		return nil, err
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		a.failedValidationResponse(w, r, map[string]string{"error": "invalid token"})
		return nil, fmt.Errorf("invalid token")
	}

	// Get the user ID from the claims
	userID := claims["user_id"].(float64) 
	return a.userModel.Get(int64(userID))
}
