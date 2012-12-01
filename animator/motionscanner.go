package animator

import (
	//"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
)

func rgbaAt(img *image.NRGBA, x, y int) color.NRGBA {
	return img.At(x, y).(color.NRGBA)
}

func rgbaAt2(img *image.NRGBA, x, y int) color.NRGBA {
	return img.At(x/2, y/2).(color.NRGBA)
}

func grayAt(img *image.Gray, x, y int) color.Gray {
	return img.At(x, y).(color.Gray)
}

type EdgeDetector func(c0, c1 color.NRGBA) color.Gray

type DeltaFunc func(c1, c2 color.NRGBA) float64

type MotionMap struct {
	Bounds image.Rectangle
	data   []image.Point
}

type Shift struct {
	Dx, Dy float64
}

type ShiftMap struct {
	Bounds image.Rectangle
	data   []Shift
}

func NewMotionMap(bounds image.Rectangle) *MotionMap {
	m := new(MotionMap)
	m.Bounds = bounds
	m.data = make([]image.Point, bounds.Dx()*bounds.Dy())
	return m
}

func GetMotion(m *MotionMap, x, y int) (v image.Point) {
	if m.Bounds.Min.X <= x && x < m.Bounds.Max.X && m.Bounds.Min.Y <= y && y < m.Bounds.Max.Y {
		v = m.data[x-m.Bounds.Min.X+(y-m.Bounds.Min.Y)*m.Bounds.Size().X]
	} else {
		panic("out of bounds")
	}
	return
}

func SetMotion(m *MotionMap, x, y int, v image.Point) {
	if m.Bounds.Min.X <= x && x < m.Bounds.Max.X && m.Bounds.Min.Y <= y && y < m.Bounds.Max.Y {
		m.data[x-m.Bounds.Min.X+(y-m.Bounds.Min.Y)*m.Bounds.Size().X] = v
	} else {
		panic("out of bounds")
	}
}

func NewShiftMap(bounds image.Rectangle) *ShiftMap {
	m := new(ShiftMap)
	m.Bounds = bounds
	m.data = make([]Shift, bounds.Dx()*bounds.Dy())
	return m
}

func GetShift(m *ShiftMap, x, y int) (s Shift) {
	if m.Bounds.Min.X <= x && x < m.Bounds.Max.X && m.Bounds.Min.Y <= y && y < m.Bounds.Max.Y {
		s = m.data[x-m.Bounds.Min.X+(y-m.Bounds.Min.Y)*m.Bounds.Size().X]
	} else {
		panic("out of bounds")
	}
	return
}

func AddShiftMap(m *ShiftMap, s *ShiftMap) {
	forRect(m.Bounds.Intersect(s.Bounds), 1, func(x, y int) {
		SetShift(m, x, y, AddShift(GetShift(m, x, y), GetShift(s, x, y)))
	})
}

func AddShift(s1, s2 Shift) Shift {
	return Shift{s1.Dx + s2.Dx, s1.Dy + s2.Dy}
}

func SetShift(m *ShiftMap, x, y int, s Shift) {
	if m.Bounds.Min.X <= x && x < m.Bounds.Max.X && m.Bounds.Min.Y <= y && y < m.Bounds.Max.Y {
		m.data[x-m.Bounds.Min.X+(y-m.Bounds.Min.Y)*m.Bounds.Size().X] = s
	} else {
		panic("out of bounds")
	}
}

func rectAt(x, y, dia int) image.Rectangle {
	return image.Rect(x-dia, y-dia, x+dia+1, y+dia+1)
}

func forRect(rect image.Rectangle, step int, f func(x, y int)) {
	for y := rect.Min.Y; y < rect.Max.Y; y += step {
		for x := rect.Min.X; x < rect.Max.X; x += step {
			f(x, y)
		}
	}
}

