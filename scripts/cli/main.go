package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	localURL = "http://127.0.0.1:4027/api/v1/aft"
)

type httpReqData struct {
	Method  string
	URL     string
	Auth    string
	Body    map[string]string
	Headers map[string]string
}

type response struct {
	StatusCode   int
	Header       http.Header
	Status, Body string
}

func main() {
	var reqNum int
	var smsBody string
	flag.IntVar(&reqNum, "n", 1, "Number of requests")
	flag.StringVar(&smsBody, "m", "", "Message to send")
	flag.Parse()

	recs := getRecs(reqNum)

	if smsBody == "" {
		smsBody = fmt.Sprintf(
			"Hello world to %v at %v",
			reqNum,
			time.Now().String()[8:19],
		)
	}
	apiURL := localURL
	envURL := os.Getenv("APISIM_URL")
	if envURL != "" {
		apiURL = envURL
	}

	log.Printf("Sending to %v recipients.", reqNum)

	body := map[string]string{
		"username": os.Getenv("APISIM_USER"),
		"to":       recs,
		"message":  smsBody,
		"from":     "ApisimCli",
	}

	httpReq := httpReqData{
		URL:    apiURL,
		Body:   body,
		Method: "POST",
		Headers: map[string]string{
			"apikey": os.Getenv("APISIM_SECRET"),
		},
	}
	log.Printf("Request: %v\n", httpReq)
	resp, err := httpReq.makeHTTPRequest()
	if err != nil {
		log.Printf("Error sending: %v\n", err)
	}
	log.Printf("Response: %v\n", resp.Body)
}

func getRecs(num int) string {
	var numbers []string
	for i := 0; i < num; i++ {
		numbers = append(numbers, getPhone())
	}
	return strings.Join(numbers, ",")
}

func getPhone() string {

	prefix := []string{"+2557", "+2536", "+2547", "+2119", "+2568", "+2541"}

	net := []string{
		"16", "17", "18", "20", "21", "22", "23", "25", "96", "27",
	}

	rand.Seed(time.Now().UnixNano())
	destCtry := prefix[rand.Intn(len(prefix))]
	destNet := net[rand.Intn(len(net))]
	randNum := strconv.Itoa(111111 + rand.Intn(999999-111111))

	return destCtry + destNet + randNum
}

func (httpReq *httpReqData) makeHTTPRequest() (response, error) {
	if len(httpReq.Body) < 0 {
		return response{}, errors.New("No form body found")
	}
	form := url.Values{}
	for key, value := range httpReq.Body {
		form.Add(key, value)
	}
	client := http.Client{Timeout: time.Second * 60 * 3}
	req, err := http.NewRequest(
		httpReq.Method, httpReq.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return response{}, fmt.Errorf("makerequest: %v", err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(form)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	if len(httpReq.Headers) > 0 {
		for key, val := range httpReq.Headers {
			req.Header.Add(key, val)
		}
	}
	if len(httpReq.Auth) > 1 {
		user := strings.Split(httpReq.Auth, ":")
		req.SetBasicAuth(user[0], user[1])
	}
	resp, err := client.Do(req)
	if err != nil {
		return response{}, fmt.Errorf("makerequest do: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, fmt.Errorf("makerequest readall: %v", err)
	}
	return response{
		Body: string(body), Header: resp.Header, Status: resp.Status, StatusCode: resp.StatusCode,
	}, nil
}
