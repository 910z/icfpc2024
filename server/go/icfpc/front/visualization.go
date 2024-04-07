package front

import (
	"icfpc/algorithms"
	"icfpc/database"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Visualization struct {
	app.Compo

	task     database.Task
	solution algorithms.Solution
}

func (v *Visualization) OnMount(ctx app.Context) {
	pixels := make([]uint8, 40000)

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			pixels[y*400+x*4] = 0
			pixels[y*400+x*4+1] = 0
			pixels[y*400+x*4+2] = 0
			pixels[y*400+x*4+3] = 255
		}
	}

	uints := app.Window().Get("Uint8ClampedArray").New(40000)
	app.CopyBytesToJS(uints, pixels)

	data := app.Window().Get("ImageData").New(uints, 100)
	app.Window().Set("imgData", data)

	ctx2d := ctx.JSSrc().Call("getContext", "2d")
	ctx2d.Call("putImageData", data, 0, 0)
}

func (v *Visualization) Render() app.UI {
	return app.Canvas().Width(2000).Height(2000)
}
