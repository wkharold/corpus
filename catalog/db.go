package catalog

import (
	"fmt"
	"os"
	"path"

	"github.com/cznic/ql"
	"github.com/wkharold/corpus/catalog"
)

type MbxMsg struct {
	ID      int64
	Mailbox string
	Folder  string
	Msgfile string
}

var (
	schema ql.List
	insmsg ql.List
)

func New(dir string) func(chan MbxMsg, chan chan int) {
	_, err := os.Open(dir)
	if err != nil {
		switch err.(*os.PathError).Err.Error() {
		case "no such file or directory":
			if os.Mkdir(ServerDir, os.ModePerm) != nil {
				panic(fmt.Sprintf("Can't create server directory %s [%v]", dir, err))
			}
		default:
			panic(fmt.Sprintf("Unexpected error [%v]", err))
		}
	}

	db, err := ql.OpenFile(path.Join(dir, "catalog.db"), &ql.Options{true, nil, nil})
	if err != nil {
		panic(fmt.Sprintf("Can't open catalog database [%v]", err))
	}

	schema = ql.MustSchema((*catalog.MbxMsg)(nil), "", nil)
	insmsg = ql.MustCompile(`
		BEGIN TRANSACTION;
			INSERT INTO MbxMsg VALUES($1, $2, $3);
		COMMIT;`)

	if _, _, err = db.Execute(ql.NewRWCtx(), schema); err != nil {
		panic(fmt.Sprintf("Can't create schema [%v]", err))
	}
	db.Close()

	return listener
}

func listener(mc chan MbxMsg, done chan chan int) {
	var rc chan int
	count := 0

loop:
	for {
		select {
		case msg := <-mc:
			if _, _, err := db.Execute(ql.NewRWCtx(), insmsg, ql.MustMarshal(&msg)...); err != nil {
				db.Close()
				panic(fmt.Sprintf("Message insert failed [%v]", err))
			}
			count++
		case rc = <-done:
			break loop
		}
	}

	db.Close()
	rc <- count
}
