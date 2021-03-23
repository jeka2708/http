package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func getResponse(url string, timeout int, dataResponses map[string][]float64, noResponses map[string]int) {

	start := time.Now()
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	result, err := client.Get(url)
	if err != nil {
		noResponses[url] = noResponses[url] + 1
		log.Fatal(err)
	}
	elapsed := time.Since(start).Seconds()
	defer result.Body.Close()
	s := fmt.Sprintf("%s %f", url, elapsed)
	log.Println(s)
	if result.StatusCode == http.StatusOK {
		appendResponse(url, elapsed, dataResponses)
	}
}

func appendResponse(url string, time float64, dataResponses map[string][]float64) {
	dataResponses[url] = append(dataResponses[url], time)
}

func httpRequest(urlChan chan string, dataResponses map[string][]float64,
	noResponses map[string]int, generalValue []int) {

	url := <-urlChan
	count := generalValue[0]
	timeout := generalValue[1]
	i := 0
	for i < count {
		go getResponse(url, timeout, dataResponses, noResponses)
		i++
	}
}

func parseArgument(item string) []string {
	return strings.Split(item, ",")
}
func findMinMaxAvg(values []float64) (min float64, max float64, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	min = values[0]
	max = values[0]
	var sum float64 = 0
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum = sum + v
	}
	var count = float64(len(values))
	avg = sum / count
	return min, max, avg
}

func printMinMaxAvg(dataResponses map[string][]float64) {
	for key, _ := range dataResponses {
		min, max, avg := findMinMaxAvg(dataResponses[key])
		fmt.Printf("url: %s, min: %f, max: %f, avg: %f \n", key, min, max, avg)

	}
}

func main() {
	urlChan := make(chan string)
	var dataResponses = map[string][]float64{}
	var noResponses = map[string]int{}
	url := flag.String("url", "", "url.")
	count := flag.Int("count", 1, "count response.")
	timeout := flag.Int("timeout", 1, "count response.")
	flag.Parse()
	var generalValue = []int{*count, *timeout}
	inputValue := parseArgument(*url)

	i := 0
	for i < len(inputValue) {
		go httpRequest(urlChan, dataResponses, noResponses, generalValue)
		urlChan <- inputValue[i]
		i++
	}
	time.Sleep(1 * time.Second)
	fmt.Println(dataResponses)
	printMinMaxAvg(dataResponses)
	fmt.Println(noResponses)
	fmt.Scanln()
}
