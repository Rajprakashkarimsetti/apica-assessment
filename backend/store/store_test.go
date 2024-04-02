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
						Key:        "key1",
						Value:      "value1",
						Expiration: 5,
						TimeStamp:  time.Now(),
					},
				},
				Head: &models.CacheData{
					Key:        "key2",
					Value:      "value2",
					Expiration: 5,
					TimeStamp:  time.Now(),
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
						Key:        "key1",
						Value:      "value1",
						Expiration: 5,
						TimeStamp:  time.Now()},
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
		input *models.CacheData
		cache *cacher.Cache
	}{
		{
			desc: "successfully inserted into cache",
			input: &models.CacheData{Key: "key1",
				Value:      "value1",
				Expiration: 5,
			},
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache:    map[string]*models.CacheData{},
				Mutex:    sync.Mutex{},
			},
		},

		{
			desc: "success, key already exists, updating the value",
			input: &models.CacheData{Key: "key1",
				Value:      "value2",
				Expiration: 5,
			},
			cache: &cacher.Cache{
				Capacity: 1024,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:        "key1",
						Value:      "value1",
						Expiration: 5,
						TimeStamp:  time.Now()},
				},
				Head: &models.CacheData{
					Key:        "key1",
					Value:      "value2",
					Expiration: 5,
					TimeStamp:  time.Now()},
				Mutex: sync.Mutex{},
			},
		},

		{
			desc: "cache capacity exceeded, removes the last element and inserts new to front",
			input: &models.CacheData{Key: "key3",
				Value:      "value3",
				Expiration: 5,
			},
			cache: &cacher.Cache{
				Capacity: 1,
				Cache: map[string]*models.CacheData{
					"key1": {
						Key:        "key1",
						Value:      "value1",
						Expiration: 5,
						TimeStamp:  time.Now()},
				},
				Tail: &models.CacheData{
					Key:        "key2",
					Value:      "value2",
					Expiration: 5,
					TimeStamp:  time.Now()},
				Mutex: sync.Mutex{},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			mockLruCacherStr := New(tc.cache)

			mockLruCacherStr.Set(tc.input)

			if _, ok := tc.cache.Cache[tc.input.Key]; !ok {
				assert.Equal(t, tc.input.Key, "", "Test failed")
			}
		})
	}
}
