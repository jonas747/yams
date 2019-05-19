package packetmappings

func init() {
	init482()
}

const (
	StateHandshake = 0 // Initial handshake state
	StateStatus    = 1 // Status is used in the server browser when the client requests the status
	StateLogin     = 2
	StatePlay      = 3
)

// map[protocol_version]map[YAMPacketID]mc_packet_id
var toMCPacket = make(map[int]map[YAMPacketID]int64)

// map[protocol_version]map[mc_packet_id]YAMPacketID
var toYAMPacket = make(map[int]map[int64]YAMPacketID)

// abstraction layer, these are mapped to actual packetID's depending on the version and mapping used
type YAMPacketID int

const (
	// Hanshaking
	HandshakingServerHandshake YAMPacketID = iota
	HandshakingServerLegacyServerListPing

	// Play
	// ClientBound
	PlayClientSpawnObject
	PlayClientSpawnExperienceOrb
	PlayClientSpawnGlobalEntity
	PlayClientSpawnMob
	PlayClientSpawnPainting
	PlayClientSpawnPlayer
	PlayClientAnimation
	PlayClientStatistics
	PlayClientBlockBreakAnimation
	PlayClientUpdateBlockEntity
	PlayClientBlockAction
	PlayClientBlockChange
	PlayClientBossBar
	PlayClientServerDifficulty
	PlayClientChatMessage
	PlayClientMultiBlockChange
	PlayClientTabComplete
	PlayClientDeclareCommands
	PlayClientConfirmTransaction
	PlayClientCloseWindow
	PlayClientOpenWindow
	PlayClientWindowItems
	PlayClientWindowProperty
	PlayClientSetSlot
	PlayClientSetCooldown
	PlayClientPluginMessage
	PlayClientNamedSoundEffect
	PlayClientDisconnect
	PlayClientEntityStatus
	PlayClientNBTQueryResponse
	PlayClientExplosion
	PlayClientUnloadChunk
	PlayClientChangeGameState
	PlayClientOpenHorseWindow
	PlayClientKeepAlive
	PlayClientChunkData
	PlayClientEffect
	PlayClientParticle
	PlayClientUpdateLight
	PlayClientJoinGame
	PlayClientMapData
	PlayClientTradeList
	PlayClientEntity
	PlayClientEntityRelativeMove
	PlayClientEntityLookAndRelativeMove
	PlayClientEntityLook
	PlayClientVehicleMove
	PlayClientOpenBook
	PlayClientOpenSignEditor
	PlayClientCraftRecipeResponse
	PlayClientPlayerAbilities
	PlayClientCombatEvent
	PlayClientPlayerInfo
	PlayClientFacePlayer
	PlayClientPlayerPositionAndLook
	PlayClientUseBed
	PlayClientUnlockRecipes
	PlayClientDestroyEntities
	PlayClientRemoveEntityEffect
	PlayClientResourcePackSend
	PlayClientRespawn
	PlayClientEntityHeadLook
	PlayClientSelectAdvancementTab
	PlayClientWorldBorder
	PlayClientCamera
	PlayClientHeldItemChange
	PlayClientUpdateViewPosition
	PlayClientUpdateViewDistance
	PlayClientDisplayScoreboard
	PlayClientEntityMetadata
	PlayClientAttachEntity
	PlayClientEntityVelocity
	PlayClientEntityEquipment
	PlayClientSetExperience
	PlayClientUpdateHealth
	PlayClientScoreboardObjective
	PlayClientSetPassengers
	PlayClientTeams
	PlayClientUpdateScore
	PlayClientSpawnPosition
	PlayClientTimeUpdate
	PlayClientTitle
	PlayClientStopSound
	PlayClientSoundEffect
	PlayClientPlayerListHeaderAndFooter
	PlayClientCollectItem
	PlayClientEntityTeleport
	PlayClientAdvancements
	PlayClientEntityProperties
	PlayClientEntityEffect
	PlayClientDeclareRecipes
	PlayClientTags
	PlayClientEntitySoundEffect

	// Serverbound
	PlayServerTeleportConfirm
	PlayServerQueryBlockNBT
	PlayServerChatMessage
	PlayServerClientStatus
	PlayServerClientSettings
	PlayServerTabComplete
	PlayServerConfirmTransaction
	PlayServerEnchantItem
	PlayServerClickWindow
	PlayServerCloseWindow
	PlayServerPluginMessage
	PlayServerEditBook
	PlayServerQueryEntityNBT
	PlayServerUseEntity
	PlayServerKeepAlive
	PlayServerPlayer
	PlayServerPlayerPosition
	PlayServerPlayerPositionAndLook
	PlayServerPlayerLook
	PlayServerVehicleMove
	PlayServerSteerBoat
	PlayServerPickItem
	PlayServerCraftRecipeRequest
	PlayServerPlayerAbilities
	PlayServerPlayerDigging
	PlayServerEntityAction
	PlayServerSteerVehicle
	PlayServerRecipeBookData
	PlayServerNameItem
	PlayServerResourcePackStatus
	PlayServerAdvancementTab
	PlayServerSelectTrade
	PlayServerSetBeaconEffect
	PlayServerHeldItemChange
	PlayServerUpdateCommandBlock
	PlayServerUpdateCommandBlockMinecart
	PlayServerCreativeInventoryAction
	PlayServerUpdateStructureBlock
	PlayServerUpdateSign
	PlayServerAnimation
	PlayServerSpectate
	PlayServerPlayerBlockPlacement
	PlayServerUseItem
	PlayServerUpdateJigsawBlock
	PlayServerSetDifficulty
	PlayServerClickWindowButton
	PlayServerLockDifficulty

	// Status
	StatusClientResponse
	StatusClientPong
	StatusServerRequest
	StatusServerPing

	// Login
	// Clientbound
	LoginClientDisconnect
	LoginClientEncryptionRequest
	LoginClientLoginSuccess
	LoginClientSetCompression
	LoginClientLoginPluginRequest
	// Serverbound
	LoginServerLoginStart
	LoginServerEncryptionResponse
	LoginServerLoginPluginResponse

	// internal ones
	Disconnected
)

