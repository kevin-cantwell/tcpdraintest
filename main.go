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
	url := os.Args[1]
	n, _ := strconv.Atoi(os.Args[2])

	client := http.DefaultClient

	d1 := run(client, url, n, false)
	d2 := run(client, url, n, true)

	fmt.Printf("Without drain (x%d): %v\n", n, d1)
	fmt.Printf("   With drain (x%d): %v\n", n, d2)

	os.Exit(0)
}

func run(client *http.Client, url string, n int, drain bool) time.Duration {
	var total time.Duration
	for i := 0; i < n; i++ {
		start := time.Now()
		resp, err := client.Get(url)
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

		if drain {
			io.Copy(ioutil.Discard, resp.Body)
		}
		resp.Body.Close()
	}
	return total
}
