package game

import (
	"encoding/json"
	"fmt"
)

// Profiles represents the 'profile' field in the JSON data
type GameProfiles struct {
	Client ClientProfile `json:"client"`
	Server ServerProfile `json:"server"`
}

// MergePlayerData merges a partial PlayerData with a template PlayerData.
func (base *ClientProfile) Merge(partial *ClientProfile) (*ClientProfile, error) {

	partialJSON, err := json.Marshal(partial)
	if err != nil {
		return nil, fmt.Errorf("error marshaling partial data: %v", err)
	}

	err = json.Unmarshal(partialJSON, &base)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling partial data: %v", err)
	}

	return base, nil
}

type ClientProfile struct {
	DisplayName string      `json:"displayname"`
	XPlatformID string      `json:"xplatformid"`
	TeamName    interface{} `json:"teamname,omitempty"`
	Weapon      string      `json:"weapon"`
	Grenade     string      `json:"grenade"`
	WeaponArm   int         `json:"weaponarm"`
	ModifyTime  int64       `json:"modifytime"`
	Ability     string      `json:"ability"`
	Legal       Legal       `json:"legal"`

	Mute          interface{}   `json:"mute"`
	NPE           NPE           `json:"npe"`
	Customization Customization `json:"customization"`
	Social        Social        `json:"social"`
	NewUnlocks    []interface{} `json:"newunlocks"`
}

type Legal struct {
	PointsPolicyVersion int `json:"points_policy_version"`
	EulaVersion         int `json:"eula_version"`
	GameAdminVersion    int `json:"game_admin_version"`
	SplashScreenVersion int `json:"splash_screen_version"`
}

// NPE represents the 'npe' field in the 'client' section
type NPE struct {
	Lobby             NPEItem   `json:"lobby"`
	FirstMatch        NPEItem   `json:"firstmatch"`
	Movement          NPEItem   `json:"movement"`
	ArenaBasics       NPEItem   `json:"arenabasics"`
	SocialTabSeen     Versioned `json:"social_tab_seen"`
	Pointer           Versioned `json:"pointer"`
	BlueTintTabSeen   Versioned `json:"blue_tint_tab_seen"`
	HeraldryTabSeen   Versioned `json:"heraldry_tab_seen"`
	OrangeTintTabSeen Versioned `json:"orange_tint_tab_seen"`
}

// NPEItem represents an item in the 'npe' section
type NPEItem struct {
	Completed bool `json:"completed"`
}
type Versioned struct {
	Version int `json:"version"`
}

// Customization represents the 'customization' field in the 'client' section
type Customization struct {
	BattlePassSeasonPOIVersion int `json:"battlepass_season_poi_version"`
	NewUnlocksPOIVersion       int `json:"new_unlocks_poi_version"`
	StoreEntryPOIVersion       int `json:"store_entry_poi_version"`
	ClearNewUnlocksVersion     int `json:"clear_new_unlocks_version"`
}

// Social represents the 'social' field in the 'client' section
type Social struct {
	CommunityValuesVersion int    `json:"community_values_version"`
	SetupVersion           int    `json:"setup_version"`
	Group                  string `json:"group"`
}

// ServerProfile represents the 'server' field in the 'profile' section
type ServerProfile struct {
	DisplayName     string         `json:"displayname"`
	XPlatformID     string         `json:"xplatformid"`
	Version         int            `json:"_version"`
	PublisherLock   string         `json:"publisher_lock"`
	PurchasedCombat int            `json:"purchasedcombat"`
	LobbyVersion    int            `json:"lobbyversion"`
	ModifyTime      int64          `json:"modifytime"`
	LoginTime       int64          `json:"logintime"`
	UpdateTime      int64          `json:"updatetime"`
	CreateTime      int64          `json:"createtime"`
	MaybeStale      interface{}    `json:"maybestale"`
	Stats           StatSection    `json:"stats"`
	Unlocks         UnlockSection  `json:"unlocks"`
	Loadout         LoadoutSection `json:"loadout"`
	Social          interface{}    `json:"social"`
	Achievements    interface{}    `json:"achievements"`
	RewardState     interface{}    `json:"reward_state"`
	Dev             DevSection     `json:"dev"`
}

type DevSection struct {
	DisableAfkTimeout bool   `json:"disable_afk_timeout"`
	XPlatformIdStr    string `json:"xplatformid"`
}

type LevelStatItem struct {
	Cnt int    `json:"cnt"`
	Op  string `json:"op"`
	Val int    `json:"val"`
}

type StatItem struct {
	Op  string  `json:"op"`
	Val float64 `json:"val"`
}

