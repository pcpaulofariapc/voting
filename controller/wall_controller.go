package controller

import (
	"errors"
	"net"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pcpaulofariapc/voting/config"
	"github.com/pcpaulofariapc/voting/model"
	"github.com/pcpaulofariapc/voting/usecase"
)

type wallController struct {
	// Usecase
	wallUsecase usecase.WallUsecase
}

func NewWallController(usecase usecase.WallUsecase) wallController {
	return wallController{
		wallUsecase: usecase,
	}
}

func (w *wallController) GetWallByID(ctx *gin.Context) {
	var (
		err    error
		id     string
		wallID uuid.UUID
	)

	id = ctx.Param("id")
	if id == "" {
		err = errors.New("id do paredão não pode ser nulo")
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	wallID, err = uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	wall, err := w.wallUsecase.GetWallByID(wallID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, wall)
}

func (w *wallController) GetAllWallsDB(ctx *gin.Context) {
	wallsDB, err := w.wallUsecase.GetAllWallsDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, wallsDB)
}

func (w *wallController) GetAllParticipantsDB(ctx *gin.Context) {
	participantsDB, err := w.wallUsecase.GetAllParticipantsDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, participantsDB)
}

func (w *wallController) GetParticipantsActiveWall(ctx *gin.Context) {

	participants, err := w.wallUsecase.GetParticipantsActiveWall()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, participants)
}

func (w *wallController) Vote(ctx *gin.Context) {
	var (
		err    error
		vote   model.Vote
		ip     string
		result *model.ResultVote
	)

	err = ctx.BindJSON(&vote)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ip, err = getClientIPByRequestRemoteAddr(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	vote.IP = ip

	result, err = w.wallUsecase.Vote(vote)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (w *wallController) GetPartialResultActiveWall(ctx *gin.Context) {
	var (
		err    error
		result *model.ActiveWall
	)
	result, err = w.wallUsecase.GetPartialResultActiveWall()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (w *wallController) CreateParticipant(ctx *gin.Context) {
	var (
		err               error
		participantCreate model.ParticipantCreate
		id                *uuid.UUID
	)

	err = ctx.BindJSON(&participantCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(participantCreate.Name) > 64 {
		err = errors.New("nome do participante não pode ter mais que 64 caracteres")
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err = w.wallUsecase.CreateParticipant(participantCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, id)
}

func (w *wallController) CreateWall(ctx *gin.Context) {
	var (
		err          error
		wallCreate   model.WallCreate
		id           *uuid.UUID
		timeLocation *time.Location
		currentTime  time.Time
	)

	timeLocation = config.GetTimeLocation()
	currentTime = time.Now().In(timeLocation)

	err = ctx.BindJSON(&wallCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	wallCreate.ParticipantsID = slices.Compact(wallCreate.ParticipantsID)
	wallCreate.StartTime = wallCreate.StartTime.In(timeLocation)
	wallCreate.EndTime = wallCreate.EndTime.In(timeLocation)

	if wallCreate.StartTime.After(wallCreate.EndTime) {
		err = errors.New("horário de início precisa ser antes do horário de término")
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if wallCreate.StartTime.Before(currentTime) {
		err = errors.New("horário de início precisa ser posterior ao horário atual")
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(wallCreate.Name) > 64 {
		err = errors.New("nome do paredão não pode ter mais que 64 caracteres")
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err = w.wallUsecase.CreateWall(wallCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func getClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {
	// Try via request
	ip, _, err = net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		err = errors.New("debug: Parsing IP from Request.RemoteAddr got nothing")
		return "", err

	}
	return userIP.String(), nil

}
