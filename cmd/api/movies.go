package main

import (
	"fmt"
	"github.com/navruz-rakhimov/greenlight/internal/data"
	"github.com/navruz-rakhimov/greenlight/internal/validator"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// If the target destination for our .Decode() is struct, then its fields must be exported.
	// key/value pairs in json are mapped to struct fields based on the struct tag names.
	// if there is no matching tag, it decodes value into field that matches the key name.
	// Other key/value pairs are ignored.
	// If you omit key/value pair in json, the field will have its zero value
	var input struct {
		Title string	`json:"title"`
		Year int32		`json:"year"`
		Runtime data.Runtime	`json:"runtime"`
		Genres []string	`json:"genres"`
	}
	// pass non-nil pointer to .Decode(v interface{})
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Runtime: 102,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
