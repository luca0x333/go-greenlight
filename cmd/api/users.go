package main

import (
	"errors"
	"github.com/luca0x333/go-greenlight/internal/data"
	"github.com/luca0x333/go-greenlight/internal/validator"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Anonymous struct to hold the data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use Set method to generate and store the hashed and plaintext password into the struct.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Validate the struct.
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Use the background helper to execute an anonymous function that sends the welcome email.
	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			// We use PrintError because by the time we encounter the errors,
			// the client will probably have already been sent a 202 Accepted response by our writeJSON() helper.
			app.logger.PrintError(err, nil)
		}
	})

	// Write a JSON response with 202 StatusAccepted status code.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
