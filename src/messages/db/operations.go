package db

import (
	"database/sql"
	"log"
)

func ExecuteReadQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return BalancerDB.Query(query, args...)
}

func ExecuteInsertQuery(query string, args ...interface{}) (string, error) {
	var id string
	err := MasterDB.QueryRow(query, args...).Scan(&id)
	if err != nil {
		log.Printf("Database write error: %v", err)
		return "", err
	}
	return id, nil
}

func ExecuteWriteQuery(query string, args ...interface{}) error {
	stmt, err := MasterDB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return err
	}

	return nil
}

func ExecuteUpdateQuery(query string, args ...interface{}) (int64, error) {
	stmt, err := MasterDB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v\n", err)
		return 0, err
	}

	return rowsAffected, nil
}
