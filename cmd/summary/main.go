package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gofrs/uuid"
	"github.com/pcpaulofariapc/voting/summary/config"
	"github.com/pcpaulofariapc/voting/summary/db"
	"github.com/pcpaulofariapc/voting/summary/service"
)

type ParticipantWall struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name_participant"`
	Votes           int       `json:"votes"`
	VotesPercentage float64   `json:"votes_percentage"`
}

type ActiveWall struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name_wall"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	TotalVotes   int                `json:"total_votes"`
	Participants []*ParticipantWall `json:"participants"`
}

func main() {

	config.Load()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	memcacheConnect, err := service.ConnectMemcache()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to memcache")
	defer memcacheConnect.Close()

	for {
		summary(dbConnection, memcacheConnect)
	}
}

func summary(dbConnection *sql.DB, memcacheConnect *memcache.Client) {
	var (
		err    error
		wallDB *ActiveWall
		wallMC *ActiveWall
	)

	config := config.Get().Application
	tk := time.NewTicker(time.Duration(config.ExecutionTime) * time.Second)
	defer tk.Stop()

	for range tk.C {
		wallDB, err = getActiveWallDB(dbConnection)
		if err != nil {
			panic(err)
		}

		wallMC, err = getActiveWallMemcache(memcacheConnect)
		if err != nil {
			panic(err)
		}

		if wallMC == nil {
			fmt.Println("wallMC NILL")
		} else {
			fmt.Println("wallMC NOT NILL")
		}

		// Atualiza o paredeão no memcache
		if wallDB != nil && wallMC != nil {
			err = updateActiveWallMemcache(memcacheConnect, wallDB)
			if err != nil {
				panic(err)
			}
		}
		// Insere o paredeão no memcache
		if wallDB != nil && wallMC == nil {
			err = insertActiveWallMemcache(memcacheConnect, wallDB)
			if err != nil {
				panic(err)
			}
		}
		// Remove o paredão do memcache
		if wallDB == nil && wallMC != nil {
			err = deleteActiveWallMemcache(memcacheConnect)
			if err != nil {
				panic(err)
			}
		}

		if wallMC != nil {
			printWall(wallMC)
		}

	}
}

func getActiveWallDB(dbConnection *sql.DB) (activeWall *ActiveWall, err error) {
	var (
		rows   *sql.Rows
		result bool
	)
	rows, err = dbConnection.Query("SELECT WALL.id, WALL.name_wall, WALL.start_time::TIMESTAMPTZ, WALL.end_time::TIMESTAMPTZ, PARTICIPANT.id AS participant_id, PARTICIPANT.name_participant, COUNT(VOTE.id) AS total_votes " +
		"FROM t_wall WALL INNER JOIN t_wall_participant WALLPARTICIPANT ON WALL.id =  WALLPARTICIPANT.wall_id INNER JOIN t_participant PARTICIPANT ON WALLPARTICIPANT.participant_id = PARTICIPANT.id LEFT JOIN t_vote VOTE ON PARTICIPANT.id = VOTE.participant_id AND WALL.ID = VOTE.wall_id " +
		"WHERE WALL.deleted_at is NULL AND WALL.start_time < NOW() AND WALL.end_time > NOW() " +
		"GROUP BY WALL.id, PARTICIPANT.id, VOTE.participant_ID ORDER BY TOTAL_VOTES ASC")
	if err != nil {
		return nil, err
	}

	activeWall = new(ActiveWall)
	for rows.Next() {
		result = true
		participant := new(ParticipantWall)
		err = rows.Scan(
			&activeWall.ID,
			&activeWall.Name,
			&activeWall.StartTime,
			&activeWall.EndTime,
			&participant.ID,
			&participant.Name,
			&participant.Votes)
		if err != nil {
			return nil, err
		}
		activeWall.Participants = append(activeWall.Participants, participant)
	}

	if result {
		calculatePercentage(activeWall)
	} else {
		return nil, nil
	}
	return
}

func getActiveWallMemcache(memcacheConnect *memcache.Client) (activeWall *ActiveWall, err error) {

	var data *memcache.Item
	data, err = memcacheConnect.Get("ActiveWall")
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if data != nil {
		var activeWallFound ActiveWall
		err = json.Unmarshal(data.Value, &activeWallFound)
		if err != nil {
			return nil, err
		}
		return &activeWallFound, nil
	}

	return nil, nil
}

func insertActiveWallMemcache(memcacheConnect *memcache.Client, activeWall *ActiveWall) (err error) {
	var dataInsert []byte
	dataInsert, err = json.Marshal(activeWall)
	if err != nil {
		return err
	}

	expiration := int32(0)
	err = memcacheConnect.Add(&memcache.Item{Key: "ActiveWall", Value: dataInsert, Expiration: expiration})
	if err != nil {
		return err
	}
	return
}

func updateActiveWallMemcache(memcacheConnect *memcache.Client, activeWall *ActiveWall) (err error) {
	var dataInsert []byte
	dataInsert, err = json.Marshal(activeWall)
	if err != nil {
		return err
	}

	expiration := int32(0)
	err = memcacheConnect.Replace(&memcache.Item{Key: "ActiveWall", Value: dataInsert, Expiration: expiration})
	if err != nil {
		return err
	}
	return
}

func deleteActiveWallMemcache(memcacheConnect *memcache.Client) (err error) {
	err = memcacheConnect.Delete("ActiveWall")
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}
	return
}

func calculatePercentage(activeWall *ActiveWall) {
	var totalVotes int
	for i := range activeWall.Participants {
		totalVotes = totalVotes + activeWall.Participants[i].Votes
	}
	activeWall.TotalVotes = totalVotes
	for i := range activeWall.Participants {
		percentage := float64(activeWall.Participants[i].Votes) / float64(totalVotes)
		activeWall.Participants[i].VotesPercentage = percentage
	}

}

func printWall(wall *ActiveWall) {
	if wall != nil {
		fmt.Printf("\n\nWall.ID: %s", wall.ID.String())
		fmt.Printf("\nWall.Name: %s", wall.Name)
		fmt.Printf("\nWall.StartTime: %s", wall.StartTime.String())
		fmt.Printf("\nWall.EndTime: %s", wall.EndTime.String())
		fmt.Printf("\nWall.TotalVotes: %d", wall.TotalVotes)
		for i := range wall.Participants {
			fmt.Printf("\nParticipant[%d].Name: %s", i, wall.Participants[i].Name)
			fmt.Printf("\nParticipant[%d].ID: %s", i, wall.Participants[i].ID.String())
			fmt.Printf("\nParticipant[%d].Votes: %d", i, wall.Participants[i].Votes)
			fmt.Printf("\nParticipant[%d].VotesPercentage: %f", i, wall.Participants[i].VotesPercentage)
		}
	}
}
