package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go"
	. "micahke/go-graphics-engine/core"
	"github.com/micahke/glfw_imgui_backend"
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

	impl := glfw_imgui_backend.ImguiGlfw3Init(window)
	defer impl.Shutdown()

	showDemoWindow := false
	showAnotherWindow := false
	counter := 0

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
	// Model Tranformation (coverting into NDC)
	model := mgl32.Translate3D(200, 200, 0)

	mvp := proj.Mul4(view).Mul4(model)

	// vp := mgl32.Vec4{100.0, 100.0, 0.0, 1.0}
	// result := proj.Mul4x1(vp) // Simulating what the shader is doing
	fmt.Println(mvp)

	shader := NewShader("vertexShader.glsl", "fragmentShader.glsl")
	shader.Bind()
	shader.SetUniform4f("u_Color", 0.8, 0.3, 0.8, 1.9)
	shader.SetUniformMat4f("u_MVP", mvp)

	texture := NewTexture("fragment.png")
	texture.Bind(0)
	shader.SetUniform1i("u_Texture", 0)

	va.Unbind()
	vb.Unbind()
	ib.Unbind()
	shader.Unbind()

	var r float32 = 0.0
	var increment float32 = 0.05

	for !window.ShouldClose() {

		impl.NewFrame()

		// 1. Show a simple window.
		// Tip: if we don't call ImGui::Begin()/ImGui::End() the widgets automatically appears in a window called "Debug".
		{
			imgui.Text("Hello, world!")

			imgui.Checkbox("Demo Window", &showDemoWindow)
			imgui.Checkbox("Another Window", &showAnotherWindow)

			if imgui.Button("Button") {
				counter++
			}
			imgui.SameLine()
			imgui.Text(fmt.Sprintf("counter = %d", counter))

		}

		// 2. Show another simple window. In most cases you will use an explicit Begin/End pair to name your windows.
		if showAnotherWindow {
			imgui.BeginV("Another Window", &showAnotherWindow, 0)
			imgui.Text("Hello from another window!")
			if imgui.Button("Close Me") {
				showAnotherWindow = false
			}
			imgui.End()
		}

		// 3. Show the ImGui demo window. Most of the sample code is in imgui.ShowDemoWindow().
		// Read its code to learn more about Dear ImGui!
		if showDemoWindow {
			imgui.ShowDemoWindow(&showDemoWindow)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.Bind()
		shader.SetUniform4f("u_Color", r, 0.3, 0.8, 1.0)

		va.Bind()
		ib.Bind()

		gl.DrawElements(gl.TRIANGLES, int32(ib.GetCount()), gl.UNSIGNED_INT, nil)

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

}
