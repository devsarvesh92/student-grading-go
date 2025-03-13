package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var students []student

	for idx, record := range records {
		//skip header row
		if idx == 0 {
			continue
		}
		parseScore := func(score string) int {
			result, _ := strconv.Atoi(score)
			return result
		}
		students = append(students, student{
			firstName:  record[0],
			lastName:   record[1],
			university: record[2],
			test1Score: parseScore(record[3]),
			test2Score: parseScore(record[4]),
			test3Score: parseScore(record[5]),
			test4Score: parseScore(record[6]),
		})
	}
	return students
}

func calculateGrade(students []student) []studentStat {
	var studentStats []studentStat

	getGrade := func(score float32) Grade {
		switch {
		case score >= 70:
			return A
		case score >= 50 && score < 70:
			return B
		case score >= 35 && score < 50:
			return C
		default:
			return F
		}
	}

	for _, student := range students {
		finalScore := (float32(student.test1Score) + float32(student.test2Score) + float32(student.test3Score) + float32(student.test4Score)) / 4
		studentStats = append(studentStats, studentStat{
			student:    student,
			finalScore: finalScore,
			grade:      getGrade(finalScore),
		})
	}
	return studentStats
}

func findOverallTopper(gradedStudents []studentStat) studentStat {

	currentTopper := gradedStudents[0]

	for i := 1; i < len(gradedStudents); i++ {
		if gradedStudents[i].finalScore > currentTopper.finalScore {
			currentTopper = gradedStudents[i]
		}
	}

	return currentTopper
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	universityToppers := make(map[string]studentStat)

	for _, studentStat := range gs {
		currentTopper := universityToppers[studentStat.university]
		if studentStat.finalScore > currentTopper.finalScore {
			universityToppers[studentStat.university] = studentStat
		}
	}
	return universityToppers
}