type ArenaStats struct {
	Level                        LevelStatItem `json:"Level"`
	Goals                        StatItem      `json:"Goals"`
	TopSpeedsTotal               StatItem      `json:"TopSpeedsTotal"`
	HighestArenaWinStreak        StatItem      `json:"HighestArenaWinStreak"`
	ArenaWinPercentage           StatItem      `json:"ArenaWinPercentage"`
	ArenaWins                    StatItem      `json:"ArenaWins"`
	GoalsPerGame                 StatItem      `json:"GoalsPerGame"`
	Points                       StatItem      `json:"Points"`
	Interceptions                StatItem      `json:"Interceptions"`
	ThreePointGoals              StatItem      `json:"ThreePointGoals"`
	Clears                       StatItem      `json:"Clears"`
	BounceGoals                  StatItem      `json:"BounceGoals"`
	PossessionTime               StatItem      `json:"PossessionTime"`
	HatTricks                    StatItem      `json:"HatTricks"`
	ShotsOnGoal                  StatItem      `json:"ShotsOnGoal"`
	HighestPoints                StatItem      `json:"HighestPoints"`
	GoalScorePercentage          StatItem      `json:"GoalScorePercentage"`
	AveragePossessionTimePerGame StatItem      `json:"AveragePossessionTimePerGame"`
	AverageTopSpeedPerGame       StatItem      `json:"AverageTopSpeedPerGame"`
	AveragePointsPerGame         StatItem      `json:"AveragePointsPerGame"`
	ArenaMVPPercentage           StatItem      `json:"ArenaMVPPercentage"`
	ArenaMVPs                    StatItem      `json:"ArenaMVPs"`
	CurrentArenaWinStreak        StatItem      `json:"CurrentArenaWinStreak"`
	CurrentArenaMVPStreak        StatItem      `json:"CurrentArenaMVPStreak"`
	HighestArenaMVPStreak        StatItem      `json:"HighestArenaMVPStreak"`
	XP                           StatItem      `json:"XP"`
	ShotsOnGoalAgainst           StatItem      `json:"ShotsOnGoalAgainst"`
	ArenaLosses                  StatItem      `json:"ArenaLosses"`
	Catches                      StatItem      `json:"Catches"`
	StunsPerGame                 StatItem      `json:"StunsPerGame"`
	HighestStuns                 StatItem      `json:"HighestStuns"`
	Steals                       StatItem      `json:"Steals"`
	Stuns                        StatItem      `json:"Stuns"`
	PunchesReceived              StatItem      `json:"PunchesReceived"`
	Passes                       StatItem      `json:"Passes"`
	Blocks                       StatItem      `json:"Blocks"`
	BlockPercentage              StatItem      `json:"BlockPercentage"`
}

type CombatStats struct {
	Level LevelStatItem `json:"Level"`
}

type StatSection struct {
	Arena  ArenaStats  `json:"arena"`
	Combat CombatStats `json:"combat"`
}

// UnlockSection represents a section in the 'unlocks' field in the 'server' section
type UnlockSection struct {
	Arena  ArenaUnlocks  `json:"arena"`
	Combat CombatUnlocks `json:"combat"`
}
type ArenaUnlocks struct {
	DecalCombatFlamingoA   bool `json:"decal_combat_flamingo_a"`
	DecalCombatLogoA       bool `json:"decal_combat_logo_a"`
	DecalDefault           bool `json:"decal_default"`
	DecalSheldonA          bool `json:"decal_sheldon_a"`
	EmoteBlinkSmileyA      bool `json:"emote_blink_smiley_a"`
	EmoteDefault           bool `json:"emote_default"`
	EmoteDizzyEyesA        bool `json:"emote_dizzy_eyes_a"`
	LoadoutNumber          bool `json:"loadout_number"`
	PatternDefault         bool `json:"pattern_default"`
	PatternLightningA      bool `json:"pattern_lightning_a"`
	RwdBannerS1Default     bool `json:"Rwd_banner_s1_default"`
	RwdBoosterDefault      bool `json:"Rwd_booster_default"`
	RwdBracerDefault       bool `json:"Rwd_bracer_default"`
	RwdChassisBodyS11A     bool `json:"Rwd_chassis_body_s11_a"`
	RwdDecalbackDefault    bool `json:"Rwd_decalback_default"`
	RwdDecalborderDefault  bool `json:"Rwd_decalborder_default"`
	RwdMedalDefault        bool `json:"Rwd_medal_default"`
	RwdTagDefault          bool `json:"Rwd_tag_default"`
	RwdTagS1ASecondary     bool `json:"Rwd_tag_s1_a_secondary"`
	RwdTitleTitleDefault   bool `json:"Rwd_title_title_default"`
	TintBlueADefault       bool `json:"tint_blue_a_default"`
	TintNeutralADefault    bool `json:"tint_neutral_a_default"`
	TintNeutralAS10Default bool `json:"tint_neutral_a_s10_default"`
	TintOrangeADefault     bool `json:"tint_orange_a_default"`
	RwdGoalFxDefault       bool `json:"Rwd_goal_fx_default"`
	EmissiveDefault        bool `json:"emissive_default"`
}

