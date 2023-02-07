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

var l_LightPosition glm.Vec3 = glm.Vec3{1.2, 1.0, 2.0}

var l_LightColor glm.Vec3 = glm.Vec3{1.0, 1.0, 1.0}

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

	positions := GetLightingMapCoords()

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)

	vao := NewVertexArray()
	lightCubeVAO := NewVertexArray()
	vbo := NewVertexBuffer(*positions)
	vbl := NewVertexBufferLayout()
	vbl.Pushf(3)
	vbl.Pushf(3)
	vbl.Pushf(2)
	// NOTE: Add a way to customize the stride so that we can ignore values in the buffer

	vao.AddBuffer(*vbo, *vbl)
	lightCubeVAO.AddBuffer(*vbo, *vbl)

	renderer := NewRenderer()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(l_mouse_callback)
	window.SetScrollCallback(l_scroll_callback)

	l_Camera = NewCamera(l_CameraPos, l_CameraFront, l_CameraUp)

	objectShader := NewShader("lsVertex.glsl", "lsObject.glsl")
	lightShader := NewShader("lsVertex.glsl", "lsLight.glsl")

	diffuseMap := NewTexture("wall.png")

	objectShader.Bind()
	objectShader.SetUniform1f("material.diffuse", 0)

	for !window.ShouldClose() {
		var currentFrame float32 = float32(glfw.GetTime())
		l_DeltaTime = currentFrame - l_LastFrame
		l_LastFrame = currentFrame

		l_process_input(window)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		objectShader.Bind()
		objectShader.SetUniform3f("objectColor", 1.0, 0.5, 0.32)
		objectShader.SetUniform3f("lightColor", 1.0, 1.0, 1.0)
		objectShader.SetUniform3f("lightPos", l_LightPosition[0], l_LightPosition[1], l_LightPosition[2])
		objectShader.SetUniform3f("viewPos", l_Camera.Position[0], l_Camera.Position[1], l_Camera.Position[2])

		// material shaders
		objectShader.SetUniform3f("material.ambient", 1.0, 0.5, 0.31)
		objectShader.SetUniform3f("material.diffuse", 1.0, 0.5, 0.31)
		objectShader.SetUniform3f("material.specular", 0.5, 0.5, 0.5)
		objectShader.SetUniform1f("material.shininess", 32.0)

		// light shaders
		objectShader.SetUniform3f("light.ambient", 0.2, 0.2, 0.2)
		objectShader.SetUniform3f("light.diffuse", 0.5, 0.5, 0.5)
		objectShader.SetUniform3f("light.specular", 1.0, 1.0, 1.0)

		// Light colors
		// l_LightColor[0] = float32(math.Sin(glfw.GetTime() * 2.0))
		// l_LightColor[1] = float32(math.Sin(glfw.GetTime() * 0.7))
		// l_LightColor[2] = float32(math.Sin(glfw.GetTime() * 1.3))

		// DIffuse and ambient colors
		diffuseColor := l_LightColor.Mul(0.5)
		ambientColor := diffuseColor.Mul(0.2)

		// Set light values in the component
		objectShader.SetUniform3f("light.ambient", ambientColor[0], ambientColor[1], ambientColor[2])
		objectShader.SetUniform3f("light.diffuse", diffuseColor[0], diffuseColor[1], diffuseColor[2])

		l_Camera.Update(&objectShader)

		// draw the object
		model := glm.Ident4()
		objectShader.SetUniformMat4f("model", model)

		diffuseMap.Bind(0)

		renderer.Draw(*vao, objectShader)

		lightShader.Bind()
		l_Camera.Update(&lightShader)

		// Draw light cube
		lightCube := glm.Ident4()
		lightTranslation := glm.Translate3D(l_LightPosition[0], l_LightPosition[1], l_LightPosition[2])
		lightScale := glm.Scale3D(0.2, 0.2, 0.2)
		lightCube = lightCube.Mul4(lightTranslation).Mul4(lightScale)
		lightShader.SetUniformMat4f("model", lightCube)
		lightShader.SetUniform3f("cubeColor", l_LightColor[0], l_LightColor[1], l_LightColor[2])
		renderer.Draw(*lightCubeVAO, lightShader)

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
