package router
import (
	f "apiGO/overFile/flower"
	c "apiGO/overFile/car"
	fu "apiGO/overFile/furniture"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

func Run(){
	router := gin.Default()
	//flowers
	router.GET("/flowers", f.GetFlowers)
	router.GET("/flowers/:id", f.GetFlowerByID)
	router.DELETE("/flowers/:id", f.DeletedById)
	router.POST("/flowers", f.PostFlowers)
	//cars
	router.GET("/cars", c.GetCars)
	router.GET("/cars/:id", c.GetCarsByID)
	router.DELETE("/cars/:id", c.DeletedById)
	//furniture
	router.GET("/furniture", fu.GetFurnitures)
	router.GET("/furniture/:id", fu.GetFurnitureByID)
	router.DELETE("/furniture/:id", fu.DeletedById)
	router.POST("/furniture", fu.PostFurnitures)
	router.Run(":8080")
	
}