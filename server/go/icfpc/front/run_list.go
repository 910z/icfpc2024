package front

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type RunList struct {
	app.Compo

	runItems []runListItem
}

func (l *RunList) OnMount(_ app.Context) {
	res, err := http.Get("/list")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if err = json.NewDecoder(res.Body).Decode(&l.runItems); err != nil {
		panic(err)
	}
}

func (l *RunList) Render() app.UI {
	return app.Div().Body(app.Table().Body(
		app.THead().Body(
			app.Tr().Body(
				app.Th().Text("Task"),
				app.Th().Text("Algorithm"),
				app.Th().Text("Score"),
			),
		),
		app.TBody().Body(l.items()),
	), &Visualization{})
}

func (l *RunList) items() app.UI {
	return app.Range(l.runItems).Slice(func(i int) app.UI {
		item := l.runItems[i]

		return app.Tr().Body(
			app.Td().Text(item.TaskID),
			app.Td().Text(fmt.Sprintf("%s %s", item.AlgorithmName, item.AlgorithmVersion)),
			app.Td().Text(fmt.Sprint(item.Score)),
		)
	})
}
