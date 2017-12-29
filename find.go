package warpcache

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// ErrNotSingleGTS is the error thrown when call on FIND doesn't return [1]
var ErrNotSingleGTS = errors.New("FIND on selector doesn't not return a single GTS")

func checkSingleGTS(config Configuration, selector string) error {

	body, err := generateFindSizeWarpScript(config.ReadToken, selector)
	if err != nil {
		return err
	}

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Post(config.HTTPProtocol+"://"+config.Endpoint+"/api/v0/exec", "", strings.NewReader(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		dump, _ := httputil.DumpResponse(resp, true)
		return errors.New("Error during WarpScript execution: " + string(dump))
	}

	findResp := make([]float64, 1)

	err = json.NewDecoder(resp.Body).Decode(&findResp)
	if err != nil {

	}

	size := findResp[0]
	if size != 1 {
		return ErrNotSingleGTS
	}

	return nil
}
