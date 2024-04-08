package algorithms

import "reflect"

func GetName(a IAlgorithm) string {
	return reflect.TypeOf(a).String()
}
