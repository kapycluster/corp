package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/exp/rand"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("db path required as first argument")
		return
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	// Generate random invite code
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	code := make([]byte, 8)
	for i := range code {
		code[i] = letters[r.Intn(len(letters))]
	}
	inviteCode := string(code)

	// Insert into database
	_, err = db.Exec("INSERT INTO invites (id, used) VALUES (?, ?)", inviteCode, 0)
	if err != nil {
		fmt.Printf("error inserting invite: %v\n", err)
		return
	}

	fmt.Printf("generated invite code: %s\n", inviteCode)
}
