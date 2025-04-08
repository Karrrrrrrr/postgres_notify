package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"postgres_notify/config"
)

func main() {
	db, err := gorm.Open(
		postgres.Open(config.ConnStr),
	)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		reader.ReadLine()
		// language=postgresql
		var t = time.Now().UnixMilli()
		marshal, _ := json.Marshal(t)
		err = db.Transaction(func(tx *gorm.DB) error {
			tx.Exec(fmt.Sprintf("NOTIFY events, '%s'", marshal))
			return nil
		})
		if err != nil {
			panic(err)
		}

	}
}
