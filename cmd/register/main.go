package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofrs/uuid"
	"github.com/pcpaulofariapc/voting/register/config"
	"github.com/pcpaulofariapc/voting/register/db"
	"github.com/pcpaulofariapc/voting/register/service"
)

type VoteRegister struct {
	ParticipantID uuid.UUID `json:"participant_id"`
	DateRegister  time.Time `json:"date_register"`
	WallID        uuid.UUID `json:"wall_id"`
	ID            uuid.UUID `json:"id"`
	IP            string    `json:"ip"`
}

func main() {

	config.Load()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	KafikaConsumerConnect, err := service.CreateConsumer()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to kafka")
	defer KafikaConsumerConnect.Close()

	err = receiveMessages(KafikaConsumerConnect, dbConnection)
	if err != nil {
		panic(err)
	}
}

func receiveMessages(consumer *kafka.Consumer, connection *sql.DB) error {
	topic := "votacaoVotoBBB"
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	for {
		ev := consumer.Poll(10)
		if ev == nil {
			continue
		}

		switch msg := ev.(type) {
		case *kafka.Message:
			err = createVote(connection, msg.Value)
			if err != nil {
				return err
			}

		case kafka.Error:
			fmt.Println("error:", err)
			return msg
		}
	}
}

func createVote(connection *sql.DB, vote []byte) error {
	var (
		voteRegister VoteRegister
		err          error
	)

	err = json.Unmarshal(vote, &voteRegister)
	if err != nil {
		return err
	}

	_, err = connection.Exec("INSERT INTO t_vote "+
		"(wall_id,participant_id,register_id,register_at,ip) "+
		"VALUES ($1, $2, $3, $4, $5) "+
		"ON CONFLICT (register_id) DO NOTHING", voteRegister.WallID, voteRegister.ParticipantID, voteRegister.ID, voteRegister.DateRegister, voteRegister.IP)
	if err != nil {
		return err
	}

	fmt.Println("registered vote ID: ", voteRegister.ID.String())

	return nil
}
