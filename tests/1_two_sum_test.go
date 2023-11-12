package main

import (
	"leetcode/problems"
	"reflect"
	"testing"
)

func TestTwo_sum(t *testing.T) {
	// Define test cases here
	type Input struct {
		nums   []int
		target int
	}
	type TestCase struct {
		inputs []Input
		want   []int
	}

	var testCases []TestCase

	for _, tc := range testCases {
		got := Two_sum(tc.inputs...)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("For inputs %v, got %v, want %v", tc.inputs, got, tc.want)
		}
	}
}

