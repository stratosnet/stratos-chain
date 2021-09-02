package types

const (
	EventTypeCreateResourceNode           = "create_resource_node"
	EventTypeRemoveResourceNode           = "remove_resource_node"
	EventTypeUpdateResourceNode           = "update_resource_node"
	EventTypeCreateIndexingNode           = "create_indexing_node"
	EventTypeRemoveIndexingNode           = "remove_indexing_node"
	EventTypeUpdateIndexingNode           = "update_indexing_node"
	EventTypeIndexingNodeRegistrationVote = "indexing_node_reg_vote"

	AttributeKeyResourceNode            = "resource_node"
	AttributeKeyIndexingNode            = "indexing_node"
	AttributeKeyNetworkAddress          = "network_address"
	AttributeKeyPubKey                  = "pub_key"
	AttributeKeyCandidateNetworkAddress = "candidate_network_address"
	AttributeKeyCandidateStatus         = "candidate_status"

	AttributeKeyOZoneLimitIncrease = "ozone_limit_increase"
	AttributeKeyOZoneLimitDecrease = "ozone_limit_decrease"

	AttributeValueCategory = ModuleName
)
