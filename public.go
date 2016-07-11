package urldispatch

import (
	"errors"
	"net/url"
	"strings"
)

type Dispatcher struct {
	root segment
}

type Outargs struct {
	amap argsMap
	ar   args2
}

func (o Outargs) Tag() int {
	return o.amap.tag
}

func (o Outargs) ParamWithName(name string) (string, error) {
	idx := 0
	pIdx := 0
	ppIdx := index(0)

	for pIdx < len(o.amap.psections) {

		if o.amap.psections[pIdx] != 0 {
			if o.amap.params[idx] == name {
				return o.Value(idx)
			}
			idx += 1
		}

		ppIdx += 1

		ok, err := o.amap.psections.isItemAtIndexBigger(pIdx, ppIdx)
		if err != nil {
			return ``, err
		}

		if !ok {
			pIdx += 1
			ppIdx = 0
		}
	}

	return ``, errors.New("paramname not found.")
}

func (o Outargs) ArrayWithName(name string) ([]string, error) {

	idx := 0

	for idx < len(o.ar.asection) {

		if len(o.amap.arrays) <= idx {
			return nil, errors.New("asections not matching with arrays.")
		}

		if o.amap.arrays[idx] == name {
			return o.Array(idx)
		}
		idx += 1

	}

	return nil, errors.New("paramname not found.")
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
		return ``, errors.New("null ptr exception.")
	}

	return o.ar.params[idx], nil
}

func (o Outargs) Array(index int) ([]string, error) {

	if len(o.ar.asection) <= index {
		return nil, errors.New("value index out of bounds 2.")
	}

	sIdx := int(o.ar.asection[index])
	end := len(o.ar.array)

	nextIdx := index + 1
	if len(o.ar.asection) > nextIdx {
		end = int(o.ar.asection[nextIdx])
	}

	return o.ar.array[sIdx:end], nil
}

func (d *Dispatcher) Dispatch(dispatch *url.URL) (Outargs, error) {
	dispPath := dispatch.Path
	dispQuery := dispatch.RawQuery

	return d.DispatchPath(dispPath, dispQuery)
}

func (d *Dispatcher) DispatchPath(dispPath string, dispQuery string) (Outargs, error) {

	dispPath = removeLeadingAndTrailingSlashes(dispPath)

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

func (d *Dispatcher) AddRoute(route *url.URL, tag int) error {

	routePath := route.Path
	routePath = removeLeadingAndTrailingSlashes(routePath)

	segs, err := parse(routePath, route.RawQuery, tag)
	if err != nil {
		return err
	}

	err = d.addRoute(segs)
	if err != nil {
		return err
	}

	return nil
}

func removeLeadingAndTrailingSlashes(path string) string {

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return path
}
