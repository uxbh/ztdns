// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

//Package ztapi implements a (partial) API client to a ZeroTier service.
package ztapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getJSON makes a request to the ZeroTier API and returns a JSON object
func getJSON(url, APIToken string, target interface{}) error {
	if APIToken == "" {
		return fmt.Errorf("API Error: No APIToken provided")
	}

	c := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "bearer "+APIToken)
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("API returned error: %s", r.Status)
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
