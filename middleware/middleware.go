package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Function upload file
func UploadFile(next http.HandlerFunc) http.HandlerFunc { // parameter dan return handhlerfunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("image") // Ambil value dari image dan simpan kedalam 3 variabel
		// Jika ada error maka kembalikan error
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieve the file.")
			return
		}
		defer file.Close()

		fmt.Printf("Uploded file : %+v\n", handler.Filename)

		tempFile, err := ioutil.TempFile("uploads", "image-*"+handler.Filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		filename := data[8:]

		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
