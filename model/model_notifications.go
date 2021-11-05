package model

import (
	"fmt"
	"github.com/Wilm0rien/omip/util"
)

const (
	NotiSndTyp_character = iota
	NotiSndTyp_corporation
	NotiSndTyp_alliance
	NotiSndTyp_faction
	NotiSndTyp_other
)

const (
	NotiMsgTyp_AcceptedAlly = iota
	NotiMsgTyp_AcceptedSurrender
	NotiMsgTyp_AgentRetiredTrigravian
	NotiMsgTyp_AllAnchoringMsg
	NotiMsgTyp_AllMaintenanceBillMsg
	NotiMsgTyp_AllStrucInvulnerableMsg
	NotiMsgTyp_AllStructVulnerableMsg
	NotiMsgTyp_AllWarCorpJoinedAllianceMsg
	NotiMsgTyp_AllWarDeclaredMsg
	NotiMsgTyp_AllWarInvalidatedMsg
	NotiMsgTyp_AllWarRetractedMsg
	NotiMsgTyp_AllWarSurrenderMsg
	NotiMsgTyp_AllianceCapitalChanged
	NotiMsgTyp_AllianceWarDeclaredV2
	NotiMsgTyp_AllyContractCancelled
	NotiMsgTyp_AllyJoinedWarAggressorMsg
	NotiMsgTyp_AllyJoinedWarAllyMsg
	NotiMsgTyp_AllyJoinedWarDefenderMsg
	NotiMsgTyp_BattlePunishFriendlyFire
	NotiMsgTyp_BillOutOfMoneyMsg
	NotiMsgTyp_BillPaidCorpAllMsg
	NotiMsgTyp_BountyClaimMsg
	NotiMsgTyp_BountyESSShared
	NotiMsgTyp_BountyESSTaken
	NotiMsgTyp_BountyPlacedAlliance
	NotiMsgTyp_BountyPlacedChar
	NotiMsgTyp_BountyPlacedCorp
	NotiMsgTyp_BountyYourBountyClaimed
	NotiMsgTyp_BuddyConnectContactAdd
	NotiMsgTyp_CharAppAcceptMsg
	NotiMsgTyp_CharAppRejectMsg
	NotiMsgTyp_CharAppWithdrawMsg
	NotiMsgTyp_CharLeftCorpMsg
	NotiMsgTyp_CharMedalMsg
	NotiMsgTyp_CharTerminationMsg
	NotiMsgTyp_CloneActivationMsg
	NotiMsgTyp_CloneActivationMsg2
	NotiMsgTyp_CloneMovedMsg
	NotiMsgTyp_CloneRevokedMsg1
	NotiMsgTyp_CloneRevokedMsg2
	NotiMsgTyp_CombatOperationFinished
	NotiMsgTyp_ContactAdd
	NotiMsgTyp_ContactEdit
	NotiMsgTyp_ContainerPasswordMsg
	NotiMsgTyp_ContractRegionChangedToPochven
	NotiMsgTyp_CorpAllBillMsg
	NotiMsgTyp_CorpAppAcceptMsg
	NotiMsgTyp_CorpAppInvitedMsg
	NotiMsgTyp_CorpAppNewMsg
	NotiMsgTyp_CorpAppRejectCustomMsg
	NotiMsgTyp_CorpAppRejectMsg
	NotiMsgTyp_CorpBecameWarEligible
	NotiMsgTyp_CorpDividendMsg
	NotiMsgTyp_CorpFriendlyFireDisableTimerCompleted
	NotiMsgTyp_CorpFriendlyFireDisableTimerStarted
	NotiMsgTyp_CorpFriendlyFireEnableTimerCompleted
	NotiMsgTyp_CorpFriendlyFireEnableTimerStarted
	NotiMsgTyp_CorpKicked
	NotiMsgTyp_CorpLiquidationMsg
	NotiMsgTyp_CorpNewCEOMsg
	NotiMsgTyp_CorpNewsMsg
	NotiMsgTyp_CorpNoLongerWarEligible
	NotiMsgTyp_CorpOfficeExpirationMsg
	NotiMsgTyp_CorpStructLostMsg
	NotiMsgTyp_CorpTaxChangeMsg
	NotiMsgTyp_CorpVoteCEORevokedMsg
	NotiMsgTyp_CorpVoteMsg
	NotiMsgTyp_CorpWarDeclaredMsg
	NotiMsgTyp_CorpWarDeclaredV2
	NotiMsgTyp_CorpWarFightingLegalMsg
	NotiMsgTyp_CorpWarInvalidatedMsg
	NotiMsgTyp_CorpWarRetractedMsg
	NotiMsgTyp_CorpWarSurrenderMsg
	NotiMsgTyp_CustomsMsg
	NotiMsgTyp_DeclareWar
	NotiMsgTyp_DistrictAttacked
	NotiMsgTyp_DustAppAcceptedMsg
	NotiMsgTyp_ESSMainBankLink
	NotiMsgTyp_EntosisCaptureStarted
	NotiMsgTyp_FWAllianceKickMsg
	NotiMsgTyp_FWAllianceWarningMsg
	NotiMsgTyp_FWCharKickMsg
	NotiMsgTyp_FWCharRankGainMsg
	NotiMsgTyp_FWCharRankLossMsg
	NotiMsgTyp_FWCharWarningMsg
	NotiMsgTyp_FWCorpJoinMsg
	NotiMsgTyp_FWCorpKickMsg
	NotiMsgTyp_FWCorpLeaveMsg
	NotiMsgTyp_FWCorpWarningMsg
	NotiMsgTyp_FacWarCorpJoinRequestMsg
	NotiMsgTyp_FacWarCorpJoinWithdrawMsg
	NotiMsgTyp_FacWarCorpLeaveRequestMsg
	NotiMsgTyp_FacWarCorpLeaveWithdrawMsg
	NotiMsgTyp_FacWarLPDisqualifiedEvent
	NotiMsgTyp_FacWarLPDisqualifiedKill
	NotiMsgTyp_FacWarLPPayoutEvent
	NotiMsgTyp_FacWarLPPayoutKill
	NotiMsgTyp_GameTimeAdded
	NotiMsgTyp_GameTimeReceived
	NotiMsgTyp_GameTimeSent
	NotiMsgTyp_GiftReceived
	NotiMsgTyp_IHubDestroyedByBillFailure
	NotiMsgTyp_IncursionCompletedMsg
	NotiMsgTyp_IndustryOperationFinished
	NotiMsgTyp_IndustryTeamAuctionLost
	NotiMsgTyp_IndustryTeamAuctionWon
	NotiMsgTyp_InfrastructureHubBillAboutToExpire
	NotiMsgTyp_InsuranceExpirationMsg
	NotiMsgTyp_InsuranceFirstShipMsg
	NotiMsgTyp_InsuranceInvalidatedMsg
	NotiMsgTyp_InsuranceIssuedMsg
	NotiMsgTyp_InsurancePayoutMsg
	NotiMsgTyp_InvasionCompletedMsg
	NotiMsgTyp_InvasionSystemLogin
	NotiMsgTyp_InvasionSystemStart
	NotiMsgTyp_JumpCloneDeletedMsg1
	NotiMsgTyp_JumpCloneDeletedMsg2
	NotiMsgTyp_KillReportFinalBlow
	NotiMsgTyp_KillReportVictim
	NotiMsgTyp_KillRightAvailable
	NotiMsgTyp_KillRightAvailableOpen
	NotiMsgTyp_KillRightEarned
	NotiMsgTyp_KillRightUnavailable
	NotiMsgTyp_KillRightUnavailableOpen
	NotiMsgTyp_KillRightUsed
	NotiMsgTyp_LocateCharMsg
	NotiMsgTyp_MadeWarMutual
	NotiMsgTyp_MercOfferRetractedMsg
	NotiMsgTyp_MercOfferedNegotiationMsg
	NotiMsgTyp_MissionCanceledTriglavian
	NotiMsgTyp_MissionOfferExpirationMsg
	NotiMsgTyp_MissionTimeoutMsg
	NotiMsgTyp_MoonminingAutomaticFracture
	NotiMsgTyp_MoonminingExtractionCancelled
	NotiMsgTyp_MoonminingExtractionFinished
	NotiMsgTyp_MoonminingExtractionStarted
	NotiMsgTyp_MoonminingLaserFired
	NotiMsgTyp_MutualWarExpired
	NotiMsgTyp_MutualWarInviteAccepted
	NotiMsgTyp_MutualWarInviteRejected
	NotiMsgTyp_MutualWarInviteSent
	NotiMsgTyp_NPCStandingsGained
	NotiMsgTyp_NPCStandingsLost
	NotiMsgTyp_OfferToAllyRetracted
	NotiMsgTyp_OfferedSurrender
	NotiMsgTyp_OfferedToAlly
	NotiMsgTyp_OfficeLeaseCanceledInsufficientStandings
	NotiMsgTyp_OldLscMessages
	NotiMsgTyp_OperationFinished
	NotiMsgTyp_OrbitalAttacked
	NotiMsgTyp_OrbitalReinforced
	NotiMsgTyp_OwnershipTransferred
	NotiMsgTyp_RaffleCreated
	NotiMsgTyp_RaffleExpired
	NotiMsgTyp_RaffleFinished
	NotiMsgTyp_ReimbursementMsg
	NotiMsgTyp_ResearchMissionAvailableMsg
	NotiMsgTyp_RetractsWar
	NotiMsgTyp_SeasonalChallengeCompleted
	NotiMsgTyp_SovAllClaimAquiredMsg
	NotiMsgTyp_SovAllClaimLostMsg
	NotiMsgTyp_SovCommandNodeEventStarted
	NotiMsgTyp_SovCorpBillLateMsg
	NotiMsgTyp_SovCorpClaimFailMsg
	NotiMsgTyp_SovDisruptorMsg
	NotiMsgTyp_SovStationEnteredFreeport
	NotiMsgTyp_SovStructureDestroyed
	NotiMsgTyp_SovStructureReinforced
	NotiMsgTyp_SovStructureSelfDestructCancel
	NotiMsgTyp_SovStructureSelfDestructFinished
	NotiMsgTyp_SovStructureSelfDestructRequested
	NotiMsgTyp_SovereigntyIHDamageMsg
	NotiMsgTyp_SovereigntySBUDamageMsg
	NotiMsgTyp_SovereigntyTCUDamageMsg
	NotiMsgTyp_StationAggressionMsg1
	NotiMsgTyp_StationAggressionMsg2
	NotiMsgTyp_StationConquerMsg
	NotiMsgTyp_StationServiceDisabled
	NotiMsgTyp_StationServiceEnabled
	NotiMsgTyp_StationStateChangeMsg
	NotiMsgTyp_StoryLineMissionAvailableMsg
	NotiMsgTyp_StructureAnchoring
	NotiMsgTyp_StructureCourierContractChanged
	NotiMsgTyp_StructureDestroyed
	NotiMsgTyp_StructureFuelAlert
	NotiMsgTyp_StructureImpendingAbandonmentAssetsAtRisk
	NotiMsgTyp_StructureItemsDelivered
	NotiMsgTyp_StructureItemsMovedToSafety
	NotiMsgTyp_StructureLostArmor
	NotiMsgTyp_StructureLostShields
	NotiMsgTyp_StructureOnline
	NotiMsgTyp_StructureServicesOffline
	NotiMsgTyp_StructureUnanchoring
	NotiMsgTyp_StructureUnderAttack
	NotiMsgTyp_StructureWentHighPower
	NotiMsgTyp_StructureWentLowPower
	NotiMsgTyp_StructuresJobsCancelled
	NotiMsgTyp_StructuresJobsPaused
	NotiMsgTyp_StructuresReinforcementChanged
	NotiMsgTyp_TowerAlertMsg
	NotiMsgTyp_TowerResourceAlertMsg
	NotiMsgTyp_TransactionReversalMsg
	NotiMsgTyp_TutorialMsg
	NotiMsgTyp_WarAdopted
	NotiMsgTyp_WarAllyInherited
	NotiMsgTyp_WarAllyOfferDeclinedMsg
	NotiMsgTyp_WarConcordInvalidates
	NotiMsgTyp_WarDeclared
	NotiMsgTyp_WarEndedHqSecurityDrop
	NotiMsgTyp_WarHQRemovedFromSpace
	NotiMsgTyp_WarInherited
	NotiMsgTyp_WarInvalid
	NotiMsgTyp_WarRetracted
	NotiMsgTyp_WarRetractedByConcord
	NotiMsgTyp_WarSurrenderDeclinedMsg
	NotiMsgTyp_WarSurrenderOfferMsg
)

