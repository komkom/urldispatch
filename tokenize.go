package urldispatch

import (
	"errors"
	"strings"
)

func tokenize(path string, rawQuery string) ([]segment, []string, error) {

	argCache := map[string]bool{}
	var segments = []segment{}
	var queryParams []string

	queryParams = strings.Split(rawQuery, "&")

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

		if strings.HasPrefix(item, ":") {
			if cseg == nil {
				return nil, nil, errors.New("trying to insert root param.")
			}
			pn := item[1:]

			err := cacheParam(pn)
			if err != nil {
				return nil, nil, err
			}

			cseg.paramNames = append(cseg.paramNames, pn)
		} else if strings.HasPrefix(item, "l:") {
			if cseg == nil {
				return nil, nil, errors.New("trying to insert root array.")
			}

			pn := item[2:]
			err := cacheParam(pn)
			if err != nil {
				return nil, nil, err
			}

			cseg.arrayName = pn
		} else {
			nseg := segment{value: item}
			segments = append(segments, nseg)
		}
	}

	// check the query params
	for _, qp := range queryParams {
		err := cacheParam(qp)
		if err != nil {
			return nil, nil, err
		}
	}

	return segments, queryParams, nil
}
