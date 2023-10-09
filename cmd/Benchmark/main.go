package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().Unix()))

func main() {
	// Get initial values
	studentCount, err := strconv.ParseUint(os.Getenv("STD_COUNT"), 10, 64)
	if err != nil {
		log.WithError(err).Fatal("cannot parse STD_COUNT")
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("please set DATABASE_URL")
	}
	authCoreURL := os.Getenv("AUTH_CORE_URL")
	if authCoreURL == "" {
		log.Fatal("please set AUTH_CORE_URL")
	}
	// Connect to database and create users
	students := setupDatabase(databaseUrl, studentCount)
	log.Info("Press enter to start...")
	fmt.Scanln()
	// Authorize users and spawn the register routines
	startRegisterBarrier := make(chan struct{}) // close this channel to register everyone
	var totalRequests, failedRequests atomic.Int32
	doneWg := new(sync.WaitGroup)
	doneWg.Add(len(students))
	log.Info("authorizing students")
	for _, std := range students {
		authorizeStudent(authCoreURL, std)
		go registerCourses(&totalRequests, &failedRequests, doneWg, startRegisterBarrier, authCoreURL, std)
	}
	// Wait for user
	runtime.GC()
	// Start the timer and wait
	log.Info("starting the registration")
	close(startRegisterBarrier)
	now := time.Now()
	doneWg.Wait()
	log.Info("Done in ", time.Since(now))
	log.Info("Total requests: ", totalRequests.Load())
	log.Info("Failed requests: ", failedRequests.Load())
}
