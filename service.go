package main

import (
	"fmt"
	"log"
	"strings"
)

var repository *SQLiteRepository

func addSentence(keyword string, sentence string) {
	k, err := repository.GetKeywordByName(keyword)
	if err != nil {
		if err == ErrNotExists {
			k, _ = repository.CreateKeyword(Keyword{Name: keyword})
		} else {
			log.Fatal(err)
		}
	}
	s, err3 := repository.CreateSentence(Sentence{
		KeywordID: k.ID,
		Value:     sentence,
	})
	if err3 != nil {
		log.Fatal(err)
	}
	log.Println(s.ID)
}

func getRandomSentence(keyword string) (string, error) {
	k, err := repository.GetKeywordByName(keyword)
	if err != nil {
		log.Println(err)
		return "", err
	}
	s, err := repository.GetRandomSentence(*k)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return s.Value, nil
}

func getSentenceCount(keyword string) (int, error) {
	k, err := repository.GetKeywordByName(keyword)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := repository.GetSentenceCount(*k)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

func processQuery(s string) (string, error) {
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
				return "", err
			}
			return sentence, nil
		case QueryInspect:
			count, err := getSentenceCount(keyword)
			if err != nil {
				log.Println(err)
				return "", err
			}
			r := fmt.Sprintf("Keyword %s\nhas %d sentences\n", keyword, count)
			return r, nil
		default:
			log.Println("How we got here? Got no args")
		}
	} else {
		q = firstWs - 1
		switch Query(s[q]) {
		case QueryAdd:
			log.Printf("add %s\n", s[firstWs:])
			addSentence(s[:q], strings.TrimSpace(s[firstWs:]))
			return "Sentence added", nil
		default:
			log.Println("How we got here?")
		}
	}
	return "Unknown command", nil
}
