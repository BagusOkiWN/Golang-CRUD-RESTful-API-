package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// DB instance to interact with the database
var DB *gorm.DB

func init() {
	// Open a database connection using GORM
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/acrud?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")

	// Run auto migration only during development to create the 'pegawai' table
	err = DB.AutoMigrate(&Pegawai{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define API routes
	e.GET("/pegawai", GetAllData)
	e.GET("/pegawai/:id", GetPegawaiByID)
	e.POST("/pegawai", CreatePegawai)
	e.PUT("/pegawai/:id", UpdatePegawai)
	e.DELETE("/pegawai/:id", DeletePegawai)

	// Start the server
	e.Start(":1324")
}

// Pegawai struct represents the employee data model
type Pegawai struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	NamaPegawai    string `json:"nama_pegawai"`
	NIK            string `json:"nik"`
	JenisPegawaiID int    `json:"jenis_pegawai_id"`
	Unit           string `json:"unit"`
	SubUnit        string `json:"sub_unit"`
	PendidikanID   int    `json:"pendidikan_id"`
	TanggalLahir   string `json:"tgl_lahir"`
	TempatLahir    string `json:"tpt_lahir"`
	JenisKelaminID int    `json:"jenkel_id"`
	AgamaID        int    `json:"agama_id"`
	Gambar         string `json:"gambar"`
}

// PegawaiRequest represents the request payload for creating or updating Pegawai
type PegawaiRequest struct {
	ID             uint   `json:"id"`
	NamaPegawai    string `json:"nama_pegawai"`
	NIK            string `json:"nik"`
	JenisPegawaiID int    `json:"jenis_pegawai_id"`
	Unit           string `json:"unit"`
	SubUnit        string `json:"sub_unit"`
	PendidikanID   int    `json:"pendidikan_id"`
	TanggalLahir   string `json:"tgl_lahir"`
	TempatLahir    string `json:"tpt_lahir"`
	JenisKelaminID int    `json:"jenkel_id"`
	AgamaID        int    `json:"agama_id"`
	Gambar         string `json:"gambar"`
}

func GetAllData(c echo.Context) error {
	// Retrieve all employee data
	var pegawaiList []Pegawai
	if err := DB.Table("pegawai").Find(&pegawaiList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pegawai not found"})
		}
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Return the Pegawai data as JSON
	return c.JSON(http.StatusOK, pegawaiList)
}

func GetPegawaiByID(c echo.Context) error {
	// Get employee ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Retrieve Pegawai by ID
	var pegawai Pegawai
	if err := DB.Table("pegawai").First(&pegawai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pegawai not found"})
		}
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Return the Pegawai data as JSON
	return c.JSON(http.StatusOK, pegawai)
}

func CreatePegawai(c echo.Context) error {
	// Parse the request payload
	var request PegawaiRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Process the uploaded image file
	file, err := c.FormFile("gambar")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Image file is required"})
	}

	// Generate a unique filename for the uploaded image
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	uploadPath := filepath.Join("upload", filename)

	// Save the image file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(uploadPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Create a new Pegawai instance with the image filename
	newPegawai := Pegawai{
		NamaPegawai:    request.NamaPegawai,
		NIK:            request.NIK,
		JenisPegawaiID: request.JenisPegawaiID,
		Unit:           request.Unit,
		SubUnit:        request.SubUnit,
		PendidikanID:   request.PendidikanID,
		TanggalLahir:   request.TanggalLahir,
		TempatLahir:    request.TempatLahir,
		JenisKelaminID: request.JenisKelaminID,
		AgamaID:        request.AgamaID,
		Gambar:         filename, // Save the image filename in the database
	}

	// Save the new Pegawai to the database
	if err := DB.Table("pegawai").Create(&newPegawai).Error; err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Return the created Pegawai as JSON
	return c.JSON(http.StatusCreated, newPegawai)
}

func UpdatePegawai(c echo.Context) error {
	// Get employee ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Parse the request payload
	var request PegawaiRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Retrieve Pegawai by ID
	var existingPegawai Pegawai
	if err := DB.Table("pegawai").First(&existingPegawai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Pegawai not found"})
		}
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Process the updated image file, if any
	file, err := c.FormFile("gambar")
	if err == nil {
		// If a new image is uploaded, generate a unique filename for it
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		uploadPath := filepath.Join("uploads", filename)

		// Save the new image file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(uploadPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Update the image filename in the database
		existingPegawai.Gambar = filename
	}

	// Update the existing Pegawai
	existingPegawai.NamaPegawai = request.NamaPegawai
	existingPegawai.NIK = request.NIK
	existingPegawai.JenisPegawaiID = request.JenisPegawaiID
	existingPegawai.Unit = request.Unit
	existingPegawai.SubUnit = request.SubUnit
	existingPegawai.PendidikanID = request.PendidikanID
	existingPegawai.TanggalLahir = request.TanggalLahir
	existingPegawai.TempatLahir = request.TempatLahir
	existingPegawai.JenisKelaminID = request.JenisKelaminID
	existingPegawai.AgamaID = request.AgamaID

	// Save the updated Pegawai to the database
	if err := DB.Table("pegawai").Save(&existingPegawai).Error; err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Return the updated Pegawai as JSON
	return c.JSON(http.StatusOK, existingPegawai)
}

func DeletePegawai(c echo.Context) error {
	// Get employee ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	// Retrieve Pegawai by ID
	var pegawai Pegawai
	if err := DB.Table("pegawai").First(&pegawai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Pegawai not found"})
		}
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Delete the Pegawai from the database
	if err := DB.Table("pegawai").Delete(&pegawai).Error; err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Delete the associated image file
	imagePath := filepath.Join("uploads", pegawai.Gambar)
	if err := os.Remove(imagePath); err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Return success message
	return c.JSON(http.StatusOK, map[string]string{"message": "Pegawai deleted successfully"})
}


