package game

import (
	"fmt"
	"testing"
)

func TestNewXPlatformId(t *testing.T) {
	platformCode := STM
	accountId := uint64(12345)

	xpi := NewXPlatformId(platformCode, accountId)

	if xpi.PlatformCode != platformCode {
		t.Errorf("NewXPlatformId() failed: expected platformCode %v, got %v", platformCode, xpi.PlatformCode)
	}

	if xpi.AccountId != accountId {
		t.Errorf("NewXPlatformId() failed: expected accountId %v, got %v", accountId, xpi.AccountId)
	}
}

func TestXPlatformID_Valid(t *testing.T) {
	xpi := &XPlatformID{PlatformCode: STM, AccountId: 12345}

	if !xpi.Valid() {
		t.Errorf("Valid() failed: expected true, got false")
	}

	xpi.PlatformCode = 0

	if xpi.Valid() {
		t.Errorf("Valid() failed: expected false, got true")
	}
}

func TestXPlatformID_Parse(t *testing.T) {
	tests := []struct {
		s           string
		expectedID  *XPlatformID
		expectedErr error
	}{
		{"STM-12345", &XPlatformID{PlatformCode: STM, AccountId: 12345}, nil},
		{"PSN-67890", &XPlatformID{PlatformCode: PSN, AccountId: 67890}, nil},
		{"XBX-54321", &XPlatformID{PlatformCode: XBX, AccountId: 54321}, nil},
		{"OVR_ORG-98765", &XPlatformID{PlatformCode: OVR_ORG, AccountId: 98765}, nil},
		{"OVR-24680", &XPlatformID{PlatformCode: OVR, AccountId: 24680}, nil},
		{"BOT-13579", &XPlatformID{PlatformCode: BOT, AccountId: 13579}, nil},
		{"DMO-86420", &XPlatformID{PlatformCode: DMO, AccountId: 86420}, nil},
		{"TEN-97531", &XPlatformID{PlatformCode: TEN, AccountId: 97531}, nil},
		{"INVALID-12345", nil, fmt.Errorf("invalid format: INVALID-12345")},
		{"STM-INVALID", nil, fmt.Errorf("failed to parse account identifier: strconv.ParseUint: parsing \"INVALID\": invalid syntax")},
	}

	for _, tt := range tests {
		xpi, err := (&XPlatformID{}).Parse(tt.s)

		if err != nil {
			if tt.expectedErr == nil {
				t.Errorf("Parse() failed: unexpected error: %v", err)
			} else if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Parse() failed: expected error %v, got %v", tt.expectedErr, err)
			}
		} else if xpi.PlatformCode != tt.expectedID.PlatformCode || xpi.AccountId != tt.expectedID.AccountId {
			t.Errorf("Parse() failed: expected %v, got %v", tt.expectedID, xpi)
		}
	}
}

func TestXPlatformID_String(t *testing.T) {
	xpi := &XPlatformID{PlatformCode: STM, AccountId: 12345}
	expectedString := "STM-12345"

	if xpi.String() != expectedString {
		t.Errorf("String() failed: expected %s, got %s", expectedString, xpi.String())
	}
}
