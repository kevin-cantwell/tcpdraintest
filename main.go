package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	n, _ := strconv.Atoi(os.Args[1])

	go func() {
		time.Sleep(1 * time.Second) // Give the web server plenty of time to begin serving requests

		client := http.DefaultClient

		d1 := testWithoutDrain(client, n)
		d2 := testWithDrain(client, n)

		fmt.Printf("Without drain (x%d): %v\n", n, d1)
		fmt.Printf("   With drain (x%d): %v\n", n, d2)

		os.Exit(0)
	}()

	http.HandleFunc("/some.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}\n\n")) // Adds extra bytes to the json response
	})
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func testWithoutDrain(client *http.Client, n int) time.Duration {
	var total time.Duration
	for i := 0; i < n; i++ {
		start := time.Now()
		resp, err := client.Get("http://localhost:" + os.Getenv("PORT") + "/some.json")
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(resp.Body)
		var body map[string]interface{}
		if err := dec.Decode(&body); err != nil {
			panic(err)
		}
		timing := time.Since(start)
		total += timing

		// Don't drain, just close
		resp.Body.Close()
	}
	return total
}

func testWithDrain(client *http.Client, n int) time.Duration {
	var total time.Duration
	for i := 0; i < n; i++ {
		start := time.Now()
		resp, err := client.Get("http://localhost:" + os.Getenv("PORT") + "/some.json")
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(resp.Body)
		var body map[string]interface{}
		if err := dec.Decode(&body); err != nil {
			panic(err)
		}
		timing := time.Since(start)
		total += timing

		// Drain and close
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return total
}
