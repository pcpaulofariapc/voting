package model

import (
	"time"

	"github.com/google/uuid"
)

type Wall struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name_wall"`
	CreatedAt time.Time `json:"created_at"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
type AllWalls struct {
	AllWalls []*Wall `json:"walls"`
}

type Participant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" `
}

type Participants struct {
	Participants []*Participant `json:"participants"`
}

type ParticipantActiveWall struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ParticipantsActiveWall struct {
	WallID       uuid.UUID                `json:"wall_id"`
	WallName     string                   `json:"wall_name"`
	Participants []*ParticipantActiveWall `json:"participants"`
}

type Vote struct {
	ParticipantID *uuid.UUID `json:"participant_id" binding:"required"`
	WallID        *uuid.UUID `json:"wall_id" binding:"required"`
	IP            string     `json:"-"`
}

type ParticipantCreate struct {
	Name string `json:"name" binding:"required"`
}

type WallCreate struct {
	Name           string      `json:"name" binding:"required"`
	StartTime      time.Time   `json:"start_time" binding:"required"`
	EndTime        time.Time   `json:"end_time" binding:"required"`
	ParticipantsID []uuid.UUID `json:"participants_id" binding:"required,min=2,dive"`
}

type VoteRegister struct {
	ParticipantID uuid.UUID `json:"participant_id"`
	DateRegister  time.Time `json:"date_register"`
	WallID        uuid.UUID `json:"wall_id"`
	ID            uuid.UUID `json:"id"`
	IP            string    `json:"ip"`
}

type PartialResult struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Votes           int       `json:"votes"`
	VotesPercentage float64   `json:"votes_percentage"`
}

type ResultVote struct {
	RegisterID     uuid.UUID       `json:"register_vote_id"`
	PartialResults []PartialResult `json:"partial_result"`
}

type ParticipantWall struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name_participant"`
	Votes           int       `json:"votes"`
	VotesPercentage float64   `json:"votes_percentage"`
}

type ActiveWall struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name_wall"`
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	TotalVotes   int                `json:"total_votes"`
	Participants []*ParticipantWall `json:"participants"`
}
