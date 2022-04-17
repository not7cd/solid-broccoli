package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

const PrefixCmd = '.'

type Query byte

const (
	QueryAdd     Query = '+'
	QueryRandom  Query = '.' // no arg
	QueryInspect Query = '?' // no arg
)

func processQuery(s string) {
	firstWs := strings.IndexByte(s, ' ')
	var q int
	fmt.Printf("Processing %s, index %s\n", s, firstWs)

	if firstWs == -1 {
		q = len(s) - 1
		keyword := s[:q]
		switch Query(s[q]) {
		case QueryRandom:
			log.Printf("Get random %s\n", keyword)
			sentence, err := getRandomSentence(keyword)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(sentence)
		case QueryInspect:
			count, err := getSentenceCount(keyword)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Keyword %s\nhas %d sentences\n", keyword, count)
		default:
			log.Println("How we got here? Got no args")
		}
	} else {
		q = firstWs - 1
		switch Query(s[q]) {
		case QueryAdd:
			log.Printf("add %s\n", s[firstWs:])
			addSentence(s[:q], strings.TrimSpace(s[firstWs:]))
		default:
			log.Println("How we got here?")
		}
	}

}

const fileName = "sqlite.db"

func main() {
	log.Println("Hey")

	in := bufio.NewReader(os.Stdin)

	_ = os.Remove(fileName)
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	repository = NewSQLiteRepository(db)
	if err := repository.Migrate(); err != nil {
		log.Fatal(err)
	}

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			log.Println("Bad input")
		}
		// sanitize it a bit
		line = strings.TrimSpace(line)
		if line[0] == PrefixCmd {
			log.Printf("Got command\n")
			processQuery(line[1:])
		}
	}
}
