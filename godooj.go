package godooj

import "fmt"
import "strings"
import "bytes"
import "time"
import "net/http"
import "net/http/cookiejar"
import "encoding/json"

type Record map[string]interface{}

type CallArgs []interface{}

type CallKWArgs map[string]interface{}

type CallResult interface{}

type List []interface{}

type Dict map[string]interface{}

/*
 * Odoo Client struct
 */
type Client struct {
	baseURL    string
	httpClient *http.Client
	Auth       *AuthResponse
	Context    map[string]interface{}
}

func (client *Client) IsValid() bool {
	return client.httpClient != nil
}

func (client *Client) Call(model, method string, args CallArgs, kwargs map[string]interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/web/dataset/call_kw/%s/%s", client.baseURL, model, method)

	kwargs["context"] = client.Context

	requestPayload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "call",
		"params": map[string]interface{}{
			"model":  model,
			"method": method,
			"args":   args,
			"kwargs": kwargs,
		},
	}

	payloadBuf, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(payloadBuf))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	httpResponse, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	response := &RPCResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC Error: %s", response.Error.Data.Message)
	}

	if response.Result == nil {
		return nil, fmt.Errorf("RPC returned no result")
	}

	return response.Result, nil
}

func (client *Client) Read(model string, ids []int, fields []string) ([]interface{}, error) {
	args := []interface{}{ids}
	kwargs := map[string]interface{}{}
	if len(fields) > 0 {
		kwargs["fields"] = fields
	}
	respIface, err := client.Call(model, "read", args, kwargs)
	if err != nil {
		return nil, err
	}

	return respIface.([]interface{}), nil
}

func (client *Client) Search(model string, domain []interface{}) ([]int, error) {
	args := []interface{}{domain}
	kwargs := map[string]interface{}{}
	respIface, err := client.Call(model, "search", args, kwargs)
	if err != nil {
		return nil, err
	}
	recs := respIface.([]interface{})
	var ids []int
	for _, rec := range recs {
		idf := rec.(float64)
		ids = append(ids, int(idf))
	}
	return ids, nil
}

func (client *Client) SearchRead(model string, domain []interface{}, fields []string) ([]interface{}, error) {
	args := []interface{}{domain}
	kwargs := map[string]interface{}{}
	if len(fields) > 0 {
		kwargs["fields"] = fields
	}
	respIface, err := client.Call(model, "search_read", args, kwargs)
	if err != nil {
		return nil, err
	}

	return respIface.([]interface{}), nil
}

func (client *Client) Write(model string, ids []int, values map[string]interface{}) (bool, error) {
	args := []interface{}{ids, values}
	kwargs := map[string]interface{}{}
	respIface, err := client.Call(model, "write", args, kwargs)
	if err != nil {
		return false, err
	}

	return respIface.(bool), nil
}

func (client *Client) Create(model string, values map[string]interface{}) (int, error) {
	args := []interface{}{values}
	kwargs := map[string]interface{}{}
	respIface, err := client.Call(model, "create", args, kwargs)
	if err != nil {
		return -1, err
	}
	return int(respIface.(float64)), nil
}

func (client *Client) Delete(model string, ids []int) (bool, error) {
	args := []interface{}{ids}
	kwargs := map[string]interface{}{}
	respIface, err := client.Call(model, "unlink", args, kwargs)
	if err != nil {
		return false, err
	}
	return respIface.(bool), nil
}

func (client *Client) WithContext(ctx map[string]interface{}) *Client {
	newContext := map[string]interface{}{}
	for k, v := range client.Context {
		newContext[k] = v
	}
	for k, v := range ctx {
		newContext[k] = v
	}
	newClient := &Client{
		baseURL:    client.baseURL,
		Auth:       client.Auth,
		httpClient: client.httpClient,
		Context:    newContext,
	}
	return newClient
}

