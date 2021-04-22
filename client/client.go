package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nerdynick/ccloud-go-sdk/logging"
	"go.uber.org/zap"
)

const (
	//DefaultRequestTimeout is the default number of seconds to wait before considering a Metrics API query/request as timedout
	DefaultRequestTimeout time.Duration = time.Second * 60
	//DefaultMaxIdleConns is the default nmax umber of Idle HTTP connections to leave open
	DefaultMaxIdleConns int = 5
	//DefaultMaxIdleConnsPerHost is the default max number of HTTP connections to leave open for a singlular HTTP Host
	DefaultMaxIdleConnsPerHost int = 5
)

//Client HTTP Client Struct
type Client struct {
	*logging.Loggable
	Context          Context
	httpClient       http.Client
	HTTPErrorHandler func(int, []byte) error
}

//Request Sends a Request synchronously
func (client *Client) Request(request *http.Request) ([]byte, error) {
	res, err := client.httpClient.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		client.Log.Error("Error returned from HTTP Request",
			zap.String("url", request.RequestURI),
			zap.Error(err),
		)
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error := NewError(res.StatusCode, request.RequestURI, err)
		client.Log.Error("Request - Failed to Read Return",
			zap.String("url", request.RequestURI),
			zap.Error(error),
			zap.Int("statusCode", res.StatusCode),
			zap.String("statusMessage", res.Status),
		)
		return nil, err
	}

	if res.StatusCode != 200 {
		err := client.HTTPErrorHandler(res.StatusCode, resBody)
		error := NewError(res.StatusCode, request.RequestURI, err)

		// err := responses.ErrorResponse{}
		// json.Unmarshal(resBody, &err)

		client.Log.Error("Request - Invalid response code",
			zap.String("url", request.RequestURI),
			zap.Int("statusCode", res.StatusCode),
			zap.String("statusMessage", res.Status),
			zap.Error(error),
		)
		return nil, error
	}

	if client.Log.Core().Enabled(logging.DebugLevel) {
		client.Log.Debug("Request - Body",
			zap.String("url", request.RequestURI),
			zap.String("results", string(resBody)),
		)
	}

	return resBody, nil
}

//RequestAsync Sends a Request asynchronously
func (client *Client) RequestAsync(request *http.Request, responseChan chan<- []byte, errorChan chan<- error) {
	res, err := client.Request(request)
	if err != nil {
		errorChan <- err
	} else {
		responseChan <- res
	}
}

//NewRequest Builds a new http Request
func (client *Client) NewRequest(method string, url string, body []byte) (*http.Request, error) {
	client.Log.Debug("Creating Request",
		zap.String("method", method),
		zap.String("url", url),
		zap.Binary("Body", body),
	)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.Context.APIKey, client.Context.APISecret)

	req.Header.Add("Content-Type", "application/json")
	if client.Context.UserAgent != "" {
		req.Header.Add("User-Agent", client.Context.UserAgent)
	}

	for header, value := range client.Context.HTTPHeaders {
		req.Header.Add(header, value)
	}

	return req, nil
}

//SendRequest Sends a Request to the given url synchronously
func (client *Client) SendRequest(method string, url string, body []byte) ([]byte, error) {
	req, err := client.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Request(req)
}

//SendRequestAsync Sends a Request to the given url asynchronously
func (client *Client) SendRequestAsync(method string, url string, body []byte, responseChan chan<- []byte, errorChan chan<- error) {
	req, err := client.NewRequest(method, url, body)
	if err != nil {
		errorChan <- err
		return
	}
	client.RequestAsync(req, responseChan, errorChan)
}

//ResponseSupplier supplier of new Struct instances to Unmarshal the JSON response into
type ResponseSupplier func() *interface{}

//SendGet send a GET request to the API
func (client *Client) Get(response interface{}, url string) error {
	res, err := client.SendRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, &response)
}

//SendGetAsync send a GET request to the API async
func (client *Client) GetAsync(responseSupplier ResponseSupplier, url string, responseChan chan<- interface{}, errorChan chan<- error) {
	go func() {
		respChan := make(chan []byte, 1)
		client.SendRequestAsync("GET", url, nil, respChan, errorChan)
		for r := range respChan {
			res := responseSupplier()
			err := json.Unmarshal(r, &res)
			if err != nil {
				errorChan <- err
			} else {
				responseChan <- res
			}
		}
	}()
}

//SendPost send a POST request with a given JSON Body
func (client Client) Post(response interface{}, url string, jsonBody interface{}) error {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	res, err := client.SendRequest("POST", url, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, &response)
}

//SendPostAsync send a POST request with a given JSON Body async
func (client *Client) PostAsync(responseSupplier ResponseSupplier, url string, jsonBody interface{}, responseChan chan<- interface{}, errorChan chan<- error) {
	go func() {
		respChan := make(chan []byte, 1)

		body, err := json.Marshal(jsonBody)
		if err != nil {
			errorChan <- err
			return
		}

		client.SendRequestAsync("POST", url, body, respChan, errorChan)
		for r := range respChan {
			res := responseSupplier()
			err := json.Unmarshal(r, &res)
			if err != nil {
				errorChan <- err
			} else {
				responseChan <- res
			}
		}
	}()
}

//New Creates a new CCloud Metrics HTTP Client
func New(apiKey string, apiSecret string, httpErrorHandler func(int, []byte) error) Client {
	log := logging.New("CCloudAPIClient")

	return Client{
		Loggable:         log,
		Context:          NewContext(apiKey, apiSecret),
		HTTPErrorHandler: httpErrorHandler,
		httpClient: http.Client{
			Timeout: DefaultRequestTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        DefaultMaxIdleConns,
				MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
			},
		},
	}
}
