package main

import (
	. "micahke/go-graphics-engine/core"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	glm "github.com/go-gl/mathgl/mgl32"
)

var l_CameraPos glm.Vec3 = glm.Vec3{0.0, 0.0, 3.0}
var l_CameraFront glm.Vec3 = glm.Vec3{0.0, 0.0, -1.0}
var l_CameraUp glm.Vec3 = glm.Vec3{0.0, 1.0, 0.0}

var l_DeltaTime float32 = 0.0
var l_LastFrame float32 = 0.0

var l_Camera *Camera

// Basically the main() function for the lighting section
func RunLighting() {

	if err := glfw.Init(); err != nil {
		panic("Error initializing GLFW")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var window, win_err = glfw.CreateWindow(960, 540, "Hello, world!", nil, nil)
	if win_err != nil {
		panic("Error creating window")
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic("Error initializing OpenGL")
	}

	positions := GetCubePositions()

	indeces := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)

	vao := NewVertexArray()
	vbo := NewVertexBuffer(*positions)
	vbl := NewVertexBufferLayout()

	vbl.Pushf(3)
	vbl.Pushf(2)
	vao.AddBuffer(*vbo, *vbl)

	ibo := NewIndexBuffer(indeces)

	shader := NewShader("vertexShader.glsl", "fragmentShader.glsl")
	shader.Bind()
	shader.SetUniform1i("u_Texture", 0)

	texture := NewTexture("wall.png")
	texture.Bind(0)
	shader.SetUniform1i("u_Texture", 0)

	renderer := NewRenderer()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(l_mouse_callback)
	window.SetScrollCallback(l_scroll_callback)

	l_Camera = NewCamera(l_CameraPos, l_CameraFront, l_CameraUp)

	for !window.ShouldClose() {
		var currentFrame float32 = float32(glfw.GetTime())
		l_DeltaTime = currentFrame - l_LastFrame
		l_LastFrame = currentFrame

		l_process_input(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		l_Camera.Update(&shader)

		model := glm.Ident4()
		shader.SetUniformMat4f("model", model)

		renderer.Draw(*vao, ibo, shader)

		window.SwapBuffers()
		glfw.PollEvents()

	}

	glfw.Terminate()

}

func l_process_input(window *glfw.Window) {
	// control the camera
  if l_Camera == nil {
    return 
  }
	if window.GetKey(glfw.KeyW) == glfw.Press {
		l_Camera.TranslateForward(l_DeltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		l_Camera.TranslateBackward(l_DeltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		l_Camera.TranslateLeft(l_DeltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		l_Camera.TranslateRight(l_DeltaTime)
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		l_Camera.TranslateUp(l_DeltaTime)
	}
	if window.GetKey(glfw.KeyZ) == glfw.Press {
		l_Camera.TranslateDown(l_DeltaTime)
	}

}

func l_mouse_callback(window *glfw.Window, xpos float64, ypos float64) {

	l_Camera.LookAtCursor(float32(xpos), float32(ypos))

}

func l_scroll_callback(window *glfw.Window, xoffset float64, yoffset float64) {
	l_Camera.StepFOV(float32(yoffset))
}
