package schema

import (
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
)

type UUIDHolder struct {
	UUID map[parser.NodeID]string
}

func (holder *UUIDHolder) Assign(nodeID parser.NodeID, generator nodeid.Generator) string {
	val, ok := holder.UUID[nodeID]
	if ok {
		return val
	}

	newUUID := generator.GenerateId()

	holder.UUID[nodeID] = newUUID

	return newUUID
}

func NewUUIDHolder() *UUIDHolder {
	return &UUIDHolder{
		UUID: make(map[parser.NodeID]string),
	}
}
