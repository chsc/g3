package g3

import (
	"os"
)

const (
	SystemQuit = iota
)

type SystemEvent struct {
	Type int
}

type FrameEvent struct {
	DeltaTime float32
}

type MouseEvent struct {
	X, Y   int32
	Dx, Dy int32
	Button int32
}

type KeyEvent struct {
	Key  uint32
	Type uint32
}

type GraphicsSettings struct {
	Width, Height int
	FullScreen    bool
	Caption       string
}

type Engine interface {
	Init(settings *GraphicsSettings /*TODO: other settings ...*/) os.Error
	Shutdown()

	GetGraphicsDevice() GraphicsDevice
	SwapBuffers()

	EnterEventLoop()
	SystemEventChan() <-chan SystemEvent
	FrameEventChan() <-chan FrameEvent
	MouseEventChan() <-chan MouseEvent
	KeyEventChan() <-chan KeyEvent
}
