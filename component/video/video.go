package video

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Display interface {
	Render()
}

type NullDisplay struct{}

func (d NullDisplay) Render() {}

type RenderFrame struct {
	X     int32
	Y     int32
	W     int32
	H     int32
	Scale int32
}

func FitInWindow(winW int32, winH int32, w int32, h int32) RenderFrame {
	deltaW, deltaH := winW-w, winH-h
	scale := int32(1)
	if deltaW < deltaH {
		scale = winW / w
	} else {
		scale = winH / h
	}
	scaledW, scaledH := w*scale, h*scale
	return RenderFrame{
		X:     (winW - scaledW) / 2,
		Y:     (winH - scaledH) / 2,
		W:     w * scale,
		H:     h * scale,
		Scale: scale,
	}
}

func ScanLines(r *sdl.Renderer, w int32, h int32, size int32) (*sdl.Texture, error) {
	tex, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, w, h)
	if err != nil {
		return nil, err
	}

	r.SetRenderTarget(tex)
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x += 2 * size {
			r.SetDrawColorArray(0, 0, 0, 0)
			for i := int32(0); i < size; i++ {
				r.DrawPoint(x+i, y)
			}
			r.SetDrawColorArray(0, 0, 0, 0x20)
			for i := int32(size); i < size*2; i++ {
				r.DrawPoint(x+i, y)
			}
		}
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)
	return tex, nil
}
