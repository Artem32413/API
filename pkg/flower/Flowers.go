package flower

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

func GetFlowers(c *gin.Context) { //Get
	flowers, _, _, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, flowers)
}
func GetFlowerByID(c *gin.Context) { //GetID
	flowers, _, _, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	id := c.Param("id")
	for _, a := range flowers {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Объект не найден"})
}
func DeleteEventByID1(events []v.Flower, id int) []v.Flower {
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

	data := data0[0].Flowers

	id := c.Param("id")
	idToDelete, err := strconv.Atoi(id)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	updatedData := DeleteEventByID1(data, idToDelete)

	data0[0].Flowers = updatedData

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
func PostFlowers(c *gin.Context) { //Post
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
	if len(items) > 0 && len(items[0].Flowers) > 0 {
		var maxID int
		for _, flower := range items[0].Flowers {
			idNum, err := strconv.Atoi(flower.ID)
			if err == nil && idNum > maxID {
				maxID = idNum
			}
		}
		nextID = maxID + 1
	}

	var updateRequest v.Flower
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	newFlower := v.Flower{
		ID:          strconv.Itoa(nextID),
		Name:        updateRequest.Name,
		Quantity:    updateRequest.Quantity,
		Price:       updateRequest.Price,
		ArrivalDate: updateRequest.ArrivalDate,
	}
	items[0].Flowers = append(items[0].Flowers, newFlower)
	c.JSON(http.StatusCreated, newFlower)

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

	flowersID := c.Param("id")
	var flowersToUpdate *v.Flower
	for i := range items[0].Flowers {
		if items[0].Flowers[i].ID == flowersID {
			flowersToUpdate = &items[0].Flowers[i]
			break
		}
	}

	var updateRequest v.Flower
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if flowersToUpdate != nil {
		flowersToUpdate.Name = updateRequest.Name
		flowersToUpdate.Quantity = updateRequest.Quantity
		flowersToUpdate.Price = updateRequest.Price
		flowersToUpdate.ArrivalDate = updateRequest.ArrivalDate

		if err := writeFile("file.json", items); err != nil {
			log.Println("Ошибка при записи в файл:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
			return
		}
		c.JSON(http.StatusOK, flowersToUpdate)
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

	flowersID := c.Param("id")
	var flowersToUpdate *v.Flower
	for i := range items[0].Flowers {
		if items[0].Flowers[i].ID == flowersID {
			flowersToUpdate = &items[0].Flowers[i]
			break
		}
	}

	var updateRequest v.Flower
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Println("Ошибка связывания данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	if flowersToUpdate != nil {
		if updateRequest.Name != ""{
			flowersToUpdate.Name = updateRequest.Name
		}
		if updateRequest.Quantity != 0{
			flowersToUpdate.Quantity = updateRequest.Quantity
		}
		if updateRequest.Price != 0{
			flowersToUpdate.Price = updateRequest.Price
		}
		if updateRequest.ArrivalDate != ""{
			flowersToUpdate.ArrivalDate = updateRequest.ArrivalDate
		}

		if err := writeFile("file.json", items); err != nil {
			log.Println("Ошибка при записи в файл:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
			return
		}
		c.JSON(http.StatusOK, flowersToUpdate)
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
