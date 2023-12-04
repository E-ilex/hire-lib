package main

import (
	"database/sql"
	"hire-test-lib/config"
	"hire-test-lib/question"
	"hire-test-lib/utils"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

const (
	component        = "component"
	messageFieldName = "message"
	timeFieldName    = "timestamp"
	timeFieldFormat  = time.RFC3339Nano
)

type server struct {
	response        *utils.Response
	config          *config.Configuration
	logging         zerolog.Logger
	router          *chi.Mux
	database        *sql.DB
	questionService question.QuestionService
	questionHandler question.QuestionHandler
}

func (s *server) loadConfiguration() {
	s.config = config.LoadConfiguration()
}

func (s *server) loadLogging() {
	zerolog.MessageFieldName = messageFieldName
	zerolog.TimestampFieldName = timeFieldName
	zerolog.TimeFieldFormat = timeFieldFormat

	s.logging = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
}

func sqlliteConnection(cfg *config.Configuration) (*sql.DB, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DBDatasource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func (s *server) loadDatabase() {
	connection, err := sqlliteConnection(s.config)
	if err != nil {
		s.logging.Fatal().Err(err).Str(component, "database").Msg("could not connect to db")
		return
	}

	s.database = connection
	s.database.SetConnMaxLifetime(60 * time.Second)
	s.database.SetMaxOpenConns(1)
}

func (s *server) loadService() {
	s.questionService = question.NewQuestionService(&s.logging, s.database)
}

func (s *server) loadRouter() {
	s.router = chi.NewRouter()
}

func (s *server) loadMiddleware() {
	s.router.Use(middleware.Compress(5))
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
}

func (s *server) loadRoutes() {
	s.questionHandler = question.NewQuestionHandler(&s.logging, s.response, s.questionService)

	s.router.Route("/v1/question", func(router chi.Router) {
		s.router.Get("/", s.questionHandler.GetQuestions)
		s.router.Post("/", s.questionHandler.CreateQuestion)
		s.router.Put("/{questionID}", s.questionHandler.UpdateQuestion)
		s.router.Delete("/{questionID}", s.questionHandler.DeleteQuestion)
	})
}

func (s *server) run() {
	s.logging.Info().Str(component, "server").Msg("starting service")
	if err := http.ListenAndServe(s.config.Port, s.router); err != nil {
		s.logging.Fatal().Err(err).Str(component, "server").Msg("could not run server")
		return
	}
}

func main() {
	s := server{}
	s.loadConfiguration()
	s.loadLogging()
	s.loadDatabase()
	s.loadService()
	s.loadRouter()
	s.loadMiddleware()
	s.loadRoutes()
	s.run()
}
