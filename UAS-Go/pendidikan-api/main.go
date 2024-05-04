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
	err = db.AutoMigrate(&Pendidikan{})
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
	pendidikanHandler := NewPendidikanHandler(db)

	e := echo.New()
	// routing
	e.GET("/pendidikan", pendidikanHandler.GetAllPendidikan)
	e.GET("/pendidikan/:id", pendidikanHandler.GetPendidikanByID)
	e.POST("/pendidikan", pendidikanHandler.CreatePendidikan)
	e.PUT("/pendidikan/:id", pendidikanHandler.UpdatePendidikan)
	e.DELETE("/pendidikan/:id", pendidikanHandler.DeletePendidikan)
	e.Logger.Fatal(e.Start(":1882"))
}

type Pendidikan struct {
	ID     		 int64  `json:"id"`
	Pendidikan   string `json:"pendidikan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Pendidikan) TableName() string {
	return "pendidikan"
}

type PendidikanHandler struct {
	db *gorm.DB
}

func NewPendidikanHandler(db *gorm.DB) *PendidikanHandler {
	return &PendidikanHandler{db: db}
}

type PendidikanRequest struct {
	ID  		  string `param:"id"`
	Pendidikan    string `json:"pendidikan"`
}

func (h *PendidikanHandler) GetAllPendidikan(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	pendidikan := make([]*Pendidikan, 0)
	query := h.db.Model(&Pendidikan{})
	if search != "" {
		query = query.Where("pendidikan LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&pendidikan).Error; err != nil { // SELECT * FROM users
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Pendidikan"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Succesfully Get All Users", "data": pendidikan, "filter": search})
}

func (h *PendidikanHandler) CreatePendidikan(ctx echo.Context) error {
	var input PendidikanRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	pendidikan := &Pendidikan{
		Pendidikan:    input.Pendidikan,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(pendidikan).Error; err != nil { // INSERT INTO users (nim, nama, alamat) VALUES('')
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Pendidikan"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Succesfully Create a Pendidikan", "data": pendidikan})
}

func (h *PendidikanHandler) GetPendidikanByID(ctx echo.Context) error {
	var input PendidikanRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	pendidikan := new(Pendidikan)

	if err := h.db.Where("id =?", input.ID).First(&pendidikan).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Pendidikan By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Get Pendidikan By ID : %s", input.ID), "data": pendidikan})
}

func (h *PendidikanHandler) UpdatePendidikan(ctx echo.Context) error {
	var input PendidikanRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	pendidikanID, _ := strconv.Atoi(input.ID)

	pendidikan := Pendidikan{
		ID:     int64(pendidikanID),
		Pendidikan:    input.Pendidikan,
		UpdatedAt: time.Now(),
	}

	query := h.db.Model(&Pendidikan{}).Where("id = ?", pendidikanID)
	if err := query.Updates(&pendidikan).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Pendidikan By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Update Pendidikan By ID : %s", input.ID), "data": input})
}

func (h *PendidikanHandler) DeletePendidikan(ctx echo.Context) error {
	var input PendidikanRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&Pendidikan{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Pendidikan By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}