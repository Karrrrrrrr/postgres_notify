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
	fmt.Println(db)
	reader := bufio.NewReader(os.Stdin)
	for {
		// language=postgresql
		reader.ReadLine()
		db.Transaction(func(tx *gorm.DB) error {
			var t = time.Now().UnixMilli()
			marshal, _ := json.Marshal(t)
			fmt.Println(string(marshal))
			tx.Exec(fmt.Sprintf("NOTIFY events, '%s'", marshal))
			if t%2 == 0 {
				return gorm.ErrRecordNotFound
			}
			return nil
		})

	}
}
