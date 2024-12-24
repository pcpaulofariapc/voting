package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pcpaulofariapc/voting/config"
	"github.com/pcpaulofariapc/voting/controller"
	"github.com/pcpaulofariapc/voting/db"
	"github.com/pcpaulofariapc/voting/repository"
	"github.com/pcpaulofariapc/voting/service"
	"github.com/pcpaulofariapc/voting/usecase"
)

func main() {

	config.Load()

	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	KafikaProducerConnect, err := service.CreateProducer()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to kafka")
	defer KafikaProducerConnect.Close()

	memcacheConnect, err := service.ConnectMemcache()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to memcache")
	defer memcacheConnect.Close()

	gin.SetMode(gin.ReleaseMode)

	//repositorys
	WallRepository := repository.NewWallRepository(dbConnection, memcacheConnect, KafikaProducerConnect)
	// usecases
	WallUsecase := usecase.NewWallUsecase(WallRepository)
	// controllers
	WallController := controller.NewWallController(WallUsecase)

	server.GET("/all-walls", WallController.GetAllWallsDB)
	server.GET("/wall/:id", WallController.GetWallByID)
	server.POST("/wall", WallController.CreateWall)
	server.GET("/all-participants", WallController.GetAllParticipantsDB)
	server.POST("/participant", WallController.CreateParticipant)
	server.GET("/partial-result-active-wall", WallController.GetPartialResultActiveWall)
	server.GET("/participants-active-wall", WallController.GetParticipantsActiveWall)
	server.POST("/vote", WallController.Vote)

	server.Run(":8000")
}
