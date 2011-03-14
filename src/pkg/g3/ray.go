package g3

type Ray3 struct {
	Pos Vec3
	Dir Vec3
}

func (ray *Ray3) NewPos(t float, d Vec3) Ray3 {
	return Ray3{Pos: ray.Pos.Add(ray.Dir.Mul(t)), Dir: d}
}
