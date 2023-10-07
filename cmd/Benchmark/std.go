package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

const jsonContentType = "application/json"

type student struct {
	id            uint64
	password      string
	auth          string
	toPickCourses []course
}

type course struct {
	courseID int32
	groupID  int32
}

type tokenResult struct {
	Token string `json:"token"`
}

func authorizeStudent(base string, student *student) {
	body := strings.NewReader(fmt.Sprintf(`{"user":%d,"password":"%s","staff":false}`, student.id, student.password))
	resp, err := http.Post(base+"/login", jsonContentType, body)
	if err != nil {
		log.WithError(err).WithField("id", student.id).Fatal("cannot start login request")
	}
	if resp.StatusCode != http.StatusOK {
		log.WithField("status", resp.Status).WithField("id", student.id).Fatal("cannot login user")
	}
	var response tokenResult
	_ = json.NewDecoder(resp.Body).Decode(&response)
	student.auth = "Bearer " + response.Token
	_ = resp.Body.Close()
}

func registerCourses(totalRequests, failedRequests *atomic.Int32, done *sync.WaitGroup, barrier chan struct{}, base string, student *student) {
	defer done.Done()
	<-barrier // wait for it
	// Create the client and do the requests
	client := http.Client{}
	for _, course := range student.toPickCourses {
		totalRequests.Add(1)
		body := strings.NewReader(fmt.Sprintf(`{"course_id":%d,"group_id":%d}`, course.courseID, course.groupID))
		req, _ := http.NewRequest("PUT", base+"/student/course", body)
		req.Header.Set("Content-Type", jsonContentType)
		req.Header.Set("Authorization", student.auth)
		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).WithField("std", student.id).WithField("course", course).Error("cannot register course")
			continue
		}
		if resp.StatusCode >= 300 {
			//body, _ := io.ReadAll(resp.Body)
			//log.WithField("body", string(body)).WithField("std", student.id).WithField("course", course).Warn("cannot register course:", resp.Status)
			failedRequests.Add(1)
		}
		_ = resp.Body.Close()
	}
}
