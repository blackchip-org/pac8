package mach

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
