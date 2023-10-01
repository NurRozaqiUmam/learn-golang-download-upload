package db

import (
	"database/sql"
	"fmt"
	"learn-golang-download-upload/config"
)

func NewInstanceDb() *sql.DB {
	// mengambil konfigurasi database dari package config
	conf := config.GetConfig()

	// membuka koneksi baru ke database PostgreSQL dengan menggunakan informasi konfigurasi
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUsername, conf.DBPassword, conf.DBName))
	if err != nil {
		panic(err) // Jika terjadi kesalahan dalam membuka koneksi, akan menghentikan program dan menampilkan error.
	}

	// melakukan ping ke database untuk memastikan koneksi berhasil
	err = db.Ping()
	if err != nil {
		panic(err) // jika ping ke database gagal, akan menghentikan program dan menampilkan error
	}

	return db // mengembalikan instance koneksi database yang telah dikonfigurasi
}
