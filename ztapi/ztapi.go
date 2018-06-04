// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

//Package ztapi implements a (partial) API client to a ZeroTier service.
package ztapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type apiTime struct {
	time time.Time
}

func (t *apiTime) UnmarshalJSON(b []byte) error {
	l := len(b)
	// if string(b) == "0" {
	// 	t.time = time.Time{}
	// 	return nil
	// }
	if l < 3 {
		n, err := strconv.ParseInt(string(b), 10, 32)
		if err != nil {
			return err
		}
		t.time = time.Unix(0, n)
		return nil
	}
	n, err := strconv.ParseInt(string(b[l-3:]), 10, 32)
	if err != nil {
		return err
	}
	s, err := strconv.ParseInt(string(b[:l-3]), 10, 32)
	if err != nil {
		return err
	}
	t.time = time.Unix(s, n)
	return nil
}

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
