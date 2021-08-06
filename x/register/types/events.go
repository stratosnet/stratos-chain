package types

const (
	EventTypeCreateResourceNode           = "create_resource_node"
	EventTypeRemoveResourceNode           = "remove_resource_node"
	EventTypeUpdateResourceNode           = "update_resource_node"
	EventTypeCreateIndexingNode           = "create_indexing_node"
	EventTypeRemoveIndexingNode           = "remove_indexing_node"
	EventTypeUpdateIndexingNode           = "update_indexing_node"
	EventTypeIndexingNodeRegistrationVote = "indexing_node_reg_vote"

	AttributeKeyResourceNode = "resource_node"
	AttributeKeyIndexingNode = "indexing_node"
	AttributeKeyNodeAddress  = "node_address"
	AttributeKeyNodePubkey   = "node_pubkey"

	AttributeValueCategory = ModuleName
)
