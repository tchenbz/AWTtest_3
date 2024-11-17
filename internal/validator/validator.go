package validator

import (
	"slices"
	//"github.com/tchenbz/AWTtest_3/internal/data"
	//"github.com/tchenbz/AWTtest_3/internal/validator"
    
)
 
type Validator struct {
    Errors map[string]string
} 

func New() *Validator {
    return &Validator {
        Errors: make(map[string]string),
    }
}

func (v *Validator) IsEmpty() bool {
    return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, message string) {
    _, exists := v.Errors[key]
    if !exists {
        v.Errors[key] = message
    }
}

func (v *Validator) Check(acceptable bool, key string, message string) {
    if !acceptable {
       v.AddError(key, message)
    }
}

func PermittedValue(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}

// func ValidateBook(v *Validator, book *data.Book) {
// 	// Ensure the title is not empty
// 	v.Check(book.Title != "", "title", "must be provided")
	
// 	// Ensure at least one author is provided
// 	v.Check(len(book.Authors) > 0, "authors", "must have at least one author")
	
// 	// ISBN should not be empty and should follow a basic pattern (e.g., length check)
// 	v.Check(book.ISBN != "", "isbn", "must be provided")
// 	v.Check(len(book.ISBN) == 13, "isbn", "must be 13 characters long") // Example for ISBN validation
	
// 	// Validate publication date (check for valid format or that it's not empty)
// 	v.Check(book.PublicationDate != "", "publication_date", "must be provided")
	
// 	// Ensure genre is not empty
// 	v.Check(book.Genre != "", "genre", "must be provided")
	
// 	// Validate average rating (e.g., should be between 0 and 5)
// 	v.Check(book.AverageRating >= 0 && book.AverageRating <= 5, "average_rating", "must be between 0 and 5")
	
// 	// Validate description (optional but could be length check)
// 	v.Check(len(book.Description) <= 1000, "description", "must not exceed 1000 characters")
// }

// // ValidateReview validates the fields of a Review struct.
// func ValidateReview(v *validator.Validator, review *Review) {
// 	// Ensure content is not empty
// 	v.Check(review.Content != "", "content", "must be provided")

// 	// Ensure author is not empty
// 	v.Check(review.Author != "", "author", "must be provided")

// 	// Ensure rating is between 1 and 5
// 	v.Check(review.Rating >= 1 && review.Rating <= 5, "rating", "must be between 1 and 5")

// 	// Ensure helpful_count is not negative
// 	v.Check(review.HelpfulCount >= 0, "helpful_count", "must not be negative")
// }