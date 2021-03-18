package job_handlers

import (
	"apisim/app/jobs"
	"apisim/app/jobs/sms_jobs"
	"apisim/app/work"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/revel/revel"
)

type httpXMLData struct {
	Body string
	URL  string
}

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

type ProcessDlrJobHandler struct {
	jobEnqueuer work.JobEnqueuer
}

func NewProcessDlrJobHandler(
	jobEnqueuer work.JobEnqueuer,
) *ProcessDlrJobHandler {
	return &ProcessDlrJobHandler{
		jobEnqueuer: jobEnqueuer,
	}
}

func (h *ProcessDlrJobHandler) Job() jobs.Job {
	return &sms_jobs.ProcessDlrJob{}
}

func (h *ProcessDlrJobHandler) PerformJob(
	ctx context.Context,
	body string,
) error {
	var theJob sms_jobs.ProcessDlrJob
	err := json.Unmarshal([]byte(body), &theJob)
	if err != nil {
		revel.AppLog.Errorf("error unmarshal send dlr task: %v", err)
		return nil
	}

	req := theJob.Request
	if req.Source != "sf" {
		var body map[string]string
		if req.Source == "at" {
			body = map[string]string{"status": req.Status, "id": req.ID, "failureReason": req.Reason}
		}
		if req.Source == "rm" {
			body = map[string]string{"sMessageId": req.ID, "sStatus": req.Status}
		}

		httpReq := httpReqData{
			URL:    req.URL,
			Body:   body,
			Method: "POST",
		}
		revel.AppLog.Debugf("at|sf dlr request: %v", httpReq)
		resp, err := httpReq.makeHTTPRequest()
		if err != nil {
			return err
		}
		revel.AppLog.Debugf("at|rm send dlr response: %v", resp.Body)
	} else {
		safDlrQuery := `<?xml version="1.0" encoding="UTF-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://www.csapi.org/schema/parlayx/sms/notification/v2_2/local" xmlns:v2="http://www.huawei.com.cn/schema/common/v2_1"><soapenv:Header><v2:NotifySOAPHeader><v2:timeStamp>%s</v2:timeStamp><v2:subReqID>%s</v2:subReqID><traceUniqueID>%s</traceUniqueID></v2:NotifySOAPHeader></soapenv:Header><soapenv:Body><loc:notifySmsDeliveryReceipt><loc:correlator>%s</loc:correlator><loc:deliveryStatus><address>tel:%s</address><deliveryStatus>%s</deliveryStatus></loc:deliveryStatus></loc:notifySmsDeliveryReceipt></soapenv:Body></soapenv:Envelope>`
		timeStamp := time.Now().Format("20060102150405")
		reqCont := fmt.Sprintf(
			safDlrQuery, timeStamp, "11111111111111", "504021503311009040428550001002",
			req.ID, req.Phone, req.Status,
		)

		httpReq := httpXMLData{
			Body: reqCont,
			URL:  req.URL,
		}
		revel.AppLog.Debugf("saf dlr request: %v", httpReq)
		resp, err := httpReq.MakeXMLRequest()
		if err != nil {
			return fmt.Errorf("send dlr sf: %v", err)
		}
		revel.AppLog.Debugf("sf send dlr sf resp: %v", resp.Body)
	}

	return nil
}

func (httpReq *httpReqData) makeHTTPRequest() (response, error) {
	if len(httpReq.Body) < 0 {
		return response{}, fmt.Errorf("No form body found")
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

func (httpReq *httpXMLData) MakeXMLRequest() (response, error) {
	httpClient := http.Client{Timeout: time.Second * 60 * 2}
	resp, err := httpClient.Post(
		httpReq.URL, "text/xml; charset=utf-8",
		bytes.NewBufferString(httpReq.Body),
	)
	if err != nil {
		return response{}, fmt.Errorf("xml post error: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, fmt.Errorf("xml response readall: %v", err)
	}
	return response{
		Body: string(body), Header: resp.Header, Status: resp.Status, StatusCode: resp.StatusCode,
	}, nil
}
