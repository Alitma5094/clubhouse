package main

import (
	"clubhouse/internal/auth"
	"clubhouse/internal/database"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
		Name:  dbUser.Name,
	}
}

func ValidateUser(email, password, name string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("%s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = ValidateUser(params.Email, params.Password, params.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return

	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("%s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Email: params.Email, Name: params.Name, HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("%s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, DatabaseUserToUser(newUser))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	dbUsers, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, DatabaseUserToUser(dbUser))
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerUsersGetSelf(w http.ResponseWriter, _ *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, DatabaseUserToUser(user))
}
