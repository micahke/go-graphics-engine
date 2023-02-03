package main

import (
	// "math"
	. "micahke/go-graphics-engine/core"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

var cameraPos mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, 3.0}
var cameraFront mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, -1.0}
var cameraUp mgl32.Vec3 = mgl32.Vec3{0.0, 1.0, 0.0}

func main() {

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
	va.AddBuffer(vb, layout)

	ib := NewIndexBuffer(indeces)

	shader := NewShader("vertexShader.glsl", "fragmentShader.glsl")
	shader.Bind()
	shader.SetUniform4f("u_Color", 0.8, 0.3, 0.8, 1.9)

	texture := NewTexture("fragment.png")
	texture.Bind(0)
	shader.SetUniform1i("u_Texture", 0)

	renderer := NewRenderer()

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), 960.0/540.0, 0.1, 100.0)
	shader.SetUniformMat4f("projection", projection)

	// RENDER LOOP
	for !window.ShouldClose() {

    processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		view := mgl32.Ident4()
		// var radius float32 = 5.0
		// camX := float32(math.Sin(glfw.GetTime())) * radius
		// camZ := float32(math.Cos(glfw.GetTime())) * radius
		cameraLookAt := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		view = view.Mul4(cameraLookAt)
		shader.SetUniformMat4f("view", view)

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

			renderer.Draw(va, ib, shader)
		}

		// glfw: swap buffers and poll IO events
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// TODO: add shutdown function to the package
	glfw.Terminate()
}

// processes the GLFW window input every frame
func processInput(window *glfw.Window) {
  // control the camera
	var cameraSpeed float32 = 0.05

	if window.GetKey(glfw.KeyW) == glfw.Press {
		translation := cameraFront.Mul(cameraSpeed)
		cameraPos = cameraPos.Add(translation)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		translation := cameraFront.Mul(cameraSpeed)
		cameraPos = cameraPos.Sub(translation)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
    crossProduct := cameraFront.Cross(cameraUp)
    crossProduct = crossProduct.Normalize()
    translation := crossProduct.Mul(cameraSpeed)
    cameraPos = cameraPos.Sub(translation)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
    crossProduct := cameraFront.Cross(cameraUp)
    crossProduct = crossProduct.Normalize()
    translation := crossProduct.Mul(cameraSpeed)
    cameraPos = cameraPos.Add(translation)
	}

}
