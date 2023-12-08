// Package main provides functionality for working with platform identifiers in the game.
package game

import (
	"fmt"
	"strconv"
	"strings"
)

// XPlatformID represents an identifier for a user on the platform.
type XPlatformID struct {
	PlatformCode PlatformCode `json:"platform_code"`
	AccountId    uint64       `json:"account_id"`
}

// Constants
//const xPlatformIdSize = 16

// NewXPlatformId initializes a new XPlatformId.
func NewXPlatformId(platformCode PlatformCode, accountId uint64) *XPlatformID {
	// Implementation for initializing a new XPlatformId goes here
	// ...

	return &XPlatformID{
		PlatformCode: platformCode,
		AccountId:    accountId,
	}
}

func (xpi *XPlatformID) Valid() bool {
	return xpi.PlatformCode != 0
}

// Parse parses a string into a given platform identifier.
func (xpi *XPlatformID) Parse(s string) (*XPlatformID, error) {
	// Obtain the position of the last dash.
	dashIndex := strings.LastIndex(s, "-")
	if dashIndex < 0 {
		return nil, fmt.Errorf("invalid format: %s", s)
	}

	// Split it there
	platformCodeStr := s[:dashIndex]
	accountIdStr := s[dashIndex+1:]

	// Determine the platform code.
	platformCode := PlatformCode(0).Parse(platformCodeStr)

	// Try to parse the account identifier
	accountId, err := strconv.ParseUint(accountIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account identifier: %v", err)
	}

	// Create the identifier
	platformId := &XPlatformID{PlatformCode: platformCode, AccountId: accountId}
	return platformId, nil
}

func (xpi *XPlatformID) String() string {
	return fmt.Sprintf("%s-%d", xpi.PlatformCode.String(), xpi.AccountId)
}
