package handler

import (
	"challenge/database"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Key      string     `gorm:"type:uuid;primary_key"`
	Filename string     `gorm:"not null"`
	ID       *uuid.UUID `gorm:"not null"`
}

// var idx = int(0)

// Crear una nueva instancia del modelo File Y lo guarda en la base de datos junto con el fileID
func saveKey(key, filename string, ID *uuid.UUID, db *gorm.DB) error {

	// idx++
	file := File{
		ID:       ID,
		Key:      key,
		Filename: filename,
	}

	// if idx == 3 {
	// 	return errors.New("SIMULATION ERROR")
	// }
	return db.Create(&file).Error
}

// Simular la carga del archivo en S3
func uploadToS3(file *multipart.FileHeader) (string, error) {
	uuid := uuid.New().String()
	key := fmt.Sprintf("files/%s-%s", uuid, file.Filename)
	return key, nil
}

func HandleEndpoint2(c *gin.Context) {
	ID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		fmt.Printf("ERROR PARSING ID TO UUID %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB máximo para el tamaño del formulario multipart
	if err != nil {
		fmt.Printf("MULTI PART FORM ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files"]

	//OPEN TRANSACTION
	tx := database.DB.Begin()
	if tx.Error != nil {
		fmt.Printf("BEGIN TRANSACTION ERROR %+v", tx.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"TRANSACTION_ERROR": tx.Error.Error()})
		return
	}

	keys := []string{}
	for _, file := range files {
		src, err := file.Open()

		if err != nil {
			fmt.Printf("CAN'T OPEN FILE %+v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
			return
		}
		defer src.Close()

		//AQUI SE PODRIA VALIDAR EL TIPO DE ARCHIVO QUE SACAMOS DEL HEADER, TAMBIEN SE PUEDE VALIDAR EL FILE.SIZE
		// if file.Header["Content-Type"][0] == "text/plain" {
		// 	fmt.Printf("SRC FILE OPEN >> %+v", file.Header["Content-Type"])
		// }

		key, err := uploadToS3(file)
		if err != nil {
			fmt.Printf("CAN'T UPLOAD FILE TO S3 %+v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
			return
		}

		err = saveKey(key, file.Filename, &ID, tx)
		if err != nil {
			fmt.Printf("CAN'T SAVE FILE %+v", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
			return
		}

		keys = append(keys, key)
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("CAN'T COMMIT THE TRANSACTION %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Files uploaded successfully", "keys": keys})
}