type CombatUnlocks struct {
	RwdBoosterS10      bool `json:"Rwd_booster_s10"`
	RwdChassisBodyS10A bool `json:"Rwd_chassis_body_s10_a"`
}

// LoadoutSection represents the 'loadout' field in the 'server' section
type LoadoutSection struct {
	Instances LoadoutInstance `json:"instances"`
	Number    int             `json:"number"`
}

// LoadoutInstance represents the 'instances' field in the 'loadout' section
type LoadoutInstance struct {
	Unified LoadoutSlots `json:"unified"`
}

// LoadoutSlots represents the 'slots' field in the 'unified' section
type LoadoutSlots struct {
	Decal          string `json:"decal"`
	DecalBody      string `json:"decal_body"`
	Emote          string `json:"emote"`
	SecondEmote    string `json:"secondemote"`
	Tint           string `json:"tint"`
	TintBody       string `json:"tint_body"`
	TintAlignmentA string `json:"tint_alignment_a"`
	TintAlignmentB string `json:"tint_alignment_b"`
	Pattern        string `json:"pattern"`
	PatternBody    string `json:"pattern_body"`
	Pip            string `json:"pip"`
	Chassis        string `json:"chassis"`
	Bracer         string `json:"bracer"`
	Booster        string `json:"booster"`
	Title          string `json:"title"`
	Tag            string `json:"tag"`
	Banner         string `json:"banner"`
	Medal          string `json:"medal"`
	GoalFX         string `json:"goal_fx"`
	Emissive       string `json:"emissive"`
}

func DefaultGameProfiles(xplatformid XPlatformID, displayname string) GameProfiles {
	return GameProfiles{
		Client: DefaultClientProfile(xplatformid, displayname),
		Server: DefaultServerProfile(xplatformid, displayname),
	}
}

