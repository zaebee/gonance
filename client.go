//Package client contains methods to make request to Binance API server.
package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

//API is a Binance API client.
type API struct {
	URL        string
	Key        string
	SecretKey  string
	HTTPClient *http.Client
	UserAgent  string
}

// BinanceError handles api errors from binance.com.
type BinanceError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error returns error message from Binance API.
func (e BinanceError) Error() string {
	return e.Msg
}

//New initializes API with given URL, api key and secret key. it also provides a way to overwrite *http.Client
func New(url, key, secretKey string, httpClient *http.Client, userAgent string) *API {
	return &API{
		URL:        url,
		Key:        key,
		SecretKey:  secretKey,
		HTTPClient: httpClient,
		UserAgent:  userAgent,
	}
}

//Making a public request to Binance API server.
func (a *API) Request(method, endpoint string, params interface{}, out interface{}) error {
	u, err := url.ParseRequestURI(a.URL)
	if err != nil {
		return err
	}
	u.Path = u.Path + endpoint

	if method == "GET" {
		//parse params to query string
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		q := u.Query()
		for k, v := range m {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		u.RawQuery = q.Encode()
	}
	log.Printf("%v %v", method, u.String())
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-MBX-APIKEY", a.Key)
	req.Header.Add("UserAgent", a.UserAgent)
	res, err := a.HTTPClient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		e := BinanceError{}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return BinanceError{Msg: err.Error()}
		}
		return e
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return BinanceError{Msg: "Invalid JSON"}
	}

	return nil
}

//Making a signed request to Binance API server.
func (a *API) SignedRequest(method, endpoint string, params interface{}, out interface{}) error {
	u, _ := url.ParseRequestURI(a.URL)
	u.Path = u.Path + endpoint

	//parse params to query string
	b, _ := json.Marshal(params)
	m := map[string]interface{}{}
	json.Unmarshal(b, &m)

	q := u.Query()
	for k, v := range m {
		q.Set(k, fmt.Sprintf("%v", v))
	}

	//timestamp is mandatory in signed request
	q.Add("timestamp", fmt.Sprintf("%v", time.Now().Unix()*1000))

	mac := hmac.New(sha256.New, []byte(a.SecretKey))
	mac.Write([]byte(q.Encode()))
	expectedMAC := mac.Sum(nil)
	signed := hex.EncodeToString(expectedMAC)
	//signature needs to be at the last param
	u.RawQuery = q.Encode() + "&signature=" + signed

	log.Printf("%v %v", method, u.String())

	req, _ := http.NewRequest(method, u.String(), nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-MBX-APIKEY", a.Key)
	req.Header.Add("UserAgent", a.UserAgent)
	res, err := a.HTTPClient.Do(req)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		type binanceError struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		e := binanceError{}
		err = json.NewDecoder(res.Body).Decode(&e)
		return errors.New(e.Msg)
	}
	defer res.Body.Close()
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(&out)
	}
	return err
}

type StreamHandler func(data []byte)

func (a *API) Stream(endpoint string, handler StreamHandler) {
	u := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s", endpoint)
	websocketClient, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}

	go func() {
		defer websocketClient.Close()
		for {
			_, m, err := websocketClient.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			go handler(m)
		}
	}()

}
