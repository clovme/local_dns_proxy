package core

import "fmt"

type routesPathMap map[string]RoutesInfo
type routesNameMap map[string]RoutesInfo

type routesMap struct {
	nameRoute routesNameMap
	pathRoute routesPathMap
}

func initRoutesMap() routesMap {
	routesName := make(routesNameMap)
	routesPath := make(routesPathMap)
	for _, route := range routesInfo {
		if _, ok := routesName[route.Name]; !ok {
			routesName[route.Name] = route
			routesPath[route.Path] = route
		} else {
			panic(fmt.Sprintf("路由名称重复: %s\n   %+v\n   %+v\n", route.Name, route, routesName[route.Name]))
		}
	}
	return routesMap{
		nameRoute: routesName,
		pathRoute: routesPath,
	}
}

func (r routesMap) Name(name string) RoutesInfo {
	return r.nameRoute[name]
}

func (r routesMap) Path(urlPath string) RoutesInfo {
	return r.pathRoute[urlPath]
}
