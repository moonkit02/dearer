package datatype

import "github.com/moonkit02/dearer/pkg/parser"

// IsParentedByNodeID checks if her or any of her parents have a given nodeID
func IsParentedByNodeID(nodeID parser.NodeID, node *parser.Node) bool {
	if node == nil {
		return false
	}
	if nodeID == node.ID() {
		return true
	}

	return IsParentedByNodeID(nodeID, node.Parent())
}
