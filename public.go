package urldispatch

import "fmt"

// TODO remove this test.
func Test(dispathPath string, queryPath string) {

	rootSeg := segment{}

	err := rootSeg.AddDispatchPath(dispathPath)
	if err != nil {
		panic(err)
	}

	args, err := rootSeg.dispatch(queryPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("args: %v", args)
}

func (s *segment) AddDispatchPath(dispatchPath string) error {
	segs, qParamNames, err := tokenize(dispatchPath)
	if err != nil {
		return err
	}

	err = s.addSegments(segs, qParamNames)
	if err != nil {
		return err
	}

	return nil
}

func (s *segment) Dispatch(queryPath string) (args, error) {
	return s.dispatch(queryPath)
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