var NotiSndTyp = map[string]int{
	"character":   NotiSndTyp_character,
	"corporation": NotiSndTyp_corporation,
	"alliance":    NotiSndTyp_alliance,
	"faction":     NotiSndTyp_faction,
	"other":       NotiSndTyp_other,
}

var NotiMsgTyp = map[string]int{
	"AcceptedAlly":                              NotiMsgTyp_AcceptedAlly,
	"AcceptedSurrender":                         NotiMsgTyp_AcceptedSurrender,
	"AgentRetiredTrigravian":                    NotiMsgTyp_AgentRetiredTrigravian,
	"AllAnchoringMsg":                           NotiMsgTyp_AllAnchoringMsg,
	"AllMaintenanceBillMsg":                     NotiMsgTyp_AllMaintenanceBillMsg,
	"AllStrucInvulnerableMsg":                   NotiMsgTyp_AllStrucInvulnerableMsg,
	"AllStructVulnerableMsg":                    NotiMsgTyp_AllStructVulnerableMsg,
	"AllWarCorpJoinedAllianceMsg":               NotiMsgTyp_AllWarCorpJoinedAllianceMsg,
	"AllWarDeclaredMsg":                         NotiMsgTyp_AllWarDeclaredMsg,
	"AllWarInvalidatedMsg":                      NotiMsgTyp_AllWarInvalidatedMsg,
	"AllWarRetractedMsg":                        NotiMsgTyp_AllWarRetractedMsg,
	"AllWarSurrenderMsg":                        NotiMsgTyp_AllWarSurrenderMsg,
	"AllianceCapitalChanged":                    NotiMsgTyp_AllianceCapitalChanged,
	"AllianceWarDeclaredV2":                     NotiMsgTyp_AllianceWarDeclaredV2,
	"AllyContractCancelled":                     NotiMsgTyp_AllyContractCancelled,
	"AllyJoinedWarAggressorMsg":                 NotiMsgTyp_AllyJoinedWarAggressorMsg,
	"AllyJoinedWarAllyMsg":                      NotiMsgTyp_AllyJoinedWarAllyMsg,
	"AllyJoinedWarDefenderMsg":                  NotiMsgTyp_AllyJoinedWarDefenderMsg,
	"BattlePunishFriendlyFire":                  NotiMsgTyp_BattlePunishFriendlyFire,
	"BillOutOfMoneyMsg":                         NotiMsgTyp_BillOutOfMoneyMsg,
	"BillPaidCorpAllMsg":                        NotiMsgTyp_BillPaidCorpAllMsg,
	"BountyClaimMsg":                            NotiMsgTyp_BountyClaimMsg,
	"BountyESSShared":                           NotiMsgTyp_BountyESSShared,
	"BountyESSTaken":                            NotiMsgTyp_BountyESSTaken,
	"BountyPlacedAlliance":                      NotiMsgTyp_BountyPlacedAlliance,
	"BountyPlacedChar":                          NotiMsgTyp_BountyPlacedChar,
	"BountyPlacedCorp":                          NotiMsgTyp_BountyPlacedCorp,
	"BountyYourBountyClaimed":                   NotiMsgTyp_BountyYourBountyClaimed,
	"BuddyConnectContactAdd":                    NotiMsgTyp_BuddyConnectContactAdd,
	"CharAppAcceptMsg":                          NotiMsgTyp_CharAppAcceptMsg,
	"CharAppRejectMsg":                          NotiMsgTyp_CharAppRejectMsg,
	"CharAppWithdrawMsg":                        NotiMsgTyp_CharAppWithdrawMsg,
	"CharLeftCorpMsg":                           NotiMsgTyp_CharLeftCorpMsg,
	"CharMedalMsg":                              NotiMsgTyp_CharMedalMsg,
	"CharTerminationMsg":                        NotiMsgTyp_CharTerminationMsg,
	"CloneActivationMsg":                        NotiMsgTyp_CloneActivationMsg,
	"CloneActivationMsg2":                       NotiMsgTyp_CloneActivationMsg2,
	"CloneMovedMsg":                             NotiMsgTyp_CloneMovedMsg,
	"CloneRevokedMsg1":                          NotiMsgTyp_CloneRevokedMsg1,
	"CloneRevokedMsg2":                          NotiMsgTyp_CloneRevokedMsg2,
	"CombatOperationFinished":                   NotiMsgTyp_CombatOperationFinished,
	"ContactAdd":                                NotiMsgTyp_ContactAdd,
	"ContactEdit":                               NotiMsgTyp_ContactEdit,
	"ContainerPasswordMsg":                      NotiMsgTyp_ContainerPasswordMsg,
	"ContractRegionChangedToPochven":            NotiMsgTyp_ContractRegionChangedToPochven,
	"CorpAllBillMsg":                            NotiMsgTyp_CorpAllBillMsg,
	"CorpAppAcceptMsg":                          NotiMsgTyp_CorpAppAcceptMsg,
	"CorpAppInvitedMsg":                         NotiMsgTyp_CorpAppInvitedMsg,
	"CorpAppNewMsg":                             NotiMsgTyp_CorpAppNewMsg,
	"CorpAppRejectCustomMsg":                    NotiMsgTyp_CorpAppRejectCustomMsg,
	"CorpAppRejectMsg":                          NotiMsgTyp_CorpAppRejectMsg,
	"CorpBecameWarEligible":                     NotiMsgTyp_CorpBecameWarEligible,
	"CorpDividendMsg":                           NotiMsgTyp_CorpDividendMsg,
	"CorpFriendlyFireDisableTimerCompleted":     NotiMsgTyp_CorpFriendlyFireDisableTimerCompleted,
	"CorpFriendlyFireDisableTimerStarted":       NotiMsgTyp_CorpFriendlyFireDisableTimerStarted,
	"CorpFriendlyFireEnableTimerCompleted":      NotiMsgTyp_CorpFriendlyFireEnableTimerCompleted,
	"CorpFriendlyFireEnableTimerStarted":        NotiMsgTyp_CorpFriendlyFireEnableTimerStarted,
	"CorpKicked":                                NotiMsgTyp_CorpKicked,
	"CorpLiquidationMsg":                        NotiMsgTyp_CorpLiquidationMsg,
	"CorpNewCEOMsg":                             NotiMsgTyp_CorpNewCEOMsg,
	"CorpNewsMsg":                               NotiMsgTyp_CorpNewsMsg,
	"CorpNoLongerWarEligible":                   NotiMsgTyp_CorpNoLongerWarEligible,
	"CorpOfficeExpirationMsg":                   NotiMsgTyp_CorpOfficeExpirationMsg,
	"CorpStructLostMsg":                         NotiMsgTyp_CorpStructLostMsg,
	"CorpTaxChangeMsg":                          NotiMsgTyp_CorpTaxChangeMsg,
	"CorpVoteCEORevokedMsg":                     NotiMsgTyp_CorpVoteCEORevokedMsg,
	"CorpVoteMsg":                               NotiMsgTyp_CorpVoteMsg,
	"CorpWarDeclaredMsg":                        NotiMsgTyp_CorpWarDeclaredMsg,
	"CorpWarDeclaredV2":                         NotiMsgTyp_CorpWarDeclaredV2,
	"CorpWarFightingLegalMsg":                   NotiMsgTyp_CorpWarFightingLegalMsg,
	"CorpWarInvalidatedMsg":                     NotiMsgTyp_CorpWarInvalidatedMsg,
	"CorpWarRetractedMsg":                       NotiMsgTyp_CorpWarRetractedMsg,
	"CorpWarSurrenderMsg":                       NotiMsgTyp_CorpWarSurrenderMsg,
	"CustomsMsg":                                NotiMsgTyp_CustomsMsg,
	"DeclareWar":                                NotiMsgTyp_DeclareWar,
	"DistrictAttacked":                          NotiMsgTyp_DistrictAttacked,
	"DustAppAcceptedMsg":                        NotiMsgTyp_DustAppAcceptedMsg,
	"ESSMainBankLink":                           NotiMsgTyp_ESSMainBankLink,
	"EntosisCaptureStarted":                     NotiMsgTyp_EntosisCaptureStarted,
	"FWAllianceKickMsg":                         NotiMsgTyp_FWAllianceKickMsg,
	"FWAllianceWarningMsg":                      NotiMsgTyp_FWAllianceWarningMsg,
	"FWCharKickMsg":                             NotiMsgTyp_FWCharKickMsg,
	"FWCharRankGainMsg":                         NotiMsgTyp_FWCharRankGainMsg,
	"FWCharRankLossMsg":                         NotiMsgTyp_FWCharRankLossMsg,
	"FWCharWarningMsg":                          NotiMsgTyp_FWCharWarningMsg,
	"FWCorpJoinMsg":                             NotiMsgTyp_FWCorpJoinMsg,
	"FWCorpKickMsg":                             NotiMsgTyp_FWCorpKickMsg,
	"FWCorpLeaveMsg":                            NotiMsgTyp_FWCorpLeaveMsg,
	"FWCorpWarningMsg":                          NotiMsgTyp_FWCorpWarningMsg,
	"FacWarCorpJoinRequestMsg":                  NotiMsgTyp_FacWarCorpJoinRequestMsg,
	"FacWarCorpJoinWithdrawMsg":                 NotiMsgTyp_FacWarCorpJoinWithdrawMsg,
	"FacWarCorpLeaveRequestMsg":                 NotiMsgTyp_FacWarCorpLeaveRequestMsg,
	"FacWarCorpLeaveWithdrawMsg":                NotiMsgTyp_FacWarCorpLeaveWithdrawMsg,
	"FacWarLPDisqualifiedEvent":                 NotiMsgTyp_FacWarLPDisqualifiedEvent,
	"FacWarLPDisqualifiedKill":                  NotiMsgTyp_FacWarLPDisqualifiedKill,
	"FacWarLPPayoutEvent":                       NotiMsgTyp_FacWarLPPayoutEvent,
	"FacWarLPPayoutKill":                        NotiMsgTyp_FacWarLPPayoutKill,
	"GameTimeAdded":                             NotiMsgTyp_GameTimeAdded,
	"GameTimeReceived":                          NotiMsgTyp_GameTimeReceived,
	"GameTimeSent":                              NotiMsgTyp_GameTimeSent,
	"GiftReceived":                              NotiMsgTyp_GiftReceived,
	"IHubDestroyedByBillFailure":                NotiMsgTyp_IHubDestroyedByBillFailure,
	"IncursionCompletedMsg":                     NotiMsgTyp_IncursionCompletedMsg,
	"IndustryOperationFinished":                 NotiMsgTyp_IndustryOperationFinished,
	"IndustryTeamAuctionLost":                   NotiMsgTyp_IndustryTeamAuctionLost,
	"IndustryTeamAuctionWon":                    NotiMsgTyp_IndustryTeamAuctionWon,
	"InfrastructureHubBillAboutToExpire":        NotiMsgTyp_InfrastructureHubBillAboutToExpire,
	"InsuranceExpirationMsg":                    NotiMsgTyp_InsuranceExpirationMsg,
	"InsuranceFirstShipMsg":                     NotiMsgTyp_InsuranceFirstShipMsg,
	"InsuranceInvalidatedMsg":                   NotiMsgTyp_InsuranceInvalidatedMsg,
	"InsuranceIssuedMsg":                        NotiMsgTyp_InsuranceIssuedMsg,
	"InsurancePayoutMsg":                        NotiMsgTyp_InsurancePayoutMsg,
	"InvasionCompletedMsg":                      NotiMsgTyp_InvasionCompletedMsg,
	"InvasionSystemLogin":                       NotiMsgTyp_InvasionSystemLogin,
	"InvasionSystemStart":                       NotiMsgTyp_InvasionSystemStart,
	"JumpCloneDeletedMsg1":                      NotiMsgTyp_JumpCloneDeletedMsg1,
	"JumpCloneDeletedMsg2":                      NotiMsgTyp_JumpCloneDeletedMsg2,
	"KillReportFinalBlow":                       NotiMsgTyp_KillReportFinalBlow,
	"KillReportVictim":                          NotiMsgTyp_KillReportVictim,
	"KillRightAvailable":                        NotiMsgTyp_KillRightAvailable,
	"KillRightAvailableOpen":                    NotiMsgTyp_KillRightAvailableOpen,
	"KillRightEarned":                           NotiMsgTyp_KillRightEarned,
	"KillRightUnavailable":                      NotiMsgTyp_KillRightUnavailable,
	"KillRightUnavailableOpen":                  NotiMsgTyp_KillRightUnavailableOpen,
	"KillRightUsed":                             NotiMsgTyp_KillRightUsed,
	"LocateCharMsg":                             NotiMsgTyp_LocateCharMsg,
	"MadeWarMutual":                             NotiMsgTyp_MadeWarMutual,
	"MercOfferRetractedMsg":                     NotiMsgTyp_MercOfferRetractedMsg,
	"MercOfferedNegotiationMsg":                 NotiMsgTyp_MercOfferedNegotiationMsg,
	"MissionCanceledTriglavian":                 NotiMsgTyp_MissionCanceledTriglavian,
	"MissionOfferExpirationMsg":                 NotiMsgTyp_MissionOfferExpirationMsg,
	"MissionTimeoutMsg":                         NotiMsgTyp_MissionTimeoutMsg,
	"MoonminingAutomaticFracture":               NotiMsgTyp_MoonminingAutomaticFracture,
	"MoonminingExtractionCancelled":             NotiMsgTyp_MoonminingExtractionCancelled,
	"MoonminingExtractionFinished":              NotiMsgTyp_MoonminingExtractionFinished,
	"MoonminingExtractionStarted":               NotiMsgTyp_MoonminingExtractionStarted,
	"MoonminingLaserFired":                      NotiMsgTyp_MoonminingLaserFired,
	"MutualWarExpired":                          NotiMsgTyp_MutualWarExpired,
	"MutualWarInviteAccepted":                   NotiMsgTyp_MutualWarInviteAccepted,
	"MutualWarInviteRejected":                   NotiMsgTyp_MutualWarInviteRejected,
	"MutualWarInviteSent":                       NotiMsgTyp_MutualWarInviteSent,
	"NPCStandingsGained":                        NotiMsgTyp_NPCStandingsGained,
	"NPCStandingsLost":                          NotiMsgTyp_NPCStandingsLost,
	"OfferToAllyRetracted":                      NotiMsgTyp_OfferToAllyRetracted,
	"OfferedSurrender":                          NotiMsgTyp_OfferedSurrender,
	"OfferedToAlly":                             NotiMsgTyp_OfferedToAlly,
	"OfficeLeaseCanceledInsufficientStandings":  NotiMsgTyp_OfficeLeaseCanceledInsufficientStandings,
	"OldLscMessages":                            NotiMsgTyp_OldLscMessages,
	"OperationFinished":                         NotiMsgTyp_OperationFinished,
	"OrbitalAttacked":                           NotiMsgTyp_OrbitalAttacked,
	"OrbitalReinforced":                         NotiMsgTyp_OrbitalReinforced,
	"OwnershipTransferred":                      NotiMsgTyp_OwnershipTransferred,
	"RaffleCreated":                             NotiMsgTyp_RaffleCreated,
	"RaffleExpired":                             NotiMsgTyp_RaffleExpired,
	"RaffleFinished":                            NotiMsgTyp_RaffleFinished,
	"ReimbursementMsg":                          NotiMsgTyp_ReimbursementMsg,
	"ResearchMissionAvailableMsg":               NotiMsgTyp_ResearchMissionAvailableMsg,
	"RetractsWar":                               NotiMsgTyp_RetractsWar,
	"SeasonalChallengeCompleted":                NotiMsgTyp_SeasonalChallengeCompleted,
	"SovAllClaimAquiredMsg":                     NotiMsgTyp_SovAllClaimAquiredMsg,
	"SovAllClaimLostMsg":                        NotiMsgTyp_SovAllClaimLostMsg,
	"SovCommandNodeEventStarted":                NotiMsgTyp_SovCommandNodeEventStarted,
	"SovCorpBillLateMsg":                        NotiMsgTyp_SovCorpBillLateMsg,
	"SovCorpClaimFailMsg":                       NotiMsgTyp_SovCorpClaimFailMsg,
	"SovDisruptorMsg":                           NotiMsgTyp_SovDisruptorMsg,
	"SovStationEnteredFreeport":                 NotiMsgTyp_SovStationEnteredFreeport,
	"SovStructureDestroyed":                     NotiMsgTyp_SovStructureDestroyed,
	"SovStructureReinforced":                    NotiMsgTyp_SovStructureReinforced,
	"SovStructureSelfDestructCancel":            NotiMsgTyp_SovStructureSelfDestructCancel,
	"SovStructureSelfDestructFinished":          NotiMsgTyp_SovStructureSelfDestructFinished,
	"SovStructureSelfDestructRequested":         NotiMsgTyp_SovStructureSelfDestructRequested,
	"SovereigntyIHDamageMsg":                    NotiMsgTyp_SovereigntyIHDamageMsg,
	"SovereigntySBUDamageMsg":                   NotiMsgTyp_SovereigntySBUDamageMsg,
	"SovereigntyTCUDamageMsg":                   NotiMsgTyp_SovereigntyTCUDamageMsg,
	"StationAggressionMsg1":                     NotiMsgTyp_StationAggressionMsg1,
	"StationAggressionMsg2":                     NotiMsgTyp_StationAggressionMsg2,
	"StationConquerMsg":                         NotiMsgTyp_StationConquerMsg,
	"StationServiceDisabled":                    NotiMsgTyp_StationServiceDisabled,
	"StationServiceEnabled":                     NotiMsgTyp_StationServiceEnabled,
	"StationStateChangeMsg":                     NotiMsgTyp_StationStateChangeMsg,
	"StoryLineMissionAvailableMsg":              NotiMsgTyp_StoryLineMissionAvailableMsg,
	"StructureAnchoring":                        NotiMsgTyp_StructureAnchoring,
	"StructureCourierContractChanged":           NotiMsgTyp_StructureCourierContractChanged,
	"StructureDestroyed":                        NotiMsgTyp_StructureDestroyed,
	"StructureFuelAlert":                        NotiMsgTyp_StructureFuelAlert,
	"StructureImpendingAbandonmentAssetsAtRisk": NotiMsgTyp_StructureImpendingAbandonmentAssetsAtRisk,
	"StructureItemsDelivered":                   NotiMsgTyp_StructureItemsDelivered,
	"StructureItemsMovedToSafety":               NotiMsgTyp_StructureItemsMovedToSafety,
	"StructureLostArmor":                        NotiMsgTyp_StructureLostArmor,
	"StructureLostShields":                      NotiMsgTyp_StructureLostShields,
	"StructureOnline":                           NotiMsgTyp_StructureOnline,
	"StructureServicesOffline":                  NotiMsgTyp_StructureServicesOffline,
	"StructureUnanchoring":                      NotiMsgTyp_StructureUnanchoring,
	"StructureUnderAttack":                      NotiMsgTyp_StructureUnderAttack,
	"StructureWentHighPower":                    NotiMsgTyp_StructureWentHighPower,
	"StructureWentLowPower":                     NotiMsgTyp_StructureWentLowPower,
	"StructuresJobsCancelled":                   NotiMsgTyp_StructuresJobsCancelled,
	"StructuresJobsPaused":                      NotiMsgTyp_StructuresJobsPaused,
	"StructuresReinforcementChanged":            NotiMsgTyp_StructuresReinforcementChanged,
	"TowerAlertMsg":                             NotiMsgTyp_TowerAlertMsg,
	"TowerResourceAlertMsg":                     NotiMsgTyp_TowerResourceAlertMsg,
	"TransactionReversalMsg":                    NotiMsgTyp_TransactionReversalMsg,
	"TutorialMsg":                               NotiMsgTyp_TutorialMsg,
	"WarAdopted ":                               NotiMsgTyp_WarAdopted,
	"WarAllyInherited":                          NotiMsgTyp_WarAllyInherited,
	"WarAllyOfferDeclinedMsg":                   NotiMsgTyp_WarAllyOfferDeclinedMsg,
	"WarConcordInvalidates":                     NotiMsgTyp_WarConcordInvalidates,
	"WarDeclared":                               NotiMsgTyp_WarDeclared,
	"WarEndedHqSecurityDrop":                    NotiMsgTyp_WarEndedHqSecurityDrop,
	"WarHQRemovedFromSpace":                     NotiMsgTyp_WarHQRemovedFromSpace,
	"WarInherited":                              NotiMsgTyp_WarInherited,
	"WarInvalid":                                NotiMsgTyp_WarInvalid,
	"WarRetracted":                              NotiMsgTyp_WarRetracted,
	"WarRetractedByConcord":                     NotiMsgTyp_WarRetractedByConcord,
	"WarSurrenderDeclinedMsg":                   NotiMsgTyp_WarSurrenderDeclinedMsg,
	"WarSurrenderOfferMsg":                      NotiMsgTyp_WarSurrenderOfferMsg,
}

