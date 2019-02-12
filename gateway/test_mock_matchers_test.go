package gateway

// mockFuncMatcher satisfies the github.com/golang/mock/gomock.Matcher interface
// for allowing to customize a mock param checker through a function
type mockFuncMatcher struct {
	// Func must return true if the argument satisfies its checks, otherwise false.
	Func func(interface{}) bool
}

// Matches returns the boolean returned by Func.
func (f mockFuncMatcher) Matches(x interface{}) bool {
	return f.Func(x)
}

// String return a statick message for this custom matcher.
func (f mockFuncMatcher) String() string {
	return "is satisfied by Func"
}
