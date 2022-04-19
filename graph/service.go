package graph

import (
	"errors"
	"fmt"
	"log"
	"solidbroccoli/config"
	"strings"
)

var Repository *SQLiteRepository

var (
	ErrBadInput = errors.New("bad input")
)

func addSentence(keyword string, sentence string) {
	k, err := Repository.GetKeywordByName(keyword)
	if err != nil {
		if err == ErrNotExists {
			k, _ = Repository.CreateKeyword(Keyword{Name: keyword})
		} else {
			log.Fatal(err)
		}
	}
	s, err3 := Repository.CreateSentence(Sentence{
		KeywordID: k.ID,
		Value:     sentence,
	})
	if err3 != nil {
		log.Fatal(err)
	}
	log.Println(s.ID)
}

func getRandomSentence(keyword string) (string, error) {
	k, err := Repository.GetKeywordByName(keyword)
	if err != nil {
		log.Println(err)
		return "", err
	}
	s, err := Repository.GetRandomSentence(*k)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return s.Value, nil
}

func getSentenceCount(keyword string) (int, error) {
	k, err := Repository.GetKeywordByName(keyword)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := Repository.GetSentenceCount(*k)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

func ProcessQuery(s string) (string, error) {
	firstWs := strings.IndexByte(s, ' ')
	var q int
	fmt.Printf("Processing %s, index %d\n", s, firstWs)

	// TODO: this won't catch everything
	if firstWs == 0 {
		return "", ErrBadInput
	}

	if firstWs == -1 {
		q = len(s) - 1
		keyword := s[:q]
		switch config.Query(s[q]) {
		case config.QueryRandom:
			log.Printf("Get random %s\n", keyword)
			sentence, err := getRandomSentence(keyword)
			if err != nil {
				log.Println(err)
				return "", err
			}
			return sentence, nil
		case config.QueryInspect:
			count, err := getSentenceCount(keyword)
			if err != nil {
				log.Println(err)
				return "", err
			}
			r := fmt.Sprintf("ℹ️ keyword %s\nhas %d sentences\n", keyword, count)
			return r, nil
		case config.QueryAdd:
			return "❌ no sentence after +", nil
		default:
			log.Println("How we got here? Got no args")
		}
	} else {
		q = firstWs - 1
		switch config.Query(s[q]) {
		case config.QueryAdd:
			log.Printf("add %s\n", s[firstWs:])
			addSentence(s[:q], strings.TrimSpace(s[firstWs:]))
			return "✅ added", nil
		default:
			log.Println("How we got here?")
		}
	}
	return "❌ unknown command, try [.?+] at the end of keyword", nil
}
