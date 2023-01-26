package core

import (
	"encoding/binary"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type VertexBufferElement struct {
	m_Type       int
	m_Count      int
	m_Normalized bool
}

func NewVertexBufferElement(typ int, count int, normalized bool) *VertexBufferElement {
	return &VertexBufferElement{
		m_Type:       typ,
		m_Count:      count,
		m_Normalized: normalized,
	}
}

func (vbl *VertexBufferElement) GetTypeSize() int {
	switch vbl.m_Type {
	case gl.FLOAT:
		return binary.Size(float32(0))
	case gl.UNSIGNED_INT:
		return binary.Size(uint32(0))
	default:
		return 0
	}
}
