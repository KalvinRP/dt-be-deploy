package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadFile(point http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			file, _, err := r.FormFile("image")

			if err != nil {
				fmt.Println(err)
				json.NewEncoder(w).Encode("Error Retrieving the File")
				return
			}
			defer file.Close()
			const MAX_UPLOAD_SIZE = 32 << 20 // 32MB

			r.ParseMultipartForm(MAX_UPLOAD_SIZE)
			if r.ContentLength > MAX_UPLOAD_SIZE {
				w.WriteHeader(http.StatusBadRequest)
				response := Result{Code: http.StatusBadRequest, Message: "Files sizes are too big!"}
				json.NewEncoder(w).Encode(response)
				return
			}

			tempFile, err := os.CreateTemp("uploads", "image-*.jpeg")
			if err != nil {
				fmt.Println(err)
				fmt.Println("path upload error")
				json.NewEncoder(w).Encode(err)
				return
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			data := tempFile.Name()
			filename := data[8:] // split uploads/

			ctx := context.WithValue(r.Context(), "dataFile", filename)
			point.ServeHTTP(w, r.WithContext(ctx))
		})
}
