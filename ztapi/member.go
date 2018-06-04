// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package ztapi

import (
	"fmt"
	"net"
	"strconv"
)

// GetMemberInfo returns a Member containing informationa about a specific member in a ZeroTier network
func GetMemberInfo(API, host, networkID, memberID string) (*Member, error) {
	resp := new(Member)
	url := fmt.Sprintf("%s/network/%s/member/%s", host, networkID, memberID)
	err := getJSON(url, API, resp)
	if err != nil {
		return nil, fmt.Errorf("Unable to get member info: %s", err.Error())
	}
	return resp, nil
}

// GetMemberList gets a Slice of Members in a ZeroTier network
func GetMemberList(API, host, networkID string) (*Members, error) {
	resp := new(Members)
	url := fmt.Sprintf("%s/network/%s/member", host, networkID)
	err := getJSON(url, API, resp)
	if err != nil {
		return nil, fmt.Errorf("Unable to get member list: %s", err.Error())
	}
	return resp, nil
}

// Members is a List of Members
type Members []Member

// Member contains data from a Member request
type Member struct {
	ID           string
	Type         string
	Clock        apiTime
	NetworkID    string
	NodeID       string
	ControllerID string
	Hidden       bool
	Name         string
	Online       bool
	Description  string
	Config       struct {
		ActiveBridge bool
		Address      string
		AuthHistory  []struct {
			A  bool
			By string
			C  string
			Ct string
			Ts int
		}
		Authorized           bool
		Capabilities         []string
		CreationTime         apiTime
		ID                   string
		Identity             string
		IPAssignments        []string
		LastAuthorizedTime   apiTime
		LastDeauthorizedTime apiTime
		NoAutoAssignIPs      bool
		Nwid                 string
		Objtype              string
		PhysicalAddr         string
		Revision             int
		Tags                 []string
		VMajor               int
		VMinor               int
		VProto               int
		VRev                 int
	}
	LastOnline             apiTime
	LastOffline            apiTime
	PhysicalAddress        string
	PhysicalLocation       []float64
	ClientVersion          string
	OfflineNotifyDelay     int
	ProtocolVersion        int
	SupportsCircuitTesting bool
	SupportsRulesEngine    bool
}

// Get6Plane returns the 6Plane address for a given network member.
// See https://support.zerotier.com/hc/en-us/articles/115001080308-ZeroTier-6PLANE-IPv6-Addressing for details.
// If the ZeroTier network assigns 6Plane addresses, this will be the device's address range.
func (m *Member) Get6Plane() net.IP {
	n, _ := strconv.ParseUint(m.NetworkID, 16, 64)
	d, _ := strconv.ParseUint(m.NodeID, 16, 64)
	s := n&(0xFFFFFFFF<<0x20)>>0x20 ^ n&0xFFFFFFFF
	ip := net.IP{
		0xfc, byte(s >> 0x18),
		byte((s >> 0x10) - ((s >> 0x18) << 0x8)),
		byte((s >> 0x8) - ((s >> 0x10) << 0x8)),
		byte((s) - ((s >> 0x8) << 0x8)),
		byte(d >> 0x20),
		byte((d >> 0x18) - ((d >> 0x20) << 0x8)),
		byte((d >> 0x10) - ((d >> 0x18) << 0x8)),
		byte((d >> 0x08) - ((d >> 0x10) << 0x8)),
		byte((d >> 0x00) - ((d >> 0x08) << 0x8)),
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	return ip
}

// GetRFC4193 returns the RFC4193 address of a given network member.
// If the ZeroTier network assigns RFC4193 addresses, this will be the device's address.
func (m *Member) GetRFC4193() net.IP {
	n, _ := strconv.ParseUint(m.NetworkID, 16, 64)
	d, _ := strconv.ParseUint(m.NodeID, 16, 64)
	ip := net.IP{
		0xfd, byte(n >> 0x38),
		byte((n >> 0x30) - ((n >> 0x38) << 0x8)),
		byte((n >> 0x28) - ((n >> 0x30) << 0x8)),
		byte((n >> 0x20) - ((n >> 0x28) << 0x8)),
		byte((n >> 0x18) - ((n >> 0x20) << 0x8)),
		byte((n >> 0x10) - ((n >> 0x18) << 0x8)),
		byte((n >> 0x08) - ((n >> 0x10) << 0x8)),
		byte((n >> 0x00) - ((n >> 0x08) << 0x8)),
		0x99, 0x93, byte(d >> 0x20),
		byte((d >> 0x18) - ((d >> 0x20) << 0x8)),
		byte((d >> 0x10) - ((d >> 0x18) << 0x8)),
		byte((d >> 0x08) - ((d >> 0x10) << 0x8)),
		byte((d >> 0x00) - ((d >> 0x08) << 0x8)),
	}
	return ip
}
