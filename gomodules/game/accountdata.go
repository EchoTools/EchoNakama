// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    gameProfiles, err := UnmarshalGameProfiles(bytes)
//    bytes, err = gameProfiles.Marshal()

package game

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// MergePlayerData merges a partial PlayerData with a template PlayerData.
func (base *EchoPlayerPreferences) Merge(partial *EchoPlayerPreferences) (*EchoPlayerPreferences, error) {

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

// Profiles represents the 'profile' field in the JSON data
type GameProfiles struct {
	Client EchoPlayerPreferences `json:"client"`
	Server ServerProfile         `json:"server"`
}

func UnmarshalGameProfiles(data []byte) (GameProfiles, error) {
	var r GameProfiles
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GameProfiles) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Structure
func customFunc(fl validator.FieldLevel) bool {

	if fl.Field().String() == "invalid" {
		return false
	}

	return true
}

type EchoPlayerPreferences struct {
	// WARNING: EchoVR dictates this struct/schema.
	DisplayName     string `json:"displayname" validate:"required,printascii,min=3,max=20"`
	EchoUserIdToken string `json:"xplatformid" validate:"required"`

	// The team name shown on the spectator scoreboard overlay
	TeamName           string            `json:"teamname,omitempty" validate:"omitempty,ascii"`
	CombatWeapon       string            `json:"weapon" validate:"omitnil,oneof=assault blaster rocket scout magnum smg chain rifle"`
	CombatGrenade      string            `json:"grenade" validate:"omitnil,oneof=arc burst det stun loc"`
	CombatDominantHand uint8             `json:"weaponarm" validate:"eq=0|eq=1"`
	ModifyTime         int64             `json:"modifytime" validate:"required"`
	CombatAbility      string            `json:"ability" validate:"required"`
	LegalConsents      LegalConsents     `json:"legal" validate:"required"`
	MutedPlayers       Players           `json:"mute" validate:"required"`
	GhostedPlayers     Players           `json:"ghost" validate:"required"`
	NewPlayerProgress  NewPlayerProgress `json:"npe" validate:"required"`
	Customization      Customization     `json:"customization" validate:"required"`
	Social             Social            `json:"social" validate:"required"`
	NewUnlocks         []int64           `json:"newunlocks" validate:"required"`
}

type Customization struct {
	// WARNING: EchoVR dictates this struct/schema.
	BattlePassSeasonPoiVersion uint16 `json:"battlepass_season_poi_version" validate:"required"` // Battle pass season point of interest version (manually set to 3246)
	NewUnlocksPoiVersion       uint16 `json:"new_unlocks_poi_version" validate:"required"`       // New unlocks point of interest version
	StoreEntryPoiVersion       uint16 `json:"store_entry_poi_version" validate:"required"`       // Store entry point of interest version
	ClearNewUnlocksVersion     uint16 `json:"clear_new_unlocks_version" validate:"required"`     // Clear new unlocks version
}

type Players struct {
	// WARNING: EchoVR dictates this struct/schema.
	UserIds []string `json:"users" validate:"required"`
}

type LegalConsents struct {
	// WARNING: EchoVR dictates this struct/schema.
	PointsPolicyVersion int64 `json:"points_policy_version" validate:"required"`
	EulaVersion         int64 `json:"eula_version" validate:"required"`
	GameAdminVersion    int64 `json:"game_admin_version" validate:"required"`
	SplashScreenVersion int64 `json:"splash_screen_version" validate:"required"`
	GroupsLegalVersion  int64 `json:"groups_legal_version" validate:"required"`
}

type NewPlayerProgress struct {
	// WARNING: EchoVR dictates this struct/schema.
	Lobby             NpeMilestone `json:"lobby" validate:"required"`                // User has completed the tutorial
	FirstMatch        NpeMilestone `json:"firstmatch" validate:"required"`           // User has completed their first match
	Movement          NpeMilestone `json:"movement" validate:"required"`             // User has completed the movement tutorial
	ArenaBasics       NpeMilestone `json:"arenabasics" validate:"required"`          // User has completed the arena basics tutorial
	SocialTabSeen     Versioned    `json:"social_tab_seen" validate:"required"`      // User has seen the social tab
	Pointer           Versioned    `json:"pointer" validate:"required"`              // User has seen the pointer ?
	BlueTintTabSeen   Versioned    `json:"blue_tint_tab_seen" validate:"required"`   // User has seen the blue tint tab in the character room
	HeraldryTabSeen   Versioned    `json:"heraldry_tab_seen" validate:"required"`    // User has seen the heraldry tab in the character room
	OrangeTintTabSeen Versioned    `json:"orange_tint_tab_seen" validate:"required"` // User has seen the orange tint tab in the character room
}

type NpeMilestone struct {
	// WARNING: EchoVR dictates this struct/schema.
	Completed bool `json:"completed" validate:"boolean required"` // User has completed the milestone
}

type Versioned struct {
	// WARNING: EchoVR dictates this struct/schema.
	Version int `json:"version" validate:"gte=0 required"` // A version number, 1 is seen, 0 is not seen ?
}

type Social struct {
	// WARNING: EchoVR dictates this struct/schema.
	CommunityValuesVersion int64  `json:"community_values_version" validate:"gte=0,required"`
	SetupVersion           int64  `json:"setup_version" validate:"gte=0,required"`
	Group                  string `json:"group" validate:"uuid4,required"`
}

type ServerProfile struct {
	// WARNING: EchoVR dictates this struct/schema.
	DisplayName       string            `json:"displayname" validate:"required,printascii,min=3,max=20"`
	EchoUserIdToken   string            `json:"xplatformid" validate:"required"`
	SchemaVersion     int16             `json:"_version" validate:"gte=0,required"`            // Version of the schema(?)
	PublisherLock     string            `json:"publisher_lock" validate:"required"`            // unused atm
	PurchasedCombat   int8              `json:"purchasedcombat" validate:"eq=0|eq=1,required"` // unused (combat was made free)
	LobbyVersion      int64             `json:"lobbyversion" validate:"gte=0,required"`        // set from the login request (no known effect)
	ModifyTime        int64             `json:"modifytime" validate:"gte=0,required"`
	LoginTime         int64             `json:"logintime" validate:"gte=0,required"`
	UpdateTime        int64             `json:"updatetime" validate:"gte=0,required"`
	CreateTime        int64             `json:"createtime" validate:"gte=0,required"`
	Statistics        PlayerStatistics  `json:"stats" validate:"required"`
	MaybeStale        *bool             `json:"maybestale" validate:"boolean,required"`
	UnlockedCosmetics UnlockedCosmetics `json:"unlocks"`
	EquippedCosmetics EquippedCosmetics `json:"loadout"`
	Social            Social            `json:"social"`
	Achievements      interface{}       `json:"achievements"`
	RewardState       interface{}       `json:"reward_state"`
	DeveloperFeatures DeveloperFeatures `json:"dev"`
}

type DeveloperFeatures struct {
	// WARNING: EchoVR dictates this struct/schema.
	DisableAfkTimeout bool   `json:"disable_afk_timeout" validate:"boolean required"`
	EchoUserIdToken   string `json:"xplatformid" validate:"required"`
}

type PlayerStatistics struct {
	Arena  ArenaStatistics  `json:"arena"`
	Combat CombatStatistics `json:"combat"`
}

type ArenaStatistics struct {
	Level                        Level               `json:"Level"`
	Stuns                        DiscreteStatistic   `json:"Stuns,omitempty"`
	TopSpeedsTotal               ContinuousStatistic `json:"TopSpeedsTotal,omitempty"`
	HighestArenaWinStreak        DiscreteStatistic   `json:"HighestArenaWinStreak"`
	ArenaWinPercentage           ContinuousStatistic `json:"ArenaWinPercentage,omitempty"`
	ArenaWINS                    DiscreteStatistic   `json:"ArenaWins,omitempty"`
	ShotsOnGoalAgainst           DiscreteStatistic   `json:"ShotsOnGoalAgainst,omitempty"`
	Clears                       DiscreteStatistic   `json:"Clears,omitempty"`
	AssistsPerGame               ContinuousStatistic `json:"AssistsPerGame,omitempty"`
	Passes                       DiscreteStatistic   `json:"Passes,omitempty"`
	AveragePossessionTimePerGame ContinuousStatistic `json:"AveragePossessionTimePerGame,omitempty"`
	Catches                      DiscreteStatistic   `json:"Catches,omitempty"`
	PossessionTime               ContinuousStatistic `json:"PossessionTime,omitempty"`
	StunsPerGame                 ContinuousStatistic `json:"StunsPerGame,omitempty"`
	ShotsOnGoal                  DiscreteStatistic   `json:"ShotsOnGoal,omitempty"`
	PunchesReceived              DiscreteStatistic   `json:"PunchesReceived,omitempty"`
	CurrentArenaWinStreak        DiscreteStatistic   `json:"CurrentArenaWinStreak,omitempty"`
	Assists                      DiscreteStatistic   `json:"Assists,omitempty"`
	Interceptions                DiscreteStatistic   `json:"Interceptions,omitempty"`
	HighestStuns                 DiscreteStatistic   `json:"HighestStuns,omitempty"`
	AverageTopSpeedPerGame       ContinuousStatistic `json:"AverageTopSpeedPerGame,omitempty"`
	XP                           DiscreteStatistic   `json:"XP,omitempty"`
	ArenaLosses                  DiscreteStatistic   `json:"ArenaLosses,omitempty"`
	SavesPerGame                 ContinuousStatistic `json:"SavesPerGame,omitempty"`
	Blocks                       DiscreteStatistic   `json:"Blocks,omitempty"`
	Saves                        DiscreteStatistic   `json:"Saves,omitempty"`
	HighestSaves                 DiscreteStatistic   `json:"HighestSaves,omitempty"`
	GoalSavePercentage           ContinuousStatistic `json:"GoalSavePercentage,omitempty"`
	BlockPercentage              ContinuousStatistic `json:"BlockPercentage,omitempty"`
	GoalsPerGame                 ContinuousStatistic `json:"GoalsPerGame,omitempty"`
	Points                       DiscreteStatistic   `json:"Points,omitempty"`
	Goals                        DiscreteStatistic   `json:"Goals,omitempty"`
	Steals                       DiscreteStatistic   `json:"Steals,omitempty"`
	TwoPointGoals                DiscreteStatistic   `json:"TwoPointGoals,omitempty"`
	HighestPoints                DiscreteStatistic   `json:"HighestPoints,omitempty"`
	GoalScorePercentage          ContinuousStatistic `json:"GoalScorePercentage,omitempty"`
	AveragePointsPerGame         ContinuousStatistic `json:"AveragePointsPerGame,omitempty"`
	ThreePointGoals              DiscreteStatistic   `json:"ThreePointGoals,omitempty"`
	BounceGoals                  DiscreteStatistic   `json:"BounceGoals,omitempty"`
	ArenaMVPPercentage           ContinuousStatistic `json:"ArenaMVPPercentage,omitempty"`
	ArenaMVPS                    DiscreteStatistic   `json:"ArenaMVPs,omitempty"`
	CurrentArenaMVPStreak        DiscreteStatistic   `json:"CurrentArenaMVPStreak,omitempty"`
	HighestArenaMVPStreak        DiscreteStatistic   `json:"HighestArenaMVPStreak,omitempty"`
	HeadbuttGoals                DiscreteStatistic   `json:"HeadbuttGoals,omitempty"`
	HatTricks                    DiscreteStatistic   `json:"HatTricks,omitempty"`
}
type CombatStatistics struct {
	Level                              Level                      `json:"Level"`
	CombatAssists                      CountedDiscreteStatistic   `json:"CombatAssists,omitempty"`
	CombatObjectiveDamage              CountedContinuousStatistic `json:"CombatObjectiveDamage,omitempty"`
	CombatEliminations                 CountedDiscreteStatistic   `json:"CombatEliminations,omitempty"`
	CombatDamageAbsorbed               ContinuousStatistic        `json:"CombatDamageAbsorbed,omitempty"`
	CombatWINS                         CountedDiscreteStatistic   `json:"CombatWins,omitempty"`
	CombatDamageTaken                  CountedContinuousStatistic `json:"CombatDamageTaken,omitempty"`
	CombatWinPercentage                CountedDiscreteStatistic   `json:"CombatWinPercentage,omitempty"`
	CombatStuns                        CountedDiscreteStatistic   `json:"CombatStuns,omitempty"`
	CombatKills                        CountedDiscreteStatistic   `json:"CombatKills,omitempty"`
	CombatPointCaptureGamesPlayed      DiscreteStatistic          `json:"CombatPointCaptureGamesPlayed,omitempty"`
	CombatPointCaptureWINS             DiscreteStatistic          `json:"CombatPointCaptureWins,omitempty"`
	CombatObjectiveTime                CountedContinuousStatistic `json:"CombatObjectiveTime,omitempty"`
	CombatAverageEliminationDeathRatio CountedContinuousStatistic `json:"CombatAverageEliminationDeathRatio,omitempty"`
	CombatPointCaptureWinPercentage    ContinuousStatistic        `json:"CombatPointCaptureWinPercentage,omitempty"`
	CombatDeaths                       CountedDiscreteStatistic   `json:"CombatDeaths,omitempty"`
	CombatDamage                       CountedContinuousStatistic `json:"CombatDamage,omitempty"`
	CombatObjectiveEliminations        CountedDiscreteStatistic   `json:"CombatObjectiveEliminations,omitempty"`
	CombatBestEliminationStreak        CountedDiscreteStatistic   `json:"CombatBestEliminationStreak,omitempty"`
	CombatSoloKills                    CountedDiscreteStatistic   `json:"CombatSoloKills,omitempty"`
	CombatHeadshotKills                CountedDiscreteStatistic   `json:"CombatHeadshotKills,omitempty"`
	XP                                 CountedDiscreteStatistic   `json:"XP,omitempty"`
	CombatMVPS                         DiscreteStatistic          `json:"CombatMVPs,omitempty"`
	CombatHillDefends                  DiscreteStatistic          `json:"CombatHillDefends,omitempty"`
	CombatPayloadWINS                  CountedDiscreteStatistic   `json:"CombatPayloadWins,omitempty"`
	CombatPayloadGamesPlayed           CountedDiscreteStatistic   `json:"CombatPayloadGamesPlayed,omitempty"`
	CombatPayloadWinPercentage         CountedDiscreteStatistic   `json:"CombatPayloadWinPercentage,omitempty"`
	CombatHillCaptures                 DiscreteStatistic          `json:"CombatHillCaptures,omitempty"`
	CombatHealing                      CountedContinuousStatistic `json:"CombatHealing,omitempty"`
	CombatTeammateHealing              CountedContinuousStatistic `json:"CombatTeammateHealing,omitempty"`
	CombatLosses                       CountedDiscreteStatistic   `json:"CombatLosses,omitempty"`
}

type CountedDiscreteStatistic struct {
	Count   int64  `json:"cnt" validate:"gte=0,required_with=Operand Value"`
	Operand string `json:"op" validate:"oneof=add rep max, required_with=Count Value"`
	Value   uint64 `json:"val" validate:"gte=0,required_with=Operand Count"`
}

type CountedContinuousStatistic struct {
	Operand string  `json:"op" validate:"oneof=add rep max,required_with=Count Value"`
	Value   float64 `json:"val" validate:"gte=0,required_with=Operand Count"`
	Count   uint64  `json:"cnt" validate:"gte=0,required_with=Operand Value"`
}

type DiscreteStatistic struct {
	Operand string `json:"op" validate:"oneof=add rep max,required_with=Value"`
	Value   uint64 `json:"val" validate:"gte=0,required_with=Operand"`
}

type ContinuousStatistic struct {
	Operand string  `json:"op" validate:"oneof=add rep max,required_with=Value"`
	Value   float64 `json:"val" validate:"gte=0,required_with=Operand"`
}

type Level struct {
	Count   uint8  `json:"cnt" validate:"gte=0,required"`
	Operand string `json:"op" validate:"oneof=add,required"`
	Value   uint8  `json:"val" validate:"gte=0,required"`
}

type DailyStats struct {
	Stuns                        DiscreteStatistic   `json:"Stuns,omitempty"`
	XP                           DiscreteStatistic   `json:"XP,omitempty"`
	TopSpeedsTotal               ContinuousStatistic `json:"TopSpeedsTotal,omitempty"`
	HighestArenaWinStreak        DiscreteStatistic   `json:"HighestArenaWinStreak,omitempty"`
	ArenaWinPercentage           DiscreteStatistic   `json:"ArenaWinPercentage,omitempty"`
	ArenaWINS                    DiscreteStatistic   `json:"ArenaWins,omitempty"`
	ShotsOnGoalAgainst           DiscreteStatistic   `json:"ShotsOnGoalAgainst,omitempty"`
	Clears                       DiscreteStatistic   `json:"Clears,omitempty"`
	AssistsPerGame               ContinuousStatistic `json:"AssistsPerGame,omitempty"`
	Passes                       DiscreteStatistic   `json:"Passes,omitempty"`
	AveragePossessionTimePerGame ContinuousStatistic `json:"AveragePossessionTimePerGame,omitempty"`
	Catches                      DiscreteStatistic   `json:"Catches,omitempty"`
	PossessionTime               ContinuousStatistic `json:"PossessionTime,omitempty"`
	StunsPerGame                 DiscreteStatistic   `json:"StunsPerGame,omitempty"`
	ShotsOnGoal                  DiscreteStatistic   `json:"ShotsOnGoal,omitempty"`
	PunchesReceived              DiscreteStatistic   `json:"PunchesReceived,omitempty"`
	CurrentArenaWinStreak        DiscreteStatistic   `json:"CurrentArenaWinStreak,omitempty"`
	Assists                      DiscreteStatistic   `json:"Assists,omitempty"`
	Interceptions                DiscreteStatistic   `json:"Interceptions,omitempty"`
	HighestStuns                 DiscreteStatistic   `json:"HighestStuns,omitempty"`
	AverageTopSpeedPerGame       ContinuousStatistic `json:"AverageTopSpeedPerGame,omitempty"`
	ArenaLosses                  DiscreteStatistic   `json:"ArenaLosses,omitempty"`
	SavesPerGame                 ContinuousStatistic `json:"SavesPerGame,omitempty"`
	Blocks                       DiscreteStatistic   `json:"Blocks,omitempty"`
	Saves                        DiscreteStatistic   `json:"Saves,omitempty"`
	HighestSaves                 DiscreteStatistic   `json:"HighestSaves,omitempty"`
	GoalSavePercentage           ContinuousStatistic `json:"GoalSavePercentage,omitempty"`
	BlockPercentage              ContinuousStatistic `json:"BlockPercentage,omitempty"`
}

type WeelkyStats struct {
	Stuns                        DiscreteStatistic   `json:"Stuns,omitempty"`
	XP                           DiscreteStatistic   `json:"XP,omitempty"`
	TopSpeedsTotal               ContinuousStatistic `json:"TopSpeedsTotal,omitempty"`
	HighestArenaWinStreak        DiscreteStatistic   `json:"HighestArenaWinStreak,omitempty"`
	ArenaWinPercentage           DiscreteStatistic   `json:"ArenaWinPercentage,omitempty"`
	ArenaWINS                    DiscreteStatistic   `json:"ArenaWins,omitempty"`
	ShotsOnGoalAgainst           DiscreteStatistic   `json:"ShotsOnGoalAgainst,omitempty"`
	Clears                       DiscreteStatistic   `json:"Clears,omitempty"`
	AssistsPerGame               ContinuousStatistic `json:"AssistsPerGame,omitempty"`
	Passes                       DiscreteStatistic   `json:"Passes,omitempty"`
	AveragePossessionTimePerGame ContinuousStatistic `json:"AveragePossessionTimePerGame,omitempty"`
	Catches                      DiscreteStatistic   `json:"Catches,omitempty"`
	PossessionTime               ContinuousStatistic `json:"PossessionTime,omitempty"`
	StunsPerGame                 ContinuousStatistic `json:"StunsPerGame,omitempty"`
	ShotsOnGoal                  DiscreteStatistic   `json:"ShotsOnGoal,omitempty"`
	PunchesReceived              DiscreteStatistic   `json:"PunchesReceived,omitempty"`
	CurrentArenaWinStreak        DiscreteStatistic   `json:"CurrentArenaWinStreak,omitempty"`
	Assists                      DiscreteStatistic   `json:"Assists,omitempty"`
	Interceptions                DiscreteStatistic   `json:"Interceptions,omitempty"`
	HighestStuns                 DiscreteStatistic   `json:"HighestStuns,omitempty"`
	AverageTopSpeedPerGame       ContinuousStatistic `json:"AverageTopSpeedPerGame,omitempty"`
	ArenaLosses                  DiscreteStatistic   `json:"ArenaLosses,omitempty"`
	SavesPerGame                 ContinuousStatistic `json:"SavesPerGame,omitempty"`
	Blocks                       DiscreteStatistic   `json:"Blocks,omitempty"`
	Saves                        DiscreteStatistic   `json:"Saves,omitempty"`
	HighestSaves                 DiscreteStatistic   `json:"HighestSaves,omitempty"`
	GoalSavePercentage           ContinuousStatistic `json:"GoalSavePercentage,omitempty"`
	BlockPercentage              ContinuousStatistic `json:"BlockPercentage,omitempty"`
	GoalsPerGame                 ContinuousStatistic `json:"GoalsPerGame,omitempty"`
	Points                       DiscreteStatistic   `json:"Points,omitempty"`
	Goals                        DiscreteStatistic   `json:"Goals,omitempty"`
	Steals                       DiscreteStatistic   `json:"Steals,omitempty"`
	TwoPointGoals                DiscreteStatistic   `json:"TwoPointGoals,omitempty"`
	HighestPoints                DiscreteStatistic   `json:"HighestPoints,omitempty"`
	GoalScorePercentage          ContinuousStatistic `json:"GoalScorePercentage,omitempty"`
	AveragePointsPerGame         ContinuousStatistic `json:"AveragePointsPerGame,omitempty"`
}

type EquippedCosmetics struct {
	Instances Instances `json:"instances"`
	Number    int64     `json:"number"`
}

type Instances struct {
	Unified Unified `json:"unified"`
}

type Unified struct {
	Slots Slots `json:"slots"`
}

type Slots struct {
	Decal          string `json:"decal"`
	DecalBody      string `json:"decal_body"`
	Emote          string `json:"emote"`
	Secondemote    string `json:"secondemote"`
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
	GoalFx         string `json:"goal_fx"`
	SecondEmote    string `json:"second_emote"`
	Emissive       string `json:"emissive"`
}

type UnlockedCosmetics struct {
	Arena  UnlocksArena  `json:"arena"`
	Combat UnlocksCombat `json:"combat"`
}

type UnlocksArena struct {
	DecalCombatFlamingoA        bool `json:"decal_combat_flamingo_a"`
	DecalCombatLogoA            bool `json:"decal_combat_logo_a"`
	DecalDefault                bool `json:"decal_default"`
	DecalSheldonA               bool `json:"decal_sheldon_a"`
	EmoteBlinkSmileyA           bool `json:"emote_blink_smiley_a"`
	EmoteDefault                bool `json:"emote_default"`
	EmoteDizzyEyesA             bool `json:"emote_dizzy_eyes_a"`
	LoadoutNumber               bool `json:"loadout_number"`
	PatternDefault              bool `json:"pattern_default"`
	PatternLightningA           bool `json:"pattern_lightning_a"`
	RWDBannerS1Default          bool `json:"rwd_banner_s1_default"`
	RWDBoosterDefault           bool `json:"rwd_booster_default"`
	RWDBracerDefault            bool `json:"rwd_bracer_default"`
	RWDChassisBodyS11A          bool `json:"rwd_chassis_body_s11_a"`
	RWDDecalbackDefault         bool `json:"rwd_decalback_default"`
	RWDDecalborderDefault       bool `json:"rwd_decalborder_default"`
	RWDMedalDefault             bool `json:"rwd_medal_default"`
	RWDTagDefault               bool `json:"rwd_tag_default"`
	RWDTagS1ASecondary          bool `json:"rwd_tag_s1_a_secondary"`
	RWDTitleTitleDefault        bool `json:"rwd_title_title_default"`
	TintBlueADefault            bool `json:"tint_blue_a_default"`
	TintNeutralADefault         bool `json:"tint_neutral_a_default"`
	TintNeutralAS10Default      bool `json:"tint_neutral_a_s10_default"`
	TintOrangeADefault          bool `json:"tint_orange_a_default"`
	RWDGoalFxDefault            bool `json:"rwd_goal_fx_default"`
	EmissiveDefault             bool `json:"emissive_default"`
	PatternTrianglesA           bool `json:"pattern_triangles_a"`
	RWDTagS1ISecondary          bool `json:"rwd_tag_s1_i_secondary"`
	PatternHoneycombTripleA     bool `json:"pattern_honeycomb_triple_a"`
	RWDTitleTitleE              bool `json:"rwd_title_title_e"`
	RWDTagS1OSecondary          bool `json:"rwd_tag_s1_o_secondary"`
	TintOrangeJDefault          bool `json:"tint_orange_j_default"`
	TintNeutralLDefault         bool `json:"tint_neutral_l_default"`
	TintOrangeEDefault          bool `json:"tint_orange_e_default"`
	PatternAnglesA              bool `json:"pattern_angles_a"`
	DecalRayGunA                bool `json:"decal_ray_gun_a"`
	RWDTagS1ESecondary          bool `json:"rwd_tag_s1_e_secondary"`
	RWDPatternSaltA             bool `json:"rwd_pattern_salt_a"`
	PatternSquigglesA           bool `json:"pattern_squiggles_a"`
	TintOrangeIDefault          bool `json:"tint_orange_i_default"`
	PatternWeaveA               bool `json:"pattern_weave_a"`
	TintBlueKDefault            bool `json:"tint_blue_k_default"`
	EmoteWinkyTongueA           bool `json:"emote_winky_tongue_a"`
	EmoteStinkyPoopA            bool `json:"emote_stinky_poop_a"`
	DecalRoseA                  bool `json:"decal_rose_a"`
	DecalMusicNoteA             bool `json:"decal_music_note_a"`
	TintBlueHDefault            bool `json:"tint_blue_h_default"`
	DecalBombA                  bool `json:"decal_bomb_a"`
	EmoteSickFaceA              bool `json:"emote_sick_face_a"`
	DecalAlienHeadA             bool `json:"decal_alien_head_a"`
	TintNeutralEDefault         bool `json:"tint_neutral_e_default"`
	RWDEmissive0011             bool `json:"rwd_emissive_0011"`
	DecalCombatPigA             bool `json:"decal_combat_pig_a"`
	DecalDiscA                  bool `json:"decal_disc_a"`
	TintOrangeDDefault          bool `json:"tint_orange_d_default"`
	HeraldryDefault             bool `json:"heraldry_default"`
	TintBlueBDefault            bool `json:"tint_blue_b_default"`
	TintNeutralDDefault         bool `json:"tint_neutral_d_default"`
	RWDEmissive0014             bool `json:"rwd_emissive_0014"`
	TintChassisDefault          bool `json:"tint_chassis_default"`
	DecalDinosaurA              bool `json:"decal_dinosaur_a"`
	RWDTagS1HSecondary          bool `json:"rwd_tag_s1_h_secondary"`
	DecalSpiderA                bool `json:"decal_spider_a"`
	DecalCombatMeteorA          bool `json:"decal_combat_meteor_a"`
	RWDPip0007                  bool `json:"rwd_pip_0007"`
	PatternInsetCubesA          bool `json:"pattern_inset_cubes_a"`
	TintOrangeFDefault          bool `json:"tint_orange_f_default"`
	PatternBananasA             bool `json:"pattern_bananas_a"`
	RWDBannerS1Basic            bool `json:"rwd_banner_s1_basic"`
	RWDBannerS1BoldStripe       bool `json:"rwd_banner_s1_bold_stripe"`
	RWDEmissive0012             bool `json:"rwd_emissive_0012"`
	RWDPatternRageWolfA         bool `json:"rwd_pattern_rage_wolf_a"`
	RWDTagS1MSecondary          bool `json:"rwd_tag_s1_m_secondary"`
	DecalCombatPulsarA          bool `json:"decal_combat_pulsar_a"`
	TintNeutralKDefault         bool `json:"tint_neutral_k_default"`
	EmoteSleepyZzzA             bool `json:"emote_sleepy_zzz_a"`
	TintNeutralNDefault         bool `json:"tint_neutral_n_default"`
	PatternDigitalCamoA         bool `json:"pattern_digital_camo_a"`
	EmoteKissyLipsA             bool `json:"emote_kissy_lips_a"`
	TintOrangeKDefault          bool `json:"tint_orange_k_default"`
	RWDPip0010                  bool `json:"rwd_pip_0010"`
	TintNeutralGDefault         bool `json:"tint_neutral_g_default"`
	RWDTagS1KSecondary          bool `json:"rwd_tag_s1_k_secondary"`
	TintOrangeGDefault          bool `json:"tint_orange_g_default"`
	RWDPip0001                  bool `json:"rwd_pip_0001"`
	DecalCombatRageBearA        bool `json:"decal_combat_rage_bear_a"`
	RWDTagS1BSecondary          bool `json:"rwd_tag_s1_b_secondary"`
	DecalLightningBoltA         bool `json:"decal_lightning_bolt_a"`
	DecalKoiFishA               bool `json:"decal_koi_fish_a"`
	RWDBannerS1Chevrons         bool `json:"rwd_banner_s1_chevrons"`
	EmoteExclamationPointA      bool `json:"emote_exclamation_point_a"`
	EmoteAngryFaceA             bool `json:"emote_angry_face_a"`
	PatternGearsA               bool `json:"pattern_gears_a"`
	TintBlueJDefault            bool `json:"tint_blue_j_default"`
	RWDEmissive0010             bool `json:"rwd_emissive_0010"`
	RWDTitleTitleD              bool `json:"rwd_title_title_d"`
	RWDGoalFx0002               bool `json:"rwd_goal_fx_0002"`
	DecalCrosshairA             bool `json:"decal_crosshair_a"`
	PatternScalesA              bool `json:"pattern_scales_a"`
	DecalBullseyeA              bool `json:"decal_bullseye_a"`
	EmoteCryingFaceA            bool `json:"emote_crying_face_a"`
	DecalCombatDemonA           bool `json:"decal_combat_demon_a"`
	TintBlueFDefault            bool `json:"tint_blue_f_default"`
	RWDTagS1JSecondary          bool `json:"rwd_tag_s1_j_secondary"`
	EmoteHourglassA             bool `json:"emote_hourglass_a"`
	RWDMedalS1ArenaSilver       bool `json:"rwd_medal_s1_arena_silver"`
	EmoteBrokenHeartA           bool `json:"emote_broken_heart_a"`
	TintBlueCDefault            bool `json:"tint_blue_c_default"`
	PatternTigerA               bool `json:"pattern_tiger_a"`
	TintNeutralFDefault         bool `json:"tint_neutral_f_default"`
	TintBlueIDefault            bool `json:"tint_blue_i_default"`
	TintOrangeBDefault          bool `json:"tint_orange_b_default"`
	RWDPip0014                  bool `json:"rwd_pip_0014"`
	EmoteWifiSymbolA            bool `json:"emote_wifi_symbol_a"`
	EmoteClockA                 bool `json:"emote_clock_a"`
	DecalRageWolfA              bool `json:"decal_rage_wolf_a"`
	RWDBoosterS11S1ARetro       bool `json:"rwd_booster_s11_s1_a_retro"`
	EmoteDeadFaceA              bool `json:"emote_dead_face_a"`
	PatternLeopardA             bool `json:"pattern_leopard_a"`
	PatternDiamondPlateA        bool `json:"pattern_diamond_plate_a"`
	EmoteHeartEyesA             bool `json:"emote_heart_eyes_a"`
	DecalEagleA                 bool `json:"decal_eagle_a"`
	PatternHawaiianA            bool `json:"pattern_hawaiian_a"`
	RWDPip0011                  bool `json:"rwd_pip_0011"`
	EmoteTearDropA              bool `json:"emote_tear_drop_a"`
	RWDPip0005                  bool `json:"rwd_pip_0005"`
	EmoteMoustacheA             bool `json:"emote_moustache_a"`
	DecalSaltShakerA            bool `json:"decal_salt_shaker_a"`
	RWDPatternPizzaA            bool `json:"rwd_pattern_pizza_a"`
	PatternPineappleA           bool `json:"pattern_pineapple_a"`
	EmotePizzaDance             bool `json:"emote_pizza_dance"`
	RWDPip0006                  bool `json:"rwd_pip_0006"`
	PatternStreaksA             bool `json:"pattern_streaks_a"`
	DecalFireballA              bool `json:"decal_fireball_a"`
	RWDTagS1VSecondary          bool `json:"rwd_tag_s1_v_secondary"`
	EmoteEyeRollA               bool `json:"emote_eye_roll_a"`
	DecalCombatPizzaA           bool `json:"decal_combat_pizza_a"`
	RWDPatternCupcakeA          bool `json:"rwd_pattern_cupcake_a"`
	RWDEmissive0007             bool `json:"rwd_emissive_0007"`
	DecalRadioactiveA           bool `json:"decal_radioactive_a"`
	RWDEmissive0025             bool `json:"rwd_emissive_0025"`
	PatternStringsA             bool `json:"pattern_strings_a"`
	EmoteStarEyesA              bool `json:"emote_star_eyes_a"`
	PatternArrowheadsA          bool `json:"pattern_arrowheads_a"`
	RWDBoosterVintageA          bool `json:"rwd_booster_vintage_a"`
	DecalSaturnA                bool `json:"decal_saturn_a"`
	PatternPawsA                bool `json:"pattern_paws_a"`
	DecalSwordsA                bool `json:"decal_swords_a"`
	RWDEmissive0004             bool `json:"rwd_emissive_0004"`
	RWDPatternHamburgerA        bool `json:"rwd_pattern_hamburger_a"`
	PatternTreadsA              bool `json:"pattern_treads_a"`
	DecalRocketA                bool `json:"decal_rocket_a"`
	RWDBracerVintageA           bool `json:"rwd_bracer_vintage_a"`
	RWDTagS1CSecondary          bool `json:"rwd_tag_s1_c_secondary"`
	TintBlueDDefault            bool `json:"tint_blue_d_default"`
	DecalCombatSkullCrossbonesA bool `json:"decal_combat_skull_crossbones_a"`
	DecalRadioactiveBioA        bool `json:"decal_radioactive_bio_a"`
	RWDPip0008                  bool `json:"rwd_pip_0008"`
	RWDPip0009                  bool `json:"rwd_pip_0009"`
	RWDPip0013                  bool `json:"rwd_pip_0013"`
	PatternDumbbellsA           bool `json:"pattern_dumbbells_a"`
	RWDGoalFx0008               bool `json:"rwd_goal_fx_0008"`
	RWDPip0015                  bool `json:"rwd_pip_0015"`
	RWDEmissive0006             bool `json:"rwd_emissive_0006"`
	RWDEmissive0001             bool `json:"rwd_emissive_0001"`
	RWDPatternSkullA            bool `json:"rwd_pattern_skull_a"`
	RWDPatternAlienA            bool `json:"rwd_pattern_alien_a"`
	DecalCombatTrexSkullA       bool `json:"decal_combat_trex_skull_a"`
	PatternCatsA                bool `json:"pattern_cats_a"`
	PatternDotsA                bool `json:"pattern_dots_a"`
	RWDEmissive0008             bool `json:"rwd_emissive_0008"`
	RWDEmissive0009             bool `json:"rwd_emissive_0009"`
	EmoteMoneyBagA              bool `json:"emote_money_bag_a"`
	RWDEmissive0002             bool `json:"rwd_emissive_0002"`
	RWDChassisS11RetroA         bool `json:"rwd_chassis_s11_retro_a"`
	RWDEmissive0003             bool `json:"rwd_emissive_0003"`
	DecalCombatFlyingSaucerA    bool `json:"decal_combat_flying_saucer_a"`
	RWDEmissive0005             bool `json:"rwd_emissive_0005"`
	RWDEmissive0013             bool `json:"rwd_emissive_0013"`
	RWDMedalS1ArenaGold         bool `json:"rwd_medal_s1_arena_gold"`
	RWDBannerS1Tritip           bool `json:"rwd_banner_s1_tritip"`
	DecalCombatMedicA           bool `json:"decal_combat_medic_a"`
	DecalCombatCometA           bool `json:"decal_combat_comet_a"`
	DecalCombatPuppyA           bool `json:"decal_combat_puppy_a"`
	RWDBoosterS11S1AFire        bool `json:"rwd_booster_s11_s1_a_fire"`
	EmoteReticleA               bool `json:"emote_reticle_a"`
	DecalCombatOctopusA         bool `json:"decal_combat_octopus_a"`
	TintBlueEDefault            bool `json:"tint_blue_e_default"`
	RWDBannerS1Squish           bool `json:"rwd_banner_s1_squish"`
	DecalHamburgerA             bool `json:"decal_hamburger_a"`
	EmoteSkullCrossbonesA       bool `json:"emote_skull_crossbones_a"`
	EmoteGgA                    bool `json:"emote_gg_a"`
	PatternCubesA               bool `json:"pattern_cubes_a"`
	PatternSwirlA               bool `json:"pattern_swirl_a"`
	DecalBearPawA               bool `json:"decal_bear_paw_a"`
	PatternStarsA               bool `json:"pattern_stars_a"`
	TintNeutralJDefault         bool `json:"tint_neutral_j_default"`
	EmoteDollarEyesA            bool `json:"emote_dollar_eyes_a"`
	RWDChassisS8BA              bool `json:"rwd_chassis_s8b_a"`
	EmoteLoadingA               bool `json:"emote_loading_a"`
	RWDChassisS11FlameA         bool `json:"rwd_chassis_s11_flame_a"`
	DecalCombatMilitaryBadgeA   bool `json:"decal_combat_military_badge_a"`
	DecalCatA                   bool `json:"decal_cat_a"`
	PatternTableclothA          bool `json:"pattern_tablecloth_a"`
	RWDBannerS1Hourglass        bool `json:"rwd_banner_s1_hourglass"`
	TintBlueGDefault            bool `json:"tint_blue_g_default"`
	RWDPatternTrexSkullA        bool `json:"rwd_pattern_trex_skull_a"`
	RWDTagS1FSecondary          bool `json:"rwd_tag_s1_f_secondary"`
	DecalCombatIceCreamA        bool `json:"decal_combat_ice_cream_a"`
	PatternDiamondsA            bool `json:"pattern_diamonds_a"`
	TintNeutralCDefault         bool `json:"tint_neutral_c_default"`
	TintNeutralIDefault         bool `json:"tint_neutral_i_default"`
	RWDGoalFx0005               bool `json:"rwd_goal_fx_0005"`
	DecalProfileWolfA           bool `json:"decal_profile_wolf_a"`
	RWDGoalFx0010               bool `json:"rwd_goal_fx_0010"`
	RWDGoalFx0011               bool `json:"rwd_goal_fx_0011"`
	RWDPatternRocketA           bool `json:"rwd_pattern_rocket_a"`
	EmoteLightbulbA             bool `json:"emote_lightbulb_a"`
	RWDTitleTitleC              bool `json:"rwd_title_title_c"`
	RWDTitleTitleA              bool `json:"rwd_title_title_a"`
	PatternChevronA             bool `json:"pattern_chevron_a"`
	TintOrangeHDefault          bool `json:"tint_orange_h_default"`
	DecalCombatNovaA            bool `json:"decal_combat_nova_a"`
	DecalCombatLionA            bool `json:"decal_combat_lion_a"`
	EmoteQuestionMarkA          bool `json:"emote_question_mark_a"`
	RWDTagS1DSecondary          bool `json:"rwd_tag_s1_d_secondary"`
	TintNeutralHDefault         bool `json:"tint_neutral_h_default"`
	DecalCupcakeA               bool `json:"decal_cupcake_a"`
	DecalSkullA                 bool `json:"decal_skull_a"`
	EmoteFlyingHeartsA          bool `json:"emote_flying_hearts_a"`
	DecalCrownA                 bool `json:"decal_crown_a"`
	DecalCombatScratchA         bool `json:"decal_combat_scratch_a"`
	RWDMedalS1ArenaBronze       bool `json:"rwd_medal_s1_arena_bronze"`
	TintNeutralBDefault         bool `json:"tint_neutral_b_default"`
	EmoteStarSparklesA          bool `json:"emote_star_sparkles_a"`
	TintOrangeCDefault          bool `json:"tint_orange_c_default"`
	EmoteSmirkFaceA             bool `json:"emote_smirk_face_a"`
	RWDChassisMakoS1A           bool `json:"rwd_chassis_mako_s1_a"`
	RWDEmoteBatteryS1A          bool `json:"rwd_emote_battery_s1_a"`
	RWDDecalPepperA             bool `json:"rwd_decal_pepper_a"`
	RWDBannerS1Digi             bool `json:"rwd_banner_s1_digi"`
	RWDBracerMakoS1A            bool `json:"rwd_bracer_mako_s1_a"`
	RWDTagS1TSecondary          bool `json:"rwd_tag_s1_t_secondary"`
	RWDTintS1CDefault           bool `json:"rwd_tint_s1_c_default"`
	RWDTitleS1A                 bool `json:"rwd_title_s1_a"`
	RWDXPBoostIndividualS0101   bool `json:"rwd_xp_boost_individual_s01_01"`
	RWDBoosterMakoS1A           bool `json:"rwd_booster_mako_s1_a"`
	RWDEmoteCoffeeS1A           bool `json:"rwd_emote_coffee_s1_a"`
	RWDXPBoostGroupS0101        bool `json:"rwd_xp_boost_group_s01_01"`
	RWDDecalGgA                 bool `json:"rwd_decal_gg_a"`
	RWDMedalS1EchoPassBronze    bool `json:"rwd_medal_s1_echo_pass_bronze"`
	RWDCurrencyS0101            bool `json:"rwd_currency_s01_01"`
	RWDXPBoostIndividualS0102   bool `json:"rwd_xp_boost_individual_s01_02"`
	RWDBannerS1Flames           bool `json:"rwd_banner_s1_flames"`
	RWDPatternS1B               bool `json:"rwd_pattern_s1_b"`
	RWDXPBoostGroupS0102        bool `json:"rwd_xp_boost_group_s01_02"`
	RWDBoosterArcadeS1A         bool `json:"rwd_booster_arcade_s1_a"`
	RWDTitleS1B                 bool `json:"rwd_title_s1_b"`
	RWDXPBoostIndividualS0103   bool `json:"rwd_xp_boost_individual_s01_03"`
	RWDEmoteMeteorS1A           bool `json:"rwd_emote_meteor_s1_a"`
	RWDTagS1QSecondary          bool `json:"rwd_tag_s1_q_secondary"`
	RWDCurrencyS0102            bool `json:"rwd_currency_s01_02"`
	RWDXPBoostGroupS0103        bool `json:"rwd_xp_boost_group_s01_03"`
	RWDMedalS1EchoPassSilver    bool `json:"rwd_medal_s1_echo_pass_silver"`
	RWDXPBoostIndividualS0104   bool `json:"rwd_xp_boost_individual_s01_04"`
	RWDDecalCherryBlossomA      bool `json:"rwd_decal_cherry_blossom_a"`
	RWDBracerArcadeS1A          bool `json:"rwd_bracer_arcade_s1_a"`
	RWDBannerS1Trex             bool `json:"rwd_banner_s1_trex"`
	RWDXPBoostGroupS0104        bool `json:"rwd_xp_boost_group_s01_04"`
	RWDTintS1DDefault           bool `json:"rwd_tint_s1_d_default"`
	RWDPatternS1C               bool `json:"rwd_pattern_s1_c"`
	RWDCurrencyS0103            bool `json:"rwd_currency_s01_03"`
	RWDDecalRamenA              bool `json:"rwd_decal_ramen_a"`
	RWDBannerS1Tattered         bool `json:"rwd_banner_s1_tattered"`
	RWDXPBoostIndividualS0105   bool `json:"rwd_xp_boost_individual_s01_05"`
	RWDBracerArcadeVarS1A       bool `json:"rwd_bracer_arcade_var_s1_a"`
	RWDBoosterTrexS1A           bool `json:"rwd_booster_trex_s1_a"`
	RWDXPBoostGroupS0105        bool `json:"rwd_xp_boost_group_s01_05"`
	RWDTagS1GSecondary          bool `json:"rwd_tag_s1_g_secondary"`
	RWDPatternS1D               bool `json:"rwd_pattern_s1_d"`
	RWDCurrencyS0104            bool `json:"rwd_currency_s01_04"`
	RWDBracerTrexS1A            bool `json:"rwd_bracer_trex_s1_a"`
	RWDBannerS1Wings            bool `json:"rwd_banner_s1_wings"`
	RWDTitleS1C                 bool `json:"rwd_title_s1_c"`
	RWDMedalS1EchoPassGold      bool `json:"rwd_medal_s1_echo_pass_gold"`
	RWDBoosterArcadeVarS1A      bool `json:"rwd_booster_arcade_var_s1_a"`
	RWDChassisTrexS1A           bool `json:"rwd_chassis_trex_s1_a"`
	RWDChassisAutomatonS2A      bool `json:"rwd_chassis_automaton_s2_a"`
	EmoteShiftyEyesS2A          bool `json:"emote_shifty_eyes_s2_a"`
	RWDTintS1ADefault           bool `json:"rwd_tint_s1_a_default"`
	RWDBannerS2Deco             bool `json:"rwd_banner_s2_deco"`
	RWDBracerAutomatonS2A       bool `json:"rwd_bracer_automaton_s2_a"`
	RWDDecalScarabS2A           bool `json:"rwd_decal_scarab_s2_a"`
	RWDPatternS1A               bool `json:"rwd_pattern_s1_a"`
	RWDTitleS2A                 bool `json:"rwd_title_s2_a"`
	RWDXPBoostIndividualS0201   bool `json:"rwd_xp_boost_individual_s02_01"`
	RWDBoosterAutomatonS2A      bool `json:"rwd_booster_automaton_s2_a"`
	EmoteSoundWaveS2A           bool `json:"emote_sound_wave_s2_a"`
	RWDXPBoostGroupS0201        bool `json:"rwd_xp_boost_group_s02_01"`
	RWDTintS2CDefault           bool `json:"rwd_tint_s2_c_default"`
	RWDMedalS2EchoPassBronze    bool `json:"rwd_medal_s2_echo_pass_bronze"`
	RWDCurrencyS0201            bool `json:"rwd_currency_s02_01"`
	RWDXPBoostIndividualS0202   bool `json:"rwd_xp_boost_individual_s02_02"`
	RWDBannerS2Gears            bool `json:"rwd_banner_s2_gears"`
	RWDPatternS2B               bool `json:"rwd_pattern_s2_b"`
	RWDXPBoostGroupS0202        bool `json:"rwd_xp_boost_group_s02_02"`
	RWDBracerLadybugS2A         bool `json:"rwd_bracer_ladybug_s2_a"`
	RWDTitleS2B                 bool `json:"rwd_title_s2_b"`
	RWDXPBoostIndividualS0203   bool `json:"rwd_xp_boost_individual_s02_03"`
	RWDTintS2BDefault           bool `json:"rwd_tint_s2_b_default"`
	RWDTagS2BSecondary          bool `json:"rwd_tag_s2_b_secondary"`
	RWDCurrencyS0202            bool `json:"rwd_currency_s02_02"`
	RWDXPBoostGroupS0203        bool `json:"rwd_xp_boost_group_s02_03"`
	RWDBannerS2Pyramids         bool `json:"rwd_banner_s2_pyramids"`
	RWDXPBoostIndividualS0204   bool `json:"rwd_xp_boost_individual_s02_04"`
	EmoteUwuS2A                 bool `json:"emote_uwu_s2_a"`
	RWDMedalS2EchoPassSilver    bool `json:"rwd_medal_s2_echo_pass_silver"`
	RWDBoosterLadybugS2A        bool `json:"rwd_booster_ladybug_s2_a"`
	RWDXPBoostGroupS0204        bool `json:"rwd_xp_boost_group_s02_04"`
	RWDTagS2GSecondary          bool `json:"rwd_tag_s2_g_secondary"`
	RWDDecalGearsS2A            bool `json:"rwd_decal_gears_s2_a"`
	RWDCurrencyS0203            bool `json:"rwd_currency_s02_03"`
	RWDPatternS2C               bool `json:"rwd_pattern_s2_c"`
	RWDBannerS2Ladybug          bool `json:"rwd_banner_s2_ladybug"`
	RWDXPBoostIndividualS0205   bool `json:"rwd_xp_boost_individual_s02_05"`
	RWDBracerBeeS2A             bool `json:"rwd_bracer_bee_s2_a"`
	RWDBoosterAnubisS2A         bool `json:"rwd_booster_anubis_s2_a"`
	RWDXPBoostGroupS0205        bool `json:"rwd_xp_boost_group_s02_05"`
	RWDTitleS2C                 bool `json:"rwd_title_s2_c"`
	RWDTagS2HSecondary          bool `json:"rwd_tag_s2_h_secondary"`
	RWDCurrencyS0204            bool `json:"rwd_currency_s02_04"`
	RWDBracerAnubisS2A          bool `json:"rwd_bracer_anubis_s2_a"`
	RWDDecalAxolotlS2A          bool `json:"rwd_decal_axolotl_s2_a"`
	RWDMedalS2EchoPassGold      bool `json:"rwd_medal_s2_echo_pass_gold"`
	RWDBannerS2Squares          bool `json:"rwd_banner_s2_squares"`
	RWDBoosterBeeS2A            bool `json:"rwd_booster_bee_s2_a"`
	RWDChassisAnubisS2A         bool `json:"rwd_chassis_anubis_s2_a"`
	RWDChassisSpartanA          bool `json:"rwd_chassis_spartan_a"`
	RWDTintS3TintA              bool `json:"rwd_tint_s3_tint_a"`
	RWDEmoteLightningA          bool `json:"rwd_emote_lightning_a"`
	RWDBannerTrianglesA         bool `json:"rwd_banner_triangles_a"`
	RWDBracerSpartanA           bool `json:"rwd_bracer_spartan_a"`
	RWDPatternCircuitBoardA     bool `json:"rwd_pattern_circuit_board_a"`
	RWDTitleGuardianA           bool `json:"rwd_title_guardian_a"`
	RWDTagDiamondsA             bool `json:"rwd_tag_diamonds_a"`
	RWDXPBoostIndividualS0301   bool `json:"rwd_xp_boost_individual_s03_01"`
	RWDBoosterSpartanA          bool `json:"rwd_booster_spartan_a"`
	RWDDecalNarwhalA            bool `json:"rwd_decal_narwhal_a"`
	RWDXPBoostGroupS0301        bool `json:"rwd_xp_boost_group_s03_01"`
	RWDTintS3TintB              bool `json:"rwd_tint_s3_tint_b"`
	RWDMedalS3EchoPassBronzeA   bool `json:"rwd_medal_s3_echo_pass_bronze_a"`
	RWDCurrencyS0301            bool `json:"rwd_currency_s03_01"`
	RWDXPBoostIndividualS0302   bool `json:"rwd_xp_boost_individual_s03_02"`
	RWDBracerLazurliteA         bool `json:"rwd_bracer_lazurlite_a"`
	RWDEmoteBattleCryA          bool `json:"rwd_emote_battle_cry_a"`
	RWDXPBoostGroupS0302        bool `json:"rwd_xp_boost_group_s03_02"`
	RWDBannerSpartanShieldA     bool `json:"rwd_banner_spartan_shield_a"`
	RWDPatternSpearShieldA      bool `json:"rwd_pattern_spear_shield_a"`
	RWDXPBoostIndividualS0303   bool `json:"rwd_xp_boost_individual_s03_03"`
	RWDBoosterLazurliteA        bool `json:"rwd_booster_lazurlite_a"`
	RWDTintS3TintC              bool `json:"rwd_tint_s3_tint_c"`
	RWDCurrencyS0302            bool `json:"rwd_currency_s03_02"`
	RWDTitleShieldBearerA       bool `json:"rwd_title_shield_bearer_a"`
	RWDXPBoostGroupS0303        bool `json:"rwd_xp_boost_group_s03_03"`
	RWDDecalSpartanA            bool `json:"rwd_decal_spartan_a"`
	RWDBracerAurumA             bool `json:"rwd_bracer_aurum_a"`
	RWDTagSpearA                bool `json:"rwd_tag_spear_a"`
	RWDXPBoostIndividualS0304   bool `json:"rwd_xp_boost_individual_s03_04"`
	RWDMedalS3EchoPassSilverA   bool `json:"rwd_medal_s3_echo_pass_silver_a"`
	RWDXPBoostGroupS0304        bool `json:"rwd_xp_boost_group_s03_04"`
	RWDBoosterAurumA            bool `json:"rwd_booster_aurum_a"`
	RWDEmoteSamuraiMaskA        bool `json:"rwd_emote_samurai_mask_a"`
	RWDTintS3TintD              bool `json:"rwd_tint_s3_tint_d"`
	RWDBannerSashimonoA         bool `json:"rwd_banner_sashimono_a"`
	RWDCurrencyS0303            bool `json:"rwd_currency_s03_03"`
	RWDPatternSeigaihaA         bool `json:"rwd_pattern_seigaiha_a"`
	RWDXPBoostIndividualS0305   bool `json:"rwd_xp_boost_individual_s03_05"`
	RWDTitleRoninA              bool `json:"rwd_title_ronin_a"`
	RWDBracerSamuraiA           bool `json:"rwd_bracer_samurai_a"`
	RWDXPBoostGroupS0305        bool `json:"rwd_xp_boost_group_s03_05"`
	RWDDecalOniA                bool `json:"rwd_decal_oni_a"`
	RWDBoosterSamuraiA          bool `json:"rwd_booster_samurai_a"`
	RWDTintS3TintE              bool `json:"rwd_tint_s3_tint_e"`
	RWDMedalS3EchoPassGoldA     bool `json:"rwd_medal_s3_echo_pass_gold_a"`
	RWDTagToriA                 bool `json:"rwd_tag_tori_a"`
	RWDCurrencyS0304            bool `json:"rwd_currency_s03_04"`
	RWDChassisSamuraiA          bool `json:"rwd_chassis_samurai_a"`
	RWDChassisStreetwearA       bool `json:"rwd_chassis_streetwear_a"`
	RWDBanner0000               bool `json:"rwd_banner_0000"`
	RWDEmote0000                bool `json:"rwd_emote_0000"`
	RWDTag0000                  bool `json:"rwd_tag_0000"`
	RWDBracerStreetwearA        bool `json:"rwd_bracer_streetwear_a"`
	RWDTint0000                 bool `json:"rwd_tint_0000"`
	RWDTitle0000                bool `json:"rwd_title_0000"`
	RWDPattern0000              bool `json:"rwd_pattern_0000"`
	RWDXPBoostIndividualS0401   bool `json:"rwd_xp_boost_individual_s04_01"`
	RWDBoosterStreetwearA       bool `json:"rwd_booster_streetwear_a"`
	RWDDecal0000                bool `json:"rwd_decal_0000"`
	RWDXPBoostGroupS0401        bool `json:"rwd_xp_boost_group_s04_01"`
	RWDBanner0001               bool `json:"rwd_banner_0001"`
	RWDMedal0000                bool `json:"rwd_medal_0000"`
	RWDCurrencyS0401            bool `json:"rwd_currency_s04_01"`
	RWDXPBoostIndividualS0402   bool `json:"rwd_xp_boost_individual_s04_02"`
	RWDBracerRoverA             bool `json:"rwd_bracer_rover_a"`
	RWDTag0001                  bool `json:"rwd_tag_0001"`
	RWDXPBoostGroupS0402        bool `json:"rwd_xp_boost_group_s04_02"`
	RWDEmote0001                bool `json:"rwd_emote_0001"`
	RWDTint0001                 bool `json:"rwd_tint_0001"`
	RWDXPBoostIndividualS0403   bool `json:"rwd_xp_boost_individual_s04_03"`
	RWDBoosterRoverA            bool `json:"rwd_booster_rover_a"`
	RWDPattern0001              bool `json:"rwd_pattern_0001"`
	RWDCurrencyS0402            bool `json:"rwd_currency_s04_02"`
	RWDTitle0001                bool `json:"rwd_title_0001"`
	RWDXPBoostGroupS0403        bool `json:"rwd_xp_boost_group_s04_03"`
	RWDBanner0002               bool `json:"rwd_banner_0002"`
	RWDBracerRoverADeco         bool `json:"rwd_bracer_rover_a_deco"`
	RWDDecal0001                bool `json:"rwd_decal_0001"`
	RWDXPBoostIndividualS0404   bool `json:"rwd_xp_boost_individual_s04_04"`
	RWDMedal0001                bool `json:"rwd_medal_0001"`
	RWDXPBoostGroupS0404        bool `json:"rwd_xp_boost_group_s04_04"`
	RWDBoosterRoverADeco        bool `json:"rwd_booster_rover_a_deco"`
	RWDTagS2C                   bool `json:"rwd_tag_s2_c"`
	RWDPattern0002              bool `json:"rwd_pattern_0002"`
	RWDEmote0002                bool `json:"rwd_emote_0002"`
	RWDCurrencyS0403            bool `json:"rwd_currency_s04_03"`
	RWDTint0002                 bool `json:"rwd_tint_0002"`
	RWDXPBoostIndividualS0405   bool `json:"rwd_xp_boost_individual_s04_05"`
	RWDTitle0002                bool `json:"rwd_title_0002"`
	RWDBracerFunkA              bool `json:"rwd_bracer_funk_a"`
	RWDXPBoostGroupS0405        bool `json:"rwd_xp_boost_group_s04_05"`
	RWDBanner0003               bool `json:"rwd_banner_0003"`
	RWDBoosterFunkA             bool `json:"rwd_booster_funk_a"`
	RWDDecal0002                bool `json:"rwd_decal_0002"`
	RWDMedal0002                bool `json:"rwd_medal_0002"`
	RWDTag0003                  bool `json:"rwd_tag_0003"`
	RWDCurrencyS0404            bool `json:"rwd_currency_s04_04"`
	RWDChassisFunkA             bool `json:"rwd_chassis_funk_a"`
	RWDCurrencyS0501            bool `json:"rwd_currency_s05_01"`
	RWDCurrencyS0502            bool `json:"rwd_currency_s05_02"`
	RWDCurrencyS0503            bool `json:"rwd_currency_s05_03"`
	RWDCurrencyS0504            bool `json:"rwd_currency_s05_04"`
	RWDXPBoostIndividualS0501   bool `json:"rwd_xp_boost_individual_s05_01"`
	RWDXPBoostIndividualS0502   bool `json:"rwd_xp_boost_individual_s05_02"`
	RWDXPBoostIndividualS0503   bool `json:"rwd_xp_boost_individual_s05_03"`
	RWDXPBoostIndividualS0504   bool `json:"rwd_xp_boost_individual_s05_04"`
	RWDXPBoostIndividualS0505   bool `json:"rwd_xp_boost_individual_s05_05"`
	RWDXPBoostGroupS0501        bool `json:"rwd_xp_boost_group_s05_01"`
	RWDXPBoostGroupS0502        bool `json:"rwd_xp_boost_group_s05_02"`
	RWDXPBoostGroupS0503        bool `json:"rwd_xp_boost_group_s05_03"`
	RWDXPBoostGroupS0504        bool `json:"rwd_xp_boost_group_s05_04"`
	RWDXPBoostGroupS0505        bool `json:"rwd_xp_boost_group_s05_05"`
	RWDChassisJunkyardA         bool `json:"rwd_chassis_junkyard_a"`
	RWDBanner0008               bool `json:"rwd_banner_0008"`
	RWDEmote0005                bool `json:"rwd_emote_0005"`
	RWDTag0012                  bool `json:"rwd_tag_0012"`
	RWDBracerJunkyardA          bool `json:"rwd_bracer_junkyard_a"`
	RWDTint0007                 bool `json:"rwd_tint_0007"`
	RWDTitle0004                bool `json:"rwd_title_0004"`
	RWDPattern0005              bool `json:"rwd_pattern_0005"`
	RWDBoosterJunkyardA         bool `json:"rwd_booster_junkyard_a"`
	RWDDecal0005                bool `json:"rwd_decal_0005"`
	RWDBanner0009               bool `json:"rwd_banner_0009"`
	RWDMedal0003                bool `json:"rwd_medal_0003"`
	RWDBracerNuclearA           bool `json:"rwd_bracer_nuclear_a"`
	RWDTag0013                  bool `json:"rwd_tag_0013"`
	RWDEmote0006                bool `json:"rwd_emote_0006"`
	RWDTint0008                 bool `json:"rwd_tint_0008"`
	RWDBoosterNuclearA          bool `json:"rwd_booster_nuclear_a"`
	RWDPattern0006              bool `json:"rwd_pattern_0006"`
	RWDTitle0005                bool `json:"rwd_title_0005"`
	RWDBanner0010               bool `json:"rwd_banner_0010"`
	RWDBracerNuclearAHydro      bool `json:"rwd_bracer_nuclear_a_hydro"`
	RWDDecal0006                bool `json:"rwd_decal_0006"`
	RWDMedal0004                bool `json:"rwd_medal_0004"`
	RWDBoosterNuclearAHydro     bool `json:"rwd_booster_nuclear_a_hydro"`
	RWDTag0014                  bool `json:"rwd_tag_0014"`
	RWDEmote0007                bool `json:"rwd_emote_0007"`
	RWDTint0009                 bool `json:"rwd_tint_0009"`
	RWDTitle0006                bool `json:"rwd_title_0006"`
	RWDBracerWastelandA         bool `json:"rwd_bracer_wasteland_a"`
	RWDBanner0011               bool `json:"rwd_banner_0011"`
	RWDBoosterWastelandA        bool `json:"rwd_booster_wasteland_a"`
	RWDDecal0007                bool `json:"rwd_decal_0007"`
	RWDMedal0005                bool `json:"rwd_medal_0005"`
	RWDTag0015                  bool `json:"rwd_tag_0015"`
	RWDChassisWastelandA        bool `json:"rwd_chassis_wasteland_a"`
	RWDPatternS2A               bool `json:"rwd_pattern_s2_a"`
	RWDCurrencyS0601            bool `json:"rwd_currency_s06_01"`
	RWDCurrencyS0602            bool `json:"rwd_currency_s06_02"`
	RWDCurrencyS0603            bool `json:"rwd_currency_s06_03"`
	RWDCurrencyS0604            bool `json:"rwd_currency_s06_04"`
	RWDXPBoostIndividualS0601   bool `json:"rwd_xp_boost_individual_s06_01"`
	RWDXPBoostIndividualS0602   bool `json:"rwd_xp_boost_individual_s06_02"`
	RWDXPBoostIndividualS0603   bool `json:"rwd_xp_boost_individual_s06_03"`
	RWDXPBoostIndividualS0604   bool `json:"rwd_xp_boost_individual_s06_04"`
	RWDXPBoostIndividualS0605   bool `json:"rwd_xp_boost_individual_s06_05"`
	RWDXPBoostGroupS0601        bool `json:"rwd_xp_boost_group_s06_01"`
	RWDXPBoostGroupS0602        bool `json:"rwd_xp_boost_group_s06_02"`
	RWDXPBoostGroupS0603        bool `json:"rwd_xp_boost_group_s06_03"`
	RWDXPBoostGroupS0604        bool `json:"rwd_xp_boost_group_s06_04"`
	RWDXPBoostGroupS0605        bool `json:"rwd_xp_boost_group_s06_05"`
	RWDBanner0015               bool `json:"rwd_banner_0015"`
	RWDBanner0014               bool `json:"rwd_banner_0014"`
	RWDBanner0016               bool `json:"rwd_banner_0016"`
	RWDTag0018                  bool `json:"rwd_tag_0018"`
	RWDTag0019                  bool `json:"rwd_tag_0019"`
	RWDTag0020                  bool `json:"rwd_tag_0020"`
	RWDBanner0017               bool `json:"rwd_banner_0017"`
	RWDTag0021                  bool `json:"rwd_tag_0021"`
	RWDDecal0010                bool `json:"rwd_decal_0010"`
	RWDDecal0011                bool `json:"rwd_decal_0011"`
	RWDDecal0012                bool `json:"rwd_decal_0012"`
	RWDMedal0009                bool `json:"rwd_medal_0009"`
	RWDMedal0010                bool `json:"rwd_medal_0010"`
	RWDMedal0011                bool `json:"rwd_medal_0011"`
	RWDTint0013                 bool `json:"rwd_tint_0013"`
	RWDTint0014                 bool `json:"rwd_tint_0014"`
	RWDTint0015                 bool `json:"rwd_tint_0015"`
	RWDPattern0009              bool `json:"rwd_pattern_0009"`
	RWDPattern0010              bool `json:"rwd_pattern_0010"`
	RWDEmote0010                bool `json:"rwd_emote_0010"`
	RWDEmote0011                bool `json:"rwd_emote_0011"`
	RWDEmote0012                bool `json:"rwd_emote_0012"`
	RWDTitle0007                bool `json:"rwd_title_0007"`
	RWDTitle0008                bool `json:"rwd_title_0008"`
	RWDTitle0009                bool `json:"rwd_title_0009"`
	RWDBracerSharkA             bool `json:"rwd_bracer_shark_a"`
	RWDBoosterSharkA            bool `json:"rwd_booster_shark_a"`
	RWDChassisSharkA            bool `json:"rwd_chassis_shark_a"`
	RWDPattern0008              bool `json:"rwd_pattern_0008"`
	RWDBoosterCovenantA         bool `json:"rwd_booster_covenant_a"`
	RWDBracerCovenantA          bool `json:"rwd_bracer_covenant_a"`
	RWDBracerCovenantAFlame     bool `json:"rwd_bracer_covenant_a_flame"`
	RWDBoosterCovenantAFlame    bool `json:"rwd_booster_covenant_a_flame"`
	RWDChassisScubaA            bool `json:"rwd_chassis_scuba_a"`
	RWDBracerScubaA             bool `json:"rwd_bracer_scuba_a"`
	RWDBoosterScubaA            bool `json:"rwd_booster_scuba_a"`
	RWDCurrencyS0701            bool `json:"rwd_currency_s07_01"`
	RWDCurrencyS0702            bool `json:"rwd_currency_s07_02"`
	RWDCurrencyS0703            bool `json:"rwd_currency_s07_03"`
	RWDCurrencyS0704            bool `json:"rwd_currency_s07_04"`
	RWDXPBoostIndividualS0701   bool `json:"rwd_xp_boost_individual_s07_01"`
	RWDXPBoostIndividualS0702   bool `json:"rwd_xp_boost_individual_s07_02"`
	RWDXPBoostIndividualS0703   bool `json:"rwd_xp_boost_individual_s07_03"`
	RWDXPBoostIndividualS0704   bool `json:"rwd_xp_boost_individual_s07_04"`
	RWDXPBoostIndividualS0705   bool `json:"rwd_xp_boost_individual_s07_05"`
	RWDXPBoostGroupS0701        bool `json:"rwd_xp_boost_group_s07_01"`
	RWDXPBoostGroupS0702        bool `json:"rwd_xp_boost_group_s07_02"`
	RWDXPBoostGroupS0703        bool `json:"rwd_xp_boost_group_s07_03"`
	RWDXPBoostGroupS0704        bool `json:"rwd_xp_boost_group_s07_04"`
	RWDXPBoostGroupS0705        bool `json:"rwd_xp_boost_group_s07_05"`
	RWDBanner0023               bool `json:"rwd_banner_0023"`
	RWDTitle0010                bool `json:"rwd_title_0010"`
	RWDDecal0015                bool `json:"rwd_decal_0015"`
	RWDPattern0016              bool `json:"rwd_pattern_0016"`
	RWDMedal0012                bool `json:"rwd_medal_0012"`
	RWDBanner0024               bool `json:"rwd_banner_0024"`
	RWDMedal0015                bool `json:"rwd_medal_0015"`
	RWDTint0020                 bool `json:"rwd_tint_0020"`
	RWDBoosterFumeA             bool `json:"rwd_booster_fume_a"`
	RWDTitle0011                bool `json:"rwd_title_0011"`
	RWDTint0021                 bool `json:"rwd_tint_0021"`
	RWDBracerFumeA              bool `json:"rwd_bracer_fume_a"`
	RWDMedal0016                bool `json:"rwd_medal_0016"`
	RWDEmote0016                bool `json:"rwd_emote_0016"`
	RWDMedal0014                bool `json:"rwd_medal_0014"`
	RWDDecal0016                bool `json:"rwd_decal_0016"`
	RWDTag0028                  bool `json:"rwd_tag_0028"`
	RWDTag0029                  bool `json:"rwd_tag_0029"`
	RWDPattern0017              bool `json:"rwd_pattern_0017"`
	RWDTint0022                 bool `json:"rwd_tint_0022"`
	RWDTag0030                  bool `json:"rwd_tag_0030"`
	RWDBracerNobleA             bool `json:"rwd_bracer_noble_a"`
	RWDMedal0013                bool `json:"rwd_medal_0013"`
	RWDTitle0012                bool `json:"rwd_title_0012"`
	RWDBoosterNobleA            bool `json:"rwd_booster_noble_a"`
	RWDBanner0025               bool `json:"rwd_banner_0025"`
	RWDEmote0017                bool `json:"rwd_emote_0017"`
	RWDChassisNobleA            bool `json:"rwd_chassis_noble_a"`
	RWDTitle0013                bool `json:"rwd_title_0013"`
	RWDTint0023                 bool `json:"rwd_tint_0023"`
	RWDTag0031                  bool `json:"rwd_tag_0031"`
	RWDDecal0018                bool `json:"rwd_decal_0018"`
	RWDPattern0018              bool `json:"rwd_pattern_0018"`
	RWDBanner0021               bool `json:"rwd_banner_0021"`
	RWDEmote0019                bool `json:"rwd_emote_0019"`
	RWDTitle0014                bool `json:"rwd_title_0014"`
	RWDEmote0018                bool `json:"rwd_emote_0018"`
	RWDTint0024                 bool `json:"rwd_tint_0024"`
	RWDDecal0019                bool `json:"rwd_decal_0019"`
	RWDTag0027                  bool `json:"rwd_tag_0027"`
	RWDDecal0020                bool `json:"rwd_decal_0020"`
	RWDPattern0019              bool `json:"rwd_pattern_0019"`
	RWDBanner0022               bool `json:"rwd_banner_0022"`
	RWDEmote0020                bool `json:"rwd_emote_0020"`
	RWDTitle0015                bool `json:"rwd_title_0015"`
	RWDTint0025                 bool `json:"rwd_tint_0025"`
	RWDChassisPlagueknightA     bool `json:"rwd_chassis_plagueknight_a"`
	RWDBracerPlagueknightA      bool `json:"rwd_bracer_plagueknight_a"`
	RWDBoosterPlagueknightA     bool `json:"rwd_booster_plagueknight_a"`
	RWDPip0017                  bool `json:"rwd_pip_0017"`
	RWDPip0018                  bool `json:"rwd_pip_0018"`
	RWDPip0019                  bool `json:"rwd_pip_0019"`
	RWDPip0020                  bool `json:"rwd_pip_0020"`
	RWDPip0021                  bool `json:"rwd_pip_0021"`
	RWDPip0022                  bool `json:"rwd_pip_0022"`
	RWDEmissive0023             bool `json:"rwd_emissive_0023"`
	RWDEmissive0024             bool `json:"rwd_emissive_0024"`
	RWDEmissive0026             bool `json:"rwd_emissive_0026"`
	RWDEmissive0028             bool `json:"rwd_emissive_0028"`
	RWDEmissive0029             bool `json:"rwd_emissive_0029"`
	RWDBracerFumeADaydream      bool `json:"rwd_bracer_fume_a_daydream"`
	RWDBoosterFumeADaydream     bool `json:"rwd_booster_fume_a_daydream"`
	RWDGoalFx0012               bool `json:"rwd_goal_fx_0012"`
	RWDGoalFx0004               bool `json:"rwd_goal_fx_0004"`
	RWDGoalFx0013               bool `json:"rwd_goal_fx_0013"`
	RWDGoalFx0015               bool `json:"rwd_goal_fx_0015"`
	RWDGoalFx0003               bool `json:"rwd_goal_fx_0003"`
	RWDDecal0017                bool `json:"rwd_decal_0017"`
	RWDEmissive0016             bool `json:"rwd_emissive_0016"`
	DecalKronosA                bool `json:"decal_kronos_a"`
	DecalOneYearA               bool `json:"decal_one_year_a"`
	EmoteOneA                   bool `json:"emote_one_a"`
	TintNeutralXmasADefault     bool `json:"tint_neutral_xmas_a_default"`
	TintNeutralXmasBDefault     bool `json:"tint_neutral_xmas_b_default"`
	TintNeutralXmasCDefault     bool `json:"tint_neutral_xmas_c_default"`
	TintNeutralXmasDDefault     bool `json:"tint_neutral_xmas_d_default"`
	TintNeutralXmasEDefault     bool `json:"tint_neutral_xmas_e_default"`
	PatternXmasLightsA          bool `json:"pattern_xmas_lights_a"`
	PatternXmasSnowflakesA      bool `json:"pattern_xmas_snowflakes_a"`
	PatternXmasMistletoeA       bool `json:"pattern_xmas_mistletoe_a"`
	PatternXmasFlourishA        bool `json:"pattern_xmas_flourish_a"`
	PatternXmasKnitA            bool `json:"pattern_xmas_knit_a"`
	PatternXmasKnitFlowersA     bool `json:"pattern_xmas_knit_flowers_a"`
	DecalPresentA               bool `json:"decal_present_a"`
	DecalBowA                   bool `json:"decal_bow_a"`
	DecalGingerbreadA           bool `json:"decal_gingerbread_a"`
	DecalPenguinA               bool `json:"decal_penguin_a"`
	DecalSnowmanA               bool `json:"decal_snowman_a"`
	DecalWreathA                bool `json:"decal_wreath_a"`
	DecalSnowflakeA             bool `json:"decal_snowflake_a"`
	DecalReindeerA              bool `json:"decal_reindeer_a"`
	EmoteSnowmanA               bool `json:"emote_snowman_a"`
	EmoteFireA                  bool `json:"emote_fire_a"`
	EmotePresentA               bool `json:"emote_present_a"`
	EmoteGingerbreadManA        bool `json:"emote_gingerbread_man_a"`
	TintNeutralSpookyADefault   bool `json:"tint_neutral_spooky_a_default"`
	TintNeutralSpookyBDefault   bool `json:"tint_neutral_spooky_b_default"`
	TintNeutralSpookyCDefault   bool `json:"tint_neutral_spooky_c_default"`
	TintNeutralSpookyDDefault   bool `json:"tint_neutral_spooky_d_default"`
	TintNeutralSpookyEDefault   bool `json:"tint_neutral_spooky_e_default"`
	PatternSpookyStitchesA      bool `json:"pattern_spooky_stitches_a"`
	PatternSpookyCobwebA        bool `json:"pattern_spooky_cobweb_a"`
	PatternSpookyBandagesA      bool `json:"pattern_spooky_bandages_a"`
	PatternSpookyPumpkinsA      bool `json:"pattern_spooky_pumpkins_a"`
	PatternSpookyBatsA          bool `json:"pattern_spooky_bats_a"`
	PatternSpookySkullsA        bool `json:"pattern_spooky_skulls_a"`
	DecalHalloweenBatA          bool `json:"decal_halloween_bat_a"`
	DecalHalloweenCatA          bool `json:"decal_halloween_cat_a"`
	DecalFangsA                 bool `json:"decal_fangs_a"`
	DecalHalloweenGhostA        bool `json:"decal_halloween_ghost_a"`
	DecalHalloweenPumpkinA      bool `json:"decal_halloween_pumpkin_a"`
	DecalHalloweenSkullA        bool `json:"decal_halloween_skull_a"`
	DecalHalloweenZombieA       bool `json:"decal_halloween_zombie_a"`
	DecalHalloweenScytheA       bool `json:"decal_halloween_scythe_a"`
	EmotePumpkinFaceA           bool `json:"emote_pumpkin_face_a"`
	EmoteScaredA                bool `json:"emote_scared_a"`
	EmoteRIPA                   bool `json:"emote_rip_a"`
	EmoteBatsA                  bool `json:"emote_bats_a"`
	RWDBannerLoneEcho2_A        bool `json:"rwd_banner_lone_echo_2_a"`
	RWDTagLoneEcho2_A           bool `json:"rwd_tag_lone_echo_2_a"`
	RWDDecalLoneEcho2_A         bool `json:"rwd_decal_lone_echo_2_a"`
	RWDMedalLoneEcho2_A         bool `json:"rwd_medal_lone_echo_2_a"`
	RWDBoosterHerosuitA         bool `json:"rwd_booster_herosuit_a"`
	RWDChassisHerosuitA         bool `json:"rwd_chassis_herosuit_a"`
	TintNeutralSummerADefault   bool `json:"tint_neutral_summer_a_default"`
	PatternSummerHawaiianA      bool `json:"pattern_summer_hawaiian_a"`
	DecalSummerPirateA          bool `json:"decal_summer_pirate_a"`
	DecalSummerSharkA           bool `json:"decal_summer_shark_a"`
	DecalSummerWhaleA           bool `json:"decal_summer_whale_a"`
	DecalSummerSubmarineA       bool `json:"decal_summer_submarine_a"`
	DecalHalloweenCauldronA     bool `json:"decal_halloween_cauldron_a"`
	DecalAnniversaryCupcakeA    bool `json:"decal_anniversary_cupcake_a"`
	DecalCombatAnniversaryA     bool `json:"decal_combat_anniversary_a"`
	DecalSantaCubesatA          bool `json:"decal_santa_cubesat_a"`
	DecalQuestLaunchA           bool `json:"decal_quest_launch_a"`
	EmoteDancingOctopusA        bool `json:"emote_dancing_octopus_a"`
	EmoteSpiderA                bool `json:"emote_spider_a"`
	EmoteCombatAnniversaryA     bool `json:"emote_combat_anniversary_a"`
	EmoteSnowGlobeA             bool `json:"emote_snow_globe_a"`
	EmoteDing                   bool `json:"emote_ding"`
	RWDMedalS1QuestLaunch       bool `json:"rwd_medal_s1_quest_launch"`
	RWDBoosterAnubisAHorus      bool `json:"rwd_booster_anubis_a_horus"`
	RWDBracerAnubisAHorus       bool `json:"rwd_bracer_anubis_a_horus"`
	RWDBoosterSharkATropical    bool `json:"rwd_booster_shark_a_tropical"`
	RWDBracerSharkATropical     bool `json:"rwd_bracer_shark_a_tropical"`
	RWDChassisAnubisAHorus      bool `json:"rwd_chassis_anubis_a_horus"`
	RWDChassisSharkATropical    bool `json:"rwd_chassis_shark_a_tropical"`
	RWDChassisSpartanAHero      bool `json:"rwd_chassis_spartan_a_hero"`
	RWDBracerSpartanAHero       bool `json:"rwd_bracer_spartan_a_hero"`
	RWDBoosterSpartanAHero      bool `json:"rwd_booster_spartan_a_hero"`
	RWDBracerSnacktimeA         bool `json:"rwd_bracer_snacktime_a"`
	RWDBoosterSnacktimeA        bool `json:"rwd_booster_snacktime_a"`
	RWDBracerHeartbreakA        bool `json:"rwd_bracer_heartbreak_a"`
	RWDBoosterHeartbreakA       bool `json:"rwd_booster_heartbreak_a"`
	RWDBracerVroomA             bool `json:"rwd_bracer_vroom_a"`
	RWDBoosterVroomA            bool `json:"rwd_booster_vroom_a"`
	RWDChassisNinjaA            bool `json:"rwd_chassis_ninja_a"`
	RWDBracerNinjaA             bool `json:"rwd_bracer_ninja_a"`
	RWDBoosterNinjaA            bool `json:"rwd_booster_ninja_a"`
	RWDPattern0021              bool `json:"rwd_pattern_0021"`
	RWDEmissive0022             bool `json:"rwd_emissive_0022"`
	RWDPip0025                  bool `json:"rwd_pip_0025"`
	RWDEmote0022                bool `json:"rwd_emote_0022"`
	RWDTint0029                 bool `json:"rwd_tint_0029"`
	RWDGoalFx0007               bool `json:"rwd_goal_fx_0007"`
	RWDDecal0022                bool `json:"rwd_decal_0022"`
	RWDBanner0028               bool `json:"rwd_banner_0028"`
	RWDEmote0023                bool `json:"rwd_emote_0023"`
	RWDEmissive0032             bool `json:"rwd_emissive_0032"`
	RWDPip0024                  bool `json:"rwd_pip_0024"`
	RWDTag0039                  bool `json:"rwd_tag_0039"`
	TintNeutralMDefault         bool `json:"tint_neutral_m_default"`
	DecalOculusA                bool `json:"decal_oculus_a"`
	EmoteDealGlassesA           bool `json:"emote_deal_glasses_a"`
	RWDChassisBodyS10A          bool `json:"rwd_chassis_body_s10_a"`
	RWDBoosterS10               bool `json:"rwd_booster_s10"`
	RWDTitleTitleB              bool `json:"rwd_title_title_b"`
	RWDMedalS1CombatBronze      bool `json:"rwd_medal_s1_combat_bronze"`
	RWDMedalS1CombatSilver      bool `json:"rwd_medal_s1_combat_silver"`
	RWDMedalS1CombatGold        bool `json:"rwd_medal_s1_combat_gold"`
	RWDChassisSportyA           bool `json:"rwd_chassis_sporty_a"`
	RWDBanner0004               bool `json:"rwd_banner_0004"`
	RWDTag0004                  bool `json:"rwd_tag_0004"`
	RWDTint0003                 bool `json:"rwd_tint_0003"`
	RWDTitle0003                bool `json:"rwd_title_0003"`
	RWDCurrencyStarterPack01    bool `json:"rwd_currency_starter_pack_01"`
	RWDPattern0003              bool `json:"rwd_pattern_0003"`
	RWDBanner0007               bool `json:"rwd_banner_0007"`
	RWDTag0007                  bool `json:"rwd_tag_0007"`
	RWDBracerReptileA           bool `json:"rwd_bracer_reptile_a"`
	RWDEmote0003                bool `json:"rwd_emote_0003"`
	RWDTint0006                 bool `json:"rwd_tint_0006"`
	RWDBoosterReptileA          bool `json:"rwd_booster_reptile_a"`
	RWDDecal0003                bool `json:"rwd_decal_0003"`
	RWDBanner0005               bool `json:"rwd_banner_0005"`
	RWDBracerRetroA             bool `json:"rwd_bracer_retro_a"`
	RWDEmote0004                bool `json:"rwd_emote_0004"`
	RWDTag0005                  bool `json:"rwd_tag_0005"`
	RWDBoosterRetroA            bool `json:"rwd_booster_retro_a"`
	RWDTint0004                 bool `json:"rwd_tint_0004"`
	RWDPattern0004              bool `json:"rwd_pattern_0004"`
	RWDBracerAvianA             bool `json:"rwd_bracer_avian_a"`
	RWDTag0006                  bool `json:"rwd_tag_0006"`
	RWDDecal0004                bool `json:"rwd_decal_0004"`
	RWDBoosterAvianA            bool `json:"rwd_booster_avian_a"`
	RWDBanner0006               bool `json:"rwd_banner_0006"`
	RWDTint0005                 bool `json:"rwd_tint_0005"`
	RWDChassisFrostA            bool `json:"rwd_chassis_frost_a"`
	RWDBracerFrostA             bool `json:"rwd_bracer_frost_a"`
	RWDBoosterFrostA            bool `json:"rwd_booster_frost_a"`
	RWDBoosterSpeedformA        bool `json:"rwd_booster_speedform_a"`
	RWDTag0016                  bool `json:"rwd_tag_0016"`
	RWDEmote0008                bool `json:"rwd_emote_0008"`
	RWDBracerSpeedformA         bool `json:"rwd_bracer_speedform_a"`
	RWDTint0010                 bool `json:"rwd_tint_0010"`
	RWDBoosterMechA             bool `json:"rwd_booster_mech_a"`
	RWDEmote0009                bool `json:"rwd_emote_0009"`
	RWDBanner0012               bool `json:"rwd_banner_0012"`
	RWDBracerMechA              bool `json:"rwd_bracer_mech_a"`
	RWDTint0011                 bool `json:"rwd_tint_0011"`
	RWDDecal0008                bool `json:"rwd_decal_0008"`
	RWDBoosterOrganicA          bool `json:"rwd_booster_organic_a"`
	RWDBanner0013               bool `json:"rwd_banner_0013"`
	RWDTag0017                  bool `json:"rwd_tag_0017"`
	RWDBracerOrganicA           bool `json:"rwd_bracer_organic_a"`
	RWDTint0012                 bool `json:"rwd_tint_0012"`
	RWDDecal0009                bool `json:"rwd_decal_0009"`
	RWDChassisExoA              bool `json:"rwd_chassis_exo_a"`
	RWDBracerExoA               bool `json:"rwd_bracer_exo_a"`
	RWDBoosterExoA              bool `json:"rwd_booster_exo_a"`
	RWDPattern0007              bool `json:"rwd_pattern_0007"`
	RWDBanner0018               bool `json:"rwd_banner_0018"`
	RWDTag0022                  bool `json:"rwd_tag_0022"`
	RWDChassisWolfA             bool `json:"rwd_chassis_wolf_a"`
	RWDBracerWolfA              bool `json:"rwd_bracer_wolf_a"`
	RWDBoosterWolfA             bool `json:"rwd_booster_wolf_a"`
	RWDPattern0011              bool `json:"rwd_pattern_0011"`
	RWDBracerFragmentA          bool `json:"rwd_bracer_fragment_a"`
	RWDTint0016                 bool `json:"rwd_tint_0016"`
	RWDEmote0013                bool `json:"rwd_emote_0013"`
	RWDBoosterFragmentA         bool `json:"rwd_booster_fragment_a"`
	RWDTag0023                  bool `json:"rwd_tag_0023"`
	RWDDecal0013                bool `json:"rwd_decal_0013"`
	RWDBracerBaroqueA           bool `json:"rwd_bracer_baroque_a"`
	RWDEmote0014                bool `json:"rwd_emote_0014"`
	RWDBanner0019               bool `json:"rwd_banner_0019"`
	RWDBoosterBaroqueA          bool `json:"rwd_booster_baroque_a"`
	RWDTint0017                 bool `json:"rwd_tint_0017"`
	RWDPattern0012              bool `json:"rwd_pattern_0012"`
	RWDBoosterLavaA             bool `json:"rwd_booster_lava_a"`
	RWDBanner0020               bool `json:"rwd_banner_0020"`
	RWDTint0018                 bool `json:"rwd_tint_0018"`
	RWDBracerLavaA              bool `json:"rwd_bracer_lava_a"`
	RWDTag0024                  bool `json:"rwd_tag_0024"`
	RWDDecal0014                bool `json:"rwd_decal_0014"`
	RWDTag0026                  bool `json:"rwd_tag_0026"`
	RWDTag0025                  bool `json:"rwd_tag_0025"`
	RWDEmote0015                bool `json:"rwd_emote_0015"`
	RWDTagS1RSecondary          bool `json:"rwd_tag_s1_r_secondary"`
	RWDTag0033                  bool `json:"rwd_tag_0033"`
	RWDPattern0020              bool `json:"rwd_pattern_0020"`
	RWDEmote0021                bool `json:"rwd_emote_0021"`
	RWDTint0028                 bool `json:"rwd_tint_0028"`
	RWDBracerHalloweenA         bool `json:"rwd_bracer_halloween_a"`
	RWDBoosterHalloweenA        bool `json:"rwd_booster_halloween_a"`
	RWDBracerFlamingoA          bool `json:"rwd_bracer_flamingo_a"`
	RWDBoosterFlamingoA         bool `json:"rwd_booster_flamingo_a"`
	RWDBracerPaladinA           bool `json:"rwd_bracer_paladin_a"`
	RWDBoosterPaladinA          bool `json:"rwd_booster_paladin_a"`
	RWDChassisOvergrownA        bool `json:"rwd_chassis_overgrown_a"`
	RWDBoosterOvergrownA        bool `json:"rwd_booster_overgrown_a"`
	RWDBracerOvergrownA         bool `json:"rwd_bracer_overgrown_a"`
	RWDTintS1BDefault           bool `json:"rwd_tint_s1_b_default"`
	RWDEmote0024                bool `json:"rwd_emote_0024"`
	RWDGoalFx0014               bool `json:"rwd_goal_fx_0014"`
	RWDGoalFx0001               bool `json:"rwd_goal_fx_0001"`
	RWDPip0023                  bool `json:"rwd_pip_0023"`
	RWDEmissive0030             bool `json:"rwd_emissive_0030"`
	RWDChassisSamuraiAOni       bool `json:"rwd_chassis_samurai_a_oni"`
	RWDBracerSamuraiAOni        bool `json:"rwd_bracer_samurai_a_oni"`
	RWDBoosterSamuraiAOni       bool `json:"rwd_booster_samurai_a_oni"`
	RWDPip0004                  bool `json:"rwd_pip_0004"`
	RWDEmissive0021             bool `json:"rwd_emissive_0021"`
	RWDEmissive0039             bool `json:"rwd_emissive_0039"`
	RWDBanner0026               bool `json:"rwd_banner_0026"`
	RWDTag0038                  bool `json:"rwd_tag_0038"`
	RWDBanner0027               bool `json:"rwd_banner_0027"`
	RWDPip0003                  bool `json:"rwd_pip_0003"`
	RWDBracerTrexASkelerex      bool `json:"rwd_bracer_trex_a_skelerex"`
	RWDBoosterTrexASkelerex     bool `json:"rwd_booster_trex_a_skelerex"`
	RWDChassisTrexASkelerex     bool `json:"rwd_chassis_trex_a_skelerex"`
	RWDDecal0021                bool `json:"rwd_decal_0021"`
	RWDGoalFx0006               bool `json:"rwd_goal_fx_0006"`
	RWDBanner0029               bool `json:"rwd_banner_0029"`
	RWDTag0034                  bool `json:"rwd_tag_0034"`
	RWDPip0002                  bool `json:"rwd_pip_0002"`
	RWDPip0016                  bool `json:"rwd_pip_0016"`
	RWDPip0012                  bool `json:"rwd_pip_0012"`
	RWDEmissive0015             bool `json:"rwd_emissive_0015"`
	RWDEmissive0018             bool `json:"rwd_emissive_0018"`
	RWDEmissive0019             bool `json:"rwd_emissive_0019"`
	RWDEmissive0020             bool `json:"rwd_emissive_0020"`
	RWDBannerS2Lines            bool `json:"rwd_banner_s2_lines"`
	RWDTintS2ADefault           bool `json:"rwd_tint_s2_a_default"`
	RWDTitle0016                bool `json:"rwd_title_0016"`
	RWDTitle0017                bool `json:"rwd_title_0017"`
	RWDTitle0018                bool `json:"rwd_title_0018"`
	RWDTitle0019                bool `json:"rwd_title_0019"`
}

type UnlocksCombat struct {
	DecalCombatFlamingoA   bool `json:"decal_combat_flamingo_a"`
	DecalCombatLogoA       bool `json:"decal_combat_logo_a"`
	EmoteDizzyEyesA        bool `json:"emote_dizzy_eyes_a"`
	PatternLightningA      bool `json:"pattern_lightning_a"`
	RWDBoosterS10          bool `json:"rwd_booster_s10"`
	RWDChassisBodyS10A     bool `json:"rwd_chassis_body_s10_a"`
	RWDMedalS1CombatBronze bool `json:"rwd_medal_s1_combat_bronze"`
	RWDMedalS1CombatGold   bool `json:"rwd_medal_s1_combat_gold"`
	RWDMedalS1CombatSilver bool `json:"rwd_medal_s1_combat_silver"`
	RWDTitleTitleB         bool `json:"rwd_title_title_b"`
}

/*
func DefaultGameProfiles(xplatformid EchoUserId, displayname string) GameProfiles {
	return GameProfiles{
		Client: DefaultClientProfile(xplatformid, displayname),
		Server: DefaultServerProfile(xplatformid, displayname),
	}
}
*/

func DefaultClientProfile(echoUserId EchoUserId, displayName string) EchoPlayerPreferences {
	return EchoPlayerPreferences{
		DisplayName:     displayName,
		EchoUserIdToken: echoUserId.String(),

		CombatWeapon:       "scout",
		CombatGrenade:      "det",
		CombatDominantHand: 1,
		CombatAbility:      "heal",
		LegalConsents: LegalConsents{
			PointsPolicyVersion: 1,
			EulaVersion:         1,
			GameAdminVersion:    1,
			SplashScreenVersion: 2,
		},
		NewPlayerProgress: NewPlayerProgress{
			Lobby: NpeMilestone{Completed: true},

			FirstMatch:        NpeMilestone{Completed: true},
			Movement:          NpeMilestone{Completed: true},
			ArenaBasics:       NpeMilestone{Completed: true},
			SocialTabSeen:     Versioned{Version: 1},
			Pointer:           Versioned{Version: 1},
			BlueTintTabSeen:   Versioned{Version: 1},
			HeraldryTabSeen:   Versioned{Version: 1},
			OrangeTintTabSeen: Versioned{Version: 1},
		},
		Customization: Customization{
			BattlePassSeasonPoiVersion: 0,
			NewUnlocksPoiVersion:       1,
			StoreEntryPoiVersion:       0,
			ClearNewUnlocksVersion:     1,
		},
		Social: Social{
			CommunityValuesVersion: 1,
			SetupVersion:           1,
			Group:                  "90DD4DB5-B5DD-4655-839E-FDBE5F4BC0BF",
		},
		NewUnlocks: []int64{},
	}
}

func DefaultServerProfile(gameUserId EchoUserId, displayName string) ServerProfile {
	distServerProfile := []byte(`
	{
		"_version": 4,
		"loadout": {
			"instances": {
				"unified": {
					"slots": {
						"emote": "emote_blink_smiley_a",
						"decal": "decal_default",
						"tint": "tint_neutral_a_default",
						"tint_alignment_a": "tint_blue_a_default",
						"tint_alignment_b": "tint_orange_a_default",
						"pattern": "pattern_default",
						"pip": "rwd_decalback_default",
						"chassis": "rwd_chassis_body_s11_a",
						"bracer": "rwd_bracer_default",
						"booster": "rwd_booster_default",
						"title": "rwd_title_title_default",
						"tag": "rwd_tag_s1_a_secondary",
						"banner": "rwd_banner_s1_default",
						"medal": "rwd_medal_default",
						"goal_fx": "rwd_goal_fx_default",
						"secondemote": "emote_blink_smiley_a",
						"emissive": "emissive_default",
						"tint_body": "tint_neutral_a_default",
						"pattern_body": "pattern_default",
						"decal_body": "decal_default"
					}
				}
			},
			"number": 1
		},
		"stats": {
			"arena": {
				"Level": {
					"cnt": 1,
					"op": "add",
					"val": 1
				}
			},
			"combat": {
				"Level": {
					"cnt": 1,
					"op": "add",
					"val": 1
				}
			}
		},
		"unlocks": {
			"arena": {
				"decal_combat_flamingo_a": true,
				"decal_combat_logo_a": true,
				"decal_default": true,
				"decal_sheldon_a": true,
				"emote_blink_smiley_a": true,
				"emote_default": true,
				"emote_dizzy_eyes_a": true,
				"loadout_number": true,
				"pattern_default": true,
				"pattern_lightning_a": true,
				"rwd_banner_s1_default": true,
				"rwd_booster_default": true,
				"rwd_bracer_default": true,
				"rwd_chassis_body_s11_a": true,
				"rwd_decalback_default": true,
				"rwd_decalborder_default": true,
				"rwd_medal_default": true,
				"rwd_tag_default": true,
				"rwd_tag_s1_a_secondary": true,
				"rwd_title_title_default": true,
				"tint_blue_a_default": true,
				"tint_neutral_a_default": true,
				"tint_neutral_a_s10_default": true,
				"tint_orange_a_default": true,
				"rwd_goal_fx_default" : true,
				"emissive_default" : true
			},
			"combat": {
				"rwd_booster_s10": true,
				"rwd_chassis_body_s10_a": true
			}
		},
		"xplatformid": "default_battlepass"
	}
	`)
	var serverProfile ServerProfile
	err := json.Unmarshal([]byte(distServerProfile), &serverProfile)
	if err != nil {
		fmt.Println("JSON decode error!")
		return ServerProfile{}
	}
	serverProfile.EchoUserIdToken = gameUserId.String()
	return serverProfile
}

func DefaultGameProfiles(xplatformid EchoUserId, displayname string) GameProfiles {
	return GameProfiles{
		Client: DefaultClientProfile(xplatformid, displayname),
		Server: DefaultServerProfile(xplatformid, displayname),
	}
}
