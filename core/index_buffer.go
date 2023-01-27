package core

import "github.com/go-gl/gl/v3.3-core/gl"

type IndexBuffer struct {
	m_RendererID uint32
	m_Count      int
}

func NewIndexBuffer(data []uint32) IndexBuffer {
	var ib uint32
	gl.GenBuffers(1, &ib)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	return IndexBuffer{m_RendererID: ib, m_Count: len(data)}
}

func (ib *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.m_RendererID)
}

func (ib *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBuffer) GetCount() int {
	return ib.m_Count
}