//  ______
// |_|_|_|
// |_|_|_|
// |_|_|_|
//
// 
func MarkEdges(src *image.NRGBA, ef EdgeDetector) (dst *image.Gray) {
	bounds := src.Bounds()
	dbounds := image.Rectangle{bounds.Min.Mul(2), bounds.Max.Mul(2).Sub(image.Pt(1, 1))}

	dst = image.NewGray(dbounds)

	xbounds := image.Rect(dbounds.Min.X+1, dbounds.Min.Y, dbounds.Max.X-1, dbounds.Max.Y)
	ybounds := image.Rect(dbounds.Min.X, dbounds.Min.Y+1, dbounds.Max.X, dbounds.Max.Y-1)

	forRect(xbounds, 2, func(x, y int) {
		c0 := rgbaAt2(src, x-1, y)
		c1 := rgbaAt2(src, x+1, y)
		dst.SetGray(x, y, ef(c0, c1))
	})

	forRect(ybounds, 2, func(x, y int) {
		c0 := rgbaAt2(src, x, y-1)
		c1 := rgbaAt2(src, x, y+1)
		dst.SetGray(x, y, ef(c0, c1))
	})

	forRect(dbounds.Inset(1), 2, func(x, y int) {
		g := uint8(0)
		g += grayAt(dst, x-1, y).Y / 128
		g += grayAt(dst, x+1, y).Y / 128
		g += grayAt(dst, x, y-1).Y / 128
		g += grayAt(dst, x, y+1).Y / 128

		if g > 1 {
			dst.SetGray(x, y, color.Gray{255})
		}
	})

	// forRect(dbounds.Inset(1), 2, func(x, y int) {
	// 	c0 := rgbaAt2(src, x-1, y-1)
	// 	c1 := rgbaAt2(src, x+1, y+1)
	// 	c2 := rgbaAt2(src, x+1, y-1)
	// 	c3 := rgbaAt2(src, x-1, y+1)
	// 	y1 := ef(c0, c1).Y
	// 	y2 := ef(c2, c3).Y

	// 	ym := y1
	// 	if y1 < y2 {
	// 		ym = y2
	// 	}

	// 	// if x == 55 && y == 33 {
	// 	// 	fmt.Println(c0, c1, c2, c3, y1, y2, ym)
	// 	// }

	// 	dst.SetGray(x, y, color.Gray{ym})
	// })

	return
}

func GradientInflate(img *image.Gray) {
	bounds := img.Bounds().Inset(1)

	for c := uint8(255); c > uint8(0); c -= 1 {
		changed := false
		forRect(bounds, 1, func(x, y int) {
			if grayAt(img, x, y).Y == 0 {
			loopNeighbours:
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if grayAt(img, x+dx, y+dy).Y == c {
							img.SetGray(x, y, color.Gray{c - 1})
							changed = true
							break loopNeighbours
						}
					}
				}
			}
		})

		if !changed {
			break
		}
	}
}

func Skelton(img *image.NRGBA, gradient *image.Gray) (ske *image.NRGBA) {
	bounds := gradient.Bounds()
	ske = image.NewNRGBA(bounds)

	forRect(bounds.Inset(1), 1, func(x, y int) {
		c := grayAt(gradient, x, y).Y

		// if (grayAt(gradient, x-1, y).Y > c && grayAt(gradient, x+1, y).Y > c) ||
		// 	(grayAt(gradient, x, y-1).Y > c && grayAt(gradient, x, y+1).Y > c) {
		// 	ske.SetNRGBA(x, y, rgbaAt2(img, x, y))
		// }

		if (grayAt(gradient, x-1, y).Y > c && grayAt(gradient, x+1, y).Y > c) ||
			(grayAt(gradient, x, y-1).Y > c && grayAt(gradient, x, y+1).Y > c) ||
			((grayAt(gradient, x-1, y).Y == grayAt(gradient, x+1, y).Y ||
				grayAt(gradient, x, y-1).Y == grayAt(gradient, x, y+1).Y) &&
				((grayAt(gradient, x-1, y-1).Y > c && grayAt(gradient, x+1, y+1).Y > c) ||
					(grayAt(gradient, x-1, y+1).Y > c && grayAt(gradient, x+1, y-1).Y > c))) {
			ske.SetNRGBA(x, y, rgbaAt2(img, x, y))
		}
	})

	return ske
}

// func mean(min, max int, f func(i int) int) float64 {
// 	sum, count := 0, 0
// 	for i := min; i <= max; i++ {
// 		v := f(i)
// 		sum += v * i
// 		count += v
// 	}
// 	if count != 0 {
// 		return float64(sum) / float64(count)
// 	}
// 	return float64(min+max) * 0.5
// }

