package car

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

func GetCars(c *gin.Context) {
	_, cars, _, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, cars)
}
func GetCarsByID(c *gin.Context) {
	_, cars, _, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	id := c.Param("id")
	for _, a := range cars {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func DeleteEventByID1(events []v.Car, id int) []v.Car {
	idInt := strconv.Itoa(id)
	for i, event := range events {
		if event.ID == idInt {
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

	data := data0[0].Cars

	id := c.Param("id")
	idToDelete, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	updatedData := DeleteEventByID1(data, idToDelete)

	data0[0].Cars = updatedData

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

	carsID := c.Param("id")
	var carsToUpdate *v.Car
	for i := range items[0].Cars {
		if items[0].Cars[i].ID == carsID {
			carsToUpdate = &items[0].Cars[i]
			break
		}
	}

	var updateRequest v.Car
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if carsToUpdate != nil {
		carsToUpdate.Brand = updateRequest.Brand
		carsToUpdate.Model = updateRequest.Model
		carsToUpdate.Mileage = updateRequest.Mileage
		carsToUpdate.Owners = updateRequest.Owners
		c.JSON(http.StatusOK, gin.H{"message": "Машина успешно обновлена"})
	} else {
		newCar := v.Car{
			ID:      carsID,
			Brand:   updateRequest.Brand,
			Model:   updateRequest.Model,
			Mileage: updateRequest.Mileage,
			Owners:  updateRequest.Owners,
		}
		items[0].Cars = append(items[0].Cars, newCar)
		c.JSON(http.StatusCreated, gin.H{"message": "Машина успешно добавлена"})
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

	carsID := c.Param("id")
	var carsToUpdate *v.Car
	for i := range items[0].Cars {
		if items[0].Cars[i].ID == carsID {
			carsToUpdate = &items[0].Cars[i]
			break
		}
	}

	var updateRequest v.Car
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if carsToUpdate != nil {
		carsToUpdate.Model = updateRequest.Model
		carsToUpdate.Owners = updateRequest.Owners

		if err := writeFile("file.json", items); err != nil {
			log.Println("Ошибка при записи в файл:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Машина успешно обновлена"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Автомобиль не найден"})
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
