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
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
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
	s, err := os.Open("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer s.Close()

	var data []v.Flower
	decoder := json.NewDecoder(s)
	if err := decoder.Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})
		return
	}

	id := c.Param("id")
	idToDelete, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}
	updatedData := DeleteEventByID1(data, idToDelete)

	s, err = os.OpenFile("file.json", os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла для записи"})
		return
	}
	defer s.Close()

	jsonData, err := json.MarshalIndent(updatedData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации данных в JSON"})
		return
	}
	if _, err := s.Write(jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})
		return
	}
	c.Status(http.StatusNoContent)
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

    flowerID := c.Param("id")  
    var flowerToUpdate *v.Flower  
    for i := range items[0].Flowers {  
        if items[0].Flowers[i].ID == flowerID {  
            flowerToUpdate = &items[0].Flowers[i]  
            break  
        }  
    }  

    if flowerToUpdate != nil {  
        flowerToUpdate.Name = "новое имя"  
        flowerToUpdate.Quantity += 1  
        c.JSON(http.StatusOK, gin.H{"message": "Цветок успешно обновлён"})  
    } else {  
        newFlower := v.Flower{ID: flowerID, Name: "новое имя", Quantity: 1, Price: 0, ArrivalDate: ""}  
        items[0].Flowers = append(items[0].Flowers, newFlower)  
        c.JSON(http.StatusCreated, gin.H{"message": "Цветок успешно добавлен"})  
    }  

    fileWrite, err := os.OpenFile("file.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)  
    if err != nil {  
        log.Println("Ошибка открытия файла для записи:", err)  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла для записи"})  
        return  
    }  
    defer fileWrite.Close()  

    updatedDataJSON, err := json.MarshalIndent(items, "", "  ")  
    if err != nil {  
        log.Println("Ошибка сериализации данных:", err)  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации данных"})  
        return  
    }  

    if _, err := fileWrite.Write(updatedDataJSON); err != nil {  
        log.Println("Ошибка записи в файл:", err)  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при записи в файл"})  
        return  
    }  
}  