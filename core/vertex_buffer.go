package core

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type VertexBuffer struct {
	m_RendererID uint32
}

func NewVertexBuffer(data []float32) *VertexBuffer {
	var vb uint32
	gl.GenBuffers(1, &vb)
	gl.BindBuffer(gl.ARRAY_BUFFER, vb)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(&(data[0])), gl.STATIC_DRAW)
	return &VertexBuffer{m_RendererID: vb}
}

func (vb *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.m_RendererID)
}

func (vb *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
