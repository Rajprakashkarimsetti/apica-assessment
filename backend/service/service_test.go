package service

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"

	"github.com/Rajprakashkarimsetti/apica-project/models"
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
		input *models.CacheData
		mock  *gomock.Call
	}{
		{
			desc: "success",
			input: &models.CacheData{Key: "key1",
				Value:      "value1",
				Expiration: 5,
			},
			mock: mockLruCacherStr.EXPECT().Set(&models.CacheData{Key: "key1",
				Value:      "value1",
				Expiration: 5,
			}),
		},
	}

	for _, tc := range testcases {
		mockLruCacherSvc.Set(tc.input)
	}
}
