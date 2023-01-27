package core

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type VertexArray struct {
	m_RendererID uint32
}

func NewVertexArray() VertexArray {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return VertexArray{m_RendererID: vao}

}

func (vao *VertexArray) AddBuffer(vb VertexBuffer, layout VertexBufferLayout) {
	vao.Bind()
	vb.Bind()
	elements := layout.GetElements()
	offset := 0
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointer(uint32(i), int32(element.m_Count), uint32(element.m_Type), element.m_Normalized, int32(layout.GetStride()), gl.PtrOffset(offset))
		offset += element.GetTypeSize() * element.m_Count
	}
}

func (vao *VertexArray) Bind() {
	gl.BindVertexArray(vao.m_RendererID)
}

func (vao *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}
