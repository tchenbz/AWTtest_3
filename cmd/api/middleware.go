package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/time/rate"
)

func (a *applicationDependencies)recoverPanic(next http.Handler)http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func ()  {
			err := recover();
			if err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *applicationDependencies) rateLimit(next http.Handler) http.Handler {
    type client struct {
        limiter  *rate.Limiter
        lastSeen time.Time
    }
    var mu sync.Mutex    
    var clients = make(map[string]*client)
    go func() {
        for {
            time.Sleep(time.Minute)
            mu.Lock() 
            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(clients, ip)
                }
            }
            mu.Unlock() 
        }
    }()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if a.config.limiter.enabled {
            ip, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                a.serverErrorResponse(w, r, err)
                return
            }

            mu.Lock() 
            _, found := clients[ip]
            if !found {
                clients[ip] = &client{limiter: rate.NewLimiter(
                    rate.Limit(a.config.limiter.rps),
                    a.config.limiter.burst),
                }
            }
            clients[ip].lastSeen = time.Now()

            if !clients[ip].limiter.Allow() {
                mu.Unlock() 
                a.rateLimitExceededResponse(w, r)
                return
            }

            mu.Unlock()
        } 
        next.ServeHTTP(w, r)
    })

}

// AuthMiddleware is a middleware function that ensures the request is authenticated.
func (a *applicationDependencies) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            a.failedValidationResponse(w, r, map[string]string{"error": "missing authorization token"})
            return
        }

        // Parse the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("your-secret-key"), nil
        })
        if err != nil || !token.Valid {
            a.failedValidationResponse(w, r, map[string]string{"error": "invalid token"})
            return
        }

        // Extract the user ID from the token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            a.failedValidationResponse(w, r, map[string]string{"error": "invalid token claims"})
            return
        }

        userID := claims["user_id"].(float64) // User ID is in the token claims

        // Store the user ID in the request context
        ctx := context.WithValue(r.Context(), "user_id", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
