package client

import (
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type out struct {
	Account int
	Address string
}

func TestAPI_Request_WhenValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		response string
		want     *out
	}{
		{
			desc:     "WhenStatusOK",
			response: `{"account": 1, "address": "test_address"}`,
			want: &out{
				Account: 1,
				Address: "test_address",
			},
		},
	}
	for _, ts := range tests {
		ts := ts
		t.Run(ts.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, ts.response)
				if err != nil {
					t.Fatalf("while write response got err: %v", err)
				}
			}))
			defer s.Close()
			client := s.Client()
			api := New(s.URL, "key", "secretKey", client, "userAgent")

			got := &out{}
			err := api.Request("GET", "/api/v1", nil, got)

			if diff := cmp.Diff(got, ts.want); diff != "" || err != nil {
				t.Errorf("%v.Request(params) got %v want %v", api, got, ts.want)
				t.Errorf("%v.Request(params) got err: %v want nil", api, err)
			}
		})
	}
}

func TestAPI_Request_WhenInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		response string
		status   int
		want     BinanceError
	}{
		{
			desc:     "WhenStatusForbidden",
			response: `{"code": 401, "msg": "forbidden"}`,
			status:   http.StatusForbidden,
			want: BinanceError{
				Code: 401,
				Msg:  "forbidden",
			},
		},
		{
			desc:     "WhenStatusOKWithInvalidJSON",
			response: `{"code": 200, "key":}`,
			status:   http.StatusOK,
			want: BinanceError{
				Msg: "Invalid JSON",
			},
		},
		{
			desc:     "WhenStatusNotFound",
			response: `{"code": 404, "msg": "url not found"}`,
			status:   http.StatusNotFound,
			want: BinanceError{
				Code: 404,
				Msg:  "url not found",
			},
		},
	}
	for _, ts := range tests {
		ts := ts
		t.Run(ts.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(ts.status)
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, ts.response)
				if err != nil {
					t.Fatalf("while write response got err: %v", err)
				}
			}))
			defer s.Close()
			client := s.Client()
			api := New(s.URL, "key", "secretKey", client, "userAgent")

			err := api.Request("GET", "/api/v1", nil, &out{})

			if err != ts.want {
				t.Errorf("%v.Request(params) got err: %v want %v", api, err, ts.want)
			}
		})
	}
}

func TestAPI_SignedRequest_WhenValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		response string
		want     *out
	}{
		{
			desc:     "WhenStatusOK",
			response: `{"account": 1, "address": "test_address"}`,
			want: &out{
				Account: 1,
				Address: "test_address",
			},
		},
	}
	for _, ts := range tests {
		ts := ts
		t.Run(ts.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, ts.response)
				if err != nil {
					t.Fatalf("while write response got err: %v", err)
				}
			}))
			defer s.Close()
			client := s.Client()
			api := New(s.URL, "key", "secretKey", client, "userAgent")

			got := &out{}
			err := api.SignedRequest("GET", "/api/v1", nil, got)

			if diff := cmp.Diff(got, ts.want); diff != "" || err != nil {
				t.Errorf("%v.Request(params) got %v want %v", api, got, ts.want)
				t.Errorf("%v.Request(params) got err: %v want nil", api, err)
			}
		})
	}
}

func TestAPI_SignedRequest_WhenInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		response string
		status   int
		want     BinanceError
	}{
		{
			desc:     "WhenStatusForbidden",
			response: `{"code": 401, "msg": "forbidden"}`,
			status:   http.StatusForbidden,
			want: BinanceError{
				Code: 401,
				Msg:  "forbidden",
			},
		},
		{
			desc:     "WhenStatusOKWithInvalidJSON",
			response: `{"code": 200, "key":}`,
			status:   http.StatusOK,
			want: BinanceError{
				Msg: "Invalid JSON",
			},
		},
		{
			desc:     "WhenStatusNotFound",
			response: `{"code": 404, "msg": "url not found"}`,
			status:   http.StatusNotFound,
			want: BinanceError{
				Code: 404,
				Msg:  "url not found",
			},
		},
	}
	for _, ts := range tests {
		ts := ts
		t.Run(ts.desc, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(ts.status)
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, ts.response)
				if err != nil {
					t.Fatalf("while write response got err: %v", err)
				}
			}))
			defer s.Close()
			client := s.Client()
			api := New(s.URL, "key", "secretKey", client, "userAgent")

			err := api.SignedRequest("GET", "/api/v1", nil, &out{})

			if err != ts.want {
				t.Errorf("%v.Request(params) got err: %v want %v", api, err, ts.want)
			}
		})
	}
}