func reverseMapping(in map[int]int) map[int]int {
	dst := make(map[int]int)

	for k, v := range in {
		dst[v] = k
	}

	return dst
}

func copyMapping(dst, source map[int64]YAMPacketID) {
	for k, v := range source {
		dst[k] = v
	}
}

func uniqueMCPacketID(client bool, state, packetID int32) int64 {
	state8 := int8(state)
	cb := 0
	if client {
		cb = 1
	}

	return int64(cb)<<40 | int64(state8)<<32 | int64(packetID)
}

type Mapping struct {
	toYAM map[int64]YAMPacketID
}

func NewMapping() *Mapping {
	return &Mapping{
		toYAM: make(map[int64]YAMPacketID),
	}
}

func (m *Mapping) Set(client bool, state, packetID int32, yam YAMPacketID) {
	m.toYAM[uniqueMCPacketID(client, state, packetID)] = yam
}

func (m *Mapping) finish(version int) {
	reversed := make(map[YAMPacketID]int64)

	for k, v := range m.toYAM {
		reversed[v] = k
	}

	toYAMPacket[version] = m.toYAM
	toMCPacket[version] = reversed
}

func GetMCPacketID(version int, yamID YAMPacketID) int32 {
	return int32(toMCPacket[version][yamID] & 0xffffffff)
}

func GetYAMPacketID(version int, state int, client bool, mcPacketID int) YAMPacketID {
	uniqueID := uniqueMCPacketID(client, int32(state), int32(mcPacketID))
	if yp, ok := toYAMPacket[version][uniqueID]; ok {
		return yp
	}
	return -1
}

