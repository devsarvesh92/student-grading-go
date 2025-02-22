package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
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
		fmt.Printf("Error %v while reading file %v \n", err, file)
	}

	defer file.Close()

	csv := csv.NewReader(file)

	records, err := csv.ReadAll()

	students := make([]student, 0, len(records))

	if err != nil {
		fmt.Printf("Error while reading csv %v \n", err)
	}

	for id, record := range records {
		if id == 0 {
			continue
		}

		test1Score, _ := strconv.Atoi(record[3])
		test2Score, _ := strconv.Atoi(record[4])
		test3Score, _ := strconv.Atoi(record[5])
		test4Score, _ := strconv.Atoi(record[6])

		students = append(students, student{
			firstName:  record[0],
			lastName:   record[1],
			university: record[2],
			test1Score: test1Score,
			test2Score: test2Score,
			test3Score: test3Score,
			test4Score: test4Score,
		})

	}
	return students
}

func calculateGrade(students []student) []studentStat {
	allstudentStats := make([]studentStat, 0, len(students))

	var getGrade = func(finalScore float32) Grade {
		switch {
		case finalScore >= 70:
			return A
		case finalScore >= 50:
			return B
		case finalScore >= 35:
			return C
		default:
			return F
		}
	}

	for _, student := range students {
		finalScore := (float32(student.test1Score) + float32(student.test2Score) + float32(student.test3Score) + float32(student.test4Score)) / 4
		grade := getGrade(finalScore)
		allstudentStats = append(allstudentStats, studentStat{
			student:    student,
			finalScore: finalScore,
			grade:      grade,
		})
	}
	return allstudentStats
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	sort.Slice(gradedStudents, func(i, j int) bool {
		return gradedStudents[i].finalScore > gradedStudents[j].finalScore
	})
	return gradedStudents[0]
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	studetsPerUniversity := make(map[string][]studentStat)
	universityTopper := make(map[string]studentStat)

	for _, st := range gs {
		studetsPerUniversity[st.student.university] = append(studetsPerUniversity[st.student.university], st)
	}

	for key, value := range studetsPerUniversity {
		sort.Slice(value, func(i, j int) bool { return value[i].finalScore > value[j].finalScore })
		universityTopper[key] = value[0]
	}

	return universityTopper
}
