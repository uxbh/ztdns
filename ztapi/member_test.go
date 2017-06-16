// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package ztapi

import (
	"testing"
)

func BenchmarkGet6Plane(b *testing.B) {
	member := Member{
		NetworkID: "8056C2E21C24673D",
		NodeID:    "22c55a1da6",
	}
	for n := 0; n < b.N; n++ {
		member.Get6Plane()
	}
}

func BenchmarkGetRFC4193(b *testing.B) {
	member := Member{
		NetworkID: "8056C2E21C24673D",
		NodeID:    "22c55a1da6",
	}
	for n := 0; n < b.N; n++ {
		member.GetRFC4193()
	}
}

func TestGetRFC4193(t *testing.T) {
	var table = []struct {
		net  string
		node string
		out  string
	}{
		{"1", "1", "fd00::199:9300:0:1"},
		{"FFFFFFFFFFFFFFFF", "ffffffffff", "fdff:ffff:ffff:ffff:ff99:93ff:ffff:ffff"},
		{"8056C2E21C24673D", "22c55a1da6", "fd80:56c2:e21c:2467:3d99:9322:c55a:1da6"},
	}
	for _, tt := range table {
		member := Member{
			NetworkID: tt.net,
			NodeID:    tt.node,
		}
		s := member.GetRFC4193().String()
		if s != tt.out {
			t.Errorf("GetRFC4193\nin(%x,%x)\n got %v\nwant %v", tt.net, tt.node, s, tt.out)
		}
	}
}
func TestGet6Plane(t *testing.T) {
	var table = []struct {
		net  string
		node string
		out  string
	}{
		{"1", "1", "fc00:0:100:0:1::1"},
		{"FFFFFFFFFFFFFFFF", "ffffffffff", "fc00:0:ff:ffff:ffff::1"},
		{"8056C2E21C24673D", "22c55a1da6", "fc9c:72a5:df22:c55a:1da6::1"},
	}
	for _, tt := range table {
		member := Member{
			NetworkID: tt.net,
			NodeID:    tt.node,
		}
		s := member.Get6Plane().String()
		if s != tt.out {
			t.Errorf("Get6Plane\nin(%x,%x)\n got %v\nwant %v", tt.net, tt.node, s, tt.out)
		}
	}
}
