package g3

type Plane struct {
	Normal   Vec3
	Distance float32
}

func (p *Plane) Normalize() {
	ilength := 1.0 / p.Normal.Length()
	p.Normal.Scale(ilength)
	p.Distance *= ilength
}

func (p *Plane) DistanceToPoint(v *Vec3) float32 {
	return p.Normal.Dot(*v) + p.Distance
}
