package service

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"

	"github.com/Rajprakashkarimsetti/apica-project/store"
)

func initializeTest(t *testing.T) (*store.MockLruCacher, Service) {
	ctrl := gomock.NewController(t)
	mockLruCacherStr := store.NewMockLruCacher(ctrl)

	mockLruCacherSvc := New(mockLruCacherStr)

	return mockLruCacherStr, mockLruCacherSvc
}

func Test_Get(t *testing.T) {
	mockLruCacherStr, mockLruCacherSvc := initializeTest(t)

	testcases := []struct {
		desc   string
		input  string
		output string
		mock   *gomock.Call
	}{
		{
			desc:   "success",
			input:  "key1",
			output: "value1",
			mock:   mockLruCacherStr.EXPECT().Get("key1").Return("value1"),
		},
	}

	for i, tc := range testcases {
		res := mockLruCacherSvc.Get(tc.input)

		assert.Equalf(t, tc.output, res, "Test[%d] failed", i)
	}

}

func Test_Set(t *testing.T) {
	mockLruCacherStr, mockLruCacherSvc := initializeTest(t)

	testcases := []struct {
		desc  string
		key   string
		value string
		mock  *gomock.Call
	}{
		{
			desc:  "success",
			key:   "key1",
			value: "value1",
			mock:  mockLruCacherStr.EXPECT().Set("key1", "value1"),
		},
	}

	for _, tc := range testcases {
		mockLruCacherSvc.Set(tc.key, tc.value)
	}
}
