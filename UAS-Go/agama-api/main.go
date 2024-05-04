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
	err = db.AutoMigrate(&Agama{})
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
	agamaHandler := NewAgamaHandler(db)

	e := echo.New()
	// routing
	e.GET("/agama", agamaHandler.GetAllAgama)
	e.GET("/agama/:id", agamaHandler.GetAgamaByID)
	e.POST("/agama", agamaHandler.CreateAgama)
	e.PUT("/agama/:id", agamaHandler.UpdateAgama)
	e.DELETE("/agama/:id", agamaHandler.DeleteAgama)
	e.Logger.Fatal(e.Start(":1882"))
}

type Agama struct {
	ID     int64  `json:"id"`
	Nama    string `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Agama) TableName() string {
	return "agama"
}

type AgamaHandler struct {
	db *gorm.DB
}

func NewAgamaHandler(db *gorm.DB) *AgamaHandler {
	return &AgamaHandler{db: db}
}

type AgamaRequest struct {
	ID      string `param:"id"`
	Nama    string `json:"nama"`
}

func (h *AgamaHandler) GetAllAgama(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	agama := make([]*Agama, 0)
	query := h.db.Model(&Agama{})
	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&agama).Error; err != nil { // SELECT * FROM users
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Agama"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Succesfully Get All Users", "data": agama, "filter": search})
}

func (h *AgamaHandler) CreateAgama(ctx echo.Context) error {
	var input AgamaRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	agama := &Agama{
		Nama:    input.Nama,
		CreatedAt: time.Now(),
	}

	if err := h.db.Create(agama).Error; err != nil { // INSERT INTO users (nim, nama, alamat) VALUES('')
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Agama"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Succesfully Create a Agama", "data": agama})
}

func (h *AgamaHandler) GetAgamaByID(ctx echo.Context) error {
	var input AgamaRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	agama := new(Agama)

	if err := h.db.Where("id =?", input.ID).First(&agama).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Task By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Get Agama By ID : %s", input.ID), "data": agama})
}

func (h *AgamaHandler) UpdateAgama(ctx echo.Context) error {
	var input AgamaRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	agamaID, _ := strconv.Atoi(input.ID)

	agama := Agama{
		ID:        int64(agamaID),
		Nama:      input.Nama,
		UpdatedAt: time.Now(),
	}

	query := h.db.Model(&Agama{}).Where("id = ?", agamaID)
	if err := query.Updates(&agama).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Agama By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Succesfully Update Agama By ID : %s", input.ID), "data": input})
}

func (h *AgamaHandler) DeleteAgama(ctx echo.Context) error {
	var input AgamaRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&Agama{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Agama By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}