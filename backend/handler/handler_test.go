package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"github.com/Rajprakashkarimsetti/apica-project/models"
	"github.com/Rajprakashkarimsetti/apica-project/service"
)

func initializeTest(t *testing.T) (*service.MockLruCacher, handler) {
	ctrl := gomock.NewController(t)

	mockLruCacher := service.NewMockLruCacher(ctrl)
	mockHandler := New(mockLruCacher)

	return mockLruCacher, mockHandler
}

func Test_Get(t *testing.T) {
	mockLruCacher, mockHandler := initializeTest(t)

	testcases := []struct {
		desc   string
		input  string
		mock   *gomock.Call
		output string
		err    error
	}{
		{
			desc:   "success",
			input:  "key1",
			mock:   mockLruCacher.EXPECT().Get("key1").Return("value1"),
			output: "value1",
			err:    nil,
		},

		{
			desc:   "No key passed",
			input:  "",
			output: "",
			err: models.Error{
				StatusCode: http.StatusBadRequest,
				Reason:     "Parameter key is required",
				DateTime:   time.Now(),
			},
		},

		{
			desc:   "No data found for the passing key",
			input:  "key2",
			mock:   mockLruCacher.EXPECT().Get("key2").Return(""),
			output: "",
			err: models.Error{
				StatusCode: http.StatusNotFound,
				Reason:     "No key found",
				DateTime:   time.Now(),
			},
		},
	}

	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get/"+tc.input, nil)
			w := httptest.NewRecorder()

			vars := map[string]string{"key": tc.input}
			req = mux.SetURLVars(req, vars)

			mockHandler.Get(w, req)

			resp := w.Result()
			b, _ := io.ReadAll(resp.Body)

			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				var e models.Error
				_ = json.Unmarshal(b, &e)

				assert.Equalf(t, tc.err.Error(), e.Error(), "Test[%d] failed", i)
			} else {
				var res string

				json.Unmarshal(b, &res)

				assert.Equalf(t, tc.output, res, "Test[%d] failed", i)
			}
		})
	}
}

func Test_Set(t *testing.T) {
	mockLruCacher, mockHandler := initializeTest(t)

	testcases := []struct {
		desc   string
		input  []byte
		output string
		mock   *gomock.Call
		err    error
	}{
		{
			desc:   "success",
			input:  []byte(`{"key":"key1","value":"value1"}`),
			output: "Successfully inserted",
			mock:   mockLruCacher.EXPECT().Set("key1", "value1"),
			err:    nil,
		},

		{
			desc:   "Invalid value passed",
			input:  []byte(`{"key":"key1","value":123"}`),
			output: "",
			err: models.Error{
				StatusCode: http.StatusBadRequest,
				Reason:     "Invalid Body",
				DateTime:   time.Now(),
			},
		},

		{
			desc:   "key doesn't passed",
			input:  []byte(`{"value":"value1"}`),
			output: "",
			err: models.Error{
				StatusCode: http.StatusBadRequest,
				Reason:     "Key is required",
				DateTime:   time.Now(),
			},
		},

		{
			desc:   "key doesn't passed",
			input:  []byte(`{"key":"key1"}`),
			output: "",
			err: models.Error{
				StatusCode: http.StatusBadRequest,
				Reason:     "Value is required",
				DateTime:   time.Now(),
			},
		},
	}

	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/set", bytes.NewBuffer(tc.input))
			w := httptest.NewRecorder()

			mockHandler.Set(w, req)

			resp := w.Result()
			b, _ := io.ReadAll(resp.Body)

			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				var e models.Error
				_ = json.Unmarshal(b, &e)

				assert.Equalf(t, tc.err.Error(), e.Error(), "Test[%d] failed", i)
			} else {
				var res string

				json.Unmarshal(b, &res)

				assert.Equalf(t, tc.output, res, "Test[%d] failed", i)
			}
		})
	}
}
