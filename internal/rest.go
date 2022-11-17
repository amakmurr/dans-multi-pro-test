package internal

import (
	"embed"
	"net/http"

	"github.com/amakmurr/dans-multi-pro-test/internal/repository"
	"github.com/amakmurr/dans-multi-pro-test/pkg/dans"
	"github.com/amakmurr/dans-multi-pro-test/pkg/jwt"
	v1 "github.com/amakmurr/dans-multi-pro-test/pkg/openapi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	srv *http.Server
}

//go:embed swagger
var swaggerFiles embed.FS

func NewServer(cfg *Config) (*Server, error) {
	// dependencies
	db, err := initDB(cfg.Database.Source)
	if err != nil {
		return nil, err
	}
	userRepo := repository.NewUserRepository(db)
	dansClient, err := dans.NewClient(cfg.Dans.BaseURL)
	if err != nil {
		return nil, err
	}
	jwtClient := jwt.NewClient(cfg.JWT.Secret)
	jwtAuth := NewJWTAuthentication(jwtClient)
	endpoint := NewEndpoint(userRepo, jwtClient, dansClient)

	// route
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/login", endpoint.login)
	r.Route("/jobs", func(r chi.Router) {
		r.Use(jwtAuth.Authenticator)
		r.Get("/", endpoint.getJobList)
		r.Get("/{id}", endpoint.getJobDetail)
	})

	// swagger
	r.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.FS(swaggerFiles)),
	))
	r.Get("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		s, e := v1.GetSwagger()
		if e != nil {
			return
		}
		b, _ := s.MarshalJSON()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})

	// server
	return &Server{
		srv: &http.Server{
			Addr:    ":8000",
			Handler: r,
		},
	}, nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
