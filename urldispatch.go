package urldispatch

import (
	"errors"
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

type outargs struct {
	amap argsMap
	ar   args2
}

func (o outargs) paramCount() int {
	return len(o.amap.params)
}

func (o outargs) value(index int) string {
	return o.ar.params[o.ar.psection[index]]
}

func (o outargs) array(index int) []string {
	sIdx := int(o.ar.asection[index])
	end := len(o.ar.asection)
	nextIdx := index + 1
	if end > nextIdx {
		end = int(o.ar.asection[index+1])
	}

	return o.ar.array[sIdx:end]
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

func (s *segment) addRoute(segs []segment) error {

	if len(segs) > 0 {

		refmap := segs[0].amap
		for _, s := range segs {
			if !s.amap.eq(refmap) {
				return errors.New("amaps on the segments are not equal.")
			}
		}

		err := s.addable2(segs, refmap, 0)
		if err != nil {
			return err
		}

		// insert the segments
		s.insertSegments(segs)
	}

	return nil
}

func (s *segment) dispatchPath(pathSegs []string, ar args2, idx int) (outargs, error) {

	var pIdx index

	for len(pathSegs) > 0 {
		ps := pathSegs[0]

		for _, cs := range s.next {
			if cs.value == ps {

				pCount := s.amap.psections[idx]

				// fix the array args
				ar.addNullPtrParams(pCount - pIdx)
				ar.nextArray()

				return cs.dispatchPath(pathSegs[1:], ar, idx+1)
			}
		}

		// if there is room for another param.
		hasRoomForParam, err := s.amap.psections.isItemAtIndexBigger(idx, pIdx)
		if err != nil {
			return outargs{}, err
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
			return outargs{}, errors.New("param overflow with segment:" + ps)
		}
	}

	if idx == 0 {
		return outargs{}, errors.New("nothing to dispatch.")
	} else {
		return outargs{amap: s.amap, ar: ar}, nil
	}
}

func (s *segment) dispatchQuery(query string, am argsMap, ar args2, idx index) (args2, error) {

	pc := index(len(am.params)) - idx
	ar.addNullPtrParams(pc)

	kvs := strings.Split(query, "&")
	for _, rkv := range kvs {
		kv := strings.Split(rkv, "=")
		if len(kv) != 2 {
			return ar, errors.New("invalid query")
		}

		for i := idx; i < pc; i++ {
			if kv[0] == am.params[i] {

				if ar.psection[i] != nullptr {
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

				err := cs.ifParamsEqualGetNext(cseg, amap, index)
				if err != nil {
					return err
				}

				return cs.addable2(segs[1:], amap, index+1)
			}
		}
	}
	return nil
}

func (s segment) ifParamsEqualGetNext(other segment, amap argsMap, index int) error {
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

		for _, cs := range s.next {
			if cs.value == nseg.value {

				cs.insertSegments(segments[1:])
				return
			}
		}

		s.next = append(s.next, nseg)

		// insert from the appended struct.
		s.next[len(s.next)-1].insertSegments(segments[1:])
		return
	}
}
