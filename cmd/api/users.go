package main

// import (
// 	//"errors"
// 	"errors"
// 	"fmt"
// 	"net/http"

// 	"golang.org/x/crypto/bcrypt"

// 	"github.com/tchenbz/AWTtest_3/internal/data"
// 	"github.com/tchenbz/AWTtest_3/internal/validator"
// )

// // Helper function to hash password before storing it in the database
// func hashPassword(password string) (string, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hashedPassword), nil
// }

// // createUserHandler handles POST requests for creating a new user.
// func (a *applicationDependencies) createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	// Read the incoming JSON request
// 	err := a.readJSON(w, r, &input)
// 	if err != nil {
// 		a.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Hash the password before saving it to the database
// 	hashedPassword, err := hashPassword(input.Password)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Create a new User object
// 	user := &data.User{
// 		Username: input.Username,
// 		Email:    input.Email,
// 		Password: hashedPassword,  // Store the hashed password
// 	}

// 	// Validate the User object
// 	v := validator.New()
// 	// Add your validation function for user (e.g., ValidateUser)
// 	// data.ValidateUser(v, user)
// 	if !v.IsEmpty() {
// 		a.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	// Insert the User into the database
// 	err = a.userModel.Insert(user)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Set the Location header for the newly created user and respond with the user data
// 	headers := make(http.Header)
// 	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

// 	data := envelope{"user": user}
// 	err = a.writeJSON(w, http.StatusCreated, data, headers)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

// // displayUserHandler handles GET requests to fetch a specific user by ID.
// func (a *applicationDependencies) displayUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := a.readIDParam(r)
// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	user, err := a.userModel.Get(id)
// 	if err != nil {
// 		switch {
// 		case err == data.ErrRecordNotFound:
// 			a.notFoundResponse(w, r)
// 		default:
// 			a.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Return the user in JSON format
// 	data := envelope{"user": user}
// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

// // updateUserHandler handles PATCH requests to update a specific user by ID.
// func (a *applicationDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := a.readIDParam(r)
// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	user, err := a.userModel.Get(id)
// 	if err != nil {
// 		switch {
// 		case err == data.ErrRecordNotFound:
// 			a.notFoundResponse(w, r)
// 		default:
// 			a.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	var input struct {
// 		Username *string `json:"username"`
// 		Email    *string `json:"email"`
// 		Password *string `json:"password"` // Optionally update the password
// 	}

// 	// Read the JSON input
// 	err = a.readJSON(w, r, &input)
// 	if err != nil {
// 		a.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Update the fields that are provided
// 	if input.Username != nil {
// 		user.Username = *input.Username
// 	}
// 	if input.Email != nil {
// 		user.Email = *input.Email
// 	}
// 	if input.Password != nil {
// 		// Hash the new password before saving
// 		hashedPassword, err := hashPassword(*input.Password)
// 		if err != nil {
// 			a.serverErrorResponse(w, r, err)
// 			return
// 		}
// 		user.Password = hashedPassword
// 	}

// 	// Validate the updated user
// 	v := validator.New()
// 	// Add your validation function for user (e.g., ValidateUser)
// 	// data.ValidateUser(v, user)
// 	if !v.IsEmpty() {
// 		a.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	// Update the user in the database
// 	err = a.userModel.Update(user)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Respond with the updated user
// 	data := envelope{"user": user}
// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

// // deleteUserHandler handles DELETE requests to delete a specific user by ID.
// func (a *applicationDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := a.readIDParam(r)
// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	err = a.userModel.Delete(id)
// 	if err != nil {
// 		switch {
// 		case err == data.ErrRecordNotFound:
// 			a.notFoundResponse(w, r)
// 		default:
// 			a.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Respond with a success message
// 	data := envelope{"message": "user successfully deleted"}
// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

// // listUsersHandler handles GET requests to list all users with pagination and filters.
// func (a *applicationDependencies) listUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username string
// 		Email    string
// 		data.Filters
// 	}

