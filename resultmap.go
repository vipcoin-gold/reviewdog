package reviewdog

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vipcoin-gold/reviewdog/filter"
	"github.com/vipcoin-gold/reviewdog/proto/rdf"
)

// ResultMap represents a concurrent-safe map to store Diagnostics generated by concurrent jobs.
type ResultMap struct {
	sm sync.Map
}

type Result struct {
	Name        string
	Level       string
	Diagnostics []*rdf.Diagnostic

	// Optional. Report an error of the command execution.
	// Non-nil CmdErr doesn't mean failure and Diagnostics still may have
	// results.
	// It is common that a linter fails with non-zero exit code when it finds
	// lint errors.
	CmdErr error
}

// CheckUnexpectedFailure returns error on unexpected failure, if any.
func (r *Result) CheckUnexpectedFailure() error {
	if r.CmdErr != nil && len(r.Diagnostics) == 0 {
		return fmt.Errorf("%s failed with zero findings: The command itself "+
			"failed (%v) or reviewdog cannot parse the results", r.Name, r.CmdErr)
	}
	return nil
}

// Store saves a new *Result into ResultMap.
func (rm *ResultMap) Store(key string, r *Result) {
	rm.sm.Store(key, r)
}

// Load fetches *Result from ResultMap
func (rm *ResultMap) Load(key string) (*Result, error) {
	v, ok := rm.sm.Load(key)
	if !ok {
		return nil, fmt.Errorf("fail to get the value of key %q from results", key)
	}

	t, ok := v.(*Result)
	if !ok {
		return nil, errors.New("stored type in ResultMap is invalid")
	}

	return t, nil
}

// Range retrieves `key` and `values` from ResultMap iteratively.
func (rm *ResultMap) Range(f func(key string, val *Result)) {
	rm.sm.Range(func(k, v interface{}) bool {
		f(k.(string), v.(*Result))
		return true
	})
}

// Len returns the length of ResultMap count. Len() is not yet officially not supported by Go. (ref: https://github.com/golang/go/issues/20680)
func (rm *ResultMap) Len() int {
	l := 0
	rm.sm.Range(func(_, _ interface{}) bool {
		l++
		return true
	})
	return l
}

type FilteredResult struct {
	Level              string
	FilteredDiagnostic []*filter.FilteredDiagnostic
}

// FilteredResultMap represents a concurrent-safe map to store Diagnostics generated by concurrent jobs.
type FilteredResultMap struct {
	sm sync.Map
}

// Store saves a new []*FilteredCheckFilteredResult into FilteredResultMap.
func (rm *FilteredResultMap) Store(key string, r *FilteredResult) {
	rm.sm.Store(key, r)
}

// Load fetches FilteredResult from FilteredResultMap
func (rm *FilteredResultMap) Load(key string) (*FilteredResult, error) {
	v, ok := rm.sm.Load(key)
	if !ok {
		return nil, fmt.Errorf("fail to get the value of key %q from results", key)
	}

	t, ok := v.(*FilteredResult)
	if !ok {
		return nil, errors.New("stored type in FilteredResultMap is invalid")
	}

	return t, nil
}

// Range retrieves `key` and `values` from FilteredResultMap iteratively.
func (rm *FilteredResultMap) Range(f func(key string, val *FilteredResult)) {
	rm.sm.Range(func(k, v interface{}) bool {
		f(k.(string), v.(*FilteredResult))
		return true
	})
}

// Len returns the length of FilteredResultMap count. Len() is not yet officially not supported by Go. (ref: https://github.com/golang/go/issues/20680)
func (rm *FilteredResultMap) Len() int {
	l := 0
	rm.sm.Range(func(_, _ interface{}) bool {
		l++
		return true
	})
	return l
}
