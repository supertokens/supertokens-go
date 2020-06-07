package core

import (
	"testing"
)

func TestVersionCompare(t *testing.T) {
	if getLargestVersionFromIntersection([]string{"1.0"}, []string{"2.0"}) != nil {
		t.Error("should have failed")
	}

	if *getLargestVersionFromIntersection([]string{"1.0", "2.0"}, []string{"2.0"}) != "2.0" {
		t.Error("should have failed")
	}

	if *getLargestVersionFromIntersection([]string{"1.0.1", "1.1.19"}, []string{"2.0", "1.1.19", "2.1"}) != "1.1.19" {
		t.Error("should have failed")
	}

	if *getLargestVersionFromIntersection([]string{"1.0", "1.1.19", "1.1.2"}, []string{"1.0", "1.1.19", "1.1.2", "1.1.3"}) != "1.1.19" {
		t.Error("should have failed")
	}

	if *getLargestVersionFromIntersection([]string{"2.0", "1.1.19", "1.1.2"}, []string{"2.0", "1.1.19", "1.1.2", "1.1.3"}) != "2.0" {
		t.Error("should have failed")
	}
}
