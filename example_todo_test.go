package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"strings"

	"github.com/crhntr/jsonapi"
)

func ExampleServeMux() {
	var mux jsonapi.ServeMux

	const issuesEndpoint = "issues"

	type (
		Issue struct {
			Type        string `json:"-"`
			Description string `json:"desc"`
			Done        bool   `json:"done"`
		}
	)

	var issues []Issue

	issueTypes := map[string]struct{}{"bug": struct{}{}, "feature": struct{}{}, "chore": struct{}{}}

	mux.HandleCreate(issuesEndpoint, jsonapi.CreateFunc(func(res jsonapi.CreateResponder, req *http.Request, _ string) {
		var body jsonapi.CreateRequestData

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			res.AppendError(jsonapi.Error{Detail: err.Error(), Status: http.StatusBadRequest})
			return
		}

		if _, ok := issueTypes[body.Data.Type]; !ok {
			detail := fmt.Sprintf("%q is not among the type(s) that constitute the collection represented by the endpoint", body.Data.Type)
			res.AppendError(jsonapi.Error{Detail: detail, Status: http.StatusConflict})
			return
		}

		var issue Issue
		if err := json.Unmarshal(body.Data.Attributes, &issue); err != nil {
			res.AppendError(jsonapi.Error{Detail: err.Error(), Status: http.StatusBadRequest})
			return
		}
		issue.Type = body.Data.Type

		issues = append(issues, issue)

		res.SetData(body.Data.Type, strconv.Itoa(len(issues)), issue, nil, nil, nil)
	}))

	mux.HandleFetchOne(issuesEndpoint, jsonapi.FetchOneFunc(func(res jsonapi.FetchOneResonder, req *http.Request, idStr string) {
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 || id >= len(issues) {
			res.AppendError(jsonapi.Error{Detail: fmt.Sprintf("%q not found", idStr), Status: http.StatusNotFound})
			return
		}
		issue := issues[id]
		res.SetData(issue.Type, strconv.Itoa(id), issue, nil, nil, nil)
	}))

	mux.HandleFetchCollection(issuesEndpoint, jsonapi.FetchCollectionFunc(func(res jsonapi.FetchCollectionResponder, req *http.Request) {
		for id, issue := range issues {
			res.AppendData(issue.Type, strconv.Itoa(id), issue, nil, nil, nil)
		}
	}))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// Do Client Stuff

	{
		fmt.Println("# Request all issues")

		req := requestJSONAPI(http.MethodGet, testServer.URL+"/"+issuesEndpoint, nil)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(printReqRes(res, nil))
	}

	{
		fmt.Println("# Create a feature")

		reqBody := `{"data": {"type": "feature","attributes": {"desc": "As a teapot, I should pour tea"}}}`

		req := requestJSONAPI(http.MethodPost, testServer.URL+"/"+issuesEndpoint, strings.NewReader(reqBody))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(printReqRes(res, nil))
	}

	{
		fmt.Println("# Create a bug")

		reqBody := []byte(`{"data": {"type": "bug","attributes": {"desc": "When tea from teapot is poured out, it is not warm enough."}}}`)

		req := requestJSONAPI(http.MethodPost, testServer.URL+"/"+issuesEndpoint, bytes.NewReader(reqBody))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(printReqRes(res, reqBody))
	}

	{
		fmt.Println("# Fetch the bug")

		reqBody := []byte(`{"data": {"type": "bug","attributes": {"desc": "When tea from teapot is poured out, it is not warm enough."}}}`)

		req := requestJSONAPI(http.MethodGet, testServer.URL+"/"+path.Join(issuesEndpoint, "0"), nil)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(printReqRes(res, reqBody))
	}

	{
		fmt.Println("# Request Issues")

		req := requestJSONAPI(http.MethodGet, testServer.URL+"/"+issuesEndpoint, nil)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(printReqRes(res, nil))
	}

	// Output:
	// # Request all issues
	// ## REQUEST
	//
	// GET /issues HTTP/1.1
	// Accept: application/vnd.api+json
	// Content-Type: application/vnd.api+json
	// ## RESONSE
	//
	// HTTP/1.1	OK	200
	// Content-Type: application/vnd.api+json
	// {
	//   "data": []
	// }
	//
	// # Create a feature
	// ## REQUEST
	//
	// POST /issues HTTP/1.1
	// Accept: application/vnd.api+json
	// Content-Type: application/vnd.api+json
	//
	// ## RESONSE
	//
	// HTTP/1.1	Created	201
	// Content-Type: application/vnd.api+json
	// {
	//   "data": {
	//     "id": "1",
	//     "type": "feature",
	//     "attributes": {
	//       "desc": "As a teapot, I should pour tea",
	//       "done": false
	//     }
	//   }
	// }
	//
	// # Create a bug
	// ## REQUEST
	//
	// POST /issues HTTP/1.1
	// Accept: application/vnd.api+json
	// Content-Type: application/vnd.api+json
	// {
	//   "data": {
	//     "type": "bug",
	//     "attributes": {
	//       "desc": "When tea from teapot is poured out, it is not warm enough."
	//     }
	//   }
	// }
	// ## RESONSE
	//
	// HTTP/1.1	Created	201
	// Content-Type: application/vnd.api+json
	// {
	//   "data": {
	//     "id": "2",
	//     "type": "bug",
	//     "attributes": {
	//       "desc": "When tea from teapot is poured out, it is not warm enough.",
	//       "done": false
	//     }
	//   }
	// }
	//
	// # Fetch the bug
	// ## REQUEST
	//
	// GET /issues/0 HTTP/1.1
	// Accept: application/vnd.api+json
	// Content-Type: application/vnd.api+json
	// ## RESONSE
	//
	// HTTP/1.1	OK	200
	// Content-Type: application/vnd.api+json
	// {
	//   "data": {
	//     "id": "0",
	//     "type": "feature",
	//     "attributes": {
	//       "desc": "As a teapot, I should pour tea",
	//       "done": false
	//     }
	//   }
	// }
	//
	// # Request Issues
	// ## REQUEST
	//
	// GET /issues HTTP/1.1
	// Accept: application/vnd.api+json
	// Content-Type: application/vnd.api+json
	// ## RESONSE
	//
	// HTTP/1.1	OK	200
	// Content-Type: application/vnd.api+json
	// {
	//   "data": [
	//     {
	//       "id": "0",
	//       "type": "feature",
	//       "attributes": {
	//         "desc": "As a teapot, I should pour tea",
	//         "done": false
	//       }
	//     },
	//     {
	//       "id": "1",
	//       "type": "bug",
	//       "attributes": {
	//         "desc": "When tea from teapot is poured out, it is not warm enough.",
	//         "done": false
	//       }
	//     }
	//   ]
	// }
}

