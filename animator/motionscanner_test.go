package animator

import (
	"code.google.com/p/draw2d/draw2d"
	"fmt"
	"image"
	"image/color"
	. "launchpad.net/gocheck"
	//"math"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type imgSuite struct{}

var _ = Suite(&imgSuite{})

func (s *imgSuite) TestMotion(c *C) {

	img0 := ReadNRGBA("test1.png")

	// 15 * 15, 
	maxDelta := float64(225)

	detectEdge := func(c0, c1 color.NRGBA) color.Gray {
		if delta(c0, c1) > maxDelta {
			return color.Gray{255}
		}
		return color.Gray{0}
	}

	edg0 := MarkEdges(img0, detectEdge)
	WritePng("edges0.png", scaleGray(edg0, 10))
	GradientInflate(edg0)
	// WritePng("gradient0.png", scaleGray(edg0, 10))
	ske0 := Skelton(img0, edg0)
	// WritePng("skeleton0.png", ske0)

	img1 := ReadNRGBA("test2.png")
	edg1 := MarkEdges(img1, detectEdge)
	WritePng("edges1.png", scaleGray(edg1, 10))
	GradientInflate(edg1)
	// WritePng("gradient1.png", scaleGray(edg1, 10))
	ske1 := Skelton(img1, edg1)
	// WritePng("skeleton1.png", ske1)

	m0 := NearestMotion(ske0, ske1, maxDelta, delta)
	// dists0 := MotionDistance(m0)
	src0, md0 := MotionSource(m0)
	// WritePng("density0.png", scaleGray(md0, 10))
	shift0 := GravityShift(ske1, md0, delta)
	// WritePng("motion0.png", drawMotionMap(ske0, ske1, m0, shift0))
	// WritePng("source0.png", drawMotionMap(ske0, ske1, m0, src0))
	// WritePng("distance0.png", scaleGray(dists0, 10))

	m1 := NearestMotion(ske1, ske0, maxDelta, delta)
	// dists1 := MotionDistance(m1)
	src1, md1 := MotionSource(m1)
	// WritePng("density1.png", scaleGray(md1, 10))
	shift1 := GravityShift(ske0, md1, delta)
	// WritePng("motion1.png", drawMotionMap(ske1, ske0, m1, shift1))
	// WritePng("source1.png", drawMotionMap(ske1, ske0, m1, src1))
	// WritePng("distance1.png", scaleGray(dists1, 10))

	skeShift0 := SmoothSkeletonShift(ske0)
	// WritePng("smoothSkeleton0.png", drawMotionMap(ske0, ske1, m0, skeShift0))
	skeShift1 := SmoothSkeletonShift(ske1)
	// WritePng("smoothSkeleton1.png", drawMotionMap(ske1, ske0, m1, skeShift1))

	// spread0 := SpreadMotion(m0, src0, shift0, skeShift0, skeShift1, ske1, delta)
	SpreadMotion(m0, src0, shift0, skeShift0, skeShift1, ske1, delta)
	// WritePng("spread0.png", drawMotionMap(ske0, ske1, m0, spread0))
	// spread1 := SpreadMotion(m1, src1, shift1, skeShift1, skeShift0, ske0, delta)
	SpreadMotion(m1, src1, shift1, skeShift1, skeShift0, ske0, delta)
	// WritePng("spread1.png", drawMotionMap(ske1, ske0, m1, spread1))

	OptimizeMotion(m0, m1, ske1, maxDelta, delta)
	OptimizeMotion(m1, m0, ske0, maxDelta, delta)

	WritePng("opt0.png", drawMotionMap(ske0, ske1, m0, nil))
	WritePng("opt1.png", drawMotionMap(ske1, ske0, m1, nil))

	bg0, bg1, pshift0, pshift1 := CreateBackground(img0, img1, m0, m1, edg0, edg1)

	WritePng("bg0.png", bg0)
	WritePng("bg1.png", bg1)
	WritePng("bgShift0.png", scaleGray(pshift0, 10))
	WritePng("bgShift1.png", scaleGray(pshift1, 10))

	motions := append(SerializeMotions(m0, edg0, edg1), InvertMotions(SerializeMotions(m1, edg1, edg0))...)
	SortMotions(motions)
	motions = RemoveDuplicates(motions)

	fmt.Println(motions)

	writeMotions("m0.txt", motions)

	//fmt.Println(motions)

	// WritePng("bg0.png", drawMotion(edg0, 10, pshift0))
	// WritePng("bg1.png", drawMotion(edg1, 10, pshift1))

	//WritePng("bg0.png", scaleGray(gm, 60))

	// OptimizeMotion(m0, m1, ske1, maxDelta, delta)
	// WritePng("opt2.png", drawMotionMap(ske0, ske1, m0, nil))

	// TransposeMotion(m0, m1, ske1, delta)
	// WritePng("transpose0.png", drawMotionMap(ske0, ske1, m0, nil))

	// TransposeMotion(m1, m0, ske0, delta)
	// WritePng("transpose1.png", drawMotionMap(ske1, ske0, m1, nil))

	// TransposeMotion(m0, m1, ske1, delta)
	// WritePng("transpose2.png", drawMotionMap(ske0, ske1, m0, nil))
}

func delta(c1, c2 color.NRGBA) float64 {

	if c1.A == 0 && c2.A == 0 {
		return 0.0
	}

	y1, cb1, cr1 := color.RGBToYCbCr(c1.R, c1.G, c1.B)
	y2, cb2, cr2 := color.RGBToYCbCr(c2.R, c2.G, c2.B)

	dy := (int(y1) - int(y2)) * 3
	dcb := int(cb1) - int(cb2)
	dcr := int(cr1) - int(cr2)
	da := int(c1.A) - int(c2.A)

	// if c2.A == 233 && c1.A == 255 {
	// 	fmt.Println(dy, dcb, dcr, da)
	// }

	// if c1.A == 255 && c2.A == 255 && c2.R > 0 {
	// 	fmt.Println(c1, c2, dy, dcb, dcr, da, float64(dy*dy+dcb*dcb+dcr*dcr+da*da)/32.0)
	// }

	return float64(dy*dy+dcb*dcb+dcr*dcr+da*da) / 32.0
}

func scaleGray(src *image.Gray, scale uint8) *image.Gray {
	bounds := src.Bounds()

	dst := image.NewGray(bounds)
	copy(dst.Pix, src.Pix)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.SetGray(x, y, color.Gray{grayAt(src, x, y).Y * scale})
		}
	}
	return dst
}

