package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/broswen/mimoto/internal/email"
	"github.com/broswen/mimoto/internal/handlers"
	"github.com/broswen/mimoto/internal/repository"
	"github.com/broswen/mimoto/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type Server struct {
	userService  user.UserService
	emailService email.EmailService
	router       chi.Router
	logger       zerolog.Logger
	tokenAuth    *jwtauth.JWTAuth
}

func New() (Server, error) {

	userRepository, err := repository.NewPostgres()
	if err != nil {
		return Server{}, fmt.Errorf("init Repository: %w", err)
	}

	emailService, err := email.NewSendGrid()
	if err != nil {
		return Server{}, fmt.Errorf("init EmailService: %w", err)
	}

	userService, err := user.New(userRepository, emailService)
	if err != nil {
		return Server{}, fmt.Errorf("init UserService: %w", err)
	}

	logger := httplog.NewLogger("mimoto", httplog.Options{
		JSON: true,
	})

	return Server{
		userService:  userService,
		emailService: emailService,
		router:       chi.NewRouter(),
		logger:       logger,
		tokenAuth:    jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil),
	}, nil
}

func (s *Server) Listen() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
}

func (s *Server) Routes() error {
	s.router.Use(httplog.RequestLogger(s.logger))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// health check
	s.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s.router.Post("/signup", handlers.SignupHandler(s.userService))
	s.router.Post("/confirm", handlers.ConfirmHandler(s.userService))
	s.router.Post("/login", handlers.LoginHandler(s.userService))
	s.router.Post("/sendreset", handlers.SendResetHandler(s.userService))
	s.router.Post("/reset", handlers.ResetHandler(s.userService))

	s.router.Group(func(r chi.Router) {
		r.Use(handlers.JWTAuthorizer)

		r.Post("/refresh", handlers.RefreshHandler(s.userService))
		r.Post("/logout", handlers.LogoutHandler(s.userService))
	})
	return nil
}
