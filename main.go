package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", Upload)
	e.GET("/download/:filename", Download)

	e.Logger.Fatal(e.Start(":1323"))
}

func Upload(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Konfigurasi Firebase
	opt := option.WithCredentialsFile("D:/first-project/service.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	// Dapatkan client Cloud Storage
	client, err := app.Storage(context.Background())
	if err != nil {
		return err
	}

	// lokasi di Firebase Cloud Storage
	bucket, err := client.Bucket("first-project-af98e.appspot.com")
	if err != nil {
		// Handle kesalahan
		log.Fatalf("Error getting bucket: %v", err)
	}

	ctx := context.Background()
	dst := fmt.Sprintf("image/%s", file.Filename)
	wc := bucket.Object(dst).NewWriter(ctx)
	if _, err = io.Copy(wc, src); err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		// Handle kesalahan
		log.Fatalf("Error error wc: %v", err)
	}
	//return c.HTML(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields", file.Filename))
	// File telah berhasil diunggah
	// Kembalikan URL gambar sebagai respons JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "File uploaded successfully",
	})
}

func Download(c echo.Context) error {
	filename := c.Param("filename")

	// Konfigurasi Firebase
	opt := option.WithCredentialsFile("D:/first-project/service.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	// Dapatkan client Cloud Storage
	client, err := app.Storage(context.Background())
	if err != nil {
		return err
	}

	// Lokasi di Firebase Cloud Storage
	bucket, err := client.Bucket("first-project-af98e.appspot.com")
	if err != nil {
		log.Fatalf("Error getting bucket: %v", err)
	}

	// Membuat objek untuk file gambar
	ctx := context.Background()
	obj := bucket.Object("image/" + filename)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Mengatur header untuk memaksa browser mengunduh file
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set("Content-Type", "image/jpeg") // Mengatur tipe konten sebagai aplikasi biner

	// Salin data dari Firebase Cloud Storage ke response HTTP sebagai data biner
	_, err = io.Copy(c.Response().Writer, reader)
	if err != nil {
		return err
	}

	// Mengembalikan respons JSON dengan pesan sukses
	response := map[string]interface{}{
		"code":    200,
		"message": "success",
	}

	return c.JSON(http.StatusOK, response)
}
