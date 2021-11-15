package walker

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/walker/providers"
	"log"
	"reflect"
	"strings"
)

func formatLabels(labels map[string]string) string {
	var d []string
	for key, value := range labels {
		d = append(d, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	return strings.Join(d, ",")
}
func PutFlags() (map[string]int, error) {
	var c providers.ConfigProviders
	err := c.Parse("checker.yml")
	if err != nil {
		return nil, err
	}
	promResult := make(map[string]int)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		return nil, dbErr
	}
	for _, team := range teams {
		for _, service := range c.Service {
			if !reflect.ValueOf(service.Put).IsZero() {
				var f database.FlagStruct
				f.Team = team.Login
				f.Service = service.Name
				for _, script := range service.Put {

					metricLabels := map[string]string{
						"team":    team.Login,
						"service": service.Name,
						"script":  script.Name,
						"action":  "put",
					}
					metricNameStr := fmt.Sprintf("checker{%s}", formatLabels(metricLabels))
					promResult[metricNameStr] = 0

					f.Flag = providers.GenerateFlag(20)
					// response, _ := script.RunScript(team.Address, flag)
					f.ID, err = script.RunScript("localhost", f.Flag)
					f.Put()
					f.PutFlag()
					if err != nil || f.ID == "-1" {
						log.Println(err, f.ID)
						break
					}
					promResult[metricNameStr] = 1
				}
			}
		}
	}
	return promResult, nil
}

func CheckFlags() (map[string]int, error) {
	var c providers.ConfigProviders
	err := c.Parse("checker.yml")
	if err != nil {
		return nil, err
	}
	promResult := make(map[string]int)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		return nil, dbErr
	}
	for _, team := range teams {
		for _, service := range c.Service {
			if !reflect.ValueOf(service.Check).IsZero() {
				var f database.FlagStruct
				f.Team = team.Login
				f.Service = service.Name
				keys, _ := f.GetKeys()
				for i, script := range service.Check {
					if len(keys) <= i {
						break
					}
					f.ID = keys[i]

					metricLabels := map[string]string{
						"team":    team.Login,
						"service": service.Name,
						"script":  script.Name,
						"action":  "check",
					}
					metricNameStr := fmt.Sprintf("checker{%s}", formatLabels(metricLabels))
					promResult[metricNameStr] = 0

					// response, _ := script.RunScript(team.Address, flag)
					response, _ := script.RunScript("localhost", f.ID)
					flag, _ := f.GetFlag()
					if response == flag {
						promResult[metricNameStr] = 1
					}
				}
			}
		}
	}
	return promResult, nil
}
func Exploitation() (map[string]int, error) {
	var r providers.RoundsStruct
	err := r.Parse("exploits.yml")
	if err != nil {
		return nil, err
	}
	promResult := make(map[string]int)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		return nil, dbErr
	}
	for _, team := range teams {
		round, _ := database.GetRound()
		if round > len(r.Rounds)-1 {
			round = len(r.Rounds) - 1
		}
		for i := 0; i <= round; i++ {
			round := r.Rounds[i]
			if !reflect.ValueOf(round.Exploits).IsZero() {
				for _, exploit := range round.Exploits {
					metricLabels := map[string]string{
						"team":    team.Login,
						"service": exploit.ServiceName,
						"script":  exploit.ScriptName,
						"action":  "exploit",
					}
					metricNameStr := fmt.Sprintf("exploit{%s}", formatLabels(metricLabels))
					promResult[metricNameStr] = 0

					// response, _ := script.RunScript(team.Address, "")
					response, _ := exploit.RunScript("localhost", "")
					if response == "1"{
						promResult[metricNameStr] = 1
						database.AddDefenceFlag(team.Login, exploit.ServiceName)
					}
				}
			}
		}
	}
	database.IncrRound()
	return promResult, nil
}
