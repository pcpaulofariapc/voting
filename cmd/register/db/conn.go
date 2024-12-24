package db

import (
	"database/sql"
	"fmt"

	"github.com/pcpaulofariapc/voting/register/config"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {

	database := config.Get().Database

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		database.Host, database.Port, database.User, database.Password, database.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to " + database.Name)

	return db, nil
}
