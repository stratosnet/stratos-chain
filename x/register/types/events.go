package types

const (
	EventTypeCompleteUnbondingResourceNode = "complete_unbonding_resource_node"
	EventTypeCompleteUnbondingMetaNode     = "complete_unbonding_meta_node"

	EventTypeCreateResourceNode                = "create_resource_node"
	EventTypeUnbondingResourceNode             = "unbonding_resource_node"
	EventTypeUpdateResourceNode                = "update_resource_node"
	EventTypeUpdateResourceNodeStake           = "update_resource_node_stake"
	EventTypeUpdateEffectiveStake              = "update_effective_stake"
	EventTypeCreateMetaNode                    = "create_meta_node"
	EventTypeUnbondingMetaNode                 = "unbonding_Meta_node"
	EventTypeUpdateMetaNode                    = "update_meta_node"
	EventTypeUpdateMetaNodeStake               = "update_meta_node_stake"
	EventTypeMetaNodeRegistrationVote          = "meta_node_reg_vote"
	EventTypeWithdrawMetaNodeRegistrationStake = "withdraw_meta_node_reg_stake"

	AttributeKeyResourceNode            = "resource_node"
	AttributeKeyMetaNode                = "meta_node"
	AttributeKeyNetworkAddress          = "network_address"
	AttributeKeyPubKey                  = "pub_key"
	AttributeKeyCandidateNetworkAddress = "candidate_network_address"
	AttributeKeyVoterNetworkAddress     = "voter_network_address"
	AttributeKeyCandidateStatus         = "candidate_status"

	AttributeKeyUnbondingMatureTime = "unbonding_mature_time"

	AttributeKeyOZoneLimitChanges    = "ozone_limit_changes"
	AttributeKeyInitialStake         = "initial_stake"
	AttributeKeyCurrentStake         = "current_stake"
	AttributeKeyAvailableTokenBefore = "available_token_before"
	AttributeKeyAvailableTokenAfter  = "available_token_after"
	AttributeKeyStakeDelta           = "stake_delta"
	AttributeKeyStakeToRemove        = "stake_to_remove"
	AttributeKeyIncrStake            = "incr_stake"
	AttributeKeyEffectiveStakeAfter  = "effective_stake_after"
	AttributeKeyIsUnsuspended        = "is_unsuspended"

	AttributeValueCategory = ModuleName
)
