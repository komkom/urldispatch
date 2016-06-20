package urldispatch

import (
	"fmt"
	"net/url"
	"strings"
)

// TODO remove this test.
func Test(dispatchURL *url.URL, u *url.URL) {

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
}

func (s *segment) AddDispatchPath(dispatchURL *url.URL) error {

	dispatchPath := dispatchURL.Path
	if strings.HasPrefix(dispatchPath, "/") {
		dispatchPath = dispatchPath[1:]
	}

	segs, qParamNames, err := tokenize(dispatchPath, dispatchURL.RawQuery)
	if err != nil {
		return err
	}

	err = s.addSegments(segs, qParamNames)
	if err != nil {
		return err
	}

	return nil
}

func (s *segment) Dispatch(u *url.URL) (args, error) {

	path := u.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return s.dispatch(path, u.RawQuery)
}

func printSegment(s segment, intend int) {

	intendAction := func(intend int) string {
		intendString := ""
		for i := 0; i < intend; i++ {
			intendString += "\t"
		}
		return intendString
	}

	fmt.Println(intendAction(intend) + fmt.Sprintf("%v params:%v array:%v", s.value, s.paramNames, s.arrayName))

	for _, cs := range s.next {
		printSegment(cs, intend+1)
	}
}