var stringedYAMIDs = map[YAMPacketID]string{

	// Hanshaking
	HandshakingServerHandshake:            "HandshakingServerHandshake",
	HandshakingServerLegacyServerListPing: "HandshakingServerLegacyServerListPing",

	// Play
	// ClientBound
	PlayClientSpawnObject:               "PlayClientSpawnObject",
	PlayClientSpawnExperienceOrb:        "PlayClientSpawnExperienceOrb",
	PlayClientSpawnGlobalEntity:         "PlayClientSpawnGlobalEntity",
	PlayClientSpawnMob:                  "PlayClientSpawnMob",
	PlayClientSpawnPainting:             "PlayClientSpawnPainting",
	PlayClientSpawnPlayer:               "PlayClientSpawnPlayer",
	PlayClientAnimation:                 "PlayClientAnimation",
	PlayClientStatistics:                "PlayClientStatistics",
	PlayClientBlockBreakAnimation:       "PlayClientBlockBreakAnimation",
	PlayClientUpdateBlockEntity:         "PlayClientUpdateBlockEntity",
	PlayClientBlockAction:               "PlayClientBlockAction",
	PlayClientBlockChange:               "PlayClientBlockChange",
	PlayClientBossBar:                   "PlayClientBossBar",
	PlayClientServerDifficulty:          "PlayClientServerDifficulty",
	PlayClientChatMessage:               "PlayClientChatMessage",
	PlayClientMultiBlockChange:          "PlayClientMultiBlockChange",
	PlayClientTabComplete:               "PlayClientTabComplete",
	PlayClientDeclareCommands:           "PlayClientDeclareCommands",
	PlayClientConfirmTransaction:        "PlayClientConfirmTransaction",
	PlayClientCloseWindow:               "PlayClientCloseWindow",
	PlayClientOpenWindow:                "PlayClientOpenWindow",
	PlayClientWindowItems:               "PlayClientWindowItems",
	PlayClientWindowProperty:            "PlayClientWindowProperty",
	PlayClientSetSlot:                   "PlayClientSetSlot",
	PlayClientSetCooldown:               "PlayClientSetCooldown",
	PlayClientPluginMessage:             "PlayClientPluginMessage",
	PlayClientNamedSoundEffect:          "PlayClientNamedSoundEffect",
	PlayClientDisconnect:                "PlayClientDisconnect",
	PlayClientEntityStatus:              "PlayClientEntityStatus",
	PlayClientNBTQueryResponse:          "PlayClientNBTQueryResponse",
	PlayClientExplosion:                 "PlayClientExplosion",
	PlayClientUnloadChunk:               "PlayClientUnloadChunk",
	PlayClientChangeGameState:           "PlayClientChangeGameState",
	PlayClientOpenHorseWindow:           "PlayClientOpenHorseWindow",
	PlayClientKeepAlive:                 "PlayClientKeepAlive",
	PlayClientChunkData:                 "PlayClientChunkData",
	PlayClientEffect:                    "PlayClientEffect",
	PlayClientParticle:                  "PlayClientParticle",
	PlayClientUpdateLight:               "PlayClientUpdateLight",
	PlayClientJoinGame:                  "PlayClientJoinGame",
	PlayClientMapData:                   "PlayClientMapData",
	PlayClientTradeList:                 "PlayClientTradeList",
	PlayClientEntity:                    "PlayClientEntity",
	PlayClientEntityRelativeMove:        "PlayClientEntityRelativeMove",
	PlayClientEntityLookAndRelativeMove: "PlayClientEntityLookAndRelativeMove",
	PlayClientEntityLook:                "PlayClientEntityLook",
	PlayClientVehicleMove:               "PlayClientVehicleMove",
	PlayClientOpenBook:                  "PlayClientOpenBook",
	PlayClientOpenSignEditor:            "PlayClientOpenSignEditor",
	PlayClientCraftRecipeResponse:       "PlayClientCraftRecipeResponse",
	PlayClientPlayerAbilities:           "PlayClientPlayerAbilities",
	PlayClientCombatEvent:               "PlayClientCombatEvent",
	PlayClientPlayerInfo:                "PlayClientPlayerInfo",
	PlayClientFacePlayer:                "PlayClientFacePlayer",
	PlayClientPlayerPositionAndLook:     "PlayClientPlayerPositionAndLook",
	PlayClientUseBed:                    "PlayClientUseBed",
	PlayClientUnlockRecipes:             "PlayClientUnlockRecipes",
	PlayClientDestroyEntities:           "PlayClientDestroyEntities",
	PlayClientRemoveEntityEffect:        "PlayClientRemoveEntityEffect",
	PlayClientResourcePackSend:          "PlayClientResourcePackSend",
	PlayClientRespawn:                   "PlayClientRespawn",
	PlayClientEntityHeadLook:            "PlayClientEntityHeadLook",
	PlayClientSelectAdvancementTab:      "PlayClientSelectAdvancementTab",
	PlayClientWorldBorder:               "PlayClientWorldBorder",
	PlayClientCamera:                    "PlayClientCamera",
	PlayClientHeldItemChange:            "PlayClientHeldItemChange",
	PlayClientUpdateViewPosition:        "PlayClientUpdateViewPosition",
	PlayClientUpdateViewDistance:        "PlayClientUpdateViewDistance",
	PlayClientDisplayScoreboard:         "PlayClientDisplayScoreboard",
	PlayClientEntityMetadata:            "PlayClientEntityMetadata",
	PlayClientAttachEntity:              "PlayClientAttachEntity",
	PlayClientEntityVelocity:            "PlayClientEntityVelocity",
	PlayClientEntityEquipment:           "PlayClientEntityEquipment",
	PlayClientSetExperience:             "PlayClientSetExperience",
	PlayClientUpdateHealth:              "PlayClientUpdateHealth",
	PlayClientScoreboardObjective:       "PlayClientScoreboardObjective",
	PlayClientSetPassengers:             "PlayClientSetPassengers",
	PlayClientTeams:                     "PlayClientTeams",
	PlayClientUpdateScore:               "PlayClientUpdateScore",
	PlayClientSpawnPosition:             "PlayClientSpawnPosition",
	PlayClientTimeUpdate:                "PlayClientTimeUpdate",
	PlayClientTitle:                     "PlayClientTitle",
	PlayClientStopSound:                 "PlayClientStopSound",
	PlayClientSoundEffect:               "PlayClientSoundEffect",
	PlayClientPlayerListHeaderAndFooter: "PlayClientPlayerListHeaderAndFooter",
	PlayClientCollectItem:               "PlayClientCollectItem",
	PlayClientEntityTeleport:            "PlayClientEntityTeleport",
	PlayClientAdvancements:              "PlayClientAdvancements",
	PlayClientEntityProperties:          "PlayClientEntityProperties",
	PlayClientEntityEffect:              "PlayClientEntityEffect",
	PlayClientDeclareRecipes:            "PlayClientDeclareRecipes",
	PlayClientTags:                      "PlayClientTags",
	PlayClientEntitySoundEffect:         "PlayClientEntitySoundEffect",

	// Serverbound
	PlayServerTeleportConfirm:            "PlayServerTeleportConfirm",
	PlayServerQueryBlockNBT:              "PlayServerQueryBlockNBT",
	PlayServerChatMessage:                "PlayServerChatMessage",
	PlayServerClientStatus:               "PlayServerClientStatus",
	PlayServerClientSettings:             "PlayServerClientSettings",
	PlayServerTabComplete:                "PlayServerTabComplete",
	PlayServerConfirmTransaction:         "PlayServerConfirmTransaction",
	PlayServerEnchantItem:                "PlayServerEnchantItem",
	PlayServerClickWindow:                "PlayServerClickWindow",
	PlayServerCloseWindow:                "PlayServerCloseWindow",
	PlayServerPluginMessage:              "PlayServerPluginMessage",
	PlayServerEditBook:                   "PlayServerEditBook",
	PlayServerQueryEntityNBT:             "PlayServerQueryEntityNBT",
	PlayServerUseEntity:                  "PlayServerUseEntity",
	PlayServerKeepAlive:                  "PlayServerKeepAlive",
	PlayServerPlayer:                     "PlayServerPlayer",
	PlayServerPlayerPosition:             "PlayServerPlayerPosition",
	PlayServerPlayerPositionAndLook:      "PlayServerPlayerPositionAndLook",
	PlayServerPlayerLook:                 "PlayServerPlayerLook",
	PlayServerVehicleMove:                "PlayServerVehicleMove",
	PlayServerSteerBoat:                  "PlayServerSteerBoat",
	PlayServerPickItem:                   "PlayServerPickItem",
	PlayServerCraftRecipeRequest:         "PlayServerCraftRecipeRequest",
	PlayServerPlayerAbilities:            "PlayServerPlayerAbilities",
	PlayServerPlayerDigging:              "PlayServerPlayerDigging",
	PlayServerEntityAction:               "PlayServerEntityAction",
	PlayServerSteerVehicle:               "PlayServerSteerVehicle",
	PlayServerRecipeBookData:             "PlayServerRecipeBookData",
	PlayServerNameItem:                   "PlayServerNameItem",
	PlayServerResourcePackStatus:         "PlayServerResourcePackStatus",
	PlayServerAdvancementTab:             "PlayServerAdvancementTab",
	PlayServerSelectTrade:                "PlayServerSelectTrade",
	PlayServerSetBeaconEffect:            "PlayServerSetBeaconEffect",
	PlayServerHeldItemChange:             "PlayServerHeldItemChange",
	PlayServerUpdateCommandBlock:         "PlayServerUpdateCommandBlock",
	PlayServerUpdateCommandBlockMinecart: "PlayServerUpdateCommandBlockMinecart",
	PlayServerCreativeInventoryAction:    "PlayServerCreativeInventoryAction",
	PlayServerUpdateStructureBlock:       "PlayServerUpdateStructureBlock",
	PlayServerUpdateSign:                 "PlayServerUpdateSign",
	PlayServerAnimation:                  "PlayServerAnimation",
	PlayServerSpectate:                   "PlayServerSpectate",
	PlayServerPlayerBlockPlacement:       "PlayServerPlayerBlockPlacement",
	PlayServerUseItem:                    "PlayServerUseItem",
	PlayServerUpdateJigsawBlock:          "PlayServerUpdateJigsawBlock",
	PlayServerSetDifficulty:              "PlayServerSetDifficulty",
	PlayServerClickWindowButton:          "PlayServerClickWindowButton",
	PlayServerLockDifficulty:             "PlayServerLockDifficulty",

	// Status
	StatusClientResponse: "StatusClientResponse",
	StatusClientPong:     "StatusClientPong",
	StatusServerRequest:  "StatusServerRequest",
	StatusServerPing:     "StatusServerPing",

	// Login
	// Clientbound
	LoginClientDisconnect:         "LoginClientDisconnect",
	LoginClientEncryptionRequest:  "LoginClientEncryptionRequest",
	LoginClientLoginSuccess:       "LoginClientLoginSuccess",
	LoginClientSetCompression:     "LoginClientSetCompression",
	LoginClientLoginPluginRequest: "LoginClientLoginPluginRequest",
	// Serverbound
	LoginServerLoginStart:          "LoginServerLoginStart",
	LoginServerEncryptionResponse:  "LoginServerEncryptionResponse",
	LoginServerLoginPluginResponse: "LoginServerLoginPluginResponse",

	Disconnected: "Disconnected",
}

func (y YAMPacketID) String() string {
	if s, ok := stringedYAMIDs[y]; ok {
		return s
	}

	return "Unknown ID"
}
