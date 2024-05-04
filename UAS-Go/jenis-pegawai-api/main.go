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
	err = db.AutoMigrate(&JenisPegawai{})
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
	jenisPegawaiHandler := NewJenisPegawaiHandler(db)

	e := echo.New()
	// routing
	e.GET("/jenispegawai", jenisPegawaiHandler.GetAllJenisPegawai)
	e.GET("/jenispegawai/:id", jenisPegawaiHandler.GetJenisPegawaiByID)
	e.POST("/jenispegawai", jenisPegawaiHandler.CreateJenisPegawai)
	e.PUT("/jenispegawai/:id", jenisPegawaiHandler.UpdateJenisPegawai)
	e.DELETE("/jenispegawai/:id", jenisPegawaiHandler.DeleteJenisPegawai)
	e.Logger.Fatal(e.Start(":1882"))
}

type JenisPegawai struct {
	ID     				int64  `json:"id"`
	Jenis_Pegawai   	string `json:"jenis_pegawai"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (JenisPegawai) TableName() string {
	return "jenis_pegawai"
}

type JenisPegawaiHandler struct {
	db *gorm.DB
}

func NewJenisPegawaiHandler(db *gorm.DB) *JenisPegawaiHandler {
	return &JenisPegawaiHandler{db: db}
}

type JenisPegawaiRequest struct {
	ID     			 string `param:"id"`
	Jenis_Pegawai    string `json:"jenis_pegawai"`
}

func (h *JenisPegawaiHandler) GetAllJenisPegawai(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	jenispegawai := make([]*JenisPegawai, 0)
	query := h.db.Model(&JenisPegawai{})
	if search != "" {
		query = query.Where("jenispegawai LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&jenispegawai).Error; err != nil { // SELECT * FROM users
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Jenis Pegawai"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Succesfully Get All Users", "data": jenispegawai, "filter": search})
}

func (h *JenisPegawaiHandler) CreateJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenispegawai := &JenisPegawai{
		Jenis_Pegawai:    input.Jenis_Pegawai,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(jenispegawai).Error; err != nil { // INSERT INTO users (nim, nama, alamat) VALUES('')
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Jenis Pegawai"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Succesfully Create a Agama", "data": jenispegawai})
}

func (h *JenisPegawaiHandler) GetJenisPegawaiByID(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenispegawai := new(JenisPegawai)

	if err := h.db.Where("id =?", input.ID).First(&jenispegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Jenis Pegawai By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Get Jenis Pegawai By ID : %s", input.ID), "data": jenispegawai})
}

func (h *JenisPegawaiHandler) UpdateJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenispegawaiID, _ := strconv.Atoi(input.ID)

	jenispegawai := JenisPegawai{
		ID:     int64(jenispegawaiID),
		Jenis_Pegawai:    input.Jenis_Pegawai,
		UpdatedAt: time.Now(),
	}

	query := h.db.Model(&JenisPegawai{}).Where("id = ?", jenispegawaiID)
	if err := query.Updates(&jenispegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Jenis Pegawai By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Update Jenis Pegawai By ID : %s", input.ID), "data": input})
}

func (h *JenisPegawaiHandler) DeleteJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&JenisPegawai{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Jenis Pegawai By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}