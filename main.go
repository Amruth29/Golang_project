package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/process", processHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form containing image files
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Get the slice of uploaded files
	files := r.MultipartForm.File["images"]

	// Process each image file concurrently
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()
			processImage(fileHeader)
		}(file)
	}
	wg.Wait()

	fmt.Fprint(w, `{"message": "All images processed successfully"}`)
}

func processImage(fileHeader *multipart.FileHeader) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening uploaded file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the file content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file content: %v\n", err)
		return
	}

	// Decode the image
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		log.Printf("Error decoding image: %v\n", err)
		return
	}

	// Resize the image
	resizedImg := resize.Resize(300, 0, img, resize.Lanczos3) // Resize to width 300 pixels (maintain aspect ratio)

	// Save the resized image
	outputDir := "C:/Users/Anirudh R/Desktop/Go_final_78/output"
	outputPath := filepath.Join(outputDir, fileHeader.Filename)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	// Encode and save the resized image as JPEG
	if err := jpeg.Encode(outputFile, resizedImg, nil); err != nil {
		log.Printf("Error encoding image: %v\n", err)
		return
	}

	log.Printf("Resized image saved at: %s\n", outputPath)
}
