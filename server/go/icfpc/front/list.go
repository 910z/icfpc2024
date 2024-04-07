package front

import (
	"encoding/json"
	"fmt"
	"icfpc/types"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type RunList struct {
	app.Compo
	runItems []types.RunListItem
}

func (h *RunList) OnMount(ctx app.Context) {
	res, err := http.Get("/list")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&h.runItems); err != nil {
		panic(err)
	}
}

func (h *RunList) items() app.UI {
	return app.Range(h.runItems).Slice(func(i int) app.UI {
		item := h.runItems[i]
		return app.Tr().Body(
			app.Td().Text(item.TaskId),
			app.Td().Text(fmt.Sprintf("%s %s", item.AlgorithmName, item.AlgorithmVersion)),
			app.Td().Text(fmt.Sprint(item.Score)),
		)
	})
}

func (h *RunList) Render() app.UI {
	return app.Div().Body(app.Table().Body(
		app.THead().Body(
			app.Tr().Body(
				app.Th().Text("Task"),
				app.Th().Text("Algorithm"),
				app.Th().Text("Score"),
			),
		),
		app.TBody().Body(h.items()),
	), &Visualization{})
}
