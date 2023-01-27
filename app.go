package main

import (
	"github.com/AllenDang/imgui-go"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/micahke/glfw_imgui_backend"
	. "micahke/go-graphics-engine/core"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

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

	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	impl := glfw_imgui_backend.ImguiGlfw3Init(window, io)
	defer impl.Shutdown()

	positions := []float32{
		100.0, 100.0, 0.0, 0.0,
		200.0, 100.0, 1.0, 0.0,
		200.0, 200.0, 1.0, 1.0,
		100.0, 200.0, 0.0, 1.0,
	}

	indeces := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	va := NewVertexArray()
	vb := NewVertexBuffer(positions)

	layout := NewVertexBufferLayout()
	layout.Pushf(2)
	layout.Pushf(2)
	va.AddBuffer(vb, layout)

	ib := NewIndexBuffer(indeces)

	// Ortho Projection
	proj := mgl32.Ortho(0, 960, 0, 540, -1.0, 1.0)
	// View Projection (camera)
	view := mgl32.Translate3D(-100, 0, 0)

	// vp := mgl32.Vec4{100.0, 100.0, 0.0, 1.0}
	// result := proj.Mul4x1(vp) // Simulating what the shader is doing

	shader := NewShader("vertexShader.glsl", "fragmentShader.glsl")
	shader.Bind()
	shader.SetUniform4f("u_Color", 0.8, 0.3, 0.8, 1.9)

	texture := NewTexture("fragment.png")
	texture.Bind(0)
	shader.SetUniform1i("u_Texture", 0)

	va.Unbind()
	vb.Unbind()
	ib.Unbind()
	shader.Unbind()

	var r float32 = 0.0
	var increment float32 = 0.05

	translationPositions := [3]float32{200, 200, 0}

  renderer := NewRenderer()

	for !window.ShouldClose() {

		impl.NewFrame()

		// Model Tranformation (coverting into NDC)
		model := mgl32.Translate3D(translationPositions[0], translationPositions[1], translationPositions[2])

		mvp := proj.Mul4(view).Mul4(model)

		// imgui.SliderFloat3("Translation", , 0.0, 960.0)
    imgui.SliderFloat3("Translation 3D", &translationPositions, 0, 960)

		gl.Clear(gl.COLOR_BUFFER_BIT)

		// shader.Bind()
		shader.SetUniform4f("u_Color", r, 0.3, 0.8, 1.0)
		shader.SetUniformMat4f("u_MVP", mvp)

    renderer.Draw(va, ib, shader)

		if r > 1.0 {
			increment = -0.5
		} else if r < 0.0 {
			increment = 0.05
		}

		r = r + increment

		imgui.Render()
		impl.Render(imgui.RenderedDrawData())

		window.SwapBuffers()

		glfw.PollEvents()
	}

	// TODO: add shutdown function to the package
	context.Destroy()
	glfw.Terminate()
}
