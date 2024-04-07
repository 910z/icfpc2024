package main

import (
	"context"
	"database/sql"
	"fmt"
	"icfpc/algorithms"
	"icfpc/types"
	"reflect"
	"time"

	"github.com/uptrace/bun"
)

type ErrorExpl struct {
	types.Explanation
	Error string
}

func runAlgorithms(db *bun.DB) {
	ctx := context.Background()
	for {
		tasks := getTasks()
		algorithms := getAlgorithms()

		for i := range tasks {
			for _, algorithm := range algorithms {
				runRes := types.RunResult{
					TaskId:           tasks[i].Id,
					AlgorithmName:    reflect.TypeOf(algorithm).String(),
					AlgorithmVersion: algorithm.Version(),
					StartedAt:        time.Now().UTC(),
					RunStatus:        types.RunStatusStarted,
				}
				err := db.NewInsert().Model(&runRes).
					Returning("id").
					Ignore().
					Scan(ctx, &runRes.Id)
				if err == sql.ErrNoRows { // уже запущено
					continue
				}
				if err != nil {
					panic(err)
				}

				go func() {
					defer func() {
						if err := recover(); err != nil {
							runRes.RunStatus = types.RunStatusError
							runRes.FinishedAt = time.Now().UTC()
							expl := ErrorExpl{Error: fmt.Sprint(err)}
							runRes.Explanation = &expl
							_, err := db.NewUpdate().Model(&runRes).WherePK().Exec(ctx)
							if err != nil {
								panic(err)
							}
						}
					}()

					task := &tasks[i]
					solution, explanation := algorithm.Solve(task)

					runRes.FinishedAt = time.Now().UTC()
					runRes.Solution = solution
					runRes.Explanation = explanation
					runRes.RunStatus = types.RunStatusFinished

					res, err := db.NewUpdate().Model(&runRes).WherePK().Exec(ctx)
					if err != nil {
						panic(err)
					}
					rows, err := res.RowsAffected()
					if err != nil {
						panic(err)
					}
					if rows == 0 {
						panic(fmt.Sprintf("no rows updated: %v", runRes))
					}
				}()
			}
		}
		time.Sleep(time.Second)
	}
}

func getTasks() []types.Task {
	return []types.Task{
		{Id: "1", Description: "Task1", Data: 5},
		{Id: "2", Description: "Task2", Data: 10},
	}
}

func getAlgorithms() []algorithms.Algorithm {
	return []algorithms.Algorithm{
		&algorithms.Doubler{},
		&algorithms.Tripler{},
	}
}
