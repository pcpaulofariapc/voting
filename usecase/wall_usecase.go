package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pcpaulofariapc/voting/config"
	"github.com/pcpaulofariapc/voting/model"
	"github.com/pcpaulofariapc/voting/repository"
)

type WallUsecase struct {
	repository repository.WallRepository
}

func NewWallUsecase(repo repository.WallRepository) WallUsecase {
	return WallUsecase{
		repository: repo,
	}
}

func (w *WallUsecase) GetParticipantsActiveWall() (participants *model.ParticipantsActiveWall, err error) {
	var activeWall *model.ActiveWall
	activeWall, err = w.repository.GetActiveWall()
	if err != nil {
		return nil, err
	}

	if activeWall != nil {
		participants = new(model.ParticipantsActiveWall)
		participants.WallID = activeWall.ID
		participants.WallName = activeWall.Name
		participants.Participants = make([]*model.ParticipantActiveWall, len(activeWall.Participants))
		for i := range activeWall.Participants {
			participants.Participants[i] = &model.ParticipantActiveWall{
				ID:   activeWall.Participants[i].ID,
				Name: activeWall.Participants[i].Name,
			}
		}
	}
	return
}

func (w *WallUsecase) Vote(vote model.Vote) (result *model.ResultVote, err error) {
	var (
		idValidate   *uuid.UUID
		activeWAll   *model.ActiveWall
		dateRegister time.Time
		timeLocation *time.Location
	)

	timeLocation = config.GetTimeLocation()
	dateRegister = time.Now().In(timeLocation)

	activeWAll, err = w.repository.GetActiveWall()
	if err != nil {
		return nil, err
	}

	if activeWAll == nil {
		err = errors.New("não existe paredão ativo no momento")
		return nil, err
	}

	if dateRegister.After(activeWAll.EndTime.In(timeLocation)) {
		err = errors.New("o paredão encerrou")
		return nil, err
	}

	if activeWAll.ID != *vote.WallID {
		err = errors.New("o paredão ativo no momento é diferente do paredão indicado no voto")
		return nil, err
	}

	for i := range activeWAll.Participants {
		if *vote.ParticipantID == activeWAll.Participants[i].ID {
			idValidate = vote.ParticipantID
		}
	}

	if idValidate == nil {
		err = errors.New("participante não encontrado")
		return nil, err
	}

	voteRegister := model.VoteRegister{
		ParticipantID: *idValidate,
		DateRegister:  time.Now().In(config.GetTimeLocation()),
		WallID:        activeWAll.ID,
		ID:            uuid.New(),
		IP:            vote.IP,
	}

	var val []byte
	val, err = json.Marshal(voteRegister)
	if err != nil {
		return nil, err
	}

	err = w.repository.RegsiterVote(val)
	if err != nil {
		return nil, err
	}

	result = new(model.ResultVote)
	result.RegisterID = voteRegister.ID
	result.PartialResults = make([]model.PartialResult, len(activeWAll.Participants))
	for i := range activeWAll.Participants {
		result.PartialResults[i].ID = activeWAll.Participants[i].ID
		result.PartialResults[i].Name = activeWAll.Participants[i].Name
		result.PartialResults[i].Votes = activeWAll.Participants[i].Votes
		result.PartialResults[i].VotesPercentage = activeWAll.Participants[i].VotesPercentage
	}

	fmt.Printf("\n Register vote, vote ip: %s, participant id: %s date: %s, \n", voteRegister.IP, voteRegister.ParticipantID, voteRegister.DateRegister)

	return result, nil
}

func (w *WallUsecase) GetPartialResultActiveWall() (activeWall *model.ActiveWall, err error) {
	activeWall, err = w.repository.GetActiveWall()
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallUsecase) GetAllWallsDB() (allWalls *model.AllWalls, err error) {
	allWalls, err = w.repository.GetAllWallsDB()
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallUsecase) GetWallByID(wallID uuid.UUID) (wall *model.ActiveWall, err error) {
	wall, err = w.repository.GetWallByID(wallID)
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallUsecase) GetAllParticipantsDB() (allParticipants *model.Participants, err error) {
	allParticipants, err = w.repository.GetAllParticipantsDB()
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallUsecase) CreateParticipant(participant model.ParticipantCreate) (id *uuid.UUID, err error) {
	id, err = w.repository.CreateParticipant(participant)
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallUsecase) CreateWall(wall model.WallCreate) (id *uuid.UUID, err error) {
	var (
		participantsDB *model.Participants
		validatesID    int
	)

	participantsDB, err = w.repository.GetAllParticipantsDB()
	if err != nil {
		return nil, err
	}

	for i := range wall.ParticipantsID {
		for j := range participantsDB.Participants {
			if wall.ParticipantsID[i] == participantsDB.Participants[j].ID {
				validatesID++
			}
		}
	}

	if len(wall.ParticipantsID) > validatesID {
		err = errors.New("nem todos os participantes foram encontrados")
		return nil, err
	}

	id, err = w.repository.CreateWall(wall)
	if err != nil {
		return nil, err
	}

	return
}
