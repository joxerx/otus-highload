package db

import (
	"database/sql"
	"fmt"
	"log"
	"otus-highload/redis"
)

func ExecuteReadQuery(query string, args ...interface{}) (*sql.Rows, error) {
	slaveIDs := []string{"slave-1", "slave-2"}

	leastLoaded, err := redis.GetLeastLoadedSlave(slaveIDs)
	if err != nil {
		log.Println("Error getting least loaded slave:", err)
		return nil, err
	}

	slaveDB, ok := SlaveDBs[leastLoaded]
	if !ok {
		log.Println("No connection found for slave", leastLoaded)
		return nil, fmt.Errorf("no connection found for slave %s", leastLoaded)
	}

	if err := redis.IncrementSlaveCounter(leastLoaded); err != nil {
		log.Println("Error incrementing slave counter:", err)
		return nil, err
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
