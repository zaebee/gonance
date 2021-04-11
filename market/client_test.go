package market

import (
	"github.com/google/go-cmp/cmp"
	"gonance/client"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarket_Prices_WhenValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		response string
		want     Prices
	}{
		{
			desc:     "WhenEmptyResponse",
			response: `[]`,
			want:     Prices{},
		},
		{
			desc:     "WhenStatusOK",
			response: `[{"symbol": "BTCETH", "price": "100"}, {"symbol": "BTCADA", "price": "200"}]`,
			want: Prices{
				"BTCETH": {
					Symbol: "BTCETH",
					Price:  "100",
				},
				"BTCADA": {
					Symbol: "BTCADA",
					Price:  "200",
				},
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
			api := client.New(s.URL, "key", "secretKey", s.Client(), "userAgent")
			market := Client{API: *api}
			got, err := market.Prices()

			if diff := cmp.Diff(got, ts.want); diff != "" || err != nil {
				t.Errorf("%v.Prices() got %v want %v", market, got, ts.want)
				t.Errorf("%v.Prices() got err: %v want nil", market, err)
			}
		})
	}
}
