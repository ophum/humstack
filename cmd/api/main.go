package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/v1/pkg/api/controller"
	"github.com/ophum/humstack/v1/pkg/api/delivery/http/router"
	"github.com/ophum/humstack/v1/pkg/api/repository/rdb"
	"github.com/ophum/humstack/v1/pkg/api/usecase"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	r := gin.Default()

	diskRepository := rdb.NewDiskRepository(db)
	diskUsecase := usecase.NewDiskUsecase(diskRepository)
	diskController := controller.NewDiskController(diskUsecase)

	diskRouter := router.NewDiskRouter(r, diskController)
	diskRouter.RegisterRoutes()

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
