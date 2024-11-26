package furniture

import (
	r "apiGO/run"
	v "apiGO/structFile"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

func GetFurnitures(c *gin.Context) { //Get
	_, _, furniture, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, furniture)
}
func GetFurnitureByID(c *gin.Context) { //GetID
	_, _, furniture, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	id := c.Param("id")
	for _, a := range furniture {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func DeleteEventByID1(events []v.Furniture, id int) []v.Furniture {
	idInt := strconv.Itoa(id)
	for i, event := range events {
		if event.ID == idInt {
			fmt.Println("Успешное удаление")
			return append(events[:i], events[i+1:]...)
		}
	}

	return events
}
func DeletedById(c *gin.Context) { //DeleteID
	s, err := os.Open("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer s.Close()

	decoder, err := io.ReadAll(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
		return
	}

	var data0 []v.Inventory

	if err := json.Unmarshal(decoder, &data0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})
		return
	}

	data := data0[0].Furniture

	id := c.Param("id")
	idToDelete, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	updatedData := DeleteEventByID1(data, idToDelete)

	data0[0].Furniture = updatedData
	s, err = os.OpenFile("file.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла для записи"})
		return
	}
	defer s.Close()

	jsonData, err := json.MarshalIndent(data0, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации данных в JSON"})
		return
	}

	if _, err := s.Write(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Успешно": "удаление получилось"})
}
func PostFurnitures(c *gin.Context) { //Post
	file, err := os.Open("file.json")
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Файл не найден"})
		return
	}
	defer file.Close()

	readFile, err := io.ReadAll(file)
	if err != nil {
		log.Println("Ошибка чтения файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
		return
	}

	var items []v.Inventory
	if err := json.Unmarshal(readFile, &items); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})
		return
	}

	nextID := 1
	if len(items) > 0 && len(items[0].Furniture) > 0 {
		var maxID int
		for _, furniture := range items[0].Furniture {
			idNum, err := strconv.Atoi(furniture.ID)
			if err == nil && idNum > maxID {
				maxID = idNum
			}
		}
		nextID = maxID + 1
	}

	var updateRequest v.Furniture
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	newfurniture := v.Furniture{
		ID:           strconv.Itoa(nextID),
		Name:         updateRequest.Name,
		Manufacturer: updateRequest.Manufacturer,
		Height:       updateRequest.Height,
		Width:        updateRequest.Width,
		Length:       updateRequest.Length,
	}
	items[0].Furniture = append(items[0].Furniture, newfurniture)
	c.JSON(http.StatusCreated, newfurniture)

	if err := writeFile("file.json", items); err != nil {
		log.Println("Ошибка при записи в файл:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}
}
func PutItem(c *gin.Context) { //Put
	file, err := os.Open("file.json")
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Файл не найден"})
		return
	}
	defer file.Close()

	readFile, err := io.ReadAll(file)
	if err != nil {
		log.Println("Ошибка чтения файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
		return
	}

	var items []v.Inventory
	if err := json.Unmarshal(readFile, &items); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})
		return
	}

	furnitureID := c.Param("id")
	var furnitureToUpdate *v.Furniture
	for i := range items[0].Furniture {
		if items[0].Furniture[i].ID == furnitureID {
			furnitureToUpdate = &items[0].Furniture[i]
			break
		}
	}

	var updateRequest v.Furniture
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if furnitureToUpdate != nil {
		furnitureToUpdate.Name = updateRequest.Name
		furnitureToUpdate.Manufacturer = updateRequest.Manufacturer
		furnitureToUpdate.Height = updateRequest.Height
		furnitureToUpdate.Width = updateRequest.Width
		furnitureToUpdate.Length = updateRequest.Length

		if err := writeFile("file.json", items); err != nil {
			log.Println("Ошибка при записи в файл:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
			return
		}
		c.JSON(http.StatusOK, furnitureToUpdate)
	} else {
		c.JSON(http.StatusNoContent, nil)
	}

	if err := writeFile("file.json", items); err != nil {
		log.Println("Ошибка при записи в файл:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}
}
func PatchItem(c *gin.Context) { //Patch
	file, err := os.Open("file.json")
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Файл не найден"})
		return
	}
	defer file.Close()

	readFile, err := io.ReadAll(file)
	if err != nil {
		log.Println("Ошибка чтения файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
		return
	}

	var items []v.Inventory
	if err := json.Unmarshal(readFile, &items); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})
		return
	}

	furnitureID := c.Param("id")
	var furnitureToUpdate *v.Furniture
	for i := range items[0].Furniture {
		if items[0].Furniture[i].ID == furnitureID {
			furnitureToUpdate = &items[0].Furniture[i]
			break
		}
	}

	var updateRequest v.Furniture
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if furnitureToUpdate != nil {
		if updateRequest.Name  != "" {
			furnitureToUpdate.Name  = updateRequest.Name 
		}
		if updateRequest.Manufacturer != "" {
			furnitureToUpdate.Manufacturer = updateRequest.Manufacturer
		}
		if updateRequest.Height != 0 {
			furnitureToUpdate.Height = updateRequest.Height
		}
		if updateRequest.Width != 0 {
			furnitureToUpdate.Width = updateRequest.Width
		}
		if updateRequest.Length != 0 {
			furnitureToUpdate.Length = updateRequest.Length
		}

		if err := writeFile("file.json", items); err != nil {
			log.Println("Ошибка при записи в файл:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
			return
		}
		c.JSON(http.StatusOK, furnitureToUpdate)
	} else {
		c.JSON(http.StatusNoContent, nil)
	}

	if err := writeFile("file.json", items); err != nil {
		log.Println("Ошибка при записи в файл:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}
}
func writeFile(filename string, data interface{}) error {
	fileWrite, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	updatedDataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if _, err := fileWrite.Write(updatedDataJSON); err != nil {
		return err
	}

	return nil
}
