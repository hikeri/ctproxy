package src

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidatorAdditionalInfo struct {
	// maxSpeed int // in KB @todo
}

type ValidatorProxyFunc func(login string, password string) (bool, *ValidatorAdditionalInfo)

var ValidatorsAvailable = map[string]ValidatorProxyFunc{
	"test":    testValidator,
	"ct_auth": censorTrackerValidator,
}

func EmptyValidator(login string, password string) (bool, *ValidatorAdditionalInfo) {
	return true, &ValidatorAdditionalInfo{}
}

func testValidator(login string, password string) (bool, *ValidatorAdditionalInfo) {
	return login == "test" && password == "", &ValidatorAdditionalInfo{}
}

const censorTrackerEndpoint = "https://dropmigrations.censortracker.org/api/authenticate/"

type censorTrackerResponse struct {
	Email         string `json:"email"`
	Active        bool   `json:"active"`
	RemainingDays int    `json:"remainingDays"`
}

func censorTrackerValidator(login string, password string) (bool, *ValidatorAdditionalInfo) {
	req, _ := http.NewRequest("GET", censorTrackerEndpoint, nil)
	req.Header.Set("Authorization", password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	var data censorTrackerResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	if data.Active != true {
		return false, nil
	}

	return true, &ValidatorAdditionalInfo{}
}
