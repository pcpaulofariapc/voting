package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pcpaulofariapc/voting/model"
)

type WallRepository struct {
	connectionDB    *sql.DB
	connectionMem   *memcache.Client
	connectionKafka *kafka.Producer
}

func NewWallRepository(connectionDB *sql.DB, connectionMem *memcache.Client, connectionKafka *kafka.Producer) WallRepository {
	return WallRepository{
		connectionDB:    connectionDB,
		connectionMem:   connectionMem,
		connectionKafka: connectionKafka,
	}
}

func (w *WallRepository) GetAllWallsDB() (walls *model.AllWalls, err error) {
	var rows *sql.Rows

	rows, err = w.connectionDB.Query("SELECT id, name_wall, created_at::TIMESTAMPTZ, start_time::TIMESTAMPTZ, end_time::TIMESTAMPTZ FROM t_wall WHERE deleted_at IS NULL ORDER BY start_time DESC")
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("não foram encontrados paredões")
			return nil, err
		}
		return nil, err
	}

	walls = new(model.AllWalls)
	for rows.Next() {
		wall := new(model.Wall)
		err = rows.Scan(
			&wall.ID,
			&wall.Name,
			&wall.CreatedAt,
			&wall.StartTime,
			&wall.EndTime,
		)
		if err != nil {
			return nil, err
		}
		walls.AllWalls = append(walls.AllWalls, wall)
	}
	return
}

func (w *WallRepository) GetAllParticipantsDB() (participants *model.Participants, err error) {
	var rows *sql.Rows

	rows, err = w.connectionDB.Query("SELECT id, name_participant, created_at::TIMESTAMPTZ FROM t_participant WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("não foram encontrados participantes")
			return nil, err
		}
		return nil, err
	}

	participants = new(model.Participants)
	for rows.Next() {
		wall := new(model.Participant)
		err = rows.Scan(
			&wall.ID,
			&wall.Name,
			&wall.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		participants.Participants = append(participants.Participants, wall)
	}
	return
}

func (w *WallRepository) RegsiterVote(message []byte) error {
	topic := "votacaoVotoBBB"
	deliveryChan := make(chan kafka.Event)

	err := w.connectionKafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Key:            []byte("votos2"),
	}, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)

	return msg.TopicPartition.Error
}

func (w *WallRepository) GetActiveWall() (activeAll *model.ActiveWall, err error) {
	var data *memcache.Item

	data, err = w.connectionMem.Get("ActiveWall")
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}

	err = json.Unmarshal(data.Value, &activeAll)
	if err != nil {
		return nil, err
	}

	return
}

func (w *WallRepository) GetWallByID(wallID uuid.UUID) (wall *model.ActiveWall, err error) {
	var (
		rows  *sql.Rows
		query *sql.Stmt
	)

	query, err = w.connectionDB.Prepare("SELECT WALL.id, WALL.name_wall, WALL.start_time::TIMESTAMPTZ, WALL.end_time::TIMESTAMPTZ, PARTICIPANT.id AS participant_id, PARTICIPANT.name_participant, COUNT(VOTE.id) AS total_votes " +
		"FROM t_wall WALL LEFT JOIN t_wall_participant WALLPARTICIPANT ON WALL.id =  WALLPARTICIPANT.wall_id LEFT JOIN t_participant PARTICIPANT ON WALLPARTICIPANT.participant_id = PARTICIPANT.id LEFT JOIN t_vote VOTE ON PARTICIPANT.id = VOTE.participant_id AND WALL.ID = VOTE.wall_id " +
		"WHERE WALL.deleted_at IS NULL AND WALL.id = $1::UUID " +
		"GROUP BY WALL.id, PARTICIPANT.id, VOTE.participant_ID ORDER BY TOTAL_VOTES ASC")
	if err != nil {
		return nil, err
	}

	rows, err = query.Query(wallID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("paredão não foi encontrado")
			return nil, err
		}
		return nil, err
	}

	wall = new(model.ActiveWall)
	for rows.Next() {
		participant := new(model.ParticipantWall)
		err = rows.Scan(
			&wall.ID,
			&wall.Name,
			&wall.StartTime,
			&wall.EndTime,
			&participant.ID,
			&participant.Name,
			&participant.Votes)
		if err != nil {
			return nil, err
		}
		wall.Participants = append(wall.Participants, participant)
	}

	calculatePercentage(wall)
	return
}

func (w *WallRepository) CreateParticipant(participant model.ParticipantCreate) (id *uuid.UUID, err error) {
	var query *sql.Stmt
	id = new(uuid.UUID)

	query, err = w.connectionDB.Prepare("INSERT INTO t_participant" +
		"(name_participant)" +
		" VALUES ($1) RETURNING id")
	if err != nil {
		return nil, err
	}

	err = query.QueryRow(participant.Name).Scan(id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (w *WallRepository) CreateWall(wall model.WallCreate) (id *uuid.UUID, err error) {
	var query *sql.Stmt
	id = new(uuid.UUID)

	query, err = w.connectionDB.Prepare("INSERT INTO t_wall" +
		"(name_wall, start_time, end_time)" +
		" VALUES ($1,$2,$3) RETURNING id")
	if err != nil {
		return nil, err
	}

	err = query.QueryRow(wall.Name, wall.StartTime, wall.EndTime).Scan(id)
	if err != nil {
		return nil, err
	}

	query, err = w.connectionDB.Prepare("INSERT INTO t_wall_participant" +
		"(wall_id, participant_id)" +
		" VALUES ($1::UUID, unnest($2::UUID[]))")
	if err != nil {
		return nil, err
	}

	_, err = query.Exec(*id, pq.Array(wall.ParticipantsID))
	if err != nil {
		return nil, err
	}

	return id, nil
}

func calculatePercentage(wall *model.ActiveWall) {
	var totalVotes int

	for i := range wall.Participants {
		totalVotes = totalVotes + wall.Participants[i].Votes
	}
	wall.TotalVotes = totalVotes

	for i := range wall.Participants {
		percentage := float64(wall.Participants[i].Votes) / float64(totalVotes)
		wall.Participants[i].VotesPercentage = percentage
	}
}
