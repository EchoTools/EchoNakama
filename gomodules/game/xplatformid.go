package game

import (
	"fmt"
	"strconv"
	"strings"
)

func NewEchoUserId(platformCode PlatformCode, accountId uint64) *EchoUserId {
	return &EchoUserId{PlatformCode: platformCode, AccountId: accountId}
}

// EchoUserId represents an identifier for a user on the platform.
type EchoUserId struct {
	PlatformCode PlatformCode `json:"platform_code"`
	AccountId    uint64       `json:"account_id"`
}

// Constants
//const xPlatformIdSize = 16

func (xpi *EchoUserId) Valid() bool {
	return xpi.PlatformCode > STM && xpi.PlatformCode < TEN && xpi.AccountId > 0
}

// Parse parses a string into a given platform identifier.
func (xpi *EchoUserId) Parse(s string) (*EchoUserId, error) {
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
	platformId := &EchoUserId{PlatformCode: platformCode, AccountId: accountId}
	return platformId, nil
}

func (xpi *EchoUserId) String() string {
	return fmt.Sprintf("%s-%d", xpi.PlatformCode.String(), xpi.AccountId)
}

func (xpi *EchoUserId) Token() string {
	return fmt.Sprintf("%s-%d", xpi.PlatformCode.String(), xpi.AccountId)
}

func (xpi *EchoUserId) Equal(other *EchoUserId) bool {
	return xpi.PlatformCode == other.PlatformCode && xpi.AccountId == other.AccountId
}

func (xpi *EchoUserId) IsEmpty() bool {
	return xpi.PlatformCode == 0 && xpi.AccountId == 0
}

func (xpi *EchoUserId) IsNotEmpty() bool {
	return xpi.PlatformCode != 0 && xpi.AccountId != 0
}
