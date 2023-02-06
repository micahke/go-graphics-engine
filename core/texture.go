package core

import (
	"fmt"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	m_RendererID      uint32
	m_FilePath        string
	m_LocalBuffer     []uint8
	m_Width, m_Height int32
	m_BPP             int
}

func NewTexture(path string) *Texture {
	texture := &Texture{}
	texture.m_RendererID = 0
	texture.m_FilePath = path
	texture.m_BPP = 0

	file, err := os.Open("res/" + path)
	if err != nil {
		fmt.Println("Error reading file")
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding file")
	}

	flippedImage := imaging.FlipV(img)
	texture.m_Height = int32(flippedImage.Rect.Dy())
	texture.m_Width = int32(flippedImage.Rect.Dx())
	texture.m_LocalBuffer = flippedImage.Pix

	gl.GenTextures(1, &texture.m_RendererID)
	gl.BindTexture(gl.TEXTURE_2D, texture.m_RendererID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, texture.m_Width, texture.m_Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&texture.m_LocalBuffer[0]))
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture
}

func (texture *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, texture.m_RendererID)
}

func (texture *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (texture *Texture) GetWidth() int32 {
	return texture.m_Width
}

func (texture *Texture) GetHeight() int32 {
	return texture.m_Height
}
