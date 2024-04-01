package store

import (
	"sync"
	"testing"
	"time"

	"github.com/bmizerany/assert"

	"github.com/Rajprakashkarimsetti/apica-project/cacher"
	"github.com/Rajprakashkarimsetti/apica-project/models"
)

func Test_Get(t *testing.T) {
	testcases := []struct {
		desc   string
		input  string
		output string
		cache  *cacher.Cache
	}{
		{
			desc:   "success",
			input:  "key1",
			output: "value1",
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:       "key1",
						Value:     "value1",
						Timestamp: time.Now(),
					},
				},
				Head: &models.CacheData{
					Key:       "key2",
					Value:     "value2",
					Timestamp: time.Now(),
				},
				Mutex: sync.Mutex{},
			},
		},

		{
			desc:   "key not found",
			input:  "key2",
			output: "",
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:       "key1",
						Value:     "value1",
						Timestamp: time.Now(),
					},
				},
				Mutex: sync.Mutex{},
			},
		},
	}

	for i, tc := range testcases {
		mockLruCacheStr := New(tc.cache)

		res := mockLruCacheStr.Get(tc.input)

		assert.Equalf(t, tc.output, res, "Test[%d] failed", i)
	}
}

func Test_Set(t *testing.T) {
	testcases := []struct {
		desc  string
		key   string
		value string
		cache *cacher.Cache
	}{
		{
			desc:  "successfully inserted into cache",
			key:   "key1",
			value: "value1",
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache:    map[string]*models.CacheData{},
				Mutex:    sync.Mutex{},
			},
		},

		{
			desc:  "success, key already exists, updating the value",
			key:   "key1",
			value: "value2",
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:       "key1",
						Value:     "value1",
						Timestamp: time.Now(),
					},
				},
				Head: &models.CacheData{
					Key:       "key1",
					Value:     "value2",
					Timestamp: time.Now(),
				},
				Mutex: sync.Mutex{},
			},
		},

		{
			desc:  "cache capacity exceeded, removes the last element and inserts new to front",
			key:   "key3",
			value: "value3",
			cache: &cacher.Cache{
				Capacity: 1,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:       "key1",
						Value:     "value1",
						Timestamp: time.Now(),
					},
				},
				Tail: &models.CacheData{
					Key:       "key2",
					Value:     "value2",
					Timestamp: time.Now(),
				},
				Mutex: sync.Mutex{},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			mockLruCacherStr := New(tc.cache)

			mockLruCacherStr.Set(tc.key, tc.value)

			if _, ok := tc.cache.Cache[tc.key]; !ok {
				assert.Equal(t, tc.key, "", "Test failed")
			}
		})
	}
}
