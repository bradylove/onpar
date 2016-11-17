//go:generate hel

package expect_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/apoydence/onpar/expect"
)

func TestToPassesActualToMatcher(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	close(mockMatcher.MatchOutput.Err)

	Expect(mockT, 101).To(mockMatcher)

	select {
	case actual := <-mockMatcher.MatchInput.Actual:
		if !reflect.DeepEqual(actual, 101) {
			t.Errorf("Expected %v to equal %v", actual, 101)
		}
	default:
		t.Errorf("Expected Match() to be invoked")
	}
}

func TestToErrorsIfMatcherFails(t *testing.T) {
	mockT := newMockT()
	mockMatcher := newMockMatcher()
	close(mockMatcher.MatchOutput.ResultValue)
	mockMatcher.MatchOutput.Err <- fmt.Errorf("some-error")

	Expect(mockT, 101).To(mockMatcher)

	select {
	case msg := <-mockT.FatalInput.Arg0:
		if !reflect.DeepEqual(msg, []interface{}{"some-error"}) {
			t.Errorf("Expected %v to equal %v", msg, "some-error")
		}
	default:
		t.Errorf("Expected Fatal() to be invoked")
	}
}
