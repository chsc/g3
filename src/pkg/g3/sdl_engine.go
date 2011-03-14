package g3

import (
	"os"
	"runtime"
	"sdl"
)

const (
	// TODO: Add all SDL key codes
	KeyW = uint32(sdl.K_w)
	KeyA = uint32(sdl.K_a)
	KeyS = uint32(sdl.K_s)
	KeyD = uint32(sdl.K_d)
	KeyL = uint32(sdl.K_l)

	KeyF1 = uint32(sdl.K_F1)
	KeyF2 = uint32(sdl.K_F2)

	KeyUp    = uint32(sdl.K_UP)
	KeyDown  = uint32(sdl.K_DOWN)
	KeyRight = uint32(sdl.K_RIGHT)
	KeyLeft  = uint32(sdl.K_LEFT)
	KeyPageUp    = uint32(sdl.K_PAGEUP)
	KeyPageDown  = uint32(sdl.K_PAGEDOWN)
)

const (
	KeyPressed  = uint32(sdl.KEYDOWN)
	KeyReleased = uint32(sdl.KEYUP)
)

type SDLEngine struct {
	screen          *sdl.Surface
	systemEventChan chan SystemEvent
	frameEventChan  chan FrameEvent
	mouseEventChan  chan MouseEvent
	keyEventChan    chan KeyEvent
	gdevice         GraphicsDevice
}

func NewSDLEngine() *SDLEngine {
	return &SDLEngine{nil,
		make(chan SystemEvent),
		make(chan FrameEvent),
		make(chan MouseEvent, 8),
		make(chan KeyEvent, 8),
		nil}
}

func (engine *SDLEngine) Init(settings *GraphicsSettings) os.Error {
	runtime.LockOSThread()

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		return os.NewError("unable to initialize sdl.")
	}

	if sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1) != 0 {
		sdl.Quit()
		return os.NewError("double buffering not available.")
	}

	engine.screen = sdl.SetVideoMode(settings.Width, settings.Height, 16, sdl.OPENGL|sdl.RESIZABLE)
	if engine.screen == nil {
		sdl.Quit()
		return os.NewError("unable to set video mode.")
	}

	sdl.WM_SetCaption(settings.Caption, settings.Caption)

	engine.gdevice = NewOpenGLGraphicsDevice()

	return nil
}

func (engine *SDLEngine) Shutdown() {
	sdl.Quit()
	runtime.UnlockOSThread()
}

func (engine *SDLEngine) SwapBuffers() {
	sdl.GL_SwapBuffers()
}

func (engine *SDLEngine) sdlRenderLoop() {
	runtime.LockOSThread()
	for {
		var event sdl.Event
		for event.Poll() {
			switch event.Type {
			case sdl.MOUSEMOTION:
				m := event.MouseMotion()
				engine.mouseEventChan <- MouseEvent{int32(m.X), int32(m.Y), int32(m.Xrel), int32(m.Yrel), int32(m.State)}
			case sdl.KEYDOWN, sdl.KEYUP:
				k := event.Keyboard()
				engine.keyEventChan <- KeyEvent{uint32(k.Keysym.Sym), uint32(k.Type)}
			case sdl.QUIT:
				engine.systemEventChan <- SystemEvent{SystemQuit}
				return
			}
		}
		engine.frameEventChan <- FrameEvent{0.0}
	}
	runtime.UnlockOSThread()
}

func (engine *SDLEngine) EnterEventLoop() {
	go engine.sdlRenderLoop()
}

func (engine *SDLEngine) GetGraphicsDevice() GraphicsDevice {
	return engine.gdevice
}

func (engine *SDLEngine) SystemEventChan() <-chan SystemEvent {
	return engine.systemEventChan
}

func (engine *SDLEngine) FrameEventChan() <-chan FrameEvent {
	return engine.frameEventChan
}

func (engine *SDLEngine) MouseEventChan() <-chan MouseEvent {
	return engine.mouseEventChan
}

func (engine *SDLEngine) KeyEventChan() <-chan KeyEvent {
	return engine.keyEventChan
}