func forXY(x0, y0, x1, y1 int, f func(x, y int)) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			f(x, y)
		}
	}
}

func SmoothSkeletonShift(ske *image.NRGBA) (skeShift *ShiftMap) {
	bounds := ske.Bounds()
	skeShift = NewShiftMap(bounds)

	//TODO: calculate shift to average center of 3x3 surrounding (x,y) of skeleton
	forRect(bounds.Inset(1), 1, func(x, y int) {

		if rgbaAt(ske, x, y).A > 0 {
			count := 0
			sx, sy := 0, 0

			forXY(-1, -1, 1, 1, func(dx, dy int) {
				if rgbaAt(ske, x+dx, y+dy).A > 0 {
					sx += dx
					sy += dy
					count++
				}
			})

			if count > 2 {
				SetShift(skeShift, x, y, Shift{
					float64(sx) / float64(count),
					float64(sy) / float64(count),
				})
			}

			// SetShift(skeShift, x, y, Shift{
			// 	mean(-1, 1, func(i int) int {
			// 		return hasSkeleton(x+i, y-1) + hasSkeleton(x+i, y+1)
			// 	}),
			// 	mean(-1, 1, func(i int) int {
			// 		return hasSkeleton(x-1, y+i) + hasSkeleton(x+1, y+i)
			// 	}),
			// })
		}
	})
	return
}

func NearestMotion(ske0, ske1 *image.NRGBA, maxdelta float64, delta DeltaFunc) (m *MotionMap) {
	bounds := ske0.Bounds().Intersect(ske1.Bounds())

	m = NewMotionMap(bounds)

	forRect(bounds, 1, func(x, y int) {
		if rgbaAt(ske0, x, y).A > 0 {
			SetMotion(m, x, y, findNearest(ske1, x, y, maxdelta, rgbaAt(ske0, x, y), delta).Sub(image.Point{x, y}))
		}
	})

	return
}

func MotionDensity(m *MotionMap) *image.Gray {
	bounds := m.Bounds
	density := image.NewGray(bounds)

	forRect(bounds, 1, func(x, y int) {
		v := GetMotion(m, x, y)
		if v.X != 0 || v.Y != 0 {
			tx, ty := x+v.X, y+v.Y
			g := grayAt(density, tx, ty).Y
			if g < 255 {
				density.SetGray(tx, ty, color.Gray{g + 1})
			}
		}
	})

	return density
}

func MotionDistance(m *MotionMap) *image.Gray {
	bounds := m.Bounds
	density := image.NewGray(bounds)

	forRect(bounds, 1, func(x, y int) {
		v := GetMotion(m, x, y)
		if v.X != 0 || v.Y != 0 {
			tx, ty := x+v.X, y+v.Y
			dmax := grayAt(density, tx, ty).Y
			d := uint8(math.Min(math.Floor(math.Sqrt(float64(v.X*v.X+v.Y*v.Y))), 255.0))

			if dmax < d {
				density.SetGray(tx, ty, color.Gray{d})
			}
		}
	})

	return density
}

func gravity(ske *image.NRGBA, px, py int, maxdelta float64, col color.NRGBA, delta DeltaFunc) float64 {
	d := int(math.Floor(math.Sqrt(maxdelta)))
	bounds := rectAt(px, py, d).Intersect(ske.Bounds())
	g := float64(0)

	forRect(bounds, 1, func(x, y int) {
		dx, dy := float64(x-px), float64(y-py)
		dt := delta(rgbaAt(ske, x, y), col) + dx*dx + dy*dy
		if dt < maxdelta {
			g += maxdelta - dt
		}
	})

	return g
}

func centerOfGravityShift(ske *image.NRGBA, px, py, d int, col color.NRGBA, delta DeltaFunc) Shift {
	s := Shift{0, 0}
	sn := float64(0)
	r := rectAt(px, py, d).Intersect(ske.Bounds())
	maxdelta := float64(d * d)

	// if px == 65 {
	// 	fmt.Println("centerof", d, maxdelta, r)
	// }

	forRect(r, 1, func(x, y int) {
		dx, dy := float64(x-px), float64(y-py)
		dt := delta(rgbaAt(ske, x, y), col) + dx*dx + dy*dy
		if dt < maxdelta {
			g := gravity(ske, x, y, maxdelta/4.0, rgbaAt(ske, x, y), delta)
			s.Dx += dx * g
			s.Dy += dy * g
			sn += g

			// if px == 65 {
			// 	fmt.Println(dx, dy, g)
			// }
		}
	})

	s.Dx /= sn
	s.Dy /= sn
	return s
}

