package types

// register module event types
const (
	EventTypeCreateResourceNode = "create_resource_node"
	EventTypeCreateIndexingNode = "create_indexing_node"

	AttributeKeyResourceNode = "resource_node"
	AttributeKeyIndexingNode = "indexing_node"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
