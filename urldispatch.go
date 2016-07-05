package urldispatch

import (
	"errors"
	"fmt"
	"strings"
)

const (
	nullptr = index(^uint8(0))
)

type segment struct {
	value string
	amap  argsMap
	next  []segment
}

type args2 struct {
	psection indexes
	params   []string
	asection indexes
	array    []string
}

func (a *args2) appendParamValue(value string) {

	pIdx := index(len(a.params))

	a.params = append(a.params, value)
	a.psection = append(a.psection, pIdx)
}

func (a *args2) appendArrayValue(value string) {
	a.array = append(a.array, value)
}

func (a *args2) nextArray() {
	aIdx := index(len(a.array))
	a.asection = append(a.asection, aIdx)
}

func (a *args2) addNullPtrParams(count index) {

	for i := index(0); i < count; i++ {
		// point to the largest index.
		a.psection = append(a.psection, nullptr)
	}
}

func (d *Dispatcher) addRoute(segs []segment) error {

	if len(segs) > 0 {

		cseg := segs[0]
		refmap := cseg.amap

		// check if all the amaps are equal
		for _, s := range segs {
			if !s.amap.eq(refmap) {
				return errors.New("amaps on the segments are not equal.")
			}
		}

		// check if the segments are addable
		err := d.root.addable2(segs, refmap, 0)
		if err != nil {
			return err
		}

		// insert the segments.
		d.root.insertSegments(segs)
	}

	return nil
}

func (d Dispatcher) dispatchPath(pathSegs []string) (Outargs, error) {

	ar := args2{}

	if len(pathSegs) > 0 {
		return d.root.dispatchPath(pathSegs, ar, -1)
	}

	return Outargs{}, errors.New("nothing to dispatch.")
}

func (s segment) dispatchPath(pathSegs []string, ar args2, idx int) (Outargs, error) {

	var pIdx index

	for len(pathSegs) > 0 {
		ps := pathSegs[0]

		for _, cs := range s.next {
			if cs.value == ps {

				if idx > -1 {
					pCount := cs.amap.psections[idx]

					// fix the array args
					ar.addNullPtrParams(pCount - pIdx)
				}

				if cs.amap.asections[idx+1] > 0 {
					ar.nextArray()
				}
				return cs.dispatchPath(pathSegs[1:], ar, idx+1)
			}
		}

		if idx == -1 {
			return Outargs{}, errors.New("nothing to dispatch for " + ps)
		}

		// if there is room for another param.
		hasRoomForParam, err := s.amap.psections.isItemAtIndexBigger(idx, pIdx)
		if err != nil {
			return Outargs{}, err
		}

		if hasRoomForParam {

			ar.appendParamValue(ps)
			pIdx += 1
			pathSegs = pathSegs[1:]

		} else if s.amap.asections[idx] > 0 {
			// has room for another array item.
			ar.appendArrayValue(ps)
			pathSegs = pathSegs[1:]
		} else {
			return Outargs{}, errors.New("param overflow with segment:" + ps + fmt.Sprintf(" idx:%v", idx))
		}
	}

	if len(s.next) > 0 {
		return Outargs{}, errors.New("only partial dispatch.")
	} else {
		return Outargs{amap: s.amap, ar: ar}, nil
	}
}

func dispatchQuery(query string, am argsMap, ar args2, idx index) (args2, error) {

	pc := index(len(am.params)) - idx
	ar.addNullPtrParams(pc)

	kvs := strings.Split(query, "&")

	for _, rkv := range kvs {
		kv := strings.Split(rkv, "=")
		if len(kv) != 2 {
			return ar, errors.New("invalid query")
		}

		for i := idx; i < pc+idx; i++ {

			if kv[0] == am.params[i] {
				if ar.psection[i] == nullptr {
					ar.psection[i] = index(len(ar.params))
					ar.params = append(ar.params, kv[1])
					continue
				}
			}
		}
	}

	return ar, nil
}

func (s segment) addable2(segs []segment, amap argsMap, index int) error {

	if len(segs) > 0 {
		for _, cs := range s.next {
			cseg := segs[0]
			if cs.value == cseg.value {

				err := cs.compareParams(amap, index)
				if err != nil {
					return err
				}

				// only check if the segment is not the root segment.
				if index != 0 {
					if len(segs) < 2 {
						return errors.New(`partial route`)
					}

					if len(cs.next) == 0 {
						return errors.New(`partial route detected`)
					}
				}

				return cs.addable2(segs[1:], amap, index+1)
			}
		}
	}
	return nil
}

func (s segment) compareParams(amap argsMap, index int) error {
	eq, err := s.amap.compareAtIndex(amap, index)
	if err != nil {
		return err
	}

	if !eq {
		return errors.New("segment is equal but params differ.")
	}

	return nil
}

func (s *segment) insertSegments(segments []segment) {

	if len(segments) > 0 {
		nseg := segments[0]

		for i, cs := range s.next {
			if cs.value == nseg.value {

				s.next[i].insertSegments(segments[1:])
				return
			}
		}

		s.next = append(s.next, nseg)

		// insert from the appended struct.
		s.next[len(s.next)-1].insertSegments(segments[1:])
		return
	}
}