// 	query := r.URL.Query()
// 	input.Username = a.getSingleQueryParameter(query, "username", "")
// 	input.Email = a.getSingleQueryParameter(query, "email", "")
// 	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
// 	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
// 	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
// 	input.Filters.SortSafeList = []string{"id", "username", "email", "-id", "-username", "-email"}

// 	// Validate the filters
// 	v := validator.New()
// 	data.ValidateFilters(v, input.Filters)
// 	if !v.IsEmpty() {
// 		a.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	// Get the users from the database
// 	users, metadata, err := a.userModel.GetAll(input.Username, input.Email, input.Filters)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Respond with the list of users
// 	data := envelope{
// 		"users":    users,
// 		"metadata": metadata,
// 	}
// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

// // Handler for fetching user profile
// func (a *applicationDependencies) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
//     id, err := a.readIDParam(r)
//     if err != nil {
//         a.notFoundResponse(w, r)
//         return
//     }

//     user, err := a.userModel.Get(id)
//     if err != nil {
//         switch {
//         case errors.Is(err, data.ErrRecordNotFound):
//             a.notFoundResponse(w, r)
//         default:
//             a.serverErrorResponse(w, r, err)
//         }
//         return
//     }

//     // Return the user profile
//     data := envelope{"user": user}
//     err = a.writeJSON(w, http.StatusOK, data, nil)
//     if err != nil {
//         a.serverErrorResponse(w, r, err)
//     }
// }

// // Handler for fetching user's reading lists
// func (a *applicationDependencies) getUserReadingListsHandler(w http.ResponseWriter, r *http.Request) {
//     id, err := a.readIDParam(r)
//     if err != nil {
//         a.notFoundResponse(w, r)
//         return
//     }

//     // Parse query parameters for filters (pagination, sorting)
//     query := r.URL.Query()
//     filters := data.Filters{
//         Page:     a.getSingleIntegerParameter(query, "page", 1, validator.New()),
//         PageSize: a.getSingleIntegerParameter(query, "page_size", 10, validator.New()),
//         Sort:     a.getSingleQueryParameter(query, "sort", "id"),
//     }

//     // Call GetAllByUser with both userID and filters
//     readingLists, metadata, err := a.readingListModel.GetAllByUser(id, filters)
//     if err != nil {
//         a.serverErrorResponse(w, r, err)
//         return
//     }

//     // Return the user's reading lists along with pagination metadata
//     data := envelope{
//         "readinglists": readingLists,
//         "metadata":     metadata,
//     }
//     err = a.writeJSON(w, http.StatusOK, data, nil)
//     if err != nil {
//         a.serverErrorResponse(w, r, err)
//     }
// }


// func (a *applicationDependencies) getUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Retrieve the user ID from the URL parameter
// 	id, err := a.readIDParam(r)
// 	if err != nil {
// 		a.notFoundResponse(w, r)
// 		return
// 	}

// 	// Parse the query parameters to get filters (pagination, sorting)
// 	query := r.URL.Query()
// 	filters := data.Filters{
// 		Page:     a.getSingleIntegerParameter(query, "page", 1, validator.New()),
// 		PageSize: a.getSingleIntegerParameter(query, "page_size", 10, validator.New()),
// 		Sort:     a.getSingleQueryParameter(query, "sort", "id"),
// 	}

// 	// Fetch the reviews for the user using the GetAllByUser method
// 	reviews, metadata, err := a.reviewModel.GetAllByUser(id, filters)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Send the response with the list of reviews and metadata for pagination
// 	data := envelope{
// 		"reviews": reviews,
// 		"metadata": metadata,
// 	}
// 	err = a.writeJSON(w, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(w, r, err)
// 	}
// }

import (
    "fmt"
    "net/http"
    "errors"
    "time"
    "github.com/tchenbz/AWTtest_3/internal/data"
    "github.com/tchenbz/AWTtest_3/internal/validator"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
)

// User-related handler and helper functions

// Helper function to hash password before storing it in the database
func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// Function to generate JWT token for the user
func generateJWTToken(user *data.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := []byte("your-secret-key") // Replace with your secret key
    signedToken, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}

