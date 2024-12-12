package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var (
	slaveCounters = map[string]int{
		"slave-1": 0,
		"slave-2": 0,
	}
	counterMu sync.Mutex
)

func getLeastLoadedSlave(slaveIDs []string) string {
	counterMu.Lock()
	defer counterMu.Unlock()

	leastLoaded := slaveIDs[0]
	minRequests := slaveCounters[leastLoaded]

	for _, slaveID := range slaveIDs[1:] {
		if slaveCounters[slaveID] < minRequests {
			leastLoaded = slaveID
			minRequests = slaveCounters[slaveID]
		}
	}

	slaveCounters[leastLoaded]++
	return leastLoaded
}

func ExecuteReadQuery(query string, args ...interface{}) (*sql.Rows, error) {
	slaveIDs := []string{"slave-1", "slave-2"}

	leastLoaded := getLeastLoadedSlave(slaveIDs)

	slaveDB, ok := SlaveDBs[leastLoaded]
	if !ok {
		log.Println("No connection found for slave", leastLoaded)
		return nil, fmt.Errorf("no connection found for slave %s", leastLoaded)
	}

	return slaveDB.Query(query, args...)
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

	// Execute the query
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
