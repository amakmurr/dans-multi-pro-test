package internal

import "github.com/jmoiron/sqlx"

func initDB(dbSource string) (*sqlx.DB, error) {
	dbConn, err := sqlx.Open("pgx", dbSource)
	if err != nil {
		return nil, err
	}
	return dbConn, dbConn.Ping()
}
