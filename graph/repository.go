package graph

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS keywords(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE
    );
    CREATE TABLE IF NOT EXISTS sentences(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        keyword_id INTEGER NOT NULL,
        sentence TEXT NOT NULL,
        FOREIGN KEY(keyword_id) REFERENCES keywords(id)
    );
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateSentence(sentence Sentence) (*Sentence, error) {
	res, err := r.db.Exec("INSERT INTO sentences(keyword_id, sentence) values(?,?)", sentence.KeywordID, sentence.Value)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	sentence.ID = id

	return &sentence, nil
}

func (r *SQLiteRepository) CreateKeyword(website Keyword) (*Keyword, error) {
	res, err := r.db.Exec("INSERT INTO keywords(name) values(?)", website.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	website.ID = id

	return &website, nil
}

func (r *SQLiteRepository) AllSentences() ([]Sentence, error) {
	rows, err := r.db.Query("SELECT * FROM sentences")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Sentence
	for rows.Next() {
		var sentence Sentence
		if err := rows.Scan(&sentence.ID, &sentence.KeywordID, &sentence.Value); err != nil {
			return nil, err
		}
		all = append(all, sentence)
	}
	return all, nil
}

func (r *SQLiteRepository) GetRandomSentence(keyword Keyword) (*Sentence, error) {
	row := r.db.QueryRow(
		"SELECT * FROM sentences WHERE id IN (SELECT id FROM sentences WHERE keyword_id = ? ORDER BY RANDOM() LIMIT 1)",
		keyword.ID)

	var sentence Sentence
	if err := row.Scan(&sentence.ID, &sentence.KeywordID, &sentence.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &sentence, nil
}

func (r *SQLiteRepository) GetSentenceCount(keyword Keyword) (int, error) {
	row := r.db.QueryRow(
		"SELECT COUNT(id) as count_sentence FROM sentences WHERE keyword_id = ?",
		keyword.ID)

	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotExists
		}
		return 0, err
	}
	return count, nil
}

func (r *SQLiteRepository) GetKeywordByName(name string) (*Keyword, error) {
	row := r.db.QueryRow("SELECT * FROM keywords WHERE name = ?", name)

	var keyword Keyword
	if err := row.Scan(&keyword.ID, &keyword.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &keyword, nil
}

func (r *SQLiteRepository) DeleteSentence(id int64) error {
	res, err := r.db.Exec("DELETE FROM sentences WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
