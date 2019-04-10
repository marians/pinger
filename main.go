package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
)

// Check represents an HTTP endpoint to be checked
type Check struct {
	// URL is the URL of the endpoint
	URL string `yaml:url`
	// Method is the HTTP method to use
	Method string `yaml:method`
}

// Checks represents our list of checks as defined in the config,
// with defaults applied.
var (
	Checks     = []Check{}
	OneDay     = time.Duration(24) * time.Hour
	TimeFormat = "2006-01-02T15:04:05-0700"
)

// ReadConfig reads our checks configuration and applies defaults where
// details are not specified.
func ReadConfig(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file %s: #%v ", path, err)
	}
	err = yaml.Unmarshal(yamlFile, &Checks)
	if err != nil {
		log.Fatalf("Could not parse YAML file %s: %v", path, err)
	}

	// Apply defaults and sanitize
	for i, c := range Checks {
		if c.Method == "" {
			Checks[i].Method = "GET"
		} else {
			Checks[i].Method = strings.ToUpper(Checks[i].Method)
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("Missing required argument. Please give config file path as first argument.")
	}

	interval, _ := time.ParseDuration("3m")

	client := &http.Client{}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Cannot connect to redis on localhost:6379.")
	}

	for {
		ReadConfig(args[0])

		for _, check := range Checks {
			go func(check Check) {
				req, err := http.NewRequest(check.Method, check.URL, nil)
				if err != nil {
					log.Printf("%s %s failed, request could not be created: %s", check.Method, check.URL, err)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("%s %s failed: %s", check.Method, check.URL, err)
				} else if resp != nil {
					if resp.StatusCode >= 400 {
						log.Printf("%s %s HTTP Status %d", check.Method, check.URL, resp.StatusCode)
					} else {
						// store latest success in DB
						timeString := time.Now().UTC().Format(TimeFormat)
						key := fmt.Sprintf("%s %s", check.Method, check.URL)
						err := redisClient.Set(key, timeString, OneDay).Err()
						if err != nil {
							log.Fatalf("Error writing data to redis: %s", err)
						}
					}
				}
			}(check)
		}

		time.Sleep(interval)
	}
}