func drawMotionMap(ske0, ske1 *image.NRGBA, m *MotionMap, shift *ShiftMap) (img *image.RGBA) {

	bounds := ske0.Bounds().Intersect(ske1.Bounds()).Intersect(m.Bounds)

	if shift != nil {
		bounds = bounds.Intersect(shift.Bounds)
	}

	s := 10

	img = image.NewRGBA(image.Rectangle{bounds.Min.Mul(s), bounds.Max.Mul(s)})
	gc := draw2d.NewGraphicContext(img)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := color.NRGBA{255, 255, 255, 255}

			if ske0.At(x/s, y/s).(color.NRGBA).A > 0 {
				c.R = 0
			}

			if ske1.At(x/s, y/s).(color.NRGBA).A > 0 {
				c.G = 0
			}

			img.Set(x, y, c)
		}
	}

	count := 0

	for y := m.Bounds.Min.Y; y < m.Bounds.Max.Y; y++ {
		for x := m.Bounds.Min.X; x < m.Bounds.Max.X; x++ {
			v := GetMotion(m, x, y)

			if v.X != 0 || v.Y != 0 {
				cx, cy := x*s+s/2, y*s+s/2
				dx, dy := v.X*s, v.Y*s
				gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
				gc.MoveTo(float64(cx), float64(cy))
				gc.LineTo(float64(cx+dx), float64(cy+dy))
				gc.Stroke()
				count++
			}

			if shift != nil {
				sv := GetShift(shift, x, y)

				if sv.Dx != 0.0 || sv.Dy != 0.0 {
					cx, cy := x*s+s/2, y*s+s/2
					sf := float64(s)
					//fmt.Println(x, y, sv)
					gc.SetStrokeColor(color.RGBA{0, 200, 0, 255})
					gc.MoveTo(float64(cx), float64(cy))
					gc.LineTo(float64(cx)+sv.Dx*sf, float64(cy)+sv.Dy*sf)
					gc.Stroke()
				}
			}
		}
	}

	fmt.Println("count:", count)

	return
}

func drawMotion(src image.Image, cs uint8, m *MotionMap) (img *image.RGBA) {

	bounds := src.Bounds().Intersect(m.Bounds)

	s := 10

	img = image.NewRGBA(image.Rectangle{bounds.Min.Mul(s), bounds.Max.Mul(s)})
	gc := draw2d.NewGraphicContext(img)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := src.At(x/s, y/s).RGBA()
			img.Set(x, y, color.RGBA{uint8(r) * cs, uint8(g) * cs, uint8(b) * cs, uint8(a) * cs})
		}
	}

	for y := m.Bounds.Min.Y; y < m.Bounds.Max.Y; y++ {
		for x := m.Bounds.Min.X; x < m.Bounds.Max.X; x++ {
			v := GetMotion(m, x, y)

			// if y == 10 {
			// 	fmt.Println(x, y, v)
			// }

			if v.X != 0 || v.Y != 0 {
				cx, cy := x*s+s/2, y*s+s/2
				dx, dy := v.X*s, v.Y*s
				gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
				gc.MoveTo(float64(cx), float64(cy))
				gc.LineTo(float64(cx+dx), float64(cy+dy))
				gc.Stroke()
			}
		}
	}

	return
}

func writeMotions(filename string, motions []Motion) {
	out, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	arr := make([]int, len(motions)*6)

	for i, m := range motions {
		arr[i*6+0] = m.X + 1
		arr[i*6+1] = m.Y + 1
		arr[i*6+2] = m.Dx
		arr[i*6+3] = m.Dy
		arr[i*6+4] = m.R0
		arr[i*6+5] = m.R1
	}

	jsonStr, err := json.Marshal(arr)

	fmt.Println(len(motions))

	if err != nil {
		log.Fatal(err)
	}

	out.Write(jsonStr)
}
