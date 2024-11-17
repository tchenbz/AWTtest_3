package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/tchenbz/AWTtest_3/internal/data"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
	limiter struct {
        rps float64                    
        burst int                        
        enabled bool                     
    }

}

type applicationDependencies struct {
	config        serverConfig
	logger        *slog.Logger
	bookModel     data.BookModel
	readingListModel data.ReadingListModel
	reviewModel   data.ReviewModel
	userModel     data.UserModel  
}


func main() {
	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", os.Getenv("TEST3_DB_DSN"), "PostgreSQL DSN")
	flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")
	flag.IntVar(&settings.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")
	flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Initialize the application dependencies with the necessary models
	appInstance := &applicationDependencies{
		config:    settings,
		logger:    logger,
		bookModel: data.BookModel{DB: db},        
		readingListModel: data.ReadingListModel{DB: db}, 
		reviewModel: data.ReviewModel{DB: db},    
		userModel: data.UserModel{DB: db},         
	}

	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(settings serverConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))


func (a *applicationDependencies) parseToken(tokenString string) (*jwt.Token, error) {
	// Parse the token using the jwt-go package
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil  // Return the secret key to validate the signature
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	// Return the parsed token
	return token, nil
}

// Generate an activation token
func generateActivationToken(userID int64) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecretKey)
}