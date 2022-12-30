package geometry

// checked

type Vectors []*Vector

func (v *Vectors) Vector(i int) *Vector {
	return (*v)[i]
}

func (v *Vectors) Add(vector *Vector) {
	*v = append(*v, vector)
}

func (v *Vectors) Insert(vector *Vector, index int) {
	if len(*v) == index {
		v.Add(vector)
	}
	*v = append((*v)[:index+1], (*v)[index:]...)
	(*v)[index] = vector
}

func (v *Vectors) Reserve() {
	first := 0
	last := len(*v) - 1
	for first < last {
		(*v)[first], (*v)[last] = (*v)[last], (*v)[first]
		first++
		last--
	}
}

func (v *Vectors) Length() int {
	return len(*v)
}
