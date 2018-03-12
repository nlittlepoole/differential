package main

import (
	"log"
	"fmt"
	"math"
	"math/rand"
	"github.com/microo8/plgo"
)

//Test method Meh prints out message to error elog
func Meh() {
	logger := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	logger.Println("meh")
}

// RandomResponse alters a binary response based on Randomized Response Algorithm
func RandomResponse(truth bool, alpha float64, beta float64) binary {
	if rand.Float64() < alpha {
		return truth
	} else {
		if rand.Float64() < beta {
			return truth
		} else {
			return !truth
		}
	}
}

// Computes real probability of true answer, given probability of true random response answers
func ProbabilityRandomResponse(probability float64, alpha float64, beta float64) float64 {
	return (probability - math.Pow(1 - alpha), beta) / alpha
} 

//VarianceRandomResponse computes the variance of the real probability
func VarianceRandomResponse(tableName string, colName string, alpha float64) float64 {
	logger := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	db, err := plgo.Open()
	if err != nil {
		logger.Fatalf("Cannot open DB: %s", err)
	}
	defer db.Close()
	query := fmt.Sprintf(
		`SELECT count(*), variance(%s) from %s`,
		colName,
		tableName,
	)
	stmt, err := db.Prepare(query, nil)
	if err != nil {
		logger.Fatalf("Cannot prepare query statement (%s): %s", query, err)
	}
	rows, err := stmt.Query()
	if err != nil {
		logger.Fatalf("Query (%s) error: %s", query, err)
	}
	var count int64
	var variance float64

	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			logger.Fatalln("Cannot get columns", err)
		}
		logger.Println(cols)
		err = rows.Scan(&count, &variance)
		if err != nil {
			logger.Fatalln("Cannot scan values", err)
		}
	}
	return variance/ (count * alpha * alpha)
}