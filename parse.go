package urldispatch

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type index uint8
type indexes []index
type params []string

func (i indexes) incrLast() error {
	if len(i) > 0 {
		i[len(i)-1] += 1
		return nil
	}

	return errors.New("indexes have 0 length.")
}

func (i indexes) last() (index, error) {
	if len(i) == 0 {
		return 0, errors.New("zero elements.")
	}

	return i[len(i)-1], nil
}

func (p params) eq(other params) bool {
	if len(p) != len(other) {
		return false
	}

	for idx, item := range p {
		if item != other[idx] {
			return false
		}
	}

	return true
}

func (i indexes) eq(other indexes) bool {
	if len(i) != len(other) {
		return false
	}

	for idx, item := range i {
		if item != other[idx] {
			return false
		}
	}

	return true
}

func (i indexes) isItemAtIndexBigger(index int, toCompare index) (bool, error) {
	if len(i) <= index {
		return false, errors.New("paramAppendable index out of bounds.")
	}

	return (i[index] > toCompare), nil
}

type argsMap struct {
	psections indexes
	params    params
	asections indexes
	arrays    params
	tag       int
}

func (am *argsMap) incrLastPsection() error {
	if len(am.psections) > 0 {
		am.psections[len(am.psections)-1] += 1
		return nil
	}

	return errors.New("psections has len 0.")
}

func (am argsMap) compareAtIndex(other argsMap, idx int) (bool, error) {

	if len(am.psections) <= idx || len(am.asections) <= idx || len(other.psections) <= idx || len(other.asections) <= idx {
		return false, errors.New("compareAtIndex out of bounds.")
	}

	return am.psections[idx] == other.psections[idx] && am.asections[idx] == other.asections[idx], nil
}

func (am argsMap) eq(other argsMap) bool {
	return am.params.eq(other.params) && am.psections.eq(other.psections) && am.arrays.eq(other.arrays) && am.asections.eq(other.asections)
}

func parse(path string, rawQuery string, tag int) ([]segment, error) {

	argCache := map[string]bool{}
	var segments = []segment{}
	var queryParams []string
	amap := argsMap{}
	amap.tag = tag

	if len(rawQuery) > 0 {
		queryParams = strings.Split(rawQuery, "&")
	}

	poke := func() *segment {
		idx := len(segments) - 1
		if idx >= 0 {
			return &segments[idx]
		}

		return nil
	}

	cacheParam := func(paramName string) error {
		_, ok := argCache[paramName]
		if ok {
			return errors.New("paramname no unique:" + paramName)
		}
		argCache[paramName] = true
		return nil
	}

	for _, item := range strings.Split(path, "/") {

		cseg := poke()

		if strings.HasPrefix(item, ":") && strings.HasSuffix(item, "...") {
			if cseg == nil {
				return nil, errors.New("trying to insert root array.")
			}

			pn := item[1 : utf8.RuneCountInString(item)-3]

			err := cacheParam(pn)
			if err != nil {
				return nil, err
			}

			amap.arrays = append(amap.arrays, pn)
			amap.asections.incrLast()

			l, err := amap.asections.last()
			if err != nil {
				return nil, err
			}

			if l > 1 {
				return nil, errors.New("multiple arrays in one segment.")
			}

		} else if strings.HasPrefix(item, ":") {

			if cseg == nil {
				return nil, errors.New("trying to insert root param.")
			}
			pn := item[1:]

			err := cacheParam(pn)
			if err != nil {
				return nil, err
			}

			amap.params = append(amap.params, pn)
			err = amap.psections.incrLast()
			if err != nil {
				return nil, err
			}

			// make sure an array has not already been added.
			l, err := amap.asections.last()
			if err != nil {
				return nil, err
			}

			if l == 1 {
				return nil, errors.New("an array has to be the last param on a segment.")
			}

		} else {
			nseg := segment{value: item}
			segments = append(segments, nseg)

			amap.psections = append(amap.psections, 0)
			amap.asections = append(amap.asections, 0)
		}
	}

	// check the query params
	if len(queryParams) > 0 {
		amap.psections = append(amap.psections, 0)

	}

	for _, qp := range queryParams {
		err := cacheParam(qp)
		if err != nil {
			return nil, err
		}

		amap.psections.incrLast()
		amap.params = append(amap.params, qp)
	}

	// add the amap to the segments.
	for idx, _ := range segments {
		segments[idx].amap = amap
	}

	return segments, nil
}