func GravityShift(ske *image.NRGBA, md *image.Gray, delta DeltaFunc) *ShiftMap {
	bounds := ske.Bounds()
	shift := NewShiftMap(bounds)

	forRect(bounds, 1, func(x, y int) {
		d := int(grayAt(md, x, y).Y)

		if d > 1 {
			SetShift(shift, x, y, centerOfGravityShift(ske, x, y, d, rgbaAt(ske, x, y), delta))
		}
	})
	return shift
}

func findNearest(ske *image.NRGBA, px, py int, maxdelta float64, col color.NRGBA, delta DeltaFunc) image.Point {
	min := maxdelta
	minp := image.Point{px, py}

	d := int(math.Floor(math.Sqrt(maxdelta)))

	bounds := rectAt(px, py, d).Intersect(ske.Bounds())

	forRect(bounds, 1, func(x, y int) {
		dx, dy := float64(x-px), float64(y-py)
		dt := delta(rgbaAt(ske, x, y), col) + dx*dx + dy*dy
		if dt < min {
			min = dt
			minp = image.Point{x, y}
		}
	})

	return minp
}

func MotionSource(m *MotionMap) (center *ShiftMap, count *image.Gray) {
	bounds := m.Bounds
	center = NewShiftMap(bounds)
	count = image.NewGray(bounds)

	forRect(bounds, 1, func(x, y int) {
		v := GetMotion(m, x, y)
		if v.X != 0 || v.Y != 0 {
			tx, ty := x+v.X, y+v.Y

			// accumulate source vectors
			sv := GetShift(center, tx, ty)
			sv.Dx -= float64(v.X)
			sv.Dy -= float64(v.Y)
			SetShift(center, tx, ty, sv)

			g := grayAt(count, tx, ty).Y
			if g < 255 {
				count.SetGray(tx, ty, color.Gray{g + 1})
			} else {
				panic("FIXME: max count reached")
			}
		}
	})

	// get average source vectors
	forRect(bounds, 1, func(x, y int) {
		c := grayAt(count, x, y).Y
		if c > 0 {
			sv := GetShift(center, x, y)
			sv.Dx /= float64(c)
			sv.Dy /= float64(c)
			SetShift(center, x, y, sv)
		}
	})

	return
}

func SpreadMotion(m *MotionMap, sourceCenter, shift, srcShift, dstShift *ShiftMap, ske *image.NRGBA, delta DeltaFunc) *ShiftMap {

	bounds := m.Bounds.Intersect(shift.Bounds).Intersect(sourceCenter.Bounds)
	spread := NewShiftMap(bounds)

	forRect(bounds, 1, func(x, y int) {

		v := GetMotion(m, x, y)

		if v.X != 0 || v.Y != 0 {
			tx, ty := x+v.X, y+v.Y

			sv := GetShift(shift, tx, ty)
			src := GetShift(sourceCenter, tx, ty)

			srcv := GetShift(srcShift, x, y)

			dx, dy := -(src.Dx + float64(v.X)), -(src.Dy + float64(v.Y))
			//cx, cy := float64(x)-src.Dx+sv.Dx, float64(y)-src.Dy+sv.Dy
			cx, cy := float64(tx)+sv.Dx+srcv.Dx, float64(ty)+sv.Dy+srcv.Dy
			sd := math.Sqrt(sv.Dx*sv.Dx + sv.Dy*sv.Dy)

			// if tx == 17 {
			// 	fmt.Println(x, y, sv, src, sd, cx, cy, dx, dy)
			// }

			n := findOffsetNearest(ske, dstShift, cx, cy, dx, dy, sd, rgbaAt(ske, tx, ty), delta)

			SetMotion(m, x, y, image.Pt(n.X-x, n.Y-y))

			sp := Shift{
				sv.Dx - src.Dx,
				sv.Dy - src.Dy,
			}

			SetShift(spread, x, y, sp)
		}
	})

	return spread
}

