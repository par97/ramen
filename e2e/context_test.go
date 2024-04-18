package e2e_test

import (
	"fmt"
	"strings"
)

type TestContext struct {
	w Workload
	d Deployer
}

var testContextList = make(map[string]TestContext)

// Based on name passed, Init the deployer and Workload and stash in a map[string]TestContext
func addTestContext(name string, w Workload, d Deployer) {
	testContextList[name] = TestContext{w, d}
}

// TODO: Search name in map for a TestContext to return, if not found go backward
// - i.e drop the last /<name> suffix form name and search till a match  is found or all suffixes are exhausted
// - e.g If name passed is "TestSuites/Exhaustive/DaemonSet/Subscription/Undeploy"
//   - Search for above name first (it will not be found as we create context at a point where we have a d+w)
//   - Search for "TestSuites/Exhaustive/DaemonSet/Subscription" (should be found)
//   - ...and so on till we find a match or exhaust suffixes
//
// I do not think need search more than twice. only last suffix need be removed.
// No recursive unless there is a valid use case,
func getTestContext(name string) (TestContext, error) {
	testCtx, ok := testContextList[name]
	if !ok {
		i := strings.LastIndex(name, "/")
		if i < 1 {
			return TestContext{}, fmt.Errorf("not a valid name in getTestContext: %v", name)
		}
		testCtx, ok = testContextList[name[0:i]]
		if !ok {
			return TestContext{}, fmt.Errorf("can not find testContext with name: %v", name)
		}
	}
	return testCtx, nil
}
