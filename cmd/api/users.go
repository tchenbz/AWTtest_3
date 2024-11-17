package main

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

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func generateJWTToken(user *data.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := []byte("your-secret-key") 
    signedToken, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}

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
		Password: hashedPassword,  
		EmailVerified: false,      // Default to false
	}

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

func (a *applicationDependencies) getUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the URL
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    // Parse query parameters for filters (pagination and sorting)
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

