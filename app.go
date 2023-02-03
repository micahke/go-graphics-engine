package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

	positions := []float32{
		// positions           // texture coords
		0.5, 0.5, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 0.0, 1.0, // top left
	}

	indeces := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	va := NewVertexArray()
	vb := NewVertexBuffer(positions)

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

	va.Unbind()
	vb.Unbind()
	ib.Unbind()
	shader.Unbind()

	renderer := NewRenderer()

	// RENDER LOOP
	for !window.ShouldClose() {

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// transformations
		// transform := mgl32.Ident4()
		// transform = transform.Mul4(mgl32.Translate3D(0.5, -0.5, 0))
		// transform = transform.Mul4(mgl32.HomogRotate3DZ(float32(glfw.GetTime())))

    model := mgl32.Ident4()
    modelRotation := mgl32.HomogRotate3D(mgl32.DegToRad(-55.0), mgl32.Vec3{1.0, 0.0, 0.0})
    model = model.Mul4(modelRotation)

    view := mgl32.Ident4()
    viewTranslation := mgl32.Translate3D(0.0, 0.0, -3.0)
    view = view.Mul4(viewTranslation)

    projection := mgl32.Perspective(mgl32.DegToRad(45.0), 960.0 / 540.0, 0.1, 100.0)

		shader.SetUniform4f("u_Color", 0.2, 0.3, 0.8, 1.0)
    shader.SetUniformMat4f("model", model)
    shader.SetUniformMat4f("view", view)
    shader.SetUniformMat4f("projection", projection)

		renderer.Draw(va, ib, shader)

		// glfw: swap buffers and poll IO events
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// TODO: add shutdown function to the package
	glfw.Terminate()
}
