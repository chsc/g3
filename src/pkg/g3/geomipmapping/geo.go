package geomipmapping

import (
	"g3"
)

const (
	crackNone        = uint(0)
	crackLeft        = uint(1 << 0)
	crackRight       = uint(1 << 1)
	crackTop         = uint(1 << 2)
	crackBottom      = uint(1 << 3)
	crackTopLeft     = crackTop | crackLeft
	crackTopRight    = crackTop | crackRight
	crackBottomLeft  = crackBottom | crackLeft
	crackBottomRight = crackBottom | crackRight
)

const (
	numCrackTypes = 9
)

type HeightMap interface {
	Size() (width, height int)
	Height(x, y float32) float32
}

type patch struct {
	center   g3.Vec3
	vertices g3.VertexBuffer
	normals  g3.VertexBuffer
}

type GeoMipMap struct {
	heightMap      HeightMap
	maxDepth       uint
	maxPatchWidth  int
	maxPatchHeight int
	maxLOD         uint
	whScale        float32
	hScale         float32
	root           g3.SpatElement
	lodIndices     [][]g3.IndexBuffer
}

type generatedLODIndices struct {
	lod        uint
	crackIndex int
	indices    []uint32
}

type generatedPatchVertices struct {
	vertices []g3.Vec3
	normals  []g3.Vec3
}

func indexToCrack(index int) uint {
	switch index {
	case 0:
		return crackNone
	case 1:
		return crackLeft
	case 2:
		return crackRight
	case 3:
		return crackTop
	case 4:
		return crackBottom
	case 5:
		return crackTopLeft
	case 6:
		return crackTopRight
	case 7:
		return crackBottomLeft
	case 8:
		return crackBottomRight
	}
	panic("invalid index")
}

func crackToIndex(crackFlag uint) int {
	switch crackFlag {
	case crackNone:
		return 0
	case crackLeft:
		return 1
	case crackRight:
		return 2
	case crackTop:
		return 3
	case crackBottom:
		return 4
	case crackTopLeft:
		return 5
	case crackTopRight:
		return 6
	case crackBottomLeft:
		return 7
	case crackBottomRight:
		return 8
	}
	panic("unknown flag")
}

// test crack flag
func isSet(mask, test uint) bool {
	return mask&test == test
}

// calculates 2^n
func pow2(n uint) uint {
	return 1 << n
}

// calculates vertex offset in patch
func offset(width, x, y uint) uint32 {
	return uint32(y*width + x)
}

// generate normal and vertex buffer
func createPatchVertices(heightMap HeightMap, x, y, width, height int, maxLod uint, scale, hscale float32) (vertices, normals []g3.Vec3) {
	div := pow2(maxLod)
	size := div + 1
	xstep := float32(width) / float32(div)
	ystep := float32(height) / float32(div)
	vertices = make([]g3.Vec3, 0, size*size)
	normals = make([]g3.Vec3, 0, size*size)
	for i := uint(0); i < size; i++ {
		for j := uint(0); j < size; j++ {
			px0, py0 := float32(x)+float32(i)*xstep, float32(y)+float32(j)*ystep
			px1, py1 := px0+xstep, py0
			px2, py2 := px0, py0+ystep

			ph0 := heightMap.Height(px0, py0)
			ph1 := heightMap.Height(px1, py1)
			ph2 := heightMap.Height(px2, py2)

			p0 := g3.Vec3{px0 * scale, py0 * scale, ph0 * hscale}
			p1 := g3.Vec3{px1 * scale, py1 * scale, ph1 * hscale}
			p2 := g3.Vec3{px2 * scale, py2 * scale, ph2 * hscale}
			tangent, bitangent := p1.Sub(p0), p2.Sub(p0)
			n := tangent.Cross(bitangent).Normalized()

			vertices = append(vertices, p0)
			normals = append(normals, n)
		}
	}
	return
}

//
// a---b
// |  /|
// | / |
// |/  |
// c---d
//
func addQuadN(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	return append(patchIndices, a, b, c, b, d, c)
}

//
// a---b
// |  /|
// e-- |
// |  \|
// c---d
//
func addQuadL(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i, j+skip/2)
	return append(patchIndices, a, b, e, b, d, e, d, c, e)
}

//
// a---b
// |\  |
// | --e
// |/  |
// c---d
//
func addQuadR(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip, j+skip/2)
	return append(patchIndices, a, b, e, a, e, c, e, d, c)
}

//
// a-e-b
// | | |
// | | |
// |/ \|
// c---d
//
func addQuadT(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j)
	return append(patchIndices, a, e, c, c, e, d, e, b, d)
}

//
// a---b
// |\ /|
// | | |
// | | |
// c-e-d
//
func addQuadB(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j+skip)
	return append(patchIndices, e, c, a, e, a, b, e, b, d)
}

//
// a-e-b
// |\| |
// f-\ |
// |  \|
// c---d
//
func addQuadTL(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j)
	f := offset(width, i, j+skip/2)
	return append(patchIndices, f, d, c, a, d, f, a, e, d, e, b, d)
}

