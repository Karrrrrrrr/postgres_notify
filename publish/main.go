package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"postgres_notify/config"
)

func Escape(s string) string {
	var result bytes.Buffer
	for _, c := range s {
		if c == '\'' {
			result.WriteString("\\")
		}
		result.WriteRune(c)
	}

	return result.String()

}
func main() {
	db, err := gorm.Open(
		postgres.Open(config.ConnStr),
	)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = line[:len(line)-1]
		if len(line) == 0 {
			continue
		}
		// language=postgresql
		//var t = time.Now().UnixMilli()

		err = db.Transaction(func(tx *gorm.DB) error {
			tx.Exec(fmt.Sprintf("NOTIFY events, '%s'", Escape(line)))
			return nil
		})
		if err != nil {
			panic(err)
		}

	}
}
