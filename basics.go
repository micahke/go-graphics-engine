package main

import (
	. "micahke/go-graphics-engine/core"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)


var b_CameraPos mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, 3.0}
var b_CameraFront mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, -1.0}
var b_CameraUp mgl32.Vec3 = mgl32.Vec3{0.0, 1.0, 0.0}

var deltaTime float32 = 0.0
var lastFrame float32 = 0.0

var camera *Camera

func RunBasics() {

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

	// positions := []float32{
	// 	// positions           // texture coords
	// 	0.5, 0.5, 0.0, 1.0, 1.0, // top right
	// 	0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
	// 	-0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
	// 	-0.5, 0.5, 0.0, 0.0, 1.0, // top left
	// }

	positions := GetCubePositions()

	indeces := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)

	va := NewVertexArray()
	vb := NewVertexBuffer(*positions)

	layout := NewVertexBufferLayout()
	layout.Pushf(3) // represents the postions
	layout.Pushf(2) // represents the texture coords
	va.AddBuffer(*vb, *layout)

	ib := NewIndexBuffer(indeces)

	shader := NewShader("vertexShader.glsl", "fragmentShader.glsl")
	shader.Bind()
	shader.SetUniform4f("u_Color", 0.8, 0.3, 0.8, 1.9)

	texture := NewTexture("fragment.png")
	texture.Bind(0)
	shader.SetUniform1i("u_Texture", 0)

	renderer := NewRenderer()

	// set up glfw camera
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(mouse_callback)
	window.SetScrollCallback(scroll_callback)

  camera = NewCamera(b_CameraPos, b_CameraFront, b_CameraUp)

	// RENDER LOOP
	for !window.ShouldClose() {

		var currentFrame float32 = float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		processInput(window, camera)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    camera.Update(&shader)

		for i := 0; i < 10; i++ {
			model := mgl32.Ident4()
			activePosition := (*GetMultiCubePositions())[i]
			modelTranslation := mgl32.Translate3D(activePosition[0], activePosition[1], activePosition[2])
			angle := 25 * i
			modelRotation := mgl32.HomogRotate3D(mgl32.DegToRad(float32(angle))*float32(glfw.GetTime()/10), mgl32.Vec3{1.0, 0.3, 0.5})
			model = model.Mul4(modelTranslation)
			model = model.Mul4(modelRotation)
			shader.SetUniformMat4f("model", model)
			shader.SetUniform4f("u_Color", 0.2, 0.3, 0.8, 1.0)
      ib.Bind()
			renderer.Draw(*va, shader)
		}

		// glfw: swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()
}

// processes the GLFW window input every frame
func processInput(window *glfw.Window, camera *Camera) {
	// control the camera

	if window.GetKey(glfw.KeyW) == glfw.Press {
    camera.TranslateForward(deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
    camera.TranslateBackward(deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
    camera.TranslateLeft(deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
    camera.TranslateRight(deltaTime)
	}

}


func mouse_callback(window *glfw.Window, xpos float64, ypos float64) {

  camera.LookAtCursor(float32(xpos), float32(ypos))

}

func scroll_callback(window *glfw.Window, xoffset float64, yoffset float64) {
  camera.StepFOV(float32(yoffset))
}
