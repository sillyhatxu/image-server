package db

import "github.com/sillyhatxu/mysql-client"

const (
	insert_sql = `
		INSERT INTO log 
		(log_type, content)
		VALUES (?, ?)
	`
)

func Log(logType, content string) error {
	_, err := dbclient.Client.Insert(insert_sql, logType, content)
	if err != nil {
		return err
	}
	return nil
}
