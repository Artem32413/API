package router
import (
	f "apiGO/pkg/flower"
	c "apiGO/pkg/car"
	fu "apiGO/pkg/furniture"

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
	router.PUT("/flowers", f.PutItem)
	//cars
	router.GET("/cars", c.GetCars)
	router.GET("/cars/:id", c.GetCarsByID)
	router.DELETE("/cars/:id", c.DeletedById)
	router.PUT("/cars", c.PutItem)
	//furniture
	router.GET("/furniture", fu.GetFurnitures)
	router.GET("/furniture/:id", fu.GetFurnitureByID)
	router.DELETE("/furniture/:id", fu.DeletedById)
	router.POST("/furniture", fu.PostFurnitures)
	router.Run(":8080")
	
}