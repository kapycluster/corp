package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/decantor/corpy/panel/config"
	"github.com/decantor/corpy/panel/routes" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func init() {
	// Configure logger

	// Equivalent of Lshortfile
	// https://github.com/rs/zerolog?tab=readme-ov-file#customize-automatic-field-names
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()

	// Colorize output for terminal
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func main() {
	// Load config
	c, err := config.Read("config.yaml")
	if err != nil {
		log.Fatal().Msgf("failed to read config: %v", err)
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Msgf("failed to connect to database: %v", err)
	}

	// App struct
	app := &config.App{
		Config: c,
		DB:     db,
	}

	// Initialize routes
	r := routes.Init(app)

	log.Info().Msg("starting server on :4000")
	err = http.ListenAndServe("localhost:4000", r)
	if err != nil {
		log.Fatal().Msgf("failed to start server: %v", err)
	}
}