func round(f float64) int {
	return int(math.Floor(f + .5))
}

func findOffsetNearest(ske *image.NRGBA, dstShift *ShiftMap, cx, cy, dx, dy, sd float64, col color.NRGBA, delta DeltaFunc) image.Point {
	d := math.Sqrt(dx*dx + dy*dy)
	px, py := round(cx), round(cy)
	tx, ty := cx+dx, cy+dy
	minp := image.Point{px, py}
	bounds := rectAt(px, py, int(math.Ceil(d*2.0+sd))).Intersect(ske.Bounds())
	min := d*d*4 + delta(col, rgbaAt(ske, px, py))

	// if px == 19 {
	// 	fmt.Println(cx, cy, dx, dy, sd, min)
	// }

	forRect(bounds, 1, func(x, y int) {

		sv := GetShift(dstShift, x, y)
		fx, fy := float64(x)+sv.Dx, float64(y)+sv.Dy

		dcx, dcy := fx-cx, fy-cy
		dtx, dty := fx-tx, fy-ty
		cd := d - math.Sqrt(dcx*dcx+dcy*dcy)
		td := math.Sqrt(dtx*dtx + dty*dty)
		d2 := math.Abs(cd) + td
		dt := delta(rgbaAt(ske, x, y), col) + d2*d2

		// if px == 19 {
		// 	fmt.Println(x, y, ":", dt, cd, td, d2)
		// }

		if dt < min {
			min = dt
			minp = image.Point{x, y}
		}
	})

	return minp
}

func TransposeMotion(m0, m1 *MotionMap, ske1 *image.NRGBA, delta DeltaFunc) {
	bounds := m0.Bounds

	forRect(bounds, 1, func(x, y int) {
		v0 := GetMotion(m0, x, y)
		if v0.X != 0 || v0.Y != 0 {
			tx, ty := x+v0.X, y+v0.Y
			v1 := GetMotion(m1, tx, ty)
			px, py := x-v1.X, y-v1.Y
			dx, dy := float64(px-tx), float64(py-ty)
			maxdelta := (dx*dx + dy*dy) * 4

			SetMotion(m0, x, y, findNearest(ske1, px, py, maxdelta, rgbaAt(ske1, tx, ty), delta).Sub(image.Point{x, y}))
		}
	})
}

