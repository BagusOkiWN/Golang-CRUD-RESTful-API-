package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/acrud?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Pegawai{})
    if err != nil {
        return nil, err
    }

	return db, nil
}

func main() {
	// initialisasi database
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	// inisialisasi handler
	pegawaiHandler := NewPegawaiHandler(db)

	e := echo.New()
	// routing
	e.GET("/pegawai", pegawaiHandler.GetAllPegawai)
	e.GET("/pegawai/:id", pegawaiHandler.GetPegawaiByID)
	e.POST("/pegawai", pegawaiHandler.CreatePegawai)
	e.PUT("/pegawai/:id", pegawaiHandler.UpdatePegawai)
	e.DELETE("/pegawai/:id", pegawaiHandler.DeletePegawai)
	e.Logger.Fatal(e.Start(":1882"))
}

type Pegawai struct {
	ID     		 		int64  `json:"id"`
	Nama_Pegawai   		string `json:"nama_pegawai"`
	NIK			   		string `json:"nik"`
	Jenis_Pegawai_ID	int64  `json:"jenis_pegawai_id"`
	Status_Pegawai_ID	int64  `json:"status_pegawai_id"`
	Unit				string `json:"unit"`
	Sub_Unit			string `json:"sub_unit"`
	Pendidikan_ID		int64  `json:"pendidikan_id"`
	Tgl_Lahir			string `json:"tgl_lahir"`
	Tpt_Lahir			string `json:"tpt_lahir"`
	Jenkel_ID			int64  `json:"jenkel_id"`
	Agama_ID			int64  `json:"agama_id"`
	Gambar				string `json:"gambar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Pegawai) TableName() string {
	return "pegawai"
}

type PegawaiHandler struct {
	db *gorm.DB
}

func NewPegawaiHandler(db *gorm.DB) *PegawaiHandler {
	return &PegawaiHandler{db: db}
}

type PegawaiRequest struct {
	ID  		  		string `param:"id"`
	Nama_Pegawai   		string `json:"nama_pegawai"`
	NIK			   		string `json:"nik"`
	Jenis_Pegawai_ID	int64  `json:"jenis_pegawai_id"`
	Status_Pegawai_ID	int64  `json:"status_pegawai_id"`
	Unit				string `json:"unit"`
	Sub_Unit			string `json:"sub_unit"`
	Pendidikan_ID		int64  `json:"pendidikan_id"`
	Tgl_Lahir			string `json:"tgl_lahir"`
	Tpt_Lahir			string `json:"tpt_lahir"`
	Jenkel_ID			int64  `json:"jenkel_id"`
	Agama_ID			int64  `json:"agama_id"`
	Gambar				string `json:"gambar"`
}

func (h *PegawaiHandler) GetAllPegawai(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	pegawai := make([]*Pegawai, 0)
	query := h.db.Model(&Pegawai{})
	if search != "" {
		query = query.Where("nama_pegawai LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&pegawai).Error; err != nil { // SELECT * FROM users
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Pegawai"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Succesfully Get All Users", "data": pegawai, "filter": search})
}

func (h *PegawaiHandler) GetPegawaiByID(ctx echo.Context) error {
	var input PegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	pegawai := new(Pegawai)

	if err := h.db.Where("id =?", input.ID).First(&pegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Pegawai By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Get Pegawai By ID : %s", input.ID), "data": pegawai})
}

func (h *PegawaiHandler) CreatePegawai(ctx echo.Context) error {
    request := new(PegawaiRequest)
    if err := ctx.Bind(request); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
    }

    pegawai := Pegawai{
        Nama_Pegawai:      request.Nama_Pegawai,
        NIK:               request.NIK,
        Jenis_Pegawai_ID:  request.Jenis_Pegawai_ID,
        Status_Pegawai_ID: request.Status_Pegawai_ID,
        Unit:              request.Unit,
        Sub_Unit:          request.Sub_Unit,
        Pendidikan_ID:     request.Pendidikan_ID,
        Tgl_Lahir:         request.Tgl_Lahir,
        Tpt_Lahir:         request.Tpt_Lahir,
        Jenkel_ID:         request.Jenkel_ID,
        Agama_ID:          request.Agama_ID,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    if err := h.db.Create(&pegawai).Error; err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create Pegawai"})
    }

    return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Pegawai created successfully", "data": pegawai})
}

func (h *PegawaiHandler) UpdatePegawai(ctx echo.Context) error {
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
    }

    request := new(PegawaiRequest)
    if err := ctx.Bind(request); err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
    }

    // Check if the Pegawai with the given ID exists
    pegawai := new(Pegawai)
    if err := h.db.Where("id = ?", id).First(&pegawai).Error; err != nil {
        return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Pegawai not found"})
    }

    // Update Pegawai attributes
    pegawai.Nama_Pegawai = request.Nama_Pegawai
    pegawai.NIK = request.NIK
    pegawai.Jenis_Pegawai_ID = request.Jenis_Pegawai_ID
    pegawai.Status_Pegawai_ID = request.Status_Pegawai_ID
    pegawai.Unit = request.Unit
    pegawai.Sub_Unit = request.Sub_Unit
    pegawai.Pendidikan_ID = request.Pendidikan_ID
    pegawai.Tgl_Lahir = request.Tgl_Lahir
    pegawai.Tpt_Lahir = request.Tpt_Lahir
    pegawai.Jenkel_ID = request.Jenkel_ID
    pegawai.Agama_ID = request.Agama_ID

    // Save the changes to the database
    if err := h.db.Save(&pegawai).Error; err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update Pegawai"})
    }

    return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Pegawai updated successfully", "data": pegawai})
}

func (h *PegawaiHandler) DeletePegawai(ctx echo.Context) error {
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
    }

    // Check if the Pegawai with the given ID exists
    pegawai := new(Pegawai)
    if err := h.db.Where("id = ?", id).First(&pegawai).Error; err != nil {
        return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Pegawai not found"})
    }

    // Delete the Pegawai from the database
    if err := h.db.Delete(&pegawai).Error; err != nil {
        return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete Pegawai"})
    }

    return ctx.JSON(http.StatusNoContent, map[string]interface{}{"message": "Pegawai deleted successfully"})
}