// Handler for creating a new user
func (a *applicationDependencies) createUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Read the incoming JSON request
    err := a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Hash the password before saving it to the database
    hashedPassword, err := hashPassword(input.Password)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,  // Store the hashed password
		EmailVerified: false,      // Default to false
	}

    // Validate the User object
    v := validator.New()
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Insert the User into the database
    err = a.userModel.Insert(user)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Set the Location header for the newly created user and respond with the user data
    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

    data := envelope{"user": user}
    err = a.writeJSON(w, http.StatusCreated, data, headers)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}



// Handler for retrieving a specific user by ID
func (a *applicationDependencies) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    user, err := a.userModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{"user": user}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// Handler for updating a user's details
func (a *applicationDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    user, err := a.userModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    var input struct {
        Username *string `json:"username"`
        Email    *string `json:"email"`
        Password *string `json:"password"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    if input.Username != nil {
        user.Username = *input.Username
    }
    if input.Email != nil {
        user.Email = *input.Email
    }
    if input.Password != nil {
        hashedPassword, err := hashPassword(*input.Password)
        if err != nil {
            a.serverErrorResponse(w, r, err)
            return
        }
        user.Password = hashedPassword
    }

    err = a.userModel.Update(user)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{"user": user}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// Handler for deleting a user by ID
func (a *applicationDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    err = a.userModel.Delete(id)
    if err != nil {
        switch {
        case err == data.ErrRecordNotFound:
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{"message": "user successfully deleted"}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// Handler for user login (authentication)
func (a *applicationDependencies) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Read the login data from the request
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Retrieve the user by email
	user, err := a.userModel.GetByEmail(input.Email)
	if err != nil {
		// If user is not found, return 404 or bad request
		a.notFoundResponse(w, r)
		return
	}

	// Compare the provided password with the stored password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// If password doesn't match, return 401 Unauthorized
		a.notFoundResponse(w, r)
		return
	}

	// Generate a JWT token for the user
	token, err := generateJWTToken(user)
	if err != nil {
		// If token generation fails, return server error
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the generated token
	data := envelope{"token": token}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Handler to activate user via email (simulated)
func (a *applicationDependencies) activateUserHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    user, err := a.userModel.Get(id)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Simulating email activation, set EmailVerified to true
    user.EmailVerified = true
    err = a.userModel.Update(user)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{"message": "user activated"}
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// Handler for fetching user's reading lists
func (a *applicationDependencies) getUserReadingListsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the URL
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Parse query parameters for filters (pagination, sorting)
    query := r.URL.Query()
    filters := data.Filters{
        Page:     a.getSingleIntegerParameter(query, "page", 1, validator.New()),
        PageSize: a.getSingleIntegerParameter(query, "page_size", 10, validator.New()),
        Sort:     a.getSingleQueryParameter(query, "sort", "id"),
    }

	// Example of SortSafeList containing the valid fields
    filters.SortSafeList = []string{"id", "name", "status", "-id", "-name", "-status"}


    // Call GetAllByUser with both userID and filters
    readingLists, metadata, err := a.readingListModel.GetAllByUser(id, filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Return the user's reading lists along with pagination metadata
    data := envelope{
        "readinglists": readingLists,
        "metadata":     metadata,
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// Handler for fetching user's reviews
func (a *applicationDependencies) getUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the URL
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Parse query parameters for filters (pagination, sorting)
    query := r.URL.Query()
    filters := data.Filters{
        Page:     a.getSingleIntegerParameter(query, "page", 1, validator.New()),
        PageSize: a.getSingleIntegerParameter(query, "page_size", 10, validator.New()),
        Sort:     a.getSingleQueryParameter(query, "sort", "id"),
    }

	filters.SortSafeList = []string{"id", "name", "status", "-id", "-name", "-status"}

    // Fetch the reviews for the user using the GetAllByUser method
    reviews, metadata, err := a.reviewModel.GetAllByUser(id, filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    // Return the user's reviews along with pagination metadata
    data := envelope{
        "reviews": reviews,
        "metadata": metadata,
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

