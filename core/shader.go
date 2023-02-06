package core

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
)

type Shader struct {
	m_VertexPath           string
	m_FragmentPath         string
	m_Renderer_ID          uint32
	m_UniformLocationCache map[string]int32
}

func NewShader(vertexPath string, fragmentPath string) Shader {
	shader := Shader{}
	vertexShader := shader.ParseShader(vertexPath)
	fragmentShader := shader.ParseShader(fragmentPath)

	shader.CreateShader(vertexShader, fragmentShader)

	shader.m_VertexPath = vertexPath
	shader.m_FragmentPath = fragmentPath
	shader.m_UniformLocationCache = make(map[string]int32)

	return shader

}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.m_Renderer_ID)
}

func (shader *Shader) Unbind() {
	gl.UseProgram(0)
}

func (shader *Shader) ParseShader(shaderPath string) string {
	contents, err := ioutil.ReadFile("res/" + shaderPath)
	if err != nil {
		panic("There was an error parsing the shader")
	}
	source := string(contents) + "\x00"
	return source
}

func (shader *Shader) CompileShader(source string, shaderType uint32) uint32 {
	shaderID := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shaderID, 1, csource, nil)
	free()
	gl.CompileShader(shaderID)

	// Error handling
	var status int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		var infoLog [512]byte
		gl.GetShaderInfoLog(shaderID, 512, &length, &infoLog[0])
		message := string(infoLog[:length])
		s_type := "fragment"
		if shaderType == gl.VERTEX_SHADER {
			s_type = "vertex"
		}
		fmt.Println("Fauled to compile ", s_type, " shader!")
		fmt.Println(message)
		gl.DeleteShader(shaderID)
		return 0
	}
	return shaderID
}

func (shader *Shader) CreateShader(vertexShader string, fragmentShader string) {
	program := gl.CreateProgram()
	vs := shader.CompileShader(vertexShader, gl.VERTEX_SHADER)
	fs := shader.CompileShader(fragmentShader, gl.FRAGMENT_SHADER)

	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	shader.m_Renderer_ID = program
}


func (shader *Shader) SetUniform1i(name string, value int32) {
	// location := gl.GetUniformLocation(shader.m_Renderer_ID, gl.Str(name+"\x00"))
	gl.Uniform1i(shader.GetUniformLocation(name), value)
}

func (shader *Shader) SetUniform1f(name string, value float32) {
	// location := gl.GetUniformLocation(shader.m_Renderer_ID, gl.Str(name+"\x00"))
	gl.Uniform1f(shader.GetUniformLocation(name), value)
}

func (shader *Shader) SetUniform4f(name string, v0 float32, v1 float32, v2 float32, v3 float32) {
	// location := gl.GetUniformLocation(shader.m_Renderer_ID, gl.Str(name+"\x00"))
	gl.Uniform4f(shader.GetUniformLocation(name), v0, v1, v2, v3)
}

func (shader *Shader) SetUniformMat4f(name string, matrix mgl32.Mat4) {
	gl.UniformMatrix4fv(shader.GetUniformLocation(name), 1, false, &matrix[0])
}

func (shader *Shader) SetUniform3f(name string, v0 float32, v1 float32, v2 float32) {
  gl.Uniform3f(shader.GetUniformLocation(name), v0, v1, v2)
}

func (shader *Shader) GetUniformLocation(name string) int32 {
	// Check the cache
	if value, ok := shader.m_UniformLocationCache[name]; ok {
		return value
	}
	location := gl.GetUniformLocation(shader.m_Renderer_ID, gl.Str(name+"\x00"))
	if location == -1 {
		fmt.Println("-1 location for uniform: ", name)
	}
	shader.m_UniformLocationCache[name] = location
	return location
}
