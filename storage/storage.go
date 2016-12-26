package storage

import (
	"bufio"
	"database/sql"
	"os"
)

var db *sql.DB

type store struct{}

func Connect(driverName, dataSourceName string) (s store, err error) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return
	}

	db = conn

	err = conn.Ping()
	if err != nil {
		return
	}
	return
}

func (s store) Close() {
	db.Close()
}

func (s store) GetMigration(fileName string) error {
	path, _ := os.Getwd()
	file, err := os.Open(path + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	f := bufio.NewReader(file)
	for {
		str, err := f.ReadString(';')
		if err != nil {
			break
		}

		_, err = db.Exec(str)
		if err != nil {
			return err

		}
	}
	return nil
}
