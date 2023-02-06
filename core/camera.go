package core

import glm "github.com/go-gl/mathgl/mgl32"
import "math"

type Camera struct {
	position glm.Vec3
	front    glm.Vec3
	up       glm.Vec3

	speed float32
	yaw   float32
	pitch float32
	fov   float32

	firstMouse bool
	lastX      float32
	lastY      float32
}

func NewCamera(position glm.Vec3, front glm.Vec3, up glm.Vec3) *Camera {

	camera := Camera{}
	camera.position = position
	camera.front = front
	camera.up = up

	camera.yaw = -90.0
	camera.pitch = 0.0
	camera.fov = 45.0
	camera.speed = 2.5

	camera.firstMouse = true
	camera.lastX = 480
	camera.lastY = 270

	return &camera

}

func (camera *Camera) SetSpeed(speed float32) {
	camera.speed = speed
}

func (camera *Camera) Update(shader *Shader) {
	projection := glm.Perspective(glm.DegToRad(camera.fov), 960.0/540.0, 0.1, 100.0)
	shader.SetUniformMat4f("projection", projection)

	view := glm.Ident4()
	// var radius float32 = 5.0 camX := float32(math.Sin(glfw.GetTime())) * radius camZ := float32(math.Cos(glfw.GetTime())) * radius
	cameraLookAt := glm.LookAtV(camera.position, camera.position.Add(camera.front), camera.up)
	view = view.Mul4(cameraLookAt)
	shader.SetUniformMat4f("view", view)
}

func (camera *Camera) TranslateForward(dt float32) {
	translation := camera.front.Mul(camera.speed * dt)
	camera.position = camera.position.Add(translation)
}

func (camera *Camera) TranslateBackward(dt float32) {
	translation := camera.front.Mul(camera.speed * dt)
	camera.position = camera.position.Sub(translation)
}

func (camera *Camera) TranslateLeft(dt float32) {
	crossProduct := camera.front.Cross(camera.up)
	crossProduct = crossProduct.Normalize()
	translation := crossProduct.Mul(camera.speed * dt)
	camera.position = camera.position.Sub(translation)
}

func (camera *Camera) TranslateRight(dt float32) {
	crossProduct := camera.front.Cross(camera.up)
	crossProduct = crossProduct.Normalize()
	translation := crossProduct.Mul(camera.speed * dt)
	camera.position = camera.position.Add(translation)
}

func (camera *Camera) TranslateUp(dt float32) {
	translation := camera.up.Mul(camera.speed * dt)
	camera.position = camera.position.Add(translation)
}

func (camera *Camera) TranslateDown(dt float32) {
	translation := camera.up.Mul(camera.speed * dt)
	camera.position = camera.position.Sub(translation)
}

func (camera *Camera) StepFOV(fov float32) {
	camera.fov -= fov
	if camera.fov < 1.0 {
		camera.fov = 1.0
	}
	if camera.fov > 45.0 {
		camera.fov = 45.0
	}
}

func (camera *Camera) UpdateFOV(fov float32) {
	camera.fov = fov
	if camera.fov < 1.0 {
		camera.fov = 1.0
	}
	if camera.fov > 45.0 {
		camera.fov = 45.0
	}
}

func (camera *Camera) LookAtCursor(x float32, y float32) {
	if camera.firstMouse {
		camera.lastX = float32(x)
		camera.lastY = float32(y)
		camera.firstMouse = false
	}

	// calculate offset from last mouse position
	var xOffset float32 = float32(x) - camera.lastX
	var yOffset float32 = camera.lastY - float32(y) // this needs to be reversed
	camera.lastX = float32(x)
	camera.lastY = float32(y)

	var sensitivity float32 = 0.1
	xOffset *= sensitivity
	yOffset *= sensitivity

	camera.yaw += xOffset
	camera.pitch += yOffset

	if camera.pitch > 89.0 {
		camera.pitch = 89.0
	}
	if camera.pitch < -89.0 {
		camera.pitch = -89.0
	}

	var direction glm.Vec3
	direction[0] = float32(math.Cos(float64(glm.DegToRad(camera.yaw))) * math.Cos(float64(glm.DegToRad(camera.pitch))))
	direction[1] = float32(math.Sin(float64(glm.DegToRad(camera.pitch))))
	direction[2] = float32(math.Sin(float64(glm.DegToRad(camera.yaw))) * math.Cos(float64(glm.DegToRad(camera.pitch))))
	camera.front = direction.Normalize()
}
