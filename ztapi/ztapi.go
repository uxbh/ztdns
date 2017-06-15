// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package ztapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
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
	defer r.Body.Close()

	if false {
		s, _ := httputil.DumpRequest(req, true)
		fmt.Printf("Request:\n%s\n", s)
		s, _ = httputil.DumpResponse(r, true)
		fmt.Printf("Response:\n%s\n", s)
	}

	return json.NewDecoder(r.Body).Decode(target)
}
