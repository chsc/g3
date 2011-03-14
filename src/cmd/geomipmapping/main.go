package main

import (
	"os"
	_ "fmt"
	"runtime"
	_ "image/png"  // Only register png/jpeg decoder, but never use it directly. 
	_ "image/jpeg" // image.Decode does all the work for us.
	"g3"
	geo "g3/geomipmapping"
)

var (
	dirBase  = g3.Vec3{0, 1, 0}
	upBase   = g3.Vec3{0, 0, 1}
	leftBase = g3.Vec3{-1, 0, 0}
)

var (
	lightPos      = g3.Vec3{10, 10, 10}
	pos           = g3.Vec3{0, 0, 0.01}
	alpha         = float32(g3.Pi/2)
	beta          = float32(0)
	speed         = float32(0)
	dir, up, left g3.Vec3
	projection    g3.Matrix4x4
	modelView     g3.Matrix4x4
)

var (
	geoMipMap *geo.GeoMipMap
	mapShader g3.Shader
	texStone  g3.Texture2D
	texGrass  g3.Texture2D
	locLight  uint
	locStone  uint
	locGrass  uint
	frustum   *g3.Frustum
	wireframe bool
)

func multiplexEvents(engine g3.Engine) {
	for {
		select {
		case me := <-engine.MouseEventChan():
			alpha += float32(-me.Dx) / 100
			beta += float32(-me.Dy) / 100
			beta = g3.Clamp(beta, -g3.Pi/2, g3.Pi/2)
			ud := g3.MakeXRotationMatrix(beta)
			lr := g3.MakeZRotationMatrix(alpha)
			m := lr.Multiply(&ud)
			dir = m.Transform(dirBase)
			up = m.Transform(upBase)
			left = m.Transform(leftBase)
		case ke := <-engine.KeyEventChan():
			//fmt.Println(ke)
			if ke.Type == g3.KeyPressed {
				switch ke.Key {
				case g3.KeyW:
					speed += 0.0001
				case g3.KeyS:
					speed -= 0.0001
				case g3.KeyA:
					pos.Accumulate(left.Scaled(0.01))
				case g3.KeyD:
					pos.Accumulate(left.Scaled(-0.01))
				case g3.KeyF1:
					if wireframe {
						engine.GetGraphicsDevice().SetFillMode(g3.FillSolid)
					} else {
						engine.GetGraphicsDevice().SetFillMode(g3.FillWireFrame)
					}
					wireframe = !wireframe
				case g3.KeyPageUp:
					pos.Z += 0.01
				case g3.KeyPageDown:
					pos.Z -= 0.01
				}
			}
		case fe := <-engine.FrameEventChan():
			update(engine, fe.DeltaTime)
			render(engine)
			engine.SwapBuffers()
		case se := <-engine.SystemEventChan():
			//fmt.Println(se)
			if se.Type == g3.SystemQuit {
				return
			}
		}
	}
}

func initialize(engine g3.Engine) os.Error {
	gdev := engine.GetGraphicsDevice()

	// Load and create heigh map
	hmap, err := geo.NewHeightMapFromImageFile("../../../data/heightmaps/map1.png")
	if err != nil {
		return err
	}
	geoMipMap = geo.NewGeoMipMap(gdev, hmap, 3, 32, 32, 5, 0.01, 0.3)

	// Load and compile shader
	sources, err := g3.ReadStringsFromFiles("../../../data/shaders/map.vs.glsl", "../../../data/shaders/map.fs.glsl")
	if err != nil {
		return err
	}
	mapShader = gdev.NewShader(sources[0], sources[1])
	gdev.SetShader(mapShader)

	// Setup textures
	images, err := g3.ReadImagesFromFiles("../../../data/textures/stone.jpg", "../../../data/textures/grass.jpg")
	if err != nil {
		return err
	}
	texStone = gdev.NewTexture2D(images[0])
	texGrass = gdev.NewTexture2D(images[1])
	gdev.SetTexture2D(texStone, 0)
	gdev.SetTexture2D(texGrass, 1)

	// Setup texture sampler
	locStone = mapShader.GetUniformLocation("textureStone")
	locGrass = mapShader.GetUniformLocation("textureGrass")

	// Setup light position
	locLight = mapShader.GetUniformLocation("lightPos")

	//fmt.Println("locs:", locStone, locGrass, locLight)

	// Setup projection matrix
	projection = g3.MakePerspectiveMatrix(45.0, 640.0/480.0, 0.001, 100.0)
	gdev.SetMatrix(g3.MatrixProjection, &projection)

	return nil
}

func shutdown(engine g3.Engine) {
	texStone.Release()
	texGrass.Release()
	mapShader.Release()
	geoMipMap.Release()
}

func update(engine g3.Engine, deltaTime float32) {
	pos.Accumulate(dir.Scaled(speed))

	center := pos.Add(dir)
	lookAt := g3.MakeLookAtMatrix(&pos, &center, &up)

	modelView = lookAt
	modelViewProjection := projection.Multiply(&modelView)
	frustum = g3.MakeFrustumFromMatrix(&modelViewProjection)
}

func render(engine g3.Engine) {
	gdev := engine.GetGraphicsDevice()
	gdev.Clear()

	gdev.SetMatrix(g3.MatrixModelView, &modelView)
	gdev.SetShader(mapShader)

	lpos := modelView.Transform(lightPos)
	mapShader.SetVec3(locLight, &lpos)
	mapShader.SetTexture(locStone, 0)
	mapShader.SetTexture(locGrass, 1)

	geoMipMap.Render(gdev, frustum)
}

func main() {
	runtime.GOMAXPROCS(2)

	engine := g3.NewSDLEngine()

	engine.Init(&g3.GraphicsSettings{Width: 640, Height: 480, Caption: "Test Application - Geo-Mipmapping"})
	if err := initialize(engine); err != nil {
		panic(err.String())
	}

	engine.EnterEventLoop()
	multiplexEvents(engine)

	shutdown(engine)
	engine.Shutdown()
}
