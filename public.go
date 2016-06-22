package urldispatch

import (
	"net/url"
	"strings"
)

// TODO remove this test.
func Test(dispatchURL *url.URL, u *url.URL) {

	root := segment{}
	err := root.AddRoute(dispatchURL)
	if err != nil {
		panic(err)
	}

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

func (s *segment) Dispatch(dispatch *url.URL) (outargs, error) {
	dispatchPath := dispatch.Path
	dispatchQuery := dispatch.RawQuery

	ar := args2{}
	oa, err := s.dispatchPath(strings.Split(dispatchPath, "/"), ar, 0)
	if err != nil {
		return oa, err
	}

	ar, err = s.dispatchQuery(dispatchQuery, oa.amap, oa.ar, index(len(oa.ar.params)))
	if err != nil {
		return outargs{}, err
	}

	oa.ar = ar
	return oa, nil
}

func (s *segment) AddRoute(route *url.URL) error {

	routePath := route.Path
	if strings.HasPrefix(routePath, "/") {
		routePath = routePath[1:]
	}

	segs, err := parse(routePath, route.RawQuery)
	if err != nil {
		return err
	}

	err = s.addRoute(segs)
	if err != nil {
		return err
	}

	//fmt.Printf("____segs \n%v\n ___amap\n%v\n", segs, amap)
	//fmt.Printf("___m_ %v\n\n", amap)

	//fmt.Printf("___d_ %v\n\n", s)

	/*
		err = s.addSegments(segs, qParamNames)
		if err != nil {
			return err
		}
	*/
	return nil
}

/*
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
*/

/*
func (s *segment) Dispatch(u *url.URL) (args, error) {

	path := u.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return s.dispatch(path, u.RawQuery)
}

*/

/*
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
*/
