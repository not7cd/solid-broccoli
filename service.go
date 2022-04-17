package main

import "log"

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