func requestJSONAPI(method string, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", jsonapi.ContentType)
	req.Header.Set("Content-Type", jsonapi.ContentType)
	return req
}

func printReqRes(res *http.Response, requestBody []byte) string {
	req := res.Request
	var buf []string

	buf = append(buf, "## REQUEST\n")

	url := fmt.Sprintf("%v %v %v", req.Method, req.URL.Path, req.Proto)
	buf = append(buf, url)

	buf = append(buf, "Accept: "+req.Header.Get("Accept"))
	buf = append(buf, "Content-Type: "+req.Header.Get("Content-Type"))

	// If this is a POST, add post data
	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
		indentBuffer := bytes.NewBuffer(nil)
		json.Indent(indentBuffer, requestBody, "", "  ")
		buf = append(buf, string(indentBuffer.Bytes()))
	}

	buf = append(buf, "## RESONSE\n")
	buf = append(buf, fmt.Sprintf("%s	%s	%d", res.Proto, http.StatusText(res.StatusCode), res.StatusCode))
	buf = append(buf, "Content-Type: "+req.Header.Get("Content-Type"))

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if bodyBytes != nil {
		indentBuffer := bytes.NewBuffer(nil)
		json.Indent(indentBuffer, bodyBytes, "", "  ")
		buf = append(buf, string(indentBuffer.Bytes()))
	}

	// Return the request as a string
	return strings.Join(buf, "\n")
}
