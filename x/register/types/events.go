package types

const (
	EventTypeCompleteUnbondingResourceNode = "complete_unbonding_resource_node"
	EventTypeCompleteUnbondingIndexingNode = "complete_unbonding_indexing_node"
	EventTypeUnbondNode                    = "unbond_node"

	EventTypeCreateResourceNode           = "create_resource_node"
	EventTypeRemoveResourceNode           = "remove_resource_node"
	EventTypeUnbondingResourceNode        = "unbonding_resource_node"
	EventTypeUpdateResourceNode           = "update_resource_node"
	EventTypeUpdateResourceNodeStake      = "update_resource_node_stake"
	EventTypeCreateIndexingNode           = "create_indexing_node"
	EventTypeRemoveIndexingNode           = "remove_indexing_node"
	EventTypeUnbondingIndexingNode        = "unbonding_indexing_node"
	EventTypeUpdateIndexingNode           = "update_indexing_node"
	EventTypeUpdateIndexingNodeStake      = "update_indexing_node_stake"
	EventTypeIndexingNodeRegistrationVote = "indexing_node_reg_vote"

	AttributeKeyResourceNode            = "resource_node"
	AttributeKeyIndexingNode            = "indexing_node"
	AttributeKeyNetworkAddress          = "network_address"
	AttributeKeyPubKey                  = "pub_key"
	AttributeKeyCandidateNetworkAddress = "candidate_network_address"
	AttributeKeyVoterNetworkAddress     = "voter_network_address"
	AttributeKeyCandidateStatus         = "candidate_status"
	AttributeKeyNetworkAddr             = "network_addr"
	AttributeKeyIsIndexingNode          = "is_indexing_node"

	AttributeKeyUnbondingMatureTime = "unbonding_mature_time"

	AttributeKeyOZoneLimitChanges = "ozone_limit_changes"
	AttributeKeyInitialStake      = "initial_stake"
	AttributeKeyStakeDelta        = "stake_delta"
	AttributeKeyStakeToRemove     = "stake_to_remove"
	AttributeKeyIncrStakeBool     = "incr_stake"

	AttributeValueCategory = ModuleName
)
