package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/IlGuilhermelI/TestNeoWay/db"
	"github.com/IlGuilhermelI/TestNeoWay/dto"
)

var dataBase *sql.DB
var wg sync.WaitGroup

func main() {
	start := time.Now()
	dataBase = db.Connect()
	fullFilePath := "./db/base_teste.txt"
	insertData(fullFilePath)
	elapsed := time.Since(start)
	wg.Wait()
	log.Printf("Time execution application: %s", elapsed)

}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Failed try open file: %s", err)
	}
}

func insertData(fullFilePath string) {
	file, err := os.Open(fullFilePath)
	checkError(err)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fullTextFile []string

	for scanner.Scan() {
		fullTextFile = append(fullTextFile, scanner.Text())
	}

	file.Close()

	//insertInDataBase(fullTextFile, 1, int64(len(fullTextFile)-1))
	defineGoRoutinesExecution(fullTextFile)
}
func defineLineRange(text []string, start int64, finish int64) (lineRange []string) {
	if start == 0 {
		start = 1
	}
	if finish != 0 {
		lineRange = text[start:finish]
	} else {
		lineRange = text[start:]
	}
	return lineRange
}
func insertInDataBase(text []string, start int64, finish int64) {
	lineRange := defineLineRange(text, start, finish)
	for _, currentLine := range lineRange {
		if dto.ValidateAllCpfAndCnpjLine(currentLine) {
			_, err := dataBase.Exec(dto.InsertClientsPurchaseInformations(currentLine))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	wg.Done()
}
func defineGoRoutinesExecution(fullTextFile []string) {
	lineRangeToExecute := len(fullTextFile) / runtime.NumCPU()
	lastRun := false
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		if i == runtime.NumCPU()-1 {
			lastRun = true
		}
		execute(fullTextFile, i, lineRangeToExecute, lastRun)
	}

}
func execute(fullTextFile []string, currentGoRoutine int, interval int, lastRun bool) {
	start := int64(currentGoRoutine * interval)
	finish := int64(currentGoRoutine*interval + interval)
	if lastRun {
		finish = 0
	}
	go insertInDataBase(fullTextFile, start, finish)
}
