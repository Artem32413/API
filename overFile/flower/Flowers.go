package flower

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
            // Удаляем элемент, используя append  
            return append(events[:i], events[i+1:]...) // Удаляет элемент по индексу  
        }  
    }  
    return events // Если ID не найден, возвращаем оригинальный срез  
} 
func DeletedById(c *gin.Context) { //DeleteID  
    // Открытие файла для чтения  
    s, err := os.Open("file.json")  
    if err != nil {  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})  
        return  
    }  
    defer s.Close()  
    
    // Чтение содержимого файла  
    decoder, err := io.ReadAll(s) // Обработка ошибки  
    if err != nil {  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})  
        return  
    }  
    
    var data0 []v.Inventory  
    // Декодирование JSON в структуру  
    if err := json.Unmarshal(decoder, &data0); err != nil {  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при декодировании JSON"})  
        return  
    }  

    // Получение списка цветов  
    data := data0[0].Flowers  

    // Получение ID из параметров запроса  
    id := c.Param("id")  
    idToDelete, err := strconv.Atoi(id)  
    if err != nil {  
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})  
        return  
    }  

    // Удаление элемента  
    updatedData := DeleteEventByID1(data, idToDelete)  

    // Обновление структуры с новыми данными  
    data0[0].Flowers = updatedData  

    // Открытие файла для записи  
    s, err = os.OpenFile("file.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)  
    if err != nil {  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла для записи"})  
        return  
    }  
    defer s.Close()  

    // Сериализация данных обратно в JSON  
    jsonData, err := json.MarshalIndent(data0, "", "  ")  
    if err != nil {  
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сериализации данных в JSON"})  
        return  
    }  

    // Запись в файл  
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
