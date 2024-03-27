package main

import (
	"clubhouse/internal/database"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerImagesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	// Parse the form data, including the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	imageID := uuid.New()
	// Get a reference to the fileHeaders
	fileHeaders := r.MultipartForm.File["file"]
	for _, fileHeader := range fileHeaders {
		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			respondWithError(w, http.StatusOK, "Unable to open file")
			return
		}
		defer file.Close()

		// Create a new file in the server
		outFile, err := os.Create("data/" + imageID.String())
		if err != nil {
			respondWithError(w, http.StatusOK, "Unable to create file")
			return
		}
		defer outFile.Close()

		// Copy the uploaded file data to the newly created file
		_, err = io.Copy(outFile, file)
		if err != nil {
			respondWithError(w, http.StatusOK, "Unable to copy file")
			return
		}
	}

	fmt.Println(imageID.String())

	respondWithJSON(w, http.StatusCreated, struct{ ID uuid.UUID }{ID: imageID})
}

func (cfg *apiConfig) handlerImagesGet(w http.ResponseWriter, r *http.Request) {
	// Extract the image filename from the request URL
	filename := "data/" + chi.URLParam(r, "image")
	log.Println(filename)

	// Open the image file
	file, err := os.Open(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	// Set the appropriate content type header
	contentType := http.DetectContentType([]byte(filename))
	w.Header().Set("Content-Type", contentType)

	// Copy the image file's content to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving image", http.StatusInternalServerError)
		return
	}
}
