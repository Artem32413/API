package furniture

import (
	r "apiGO/run"
	v "apiGO/structFile"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

func GetFurnitures(c *gin.Context) {//Get
	_, _, furniture, err := r.ReadFileGet("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON"})
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, furniture)
}
func GetFurnitureByID(c *gin.Context) {//GetID
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
func PostFurnitures(c *gin.Context) {//Post
	s, err := os.Open("file.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer s.Close()

	var data []v.Furniture
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
