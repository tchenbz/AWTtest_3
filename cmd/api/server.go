package main

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func (a *applicationDependencies) serve() error {
    apiServer := &http.Server{
        Addr:         fmt.Sprintf(":%d", a.config.port),
        Handler:      a.routes(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError), // Ensure the logger conforms to log.Logger interface
    }

    shutdownError := make(chan error)
    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        s := <-quit
        a.logger.Info("shutting down server", "signal", s.String())

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        shutdownError <- apiServer.Shutdown(ctx)
    }()

    a.logger.Info("starting server", "address", apiServer.Addr, "environment", a.config.environment)

    err := apiServer.ListenAndServe()
    if !errors.Is(err, http.ErrServerClosed) {
        return err
    }

    err = <-shutdownError
    if err != nil {
        a.logger.Error("error during server shutdown", "error", err)
        return err
    }

    a.logger.Info("server stopped", "address", apiServer.Addr)

    return nil
}