func DefaultClientProfile(xplatformid XPlatformID, displayname string) ClientProfile {
	return ClientProfile{
		DisplayName: displayname,
		XPlatformID: xplatformid.String(),
		Weapon:      "scout",
		Grenade:     "det",
		WeaponArm:   1,
		Ability:     "heal",
		Legal: Legal{
			PointsPolicyVersion: 1,
			EulaVersion:         1,
			GameAdminVersion:    1,
			SplashScreenVersion: 2,
		},
		NPE: NPE{
			Lobby: NPEItem{Completed: true},

			FirstMatch:        NPEItem{Completed: true},
			Movement:          NPEItem{Completed: true},
			ArenaBasics:       NPEItem{Completed: true},
			SocialTabSeen:     Versioned{Version: 1},
			Pointer:           Versioned{Version: 1},
			BlueTintTabSeen:   Versioned{Version: 1},
			HeraldryTabSeen:   Versioned{Version: 1},
			OrangeTintTabSeen: Versioned{Version: 1},
		},
		Customization: Customization{
			BattlePassSeasonPOIVersion: 0,
			NewUnlocksPOIVersion:       1,
			StoreEntryPOIVersion:       0,
			ClearNewUnlocksVersion:     1,
		},
		Social: Social{
			CommunityValuesVersion: 1,
			SetupVersion:           1,
			Group:                  "90DD4DB5-B5DD-4655-839E-FDBE5F4BC0BF",
		},
		NewUnlocks: []interface{}{},
	}
}
func DefaultServerProfile(xplatformid XPlatformID, displayname string) ServerProfile {
	return ServerProfile{
		Version:         4,
		PublisherLock:   "rad15_live",
		PurchasedCombat: 1,
		LobbyVersion:    1680630467,
		Loadout: LoadoutSection{
			Instances: LoadoutInstance{
				Unified: LoadoutSlots{
					Emote:          "emote_blink_smiley_a",
					Decal:          "decal_default",
					Tint:           "tint_neutral_a_default",
					TintAlignmentA: "tint_blue_a_default",
					TintAlignmentB: "tint_orange_a_default",
					Pattern:        "pattern_default",
					Pip:            "Rwd_decalback_default",
					Chassis:        "Rwd_chassis_body_s11_a",
					Bracer:         "Rwd_bracer_default",
					Booster:        "Rwd_booster_default",
					Title:          "Rwd_title_title_default",
					Tag:            "Rwd_tag_s1_a_secondary",
					Banner:         "Rwd_banner_s1_default",
					Medal:          "Rwd_medal_default",
					GoalFX:         "Rwd_goal_fx_default",
					SecondEmote:    "emote_blink_smiley_a",
					Emissive:       "emissive_default",
					TintBody:       "tint_neutral_a_default",
					PatternBody:    "pattern_default",
					DecalBody:      "decal_default",
				},
			},
			Number: 1,
		},
		Stats: StatSection{
			Arena: ArenaStats{
				Level: LevelStatItem{
					Cnt: 1,
					Op:  "add",
					Val: 1,
				},
				Goals:                        StatItem{Op: "add", Val: 0},
				TopSpeedsTotal:               StatItem{Op: "add", Val: 0},
				HighestArenaWinStreak:        StatItem{Op: "max", Val: 0},
				ArenaWinPercentage:           StatItem{Op: "rep", Val: 0},
				ArenaWins:                    StatItem{Op: "add", Val: 0},
				GoalsPerGame:                 StatItem{Op: "rep", Val: 0},
				Points:                       StatItem{Op: "add", Val: 0},
				Interceptions:                StatItem{Op: "add", Val: 0},
				ThreePointGoals:              StatItem{Op: "add", Val: 0},
				Clears:                       StatItem{Op: "add", Val: 0},
				BounceGoals:                  StatItem{Op: "add", Val: 0},
				PossessionTime:               StatItem{Op: "add", Val: 0},
				HatTricks:                    StatItem{Op: "add", Val: 0},
				ShotsOnGoal:                  StatItem{Op: "add", Val: 0},
				HighestPoints:                StatItem{Op: "max", Val: 0},
				GoalScorePercentage:          StatItem{Op: "rep", Val: 0},
				AveragePossessionTimePerGame: StatItem{Op: "rep", Val: 0},
				AverageTopSpeedPerGame:       StatItem{Op: "rep", Val: 0},
				AveragePointsPerGame:         StatItem{Op: "rep", Val: 0},
				ArenaMVPPercentage:           StatItem{Op: "rep", Val: 0},
				ArenaMVPs:                    StatItem{Op: "add", Val: 0},
				CurrentArenaWinStreak:        StatItem{Op: "add", Val: 0},
				CurrentArenaMVPStreak:        StatItem{Op: "add", Val: 0},
				HighestArenaMVPStreak:        StatItem{Op: "max", Val: 0},
				XP:                           StatItem{Op: "add", Val: 0},
				ShotsOnGoalAgainst:           StatItem{Op: "add", Val: 0},
				ArenaLosses:                  StatItem{Op: "add", Val: 0},
				Catches:                      StatItem{Op: "add", Val: 0},
				StunsPerGame:                 StatItem{Op: "rep", Val: 0},
				HighestStuns:                 StatItem{Op: "max", Val: 0},
				Steals:                       StatItem{Op: "add", Val: 0},
				Stuns:                        StatItem{Op: "add", Val: 0},
				PunchesReceived:              StatItem{Op: "add", Val: 0},
				Passes:                       StatItem{Op: "add", Val: 0},
				Blocks:                       StatItem{Op: "add", Val: 0},
				BlockPercentage:              StatItem{Op: "rep", Val: 0},
			},
			Combat: CombatStats{
				Level: LevelStatItem{
					Cnt: 1,
					Op:  "add",
					Val: 1,
				},
			},
		},
		Unlocks: UnlockSection{
			Arena: ArenaUnlocks{
				DecalCombatFlamingoA:   true,
				DecalCombatLogoA:       true,
				DecalDefault:           true,
				DecalSheldonA:          true,
				EmoteBlinkSmileyA:      true,
				EmoteDefault:           true,
				EmoteDizzyEyesA:        true,
				LoadoutNumber:          true,
				PatternDefault:         true,
				PatternLightningA:      true,
				RwdBannerS1Default:     true,
				RwdBoosterDefault:      true,
				RwdBracerDefault:       true,
				RwdChassisBodyS11A:     true,
				RwdDecalbackDefault:    true,
				RwdDecalborderDefault:  true,
				RwdMedalDefault:        true,
				RwdTagDefault:          true,
				RwdTagS1ASecondary:     true,
				RwdTitleTitleDefault:   true,
				TintBlueADefault:       true,
				TintNeutralADefault:    true,
				TintNeutralAS10Default: true,
				TintOrangeADefault:     true,
				RwdGoalFxDefault:       true,
				EmissiveDefault:        true,
			},
			Combat: CombatUnlocks{
				RwdBoosterS10:      true,
				RwdChassisBodyS10A: true,
			},
		},
		XPlatformID: xplatformid.String(),
		DisplayName: displayname,
	}
}
