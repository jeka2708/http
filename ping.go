package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags
var dataResponses = map[string][]float64{}
var noResponses = map[string]int{}

func getResponse(url string, timeout int) {
	mils := float64(timeout) / float64(1000)
	start := time.Now()
	result, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if time.Since(start).Seconds() > mils {
		noResponses[url] = noResponses[url] + 1
		return
	}
	elapsed := time.Since(start).Seconds()
	defer result.Body.Close()
	s := fmt.Sprintf("%s %f", url, elapsed)
	log.Println(s)
	if result.StatusCode == http.StatusOK {
		appendResponse(url, elapsed)
	}
}

func appendResponse(url string, time float64) {
	dataResponses[url] = append(dataResponses[url], time)
}

func httpRequest(c chan string) {
	args := parseArgument(<-c)
	url := args[0]
	count, _ := strconv.Atoi(args[1])
	timout, _ := strconv.Atoi(args[2])
	i := 0
	for i < count {
		go getResponse(url, timout)
		i++
	}
}

func parseArgument(item string) []string {
	return strings.Split(item, ",")
}
func Min(values []float64) (min float64, e error) {
	if len(values) == 0 {
		return 0, errors.New("Cannot detect a minimum value in an empty slice")
	}

	min = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min, nil
}
func Max(values []float64) (max float64, e error) {
	if len(values) == 0 {
		return 0, errors.New("Cannot detect a maximum value in an empty slice")
	}

	max = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max, nil
}
func Avg(values []float64) (avg float64, e error) {
	if len(values) == 0 {
		return 0, errors.New("Cannot detect a average value in an empty slice")
	}

	var sum float64 = 0
	for _, v := range values {
		sum = sum + v
	}
	var count = float64(len(values))
	avg = sum / count
	return avg, nil
}
func printMin() {
	for key, _ := range dataResponses {
		min, _ := Min(dataResponses[key])
		fmt.Println(key, min, " min")

	}
}
func printMax() {
	for key, _ := range dataResponses {
		min, _ := Max(dataResponses[key])
		fmt.Println(key, min, " max")

	}
}
func printAvg() {
	for key, _ := range dataResponses {
		min, _ := Avg(dataResponses[key])
		fmt.Println(key, min, " avg")

	}
}
func main() {
	c := make(chan string)
	flag.Var(&myFlags, "list1", "List of arguments.")
	flag.Parse()
	i := 0
	for i < len(myFlags) {
		go httpRequest(c)
		c <- myFlags[i]
		i++
	}
	defer fmt.Println(dataResponses)
	defer fmt.Println(noResponses)
	defer printMin()
	defer printMax()
	defer printAvg()

	fmt.Scanln()

}
