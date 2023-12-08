package game

import "strings"

// PlatformCode represents the platforms on which a client may be operating.
type PlatformCode int

const (
	// Steam
	STM PlatformCode = iota + 1

	// Playstation
	PSN

	// Xbox
	XBX

	// Oculus VR user
	OVR_ORG

	// Oculus VR
	OVR

	// Bot/AI
	BOT

	// Demo (no ovr)
	DMO

	// Tencent
	TEN
)

// GetPrefix obtains a platform prefix string for a given PlatformCode.
func (code PlatformCode) GetPrefix() string {
	// Try to obtain a name for this platform code.
	name := code.String()

	// If we could obtain one, the prefix should just be the same as the name, but with underscores represented as dashes.
	if name != "" {
		return name
	}

	// An unknown/invalid platform is denoted with the value returned below.
	return "???"
}

// GetDisplayName obtains a display name for a given PlatformCode.
func (code PlatformCode) GetDisplayName() string {
	// Switch on the provided platform code and return a display name.
	switch code {
	case STM:
		return "Steam"
	case PSN:
		return "Playstation"
	case XBX:
		return "Xbox"
	case OVR_ORG:
		return "Oculus VR (ORG)"
	case OVR:
		return "Oculus VR"
	case BOT:
		return "Bot"
	case DMO:
		return "Demo"
	case TEN:
		return "Tencent" // TODO: Verify, this is only suspected to be the target of "TEN".
	default:
		return "Unknown"
	}
}

// Parse parses a string generated from PlatformCode's String() method back into a PlatformCode.
func (code PlatformCode) Parse(s string) PlatformCode {
	// Convert any underscores in the string to dashes.
	s = strings.ReplaceAll(s, "-", "_")

	// Get the enum option to represent this.
	if result, ok := platformCodeFromString(s); ok {
		return result
	}
	return 0
}

// platformCodeToString converts a PlatformCode to its string representation.
func (code PlatformCode) String() string {
	switch code {
	case STM:
		return "STM"
	case PSN:
		return "PSN"
	case XBX:
		return "XBX"
	case OVR_ORG:
		return "OVR_ORG"
	case OVR:
		return "OVR"
	case BOT:
		return "BOT"
	case DMO:
		return "DMO"
	case TEN:
		return "TEN"
	default:
		return ""
	}
}

// platformCodeFromString converts a string to its PlatformCode representation.
func platformCodeFromString(s string) (PlatformCode, bool) {
	switch s {
	case "STM":
		return STM, true
	case "PSN":
		return PSN, true
	case "XBX":
		return XBX, true
	case "OVR_ORG":
		return OVR_ORG, true
	case "OVR":
		return OVR, true
	case "BOT":
		return BOT, true
	case "DMO":
		return DMO, true
	case "TEN":
		return TEN, true
	default:
		return 0, false
	}
}
