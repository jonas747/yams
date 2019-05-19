package connection

type Gamemode uint8

const (
	GameModeSurvival Gamemode = iota
	GameModeCreative
	GameModeAdventure
	GameModeSpectator
)

// State represents the connection state
type State int

const (
	StateHandshake = 0 // Initial handshake state
	StateStatus    = 1 // Status is used in the server browser when the client requests the status
	StateLogin     = 2
	StatePlay      = 3
)

func (s State) String() string {
	switch s {
	case StateHandshake:
		return "StateHandshake"
	case StateStatus:
		return "StateStatus"
	case StateLogin:
		return "StateLogin"
	case StatePlay:
		return "StatePlay"
	}

	return "UnknownState"
}

type Dimension uint32

const (
	DimensionNether    Dimension = 0xFF
	DimensionOverworld Dimension = 0
	DimensionEnd       Dimension = 1
)

type Difficulty uint8

const (
	DifficultyPeaceful Difficulty = iota
	DifficultyEasy
	DifficultyNormal
	DifficultyHard
)

type ChatPosition uint8

const (
	ChatChatbox ChatPosition = iota
	ChatSystem
	ChatActionBar
)

type ScoreboardPosition uint8

const (
	ScoreBoardList ScoreboardPosition = iota
	ScoreBoardSidebar
	ScoreBoardBelowName
)

type BossBarAction int

const (
	BossBarActionAdd BossBarAction = iota
	BossBarActionRemove
	BossBarActionUpdateHealth
	BossBarActionUpdateTitle
	BossBarActionUpdateStyle
	BossBarActionFlags
)

type BossBarColor int

const (
	BossBarColorPink BossBarColor = iota
	BossBarColorBlue
	BossBarColorRed
	BossBarColorGreen
	BossBarColorYellow
	BossBarColorPurple
	BossBarColorWhite
)

type BossBarDivision int

const (
	BossBarDivisionNoDivision BossBarDivision = iota
	BossBarDivision6Notches
	BossBarDivision10Notches
	BossBarDivision12Notches
	BossBarDivision20Notches
)

type LevelType string

const (
	LevelTypeDefault     LevelType = "default"
	LevelTypeFlat        LevelType = "flat"
	LevelTypeLargeBiomes LevelType = "largeBiomes"
	LevelTypeAmplified   LevelType = "amplified"
	LevelTypeDefault_1_1 LevelType = "default_1_1"
)

type ClientStatusAction uint8

const (
	ClientStatusActionPerformRespawn ClientStatusAction = iota
	ClientStatusActionRequestStats
	ClientStatusActionOpenInventory
)

type Position struct {
	X int
	Y int
	Z int
}

type Protocol uint16

const (
	V1_7_2  Protocol = 4
	V1_7_6  Protocol = 5
	V1_8    Protocol = 47
	V1_9    Protocol = 107
	V1_9_1  Protocol = 108
	V1_9_2  Protocol = 109
	V1_9_3  Protocol = 110
	V1_10   Protocol = 210
	V1_11   Protocol = 315
	V1_11_1 Protocol = 316
	V1_12   Protocol = 335
	V1_12_1 Protocol = 338
	V1_12_2 Protocol = 340
	V1_13   Protocol = 393
	V1_13_1 Protocol = 401
	V1_13_2 Protocol = 404
	V1_14   Protocol = 477
	V1_14_1 Protocol = 480
)