func Connect(baseURL, db, login, password string) (*Client, error) {
	url := baseURL + "/web/session/authenticate"

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}

	authParams := AuthParams{
		DB:       db,
		Login:    login,
		Password: password,
	}

	authRequest := AuthRequest{
		Id:      "1",
		JsonRPC: "2.0",
		Method:  "call",
		Params:  authParams,
	}

	payload := &strings.Builder{}
	err = json.NewEncoder(payload).Encode(authRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, strings.NewReader(payload.String()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := AuthResponse{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Error != nil {
		return nil, fmt.Errorf("Auth error: %s", data.Error.Data.Message)
	}

	if data.Result == nil {
		return nil, fmt.Errorf("Auth returned no result")
	}

	odooClient := &Client{
		baseURL:    baseURL,
		Auth:       &data,
		httpClient: client,
		Context:    make(map[string]interface{}),
	}

	return odooClient, nil
}

type OString string

func (s *OString) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case string:
		*s = OString(v)
	case bool:
		*s = ""
	}
	return nil
}

type OFloat float64

func (f *OFloat) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case float64:
		*f = OFloat(v)
	case bool:
		*f = 0.0
	}
	return nil
}

type OInt int

func (i *OInt) UnmarshalJSON(b []byte) error {
	var ie interface{}
	if err := json.Unmarshal(b, &ie); err != nil {
		return err
	}
	switch v := ie.(type) {
	case int:
		*i = OInt(v)
	case bool:
		*i = 0
	}
	return nil
}

type OMany2one struct {
	ID   int
	Name string
}

func (m *OMany2one) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch i.(type) {
	case []interface{}:
		slice := i.([]interface{})
		*m = OMany2one{
			ID:   slice[0].(int),
			Name: slice[1].(string),
		}
	case bool:
		*m = OMany2one{}
	}
	return nil
}

type ODateTime time.Time

func (t *ODateTime) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch v := i.(type) {
	case string:
		ptime, err := time.Parse(time.RFC3339, v)
		if err == nil {
			*t = ODateTime(ptime)
		} else {
			*t = ODateTime(time.Time{})
		}
	case bool:
		*t = ODateTime(time.Time{})
	}
	return nil
}

type AuthParams struct {
	DB       string `json:"db"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Id      string     `json:"id"`
	JsonRPC string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  AuthParams `json:"params"`
}

type AuthErrorData struct {
	Name          OString       `json:"name,omitempty"`
	Debug         OString       `json:"debug,omitempty"`
	Message       OString       `json:"message,omitempty"`
	Arguments     []interface{} `json:"arguments,omitempty"`
	ExceptionType OString       `json:"exception_type,omitempty"`
}

type AuthError struct {
	Code    int            `json:"code,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    *AuthErrorData `json:"data,omitempty"`
}

type AuthUserContext struct {
	Language OString `json:"lang,omitempty"`
	TimeZone OString `json:"tz,omitempty"`
	UID      OInt    `json:"uid,omitempty"`
}

type AuthResult struct {
	UID                OInt             `json:"uid"`
	IsSystem           bool             `json:"is_system"`
	IsAdmin            bool             `json:"is_admin"`
	UserContext        *AuthUserContext `json:"user_context"`
	DB                 OString          `json:"db"`
	ServerVersion      OString          `json:"server_version"`
	ServerVersionInfo  []interface{}    `json:"server_version_info"`
	Name               OString          `json:"name"`
	UserName           OString          `json:"username"`
	PartnerDisplayName OString          `json:"partner_display_name"`
	CompanyID          OInt             `json:"company_id"`
	PartnerID          OInt             `json:"partner_id"`
	WebBaseURL         OString          `json:"web.base.url"`
}

type AuthResponse struct {
	JsonRPC OString     `json:"jsonrpc"`
	Id      OString     `json:"id"`
	Result  *AuthResult `json:"result"`
	Error   *AuthError  `json:"error"`
}

type RPCErrorData struct {
	Name          OString       `json:"name"`
	Debug         OString       `json:"debug"`
	Message       OString       `json:"message"`
	Arguments     []interface{} `json:"arguments"`
	ExceptionType OString       `json:"exception_type"`
}

type RPCError struct {
	Code    OInt         `json:"code"`
	Message OString      `json:"message"`
	Data    RPCErrorData `json:"data"`
}

type RPCResponse struct {
	JsonRPC OString     `json:"jsonrpc"`
	Id      OString     `json:"id"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
}
