package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Word struct {
	Word      string `json:"word"`
	Translate string `json:"translate"`
}

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func GetRandomWord(table string) (Word, error) {
	var word Word
	err := DB.Raw("SELECT word, translate FROM " + table + " ORDER BY RAND() LIMIT 1").Scan(&word).Error
	return word, err
}

func GetRandomDistractors(table string, excludeTranslate string, count int) ([]string, error) {
	var distractors []string
	err := DB.Raw(`
        SELECT DISTINCT translate FROM `+table+`
        WHERE translate != ? ORDER BY RAND() LIMIT ?`, excludeTranslate, count).
		Scan(&distractors).Error
	return distractors, err
}
