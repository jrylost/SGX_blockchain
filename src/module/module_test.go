package module

import (
	"testing"
)

func TestOutput(t *testing.T) {
	if ans := Output(); ans != 5 {
		t.Errorf("Test do not pass %d", ans)
	}

}
