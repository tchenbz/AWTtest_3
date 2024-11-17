package main

import (
	"errors"
	"log"
	"net/http"

	//"fmt"
	"github.com/tchenbz/AWTtest_3/internal/data"
	"github.com/tchenbz/AWTtest_3/internal/validator"
)

// createReadingListHandler handles POST requests for creating a new reading list.
func (a *applicationDependencies) createReadingListHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        CreatedBy   int64  `json:"created_by"`
        Status      string `json:"status"`
    }

    // Parse the request body
    err := a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Create a new ReadingList object
    readingList := &data.ReadingList{
        Name:        input.Name,
        Description: input.Description,
        CreatedBy:   input.CreatedBy,
        Status:      input.Status,
    }

    // Insert the reading list into the database
    err = a.readingListModel.Insert(readingList)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Respond with the newly created reading list
    data := envelope{"readinglist": readingList}
    err = a.writeJSON(w, http.StatusCreated, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}


// displayReadingListHandler handles GET requests to fetch a specific reading list by ID.
func (a *applicationDependencies) displayReadingListHandler(w http.ResponseWriter, r *http.Request) {
  // Extract the ID from the URL parameter
  id, err := a.readIDParam(r)
  if err != nil {
	  log.Printf("Error reading ID from URL: %v", err)  // Log error if ID is not found
	  a.notFoundResponse(w, r)
	  return
  }

  // Fetch the reading list from the database
  readingList, err := a.readingListModel.Get(id)
  if err != nil {
	  // Log the error before handling it
	  log.Printf("Error fetching reading list with ID %d: %v", id, err)
	  switch {
	  case errors.Is(err, data.ErrRecordNotFound):
		  a.notFoundResponse(w, r)
	  default:
		  a.serverErrorResponse(w, r, err)
	  }
	  return
  }

  // Return the reading list in JSON format
  data := envelope{"readinglist": readingList}
  err = a.writeJSON(w, http.StatusOK, data, nil)
  if err != nil {
	  // Log error if response writing fails
	  log.Printf("Error writing response for reading list ID %d: %v", id, err)
	  a.serverErrorResponse(w, r, err)
  }
}

// updateReadingListHandler handles PATCH requests to update a specific reading list by ID.
func (a *applicationDependencies) updateReadingListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	readingList, err := a.readingListModel.Get(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Status      *string  `json:"status"`
		Books       *[]int64 `json:"books"` // Books associated with the reading list
	}

	// Read the JSON input
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the fields that are provided
	if input.Name != nil {
		readingList.Name = *input.Name
	}
	if input.Description != nil {
		readingList.Description = *input.Description
	}
	if input.Status != nil {
		readingList.Status = *input.Status
	}
	if input.Books != nil {
		readingList.Books = *input.Books
	}

	// Validate the updated reading list
	v := validator.New()
	// Add your validation function (e.g., ValidateReadingList) if needed
	// data.ValidateReadingList(v, readingList)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the reading list in the database
	err = a.readingListModel.Update(readingList)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the updated reading list
	data := envelope{"readinglist": readingList}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// deleteReadingListHandler handles DELETE requests to delete a specific reading list by ID.
func (a *applicationDependencies) deleteReadingListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.readingListModel.Delete(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with a success message
	data := envelope{"message": "reading list successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// listReadingListsHandler handles GET requests to list all reading lists with pagination and filters.
func (a *applicationDependencies) listReadingListsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string
		Status string
		data.Filters
	}

	query := r.URL.Query()
	input.Name = a.getSingleQueryParameter(query, "name", "")
	input.Status = a.getSingleQueryParameter(query, "status", "")
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "status", "-id", "-name", "-status"}

	// Validate the filters
	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get the reading lists from the database
	readingLists, metadata, err := a.readingListModel.GetAll(input.Name, input.Status, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the list of reading lists
	data := envelope{
		"readinglists": readingLists,
		"metadata":     metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
