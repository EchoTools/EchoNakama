package game

import "testing"

func TestPlatformCode_GetPrefix(t *testing.T) {
	tests := []struct {
		code   PlatformCode
		prefix string
	}{
		{STM, "STM"},
		{PSN, "PSN"},
		{XBX, "XBX"},
		{OVR_ORG, "OVR_ORG"},
		{OVR, "OVR"},
		{BOT, "BOT"},
		{DMO, "DMO"},
		{TEN, "TEN"},
		{PlatformCode(100), "???"}, // Test unknown platform code
	}

	for _, tt := range tests {
		if got := tt.code.GetPrefix(); got != tt.prefix {
			t.Errorf("PlatformCode.GetPrefix() = %v, want %v", got, tt.prefix)
		}
	}
}

func TestPlatformCode_GetDisplayName(t *testing.T) {
	tests := []struct {
		code        PlatformCode
		displayName string
	}{
		{STM, "Steam"},
		{PSN, "Playstation"},
		{XBX, "Xbox"},
		{OVR_ORG, "Oculus VR (ORG)"},
		{OVR, "Oculus VR"},
		{BOT, "Bot"},
		{DMO, "Demo"},
		{TEN, "Tencent"},
		{PlatformCode(100), "Unknown"}, // Test unknown platform code
	}

	for _, tt := range tests {
		if got := tt.code.GetDisplayName(); got != tt.displayName {
			t.Errorf("PlatformCode.GetDisplayName() = %v, want %v", got, tt.displayName)
		}
	}
}

func TestPlatformCode_Parse(t *testing.T) {
	tests := []struct {
		s    string
		code PlatformCode
	}{
		{"STM", STM},
		{"PSN", PSN},
		{"XBX", XBX},
		{"OVR_ORG", OVR_ORG},
		{"OVR", OVR},
		{"BOT", BOT},
		{"DMO", DMO},
		{"TEN", TEN},
		{"UNKNOWN", 0}, // Test unknown platform code
	}

	for _, tt := range tests {
		if got := PlatformCode(0).Parse(tt.s); got != tt.code {
			t.Errorf("PlatformCode.Parse() = %v, want %v", got, tt.code)
		}
	}
}
