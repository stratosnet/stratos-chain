package types

const (
	EventTypeCompleteUnbondingResourceNode = "complete_unbonding_resource_node"
	EventTypeCompleteUnbondingMetaNode     = "complete_unbonding_meta_node"

	EventTypeCreateResourceNode                  = "create_resource_node"
	EventTypeUnbondingResourceNode               = "unbonding_resource_node"
	EventTypeUpdateResourceNode                  = "update_resource_node"
	EventTypeUpdateResourceNodeDeposit           = "update_resource_node_deposit"
	EventTypeUpdateEffectiveDeposit              = "update_effective_deposit"
	EventTypeCreateMetaNode                      = "create_meta_node"
	EventTypeUnbondingMetaNode                   = "unbonding_Meta_node"
	EventTypeUpdateMetaNode                      = "update_meta_node"
	EventTypeUpdateMetaNodeDeposit               = "update_meta_node_deposit"
	EventTypeMetaNodeRegistrationVote            = "meta_node_reg_vote"
	EventTypeWithdrawMetaNodeRegistrationDeposit = "withdraw_meta_node_reg_deposit"

	AttributeKeyResourceNode            = "resource_node"
	AttributeKeyMetaNode                = "meta_node"
	AttributeKeyNetworkAddress          = "network_address"
	AttributeKeyPubKey                  = "pub_key"
	AttributeKeyCandidateNetworkAddress = "candidate_network_address"
	AttributeKeyVoterNetworkAddress     = "voter_network_address"
	AttributeKeyCandidateStatus         = "candidate_status"

	AttributeKeyUnbondingMatureTime = "unbonding_mature_time"

	AttributeKeyOZoneLimitChanges     = "ozone_limit_changes"
	AttributeKeyInitialDeposit        = "initial_deposit"
	AttributeKeyCurrentDeposit        = "current_deposit"
	AttributeKeyAvailableTokenBefore  = "available_token_before"
	AttributeKeyAvailableTokenAfter   = "available_token_after"
	AttributeKeyDepositDelta          = "deposit_delta"
	AttributeKeyDepositToRemove       = "deposit_to_remove"
	AttributeKeyIncrDeposit           = "incr_deposit"
	AttributeKeyEffectiveDepositAfter = "effective_deposit_after"
	AttributeKeyIsUnsuspended         = "is_unsuspended"

	AttributeValueCategory = ModuleName
)
