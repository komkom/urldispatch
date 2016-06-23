package urldispatch

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Dispatcher struct {
	segments []segment
}

type Outargs struct {
	amap argsMap
	ar   args2
}

func (o Outargs) ParamCount() int {
	return len(o.amap.params)
}

func (o Outargs) Value(index int) (string, error) {

	if len(o.ar.psection) <= index {
		return ``, errors.New("value index out of bounds.")
	}

	idx := o.ar.psection[index]

	if len(o.ar.params) <= int(idx) {
		return ``, errors.New("value index out of bounds 2.")
	}

	return o.ar.params[idx], nil
}

func (o Outargs) Array(index int) ([]string, error) {

	if len(o.ar.asection) <= index {
		return nil, errors.New("value index out of bounds 2.")
	}

	sIdx := int(o.ar.asection[index])
	end := len(o.ar.asection)
	nextIdx := index + 1
	if end > nextIdx {
		end = int(o.ar.asection[index+1])
	}

	return o.ar.array[sIdx:end], nil
}

// TODO remove this test.
func Test(dispatchURL *url.URL, u *url.URL) {

	d := Dispatcher{}
	err := d.AddRoute(dispatchURL)
	if err != nil {
		panic(err)
	}

	oa, err := d.Dispatch(u)
	if err != nil {
		panic(err)
	}

	fmt.Printf("__args %v\n", oa.ar)
	/*
		rootSeg := segment{}

		err := rootSeg.AddDispatchPath(dispatchURL)
		if err != nil {
			panic(err)
		}

		args, err := rootSeg.Dispatch(u)
		if err != nil {
			panic(err)
		}

		fmt.Printf("args: %v", args)
	*/
}

func (d *Dispatcher) Dispatch(dispatch *url.URL) (Outargs, error) {
	dispPath := dispatch.Path
	dispQuery := dispatch.RawQuery

	if strings.HasPrefix(dispPath, "/") {
		dispPath = dispPath[1:]
	}

	oa, err := d.dispatchPath(strings.Split(dispPath, "/"))
	if err != nil {
		return oa, err
	}

	if len(dispQuery) > 0 {
		ar, err := dispatchQuery(dispQuery, oa.amap, oa.ar, index(len(oa.ar.params)))
		if err != nil {
			return Outargs{}, err
		}

		oa.ar = ar
	}

	return oa, nil
}

func (d *Dispatcher) AddRoute(route *url.URL) error {

	routePath := route.Path
	if strings.HasPrefix(routePath, "/") {
		routePath = routePath[1:]
	}

	segs, err := parse(routePath, route.RawQuery)
	if err != nil {
		return err
	}

	err = d.addRoute(segs)
	if err != nil {
		return err
	}

	return nil
}