//
// a-e-b
// | |/|
// | /-f
// |/  |
// c---d
//
func addQuadTR(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j)
	f := offset(width, i+skip, j+skip/2)
	return append(patchIndices, a, e, c, e, b, c, b, f, c, f, d, c)
}

//
// a---b
// |  /|
// f-/ |
// |/| |
// c-e-d
//
func addQuadBL(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j+skip)
	f := offset(width, i, j+skip/2)
	return append(patchIndices, a, b, f, f, b, c, c, b, e, e, b, d)
}

//
// a---b
// |\  |
// | \-f
// | |\|
// c-e-d
//
func addQuadBR(patchIndices []uint32, width, i, j, skip uint) []uint32 {
	a := offset(width, i, j)
	b := offset(width, i+skip, j)
	c := offset(width, i, j+skip)
	d := offset(width, i+skip, j+skip)
	e := offset(width, i+skip/2, j+skip)
	f := offset(width, i+skip, j+skip/2)
	return append(patchIndices, a, b, f, a, f, d, a, d, e, a, e, c)
}

func inTopLeftCorner(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackTopLeft) {
		return addQuadTL(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackLeft) {
		return addQuadL(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackTop) {
		return addQuadT(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inTopEdge(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackTop) {
		return addQuadT(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inTopRightCorner(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackTopRight) {
		return addQuadTR(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackRight) {
		return addQuadR(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackTop) {
		return addQuadT(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inLeftEdge(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackLeft) {
		return addQuadL(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inRightEdge(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackRight) {
		return addQuadR(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inBottomRightCorner(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackBottomRight) {
		return addQuadBR(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackRight) {
		return addQuadR(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackBottom) {
		return addQuadB(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inBottomLeftCorner(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackBottomLeft) {
		return addQuadBL(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackLeft) {
		return addQuadL(patchIndices, width, i, j, skip)
	}
	if isSet(crackMask, crackBottom) {
		return addQuadB(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func inBottomEdge(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if isSet(crackMask, crackBottom) {
		return addQuadB(patchIndices, width, i, j, skip)
	}
	return addQuadN(patchIndices, width, i, j, skip)
}

func addNewQuad(patchIndices []uint32, width, i, j, skip, crackMask uint) []uint32 {
	if skip == 1 {
		return addQuadN(patchIndices, width, i, j, skip)
	}
	if i == 0 {
		if j == 0 {
			return inTopLeftCorner(patchIndices, width, i, j, skip, crackMask)
		}
		if j == width-skip-1 {
			return inBottomLeftCorner(patchIndices, width, i, j, skip, crackMask)
		}
		return inLeftEdge(patchIndices, width, i, j, skip, crackMask)
	} else if i == width-skip-1 {
		if j == 0 {
			return inTopRightCorner(patchIndices, width, i, j, skip, crackMask)
		}
		if j == width-skip-1 {
			return inBottomRightCorner(patchIndices, width, i, j, skip, crackMask)
		}
		return inRightEdge(patchIndices, width, i, j, skip, crackMask)
	} else {
		if j == 0 {
			return inTopEdge(patchIndices, width, i, j, skip, crackMask)
		}
		if j == width-skip-1 {
			return inBottomEdge(patchIndices, width, i, j, skip, crackMask)
		}
		return addQuadN(patchIndices, width, i, j, skip)
	}
	return nil
}

func createLODIndices(lod, maxLod, crackMask uint) (patchIndices []uint32) {
	skip := pow2(lod)
	max := pow2(maxLod)
	patchIndices = make([]uint32, 0, 16) // grow as needed
	for i := uint(0); i < max; i += skip {
		for j := uint(0); j < max; j += skip {
			patchIndices = addNewQuad(patchIndices, max+1, i, j, skip, crackMask)
		}
	}
	return
}

// Return number of index lists generated.
func createAllLODIndices(maxLod uint, indicesChannel chan generatedLODIndices) uint {
	for crackIndex := 0; crackIndex < numCrackTypes; crackIndex++ {
		for lod := uint(0); lod < maxLod; lod++ {
			go func(ci int, l uint) {
				indices := createLODIndices(l, maxLod, indexToCrack(ci))
				indicesChannel <- generatedLODIndices{l, ci, indices}
			}(crackIndex, lod)
		}
	}
	return numCrackTypes * maxLod
}

// TODO: How can this be done in parallel with goroutines?
func (gm *GeoMipMap) buildQuadTreeRec(dev g3.GraphicsDevice, depth uint, x, y, w, h int, parent g3.SpatElement) (g3.SpatElement, int) {
	// TODO: remove magic numbers; calc from LOD
	if depth >= gm.maxDepth /*|| w <= gm.maxPatchWidth || h <= gm.maxPatchHeight*/ {
		//println(w, h)
		vertices, normals := createPatchVertices(gm.heightMap, x, y, w, h, gm.maxLOD, gm.whScale, gm.hScale)
		bbox := g3.MakeBoundingBoxFromPoints(vertices)
		center := bbox.CalculateCenter()
		patch := patch{center, dev.NewVertexBufferVec3(vertices), dev.NewVertexBufferVec3(normals)}
		return &g3.SpatLeaf{g3.SpatElementData{&bbox, parent, patch}}, 1
	}
	p := &g3.SpatNode{g3.SpatElementData{nil, parent, nil}, make([]g3.SpatElement, 4)}
	c1, l1 := gm.buildQuadTreeRec(dev, depth+1, x, y, w/2, h/2, p)
	c2, l2 := gm.buildQuadTreeRec(dev, depth+1, x+w/2, y, w/2, h/2, p)
	c3, l3 := gm.buildQuadTreeRec(dev, depth+1, x, y+h/2, w/2, h/2, p)
	c4, l4 := gm.buildQuadTreeRec(dev, depth+1, x+w/2, y+h/2, w/2, h/2, p)
	bbox := g3.MakeBoundingBoxFromBoxes([]g3.BoundingBox{
		*c1.GetBoundingVolume().(*g3.BoundingBox),
		*c2.GetBoundingVolume().(*g3.BoundingBox),
		*c3.GetBoundingVolume().(*g3.BoundingBox),
		*c4.GetBoundingVolume().(*g3.BoundingBox)})
	p.BVolume = &bbox
	p.Children = []g3.SpatElement{c1, c2, c3, c4}
	return p, l1 + l2 + l3 + l4
}

func (gm *GeoMipMap) buildLODIndices(dev g3.GraphicsDevice) {
	gm.lodIndices = make([][]g3.IndexBuffer, numCrackTypes)
	for i, _ := range gm.lodIndices {
		gm.lodIndices[i] = make([]g3.IndexBuffer, gm.maxLOD)
	}
	indicesChannel := make(chan generatedLODIndices)
	numMipMaps := createAllLODIndices(gm.maxLOD, indicesChannel)
	for i := uint(0); i < numMipMaps; i++ {
		c := <-indicesChannel
		gm.lodIndices[c.crackIndex][c.lod] = dev.NewIndexBuffer(c.indices)
	}
}

func (gm *GeoMipMap) buildPatchVertices(dev g3.GraphicsDevice) {
	w, h := gm.heightMap.Size()
	root, _ := gm.buildQuadTreeRec(dev, 0, 0, 0, w, h, nil)
	gm.root = root;
}

func NewGeoMipMap(dev g3.GraphicsDevice, heightMap HeightMap, maxDepth uint, maxPatchWidth, maxPatchHeight int, maxLOD uint, whScale, hScale float32) *GeoMipMap {
	gmap := &GeoMipMap{heightMap, maxDepth, maxPatchWidth, maxPatchHeight, maxLOD, whScale, hScale, nil, nil}
	gmap.buildLODIndices(dev)
	gmap.buildPatchVertices(dev)
	return gmap
}

func (gm *GeoMipMap) Release() {
	for _, lodBuffers := range gm.lodIndices {
		for _, crackBuffer := range lodBuffers {
			crackBuffer.Release()
		}
	}
}

func calculateLOD(frustum *g3.Frustum, center *g3.Vec3, width float32, maxLOD uint) int {
	// TODO: Use actual position instead of near clipping plane for distance calculation?
	d := int(g3.Clamp(frustum.Near.DistanceToPoint(center) / (width), 0.0, float32(maxLOD-1)))
	return d
}

// TODO: Hack, use correct distance and neighbour coordinates from height map
func  (gm *GeoMipMap) selectIndexBuffer(frustum *g3.Frustum, center *g3.Vec3, maxLOD uint) (int, int) {
	w := float32(32.0*0.01)

	n := g3.Vec3{center.X, center.Y, 0}
	d := calculateLOD(frustum, &n, w, maxLOD)

	n = g3.Vec3{center.X-w, center.Y, 0}
	td := calculateLOD(frustum, &n, w, maxLOD)

	n = g3.Vec3{center.X, center.Y-w, 0}
	ld := calculateLOD(frustum, &n, w, maxLOD)

	n = g3.Vec3{center.X, center.Y+w, 0}
	rd := calculateLOD(frustum, &n, w, maxLOD)

	n = g3.Vec3{center.X+w, center.Y, 0}
	bd := calculateLOD(frustum, &n, w, maxLOD)

	mask := crackNone
	if td < d {
		mask |= crackTop
	}
	if ld < d {
		mask |= crackLeft
	}
	if rd < d {
		mask |= crackRight
	}
	if bd < d {
		mask |= crackBottom
	}

	return crackToIndex(mask), d
}

func (gm *GeoMipMap) Render(dev g3.GraphicsDevice, frustum *g3.Frustum) {
	gm.root.TraverseFrustum(frustum, func(element g3.SpatElement, leaf bool) bool {
		if leaf {
			patch := element.GetData().(patch)
			dev.SetVertices(patch.vertices)
			dev.SetNormals(patch.normals)
			lodIndex, distIndex := gm.selectIndexBuffer(frustum, &patch.center, gm.maxLOD)
			dev.DrawIndexed(gm.lodIndices[lodIndex][distIndex])
		}
		return true
	})
}
