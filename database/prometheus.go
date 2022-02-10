package database

import (
	"encoding/json"
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PrometheusQueryRange struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Proto    string `json:"proto"`
				Route    string `json:"route"`
				Service  string `json:"service"`
				Team     string `json:"team"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
			Value  []interface{}   `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func getPrometheusState(query string, time string) (PrometheusQueryRange, error) {
	promAddr := os.Getenv("PROMETHEUS")
	if promAddr == "" {
		promAddr = "http://localhost:9090"
	}
	urlAddr := promAddr + "/api/v1/query"
	req, err := http.NewRequest("GET", urlAddr, nil)
	if err != nil {
		log.Print(err)
		return PrometheusQueryRange{}, err
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("time", time)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	client := &http.Client{}
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Println(reqErr)
		return PrometheusQueryRange{}, reqErr
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return PrometheusQueryRange{}, err
	}
	var queryRanges PrometheusQueryRange
	jsonErr := json.Unmarshal(body, &queryRanges)
	if jsonErr != nil {
		log.Println(jsonErr)
		return PrometheusQueryRange{}, err
	}
	return queryRanges, nil
}

func getPrometheusRange(query string, start string, end string) (PrometheusQueryRange, error) {
	promAddr := os.Getenv("PROMETHEUS")
	if promAddr == "" {
		promAddr = "http://localhost:9090"
	}
	urlAddr := promAddr + "/api/v1/query_range"
	req, err := http.NewRequest("GET", urlAddr, nil)
	if err != nil {
		log.Print(err)
		return PrometheusQueryRange{}, err
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	q.Add("end", end)
	q.Add("step", config.Conf.RoundInterval)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Println("here")
		return PrometheusQueryRange{}, reqErr
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return PrometheusQueryRange{}, err
	}

	var queryRanges PrometheusQueryRange
	jsonErr := json.Unmarshal(body, &queryRanges)
	if jsonErr != nil {
		log.Println(jsonErr)
		return PrometheusQueryRange{}, jsonErr
	}
	return queryRanges, nil
}

type TeamHistory struct {
	RoundsHistory map[string]map[int]float64
	Sources       float64
	TotalRounds   float64
}

func GetTeamHistory(teamName string) (TeamHistory, error) {
	query := fmt.Sprintf("checker{team=\"%s\"}", teamName)
	startTime, _ := GetStartTimeStamp()
	shiftTime, _ := time.Parse(time.RFC3339, startTime)
	log.Println(shiftTime)
	shiftTime = shiftTime.Add(time.Second)
	queryRanges, err := getPrometheusRange(query, shiftTime.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	if err != nil {
		return TeamHistory{}, err
	}

	teamHistory := TeamHistory{RoundsHistory: map[string]map[int]float64{}}
	teamHistory.Sources = float64(len(queryRanges.Data.Result))
	if teamHistory.Sources < 1 {
		log.Println("No Data")
		return TeamHistory{}, nil
	}
	teamHistory.TotalRounds = float64(len(queryRanges.Data.Result[0].Values))
	log.Println(queryRanges.Data.Result[0].Metric.Team, queryRanges.Data.Result[0].Metric.Service)
	for _, result := range queryRanges.Data.Result {
		for i, v := range result.Values {
			value, _ := strconv.Atoi(v[1].(string))
			if teamHistory.RoundsHistory[result.Metric.Service] != nil {
				teamHistory.RoundsHistory[result.Metric.Service][i] += float64(value)
			} else {
				teamHistory.RoundsHistory[result.Metric.Service] = map[int]float64{i: float64(value)}
			}
		}
		if float64(len(result.Values)) > teamHistory.TotalRounds {
			teamHistory.TotalRounds = float64(len(result.Values))
		}
	}
	teamHistory.TotalRounds -= 1
	return teamHistory, nil
}

func GetTeamStatus(teamName string) (map[string]float64, float64, error) {
	query := fmt.Sprintf("checker{team=\"%s\"}", teamName)
	//lastTime, _ := GetLastTimeStamp()
	queryRanges, err := getPrometheusState(query, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(queryRanges)
	serviceStatus := map[string]float64{}
	sources := float64(len(queryRanges.Data.Result))
	if sources < 1 {
		log.Println("No Data")
		return nil, 0, nil
	}
	log.Println(queryRanges.Data.Result[0].Metric.Team, queryRanges.Data.Result[0].Metric.Service)
	for _, result := range queryRanges.Data.Result {
		value, _ := strconv.Atoi(result.Value[1].(string))
		serviceStatus[result.Metric.Service] += float64(value)
	}

	return serviceStatus, sources, nil
}
