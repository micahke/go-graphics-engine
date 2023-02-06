package core

import "github.com/go-gl/gl/v3.3-core/gl"



type VertexBufferLayout struct {
    m_Elements []VertexBufferElement
    m_Stride   int
}

func NewVertexBufferLayout() *VertexBufferLayout {
  return &VertexBufferLayout{}
}

func (vbl *VertexBufferLayout) Pushf(count int) {
    vbl.m_Elements = append(vbl.m_Elements, VertexBufferElement{gl.FLOAT, count, false})
    vbl.m_Stride += 4 * count
}

func (vbl *VertexBufferLayout) Pushi(count int) {
    vbl.m_Elements = append(vbl.m_Elements, VertexBufferElement{gl.UNSIGNED_INT, count, false})
    vbl.m_Stride += 4 * count
}

func (vbl *VertexBufferLayout) Pushc(count int) {
    vbl.m_Elements = append(vbl.m_Elements, VertexBufferElement{gl.UNSIGNED_BYTE, count, true})
    vbl.m_Stride += 1 * count
}

func (vbl *VertexBufferLayout) GetElements() []VertexBufferElement {
    return vbl.m_Elements
}

func (vbl *VertexBufferLayout) GetStride() int {
    return vbl.m_Stride
}