type DBNotification struct {
	CharId         int
	IsRead         int
	NotificationId int64
	SenderId       int32
	SenderType     int
	TextRef        int
	TimeStamp      int64
	Type           int
}

func (obj *Model) createNotificationTable() {
	if !obj.checkTableExists("notifications") {
		_, err := obj.DB.Exec(`
		CREATE TABLE "notifications" (
			"charId" INT,
			"is_read" INT,
			"notification_id" INT,
			"sender_id" INT,
			"sender_type" INT,
			"text" INT,
			"timestamp" INT,
			"type" INT
		);`)
		util.CheckErr(err)
	}
}

func (obj *Model) AddNotificationEntry(notiItem *DBNotification) DBresult {
	whereClause := fmt.Sprintf(`notification_id="%d"`, notiItem.NotificationId)
	retval := DBR_Undefined
	num := obj.getNumEntries("notifications", whereClause)
	if num == 0 {
		stmt, err := obj.DB.Prepare(`
			INSERT INTO "notifications" (
			    charId,
				is_read,
				notification_id,
				sender_id,
				sender_type,
				text,
				timestamp,
				type)
			VALUES (?,?,?,?,?,?,?,?);`)
		util.CheckErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(
			notiItem.CharId,
			notiItem.IsRead,
			notiItem.NotificationId,
			notiItem.SenderId,
			notiItem.SenderType,
			notiItem.TextRef,
			notiItem.TimeStamp,
			notiItem.Type)
		util.CheckErr(err)
		affect, err := res.RowsAffected()
		util.CheckErr(err)
		if affect > 0 {
			retval = DBR_Inserted
		}
	}
	return retval
}

func (obj *Model) GetCharNotifications(charID int) (retval []*DBNotification) {
	retval = make([]*DBNotification, 0, 5)
	queryString := fmt.Sprintf(`
		SELECT
			charId,
			is_read,
			notification_id,
			sender_id,
			sender_type,
			text,
			timestamp,
			type
		FROM notifications 
		WHERE charId=%d
		ORDER BY timestamp DESC;`, charID)
	rows, err := obj.DB.Query(queryString)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var newNoti DBNotification
		rows.Scan(
			&newNoti.CharId,
			&newNoti.IsRead,
			&newNoti.NotificationId,
			&newNoti.SenderId,
			&newNoti.SenderType,
			&newNoti.TextRef,
			&newNoti.TimeStamp,
			&newNoti.Type)
		retval = append(retval, &newNoti)
	}
	return retval
}
