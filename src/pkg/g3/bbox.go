package g3

type BoundingBox struct {
	Min, Max Vec3
}

type BoundingSphere struct {
	Position Vec3
	Radius   float32
}

type BoundingVolume interface {
	ClassifyPoint(v *Vec3) bool
	ClassifyPlane(p *Plane) int
}

type HasBoundingBox interface {
	GetBoundingBox() *BoundingBox
}

type HasBoundingSphere interface {
	GetBoundingSphere() *BoundingSphere
}

type HasBoundingVolume interface {
	GetBoundingVolume() BoundingVolume
}

func MakeUndefinedBoundingBox() BoundingBox {
	return BoundingBox{Vec3{MathMax, MathMax, MathMax}, Vec3{-MathMax, -MathMax, -MathMax}}
}

func MakeBoundingBoxFromBoxes(boxes []BoundingBox) BoundingBox {
	bbox := MakeUndefinedBoundingBox()
	for _, b := range boxes {
		bbox.Min.X = Min(bbox.Min.X, b.Min.X)
		bbox.Min.Y = Min(bbox.Min.Y, b.Min.Y)
		bbox.Min.Z = Min(bbox.Min.Z, b.Min.Z)
		bbox.Max.X = Max(bbox.Max.X, b.Max.X)
		bbox.Max.Y = Max(bbox.Max.Y, b.Max.Y)
		bbox.Max.Z = Max(bbox.Max.Z, b.Max.Z)
	}
	return bbox
}

func MakeBoundingBoxFromPoints(points []Vec3) BoundingBox {
	bbox := MakeUndefinedBoundingBox()
	for _, p := range points {
		bbox.Min.X = Min(bbox.Min.X, p.X)
		bbox.Min.Y = Min(bbox.Min.Y, p.Y)
		bbox.Min.Z = Min(bbox.Min.Z, p.Z)
		bbox.Max.X = Max(bbox.Max.X, p.X)
		bbox.Max.Y = Max(bbox.Max.Y, p.Y)
		bbox.Max.Z = Max(bbox.Max.Z, p.Z)
	}
	return bbox
}

func (bbox *BoundingBox) ClassifyPoint(v *Vec3) bool {
	if v.X > bbox.Min.X && v.X < bbox.Max.X &&
		v.Y > bbox.Min.Y && v.Y < bbox.Max.Y &&
		v.Z > bbox.Min.Z && v.Z < bbox.Max.Z {
		return true
	}
	return false
}

func (bbox *BoundingBox) ClassifyPlane(p *Plane) int {
	inside := 0

	v := Vec3{bbox.Min.X, bbox.Min.Y, bbox.Min.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Min.X, bbox.Min.Y, bbox.Max.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Min.X, bbox.Max.Y, bbox.Min.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Min.X, bbox.Max.Y, bbox.Max.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Max.X, bbox.Min.Y, bbox.Min.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Max.X, bbox.Min.Y, bbox.Max.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Max.X, bbox.Max.Y, bbox.Min.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}
	v = Vec3{bbox.Max.X, bbox.Max.Y, bbox.Max.Z}
	if p.DistanceToPoint(&v) >= 0.0 {
		inside++
	}

	if inside == 8 {
		return 1
	} else if inside == 0 {
		return 0
	}
	return -1
}

func (bbox *BoundingBox) CalculateCenter() Vec3 {
	return Vec3{(bbox.Min.X+bbox.Max.X)/2, (bbox.Min.Y+bbox.Max.Y)/2, (bbox.Min.Z+bbox.Max.Z)/2}
}

