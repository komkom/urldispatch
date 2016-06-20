package urldispatch

import (
	"errors"
	"fmt"
	"strings"
)

type segment struct {
	value           string
	paramNames      []string
	queryParamNames []string
	arrayName       string
	paramValues     []string
	arrayValues     []string
	next            []segment
	tag             string
}

type keyValue struct {
	key   string
	value string
}

type offset struct {
	key    string
	length int
}

type args struct {
	tag          string
	params       []keyValue
	arrayOffsets []offset
	arrayValues  []string
}

func (a args) merge(toAppend args) args {

	a.params = append(a.params, toAppend.params...)
	a.arrayOffsets = append(a.arrayOffsets, toAppend.arrayOffsets...)
	a.arrayValues = append(a.arrayValues, toAppend.arrayValues...)
	return a
}

func (s *segment) addSegments(segments []segment, queryParamNames []string) error {
	err := s.addable(segments)
	if err != nil {
		return err
	}

	s.queryParamNames = queryParamNames

	s.insertSegments(segments)
	return nil
}

func (s *segment) dispatch(path string, query string) (args, error) {
	cargs := args{}

	pargs, err := s.dispatchPath(strings.Split(path, "/"))
	if err != nil {
		return pargs, err
	}
	cargs = cargs.merge(pargs)

	qargs, err := s.dispatchQuery(query, cargs)
	if err != nil {
		return cargs, err
	}
	cargs = qargs

	return cargs, nil
}

func (s *segment) dispatchPath(path []string) (args, error) {

	cargs := args{}

	for len(path) > 0 {
		pathSeg := path[0]

		for _, cs := range s.next {

			if cs.value == path[0] {
				ccargs, err := cs.dispatchPath(path[1:])
				if err != nil {
					return ccargs, err
				} else {

					return cargs.merge(ccargs), nil
				}
			}
		}

		margs, err := s.merge(pathSeg, cargs)

		if err != nil {
			return cargs, err
		}

		cargs = margs
		path = path[1:]
	}

	// check if recursion needs to terminate.
	if len(path) == 0 {
		if s.next == nil {
			// success dispatch candidate found.
			cargs.tag = s.tag
			return cargs, nil
		} else {
			return cargs, errors.New("not the final segment.")
		}
	}

	return args{}, errors.New("path has len 0.")
}

func (s *segment) dispatchQuery(query string, cargs args) (args, error) {
	kvs := strings.Split(query, "&")

	for len(kvs) > 0 {
		kv := strings.Split(kvs[0], "=")
		if len(kv) != 2 {
			return cargs, errors.New("query syntax invalid")
		}

		for _, qp := range s.queryParamNames {
			if qp == kv[0] {
				cargs.params = append(cargs.params, keyValue{key: qp, value: kv[1]})
			}
		}

		kvs = kvs[1:]
	}

	return cargs, nil
}

func (s segment) merge(arg string, cargs args) (args, error) {
	cpCount := len(cargs.params)

	if len(s.paramNames) > cpCount {

		cargs.params = append(cargs.params, keyValue{key: s.paramNames[cpCount], value: arg})

		return cargs, nil
	}

	if len(s.arrayName) > 0 {

		// add a new offset if needed
		lastOffsetIdx := len(cargs.arrayOffsets) - 1
		if lastOffsetIdx >= 0 {
			if cargs.arrayOffsets[lastOffsetIdx].key != s.arrayName {
				cargs.arrayOffsets = append(cargs.arrayOffsets, offset{key: s.arrayName, length: 0})
			}
		} else {
			cargs.arrayOffsets = append(cargs.arrayOffsets, offset{key: s.arrayName, length: 0})
		}

		// fix offset
		offestIdx := len(cargs.arrayOffsets) - 1
		cargs.arrayOffsets[offestIdx].length += 1

		// add argument.
		cargs.arrayValues = append(cargs.arrayValues, arg)

		return cargs, nil
	}

	return args{}, errors.New(fmt.Sprintf("merge failed for %v, with %v", s.value, arg))
}

func (s *segment) addable(segments []segment) error {

	if len(segments) > 0 {
		nseg := segments[0]

		for _, cs := range s.next {
			if cs.value == nseg.value {
				return cs.addable(segments[1:])
			}
		}

		return nil
	}

	return errors.New("addable failed.")
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
