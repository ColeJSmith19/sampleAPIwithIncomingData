package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/donovanhide/eventsource"
)

var allStudentData []studentData

func main() {

	go consumeSSE()

	fmt.Println("Listening on port 18080...")
	http.HandleFunc("/students/", getStudents)
	http.HandleFunc("/exams/", getExams)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func consumeSSE() {
	stream, err := eventsource.Subscribe("http://live-test-scores.herokuapp.com/scores", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 150; i++ {
		var data studentData
		ev := <-stream.Events
		e := json.Unmarshal([]byte(ev.Data()), &data)
		if e != nil {
			fmt.Println(e)
			return
		}
		allStudentData = append(allStudentData, data)
	}
}

func getExams(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/exams/")
	params := strings.Split(path, "/")

	if len(params) == 0 || params[0] == "" {
		getAllExamsIDs(w)
	} else if len(params) == 1 && params[0] != "" {
		getExamByID(w, params[0])
	} else {
		http.Error(w, "This endpoint accepts up to one path parameter", 400)
		return
	}

}

func getStudents(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/students/")
	params := strings.Split(path, "/")

	if len(params) == 0 || params[0] == "" {
		getAllStudentsNames(w)
	} else if len(params) == 1 && params[0] != "" {
		getStudentByID(w, params[0])
	} else {
		http.Error(w, "This endpoint accepts up to one path parameter", 400)
		return
	}

}

func getStudentByID(w http.ResponseWriter, id string) {

	var studentDataToReturn studentAverage
	var examAverage float64

	relevantStudentData := make([]studentData, 0)

	for _, data := range allStudentData {
		if data.StudentID == id {
			relevantStudentData = append(relevantStudentData, data)
			examAverage += data.Score
		}
	}

	studentDataToReturn.StudentData = relevantStudentData
	studentDataToReturn.Average = examAverage / float64(len(relevantStudentData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentDataToReturn)
}

func getAllStudentsNames(w http.ResponseWriter) {
	name := make([]string, 0)

	for _, data := range allStudentData {
		if !stringInSlice(data.StudentID, name) {
			name = append(name, data.StudentID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(name)
}

func getExamByID(w http.ResponseWriter, id string) {

	var examDataToReturn studentAverage
	var examAverage float64

	eid, e := strconv.Atoi(id)
	if e != nil {
		http.Error(w, "Cannot convert path parameter to int", 400)
		return
	}

	relevantExamData := make([]studentData, 0)

	for _, data := range allStudentData {
		if data.Exam == eid {
			relevantExamData = append(relevantExamData, data)
			examAverage += data.Score
		}
	}

	examDataToReturn.StudentData = relevantExamData
	examDataToReturn.Average = examAverage / float64(len(relevantExamData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(examDataToReturn)
}

func getAllExamsIDs(w http.ResponseWriter) {
	ids := make([]int, 0)

	for _, data := range allStudentData {
		if !intInSlice(data.Exam, ids) {
			ids = append(ids, data.Exam)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type studentData struct {
	StudentID string  `json:"studentId"`
	Exam      int     `json:"exam"`
	Score     float64 `json:"score"`
}

type studentAverage struct {
	Average     float64       `json:"average"`
	StudentData []studentData `json:"student_data"`
}