func OptimizeMotion(m0, m1 *MotionMap, ske1 *image.NRGBA, maxDelta float64, delta DeltaFunc) {
	// TODO: find minimum of "distance to ske1" + "distance to ske1's motion target"

	bounds := m0.Bounds
	maxd := int(math.Floor(math.Sqrt(maxDelta)))

	forRect(bounds, 1, func(x, y int) {
		v0 := GetMotion(m0, x, y)

		if v0.X != 0 || v0.Y != 0 {

			tx, ty := x+v0.X, y+v0.Y
			v1 := GetMotion(m1, tx, ty)
			vd := v0.Add(v1)

			len0 := math.Sqrt(float64(v0.X*v0.X + v0.Y*v0.Y))
			len2 := math.Sqrt(float64(vd.X*vd.X + vd.Y*vd.Y))

			d := int(math.Floor(len0 + len2))

			if d > maxd {
				d = maxd
			}

			col := rgbaAt(ske1, tx, ty)
			min := float64(d)

			forRect(rectAt(x, y, d).Intersect(ske1.Bounds()), 1, func(nx, ny int) {
				dx, dy := nx-x, ny-y

				if delta(rgbaAt(ske1, nx, ny), col)+float64(dx*dx+dy*dy) < maxDelta {

					vn := GetMotion(m1, nx, ny)
					//dt := delta(rgbaAt(ske1, nx, ny), col) + float64(dx*dx+dy*dy+vn.X*vn.X+vn.Y*vn.Y)
					drx, dry := dx+vn.X, dy+vn.Y
					//dt := delta(rgbaAt(ske1, nx, ny), col) + float64(dx*dx+dy*dy+drx*drx+dry*dry)
					rlen := math.Sqrt(float64(drx*drx + dry*dry))
					vlen := math.Sqrt(float64(dx*dx + dy*dy))
					dt := rlen + vlen

					if dt < min {
						tx, ty = nx, ny
						min = dt
					}
				}
			})

			SetMotion(m0, x, y, image.Point{tx - x, ty - y})
		}
	})
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func signInt(i int) int {
	if i < 0 {
		return -1
	} else if i > 0 {
		return 1
	}
	return 0
}

func CreateBackground(img0, img1 *image.NRGBA, m0, m1 *MotionMap, g0, g1 *image.Gray) (bg0, bg1 *image.NRGBA, msrc, mdst *image.Gray) {
	ibounds := img0.Bounds().Intersect(img1.Bounds())
	bounds := m0.Bounds.Intersect(m1.Bounds)
	bg0 = image.NewNRGBA(ibounds)
	bg1 = image.NewNRGBA(ibounds)

	msrc = image.NewGray(bounds)
	mdst = image.NewGray(bounds)

	draw.Draw(msrc, bounds, image.NewUniform(color.Gray{255}), image.Point{}, draw.Src)
	draw.Draw(mdst, bounds, image.NewUniform(color.Gray{255}), image.Point{}, draw.Src)

	set := func(x, y int, gr *image.Gray, v uint8, m *image.Gray) {
		d := 255 - int(grayAt(gr, x, y).Y)

		forRect(rectAt(x, y, d), 1, func(rx, ry int) {
			c := grayAt(m, rx, ry).Y
			if v < c {
				m.SetGray(rx, ry, color.Gray{v})
			}

		})
	}

	forRect(bounds, 1, func(x, y int) {
		v := GetMotion(m0, x, y)
		d := uint8(math.Min(255., math.Ceil(math.Sqrt(float64(v.X*v.X+v.Y*v.Y)))))

		if d > 0 {
			set(x, y, g0, d, msrc)
		}

		v = GetMotion(m1, x, y)
		d = uint8(math.Min(255., math.Ceil(math.Sqrt(float64(v.X*v.X+v.Y*v.Y)))))

		if d > 0 {
			set(x, y, g1, d, mdst)
		}
	})

	forRect(bounds, 2, func(x, y int) {
		v0, v1 := grayAt(msrc, x, y).Y, grayAt(mdst, x, y).Y
		var c0, c1 color.Color

		switch {
		case v0 > v1:
			c0 = img0.At(x/2, y/2)
			c1 = c0
		case v0 < v1:
			c0 = img1.At(x/2, y/2)
			c1 = c0
		default:
			c0 = img0.At(x/2, y/2)
			c1 = img1.At(x/2, y/2)
		}

		bg0.Set(x/2, y/2, c0)
		bg1.Set(x/2, y/2, c1)
	})

	return

	//TODO: 
	// 1) if src moves but dst is missing mark with img1
	// 2) if dst moves but src is missing mark with img0

	// two backgrounds3
}

// func CreateBackground(img0, img1 *image.NRGBA, m0, m1 *MotionMap, g0, g1 *image.Gray) (bg *image.NRGBA, msrc, mdst *MotionMap) {
// 	ibounds := img0.Bounds().Intersect(img1.Bounds())
// 	bounds := m0.Bounds.Intersect(m1.Bounds)
// 	bg = image.NewNRGBA(ibounds)

// 	add := func(x, y int, gr *image.Gray, v image.Point, m *MotionMap) {
// 		d := 255 - int(grayAt(gr, x, y).Y)

// 		forRect(rectAt(x, y, d), 1, func(rx, ry int) {
// 			SetMotion(m, rx, ry, GetMotion(m, rx, ry).Add(v))
// 		})
// 	}

// 	msrc = NewMotionMap(bounds)
// 	mdst = NewMotionMap(bounds)

// 	forRect(bounds, 1, func(x, y int) {
// 		v0 := GetMotion(m0, x, y)

// 		if v0.X != 0 || v0.Y != 0 {
// 			add(x, y, g0, v0, msrc)
// 		}

// 		v1 := GetMotion(m1, x, y)

// 		if v1.X != 0 || v1.Y != 0 {
// 			add(x, y, g1, v1, mdst)
// 		}
// 	})

// 	findNext := func(x, y int, g *image.Gray, m *MotionMap) {

// 		if x&1 != 0 || y&1 != 0 {
// 			SetMotion(m, x, y, image.Point{})
// 			return
// 		}

// 		v := GetMotion(m, x, y)
// 		d := maxInt(absInt(v.X), absInt(v.Y))

// 		if d < 1.0 {
// 			v = image.Point{}
// 		} else {
// 			// invert
// 			v.X, v.Y = -v.X, -v.Y

// 			if x == 40 && m == mdst {
// 				fmt.Println("track", x, y, v.X, v.Y, d)
// 			}

// 			p := image.Pt(x, y)
// 			r0 := uint8(0)

// 			for i := 0; p.In(g0.Bounds()); i++ {

// 				dx := (v.X * i) / d
// 				dy := (v.Y * i) / d

// 				r1 := grayAt(g, x+dx, y+dy).Y

// 				p.X, p.Y = x+dx, y+dy

// 				if x == 40 && m == mdst {
// 					fmt.Println(p.X, p.Y, dx, dy, r0, r1)
// 				}

// 				if r1 < r0 && r0 == uint8(255) {

// 					//dx = (dx + 1) &^ 1
// 					//dy = (dy + 1) &^ 1

// 					v.X, v.Y = dx, dy
// 					bg.SetNRGBA(x, y, rgbaAt2(img0, x+dx, y+dy))

// 					break
// 				}
// 				r0 = r1
// 			}

// 			if x == 40 && m == mdst {
// 				fmt.Println("set", x, y, v)
// 			}
// 		}

// 		SetMotion(m, x, y, v)
// 	}

// 	forRect(bounds.Inset(1), 1, func(x, y int) {
// 		findNext(x, y, g0, msrc)
// 		//fmt.Println(x, y, GetMotion(msrc, 86, 10))
// 		findNext(x, y, g1, mdst)
// 	})

// 	return

// 	//TODO: 
// 	// 1) if src moves but dst is missing mark with img1
// 	// 2) if dst moves but src is missing mark with img0

// 	// two backgrounds3
// }

type Motion struct {
	X, Y   int
	Dx, Dy int
	R0, R1 int
}

func SerializeMotions(m *MotionMap, g0, g1 *image.Gray) []Motion {

	result := make([]Motion, 0)

	forRect(m.Bounds.Intersect(g0.Bounds().Intersect(g1.Bounds())), 1, func(x, y int) {
		v := GetMotion(m, x, y)
		if v.X != 0 || v.Y != 0 {
			result = append(result,
				Motion{
					x, y,
					v.X, v.Y,
					int(255 - grayAt(g0, x, y).Y),
					int(255 - grayAt(g1, x+v.X, y+v.Y).Y),
				})
		}
	})

	return result
}

func InvertMotions(motions []Motion) []Motion {
	for i, m := range motions {
		motions[i] = Motion{
			m.X + m.Dx, m.Y + m.Dy,
			-m.Dx, -m.Dy,
			m.R1, m.R0,
		}
	}
	return motions
}

func RemoveDuplicates(motions []Motion) []Motion {
	result := make([]Motion, 0, len(motions))
	last := Motion{}
	for i, j := 0, 0; i < len(motions); i++ {
		m := motions[i]
		if last != m {
			result = append(result, m)
			j++
			last = m
		}
	}
	return result
}

type motions []Motion

func (m motions) Len() int      { return len(m) }
func (m motions) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m motions) Less(i, j int) bool {
	mi, mj := m[i], m[j]
	li := mi.Dx*mi.Dx + mi.Dy*mi.Dy
	lj := mj.Dx*mj.Dx + mj.Dy*mj.Dy
	return li > lj || (li == lj && mi.R0+mi.R1 > mj.R0+mj.R1)
}

func SortMotions(m []Motion) {
	sort.Sort(motions(m))
}

func ReadNRGBA(filename string) *image.NRGBA {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	pixmap := image.NewNRGBA(img.Bounds())
	draw.Draw(pixmap, img.Bounds(), img, img.Bounds().Min, draw.Src)

	return pixmap
}

func WritePng(filename string, img image.Image) {
	out, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	png.Encode(out, img)
}
