package core

import "github.com/go-gl/gl/v3.3-core/gl"

type Renderer struct {
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (renderer *Renderer) Draw(va VertexArray, shader Shader) {
	shader.Bind()
	va.Bind()
  gl.DrawArrays(gl.TRIANGLES, 0, 36)
	// gl.DrawElements(gl.TRIANGLES, int32(ib.GetCount()), gl.UNSIGNED_INT, nil)
}
