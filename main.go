package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/donovanhide/eventsource"
)

var allStudentData []studentData

func main() {

	go consumeSSE()

	fmt.Println("Listening on port 18080...")
	http.HandleFunc("/students/", getStudents)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func consumeSSE() {
	stream, err := eventsource.Subscribe("http://live-test-scores.herokuapp.com/scores", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 50; i++ {
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

func stringInSlice(a string, list []string) bool {
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
