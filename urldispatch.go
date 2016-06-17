package urldispatch

import (
	"errors"
	"fmt"
)

type segment struct {
	hash        string
	value       string
	paramNames  []string
	arrayName   string
	paramValues []string
	arrayValues []string
	next        []segment
	tag         string
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

func rootSegment() segment {
	return segment{}
}

func (s *segment) addSegments(segments []segment) error {
	preparedSegs, err := s.prepare(segments)
	if err != nil {
		return err
	}

	s.insertSegments(preparedSegs)
	return nil
}

func (s *segment) dispatch(path []string) (args, error) {

	cargs := args{}

	for _, pathSeg := range path {
		for _, cs := range s.next {
			if cs.value == path[0] {
				ccargs, err := cs.dispatch(path[1:])
				if err != nil {
					return cargs.merge(ccargs), err
				}
			}
		}

		cargs, err := s.merge(pathSeg, cargs)
		if err != nil {
			return cargs, err
		}
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

func (s segment) merge(arg string, cargs args) (args, error) {
	cpCount := len(cargs.params)

	if len(s.paramNames) > cpCount {

		cargs.params = append(cargs.params, keyValue{key: s.paramNames[cpCount], value: arg})

		return cargs, nil
	}

	if len(s.arrayName) > 0 {

		// add a new offset if needed
		lastOffsetIdx := len(cargs.arrayOffsets) - 1
		if lastOffsetIdx > 0 {
			if cargs.arrayOffsets[lastOffsetIdx].key != s.arrayName {
				cargs.arrayOffsets = append(cargs.arrayOffsets, offset{key: s.arrayName, length: 0})
			}
		}

		// fix offset
		offestIdx := len(cargs.arrayOffsets) - 1
		cargs.arrayOffsets[offestIdx].length += 1

		// add argument.
		cargs.arrayValues = append(cargs.arrayValues, arg)

		return cargs, nil
	}

	return args{}, errors.New(fmt.Sprintf("merge failed for %v, with %v", s.hash, arg))
}

func hash(s segment) (string, error) {

	if len(s.value) == 0 {
		return ``, errors.New(`value has len 0.`)
	}

	return fmt.Sprintf("%v#%v#%v", s.value, len(s.paramNames), len(s.arrayName) > 0), nil
}

func (s *segment) prepare(segments []segment) ([]segment, error) {

	ok := checkUniquenessOfParamNames(segments)
	if !ok {
		return nil, errors.New("param names not unique.")
	}

	if len(segments) > 0 {
		nseg := segments[0]
		h, err := hash(nseg)
		if err != nil {
			return nil, err
		}

		nseg.hash = h
		nseg.next = []segment{}

		for _, cs := range s.next {
			if cs.hash == nseg.hash {
				rsegs, err := cs.prepare(segments[1:])
				if err != nil {
					return nil, err
				}

				return append([]segment{cs}, rsegs...), nil
			}
		}

		return []segment{nseg}, nil
	}

	return nil, errors.New("prepare failed.")
}

// TODO implement.
func checkUniquenessOfParamNames(segments []segment) bool {
	return true
}

func (s *segment) insertSegments(segments []segment) {
	if len(segments) > 0 {
		nseg := segments[0]
		if len(nseg.hash) == 0 {
			panic(errors.New("hash has len 0."))
		}

		for _, cs := range s.next {
			if cs.hash == nseg.hash {
				cs.insertSegments(segments[1:])
				return
			}
		}

		s.next = append(s.next, nseg)
		nseg.insertSegments(segments[1:])
		return
	}
}
