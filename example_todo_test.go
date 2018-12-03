package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// {
	// 	fmt.Print("# Request Issues\n")
	// 	fmt.Print("\n## REQUEST\n")
	// 	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/"+path.Join(issuesEndpoint), nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	req.Header.Set("Accept", jsonapi.ContentType)
	// 	req.Header.Set("Content-Type", jsonapi.ContentType)
	//
	// 	fmt.Printf("%s	%s	%s\n", req.Method, req.URL.Path, req.Proto)
	//
	// 	res, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	fmt.Print("\n## RESPONSE\n")
	// 	fmt.Printf("%s	%s	%d\n", res.Proto, http.StatusText(res.StatusCode), res.StatusCode)
	// 	fmt.Println("json\n   ")
	// 	io.Copy(os.Stdout, res.Body)
	// 	fmt.Println("\n")
	// }
	//
	// {
	// 	fmt.Print("\n# Create an feature\n")
	//
	// 	fmt.Print("\n## REQUEST\n")
	// 	reqBody := `{"data": {"type": "feature","attributes": {"desc": "As a teapot, I should pour tea"}}}`
	//
	// 	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/"+path.Join(issuesEndpoint), strings.NewReader(reqBody))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	req.Header.Set("Accept", jsonapi.ContentType)
	// 	req.Header.Set("Content-Type", jsonapi.ContentType)
	//
	// 	fmt.Printf("%s	%s	%s\n", req.Method, req.URL.Path, req.Proto)
	// 	fmt.Println(reqBody)
	//
	// 	res, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Print("\n## RESPONSE\n")
	// 	fmt.Printf("%s	%s	%d\n", res.Proto, http.StatusText(res.StatusCode), res.StatusCode)
	// 	fmt.Println("json\n   ")
	// 	io.Copy(os.Stdout, res.Body)
	// 	fmt.Println("\n")
	// }
	//
	{
		fmt.Print("# Create an bug\n")

		reqBody := []byte(`{"data": {"type": "bug","attributes": {"desc": "When tea from teapot is poured out, it is not warm enough."}}}`)

		req, err := http.NewRequest(http.MethodPost, testServer.URL+"/"+path.Join(issuesEndpoint), bytes.NewReader(reqBody))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(printReqRes(res, reqBody))
	}
	//
	// {
	// 	fmt.Print("\n# Fetch one issue\n")
	//
	// 	fmt.Print("\n## REQUEST\n")
	// 	reqBody := `{"data": {"type": "bug","attributes": {"desc": "When tea from teapot is poured out, it is not warm enough."}}}`
	//
	// 	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/"+path.Join(issuesEndpoint, "0"), nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	req.Header.Set("Accept", jsonapi.ContentType)
	// 	req.Header.Set("Content-Type", jsonapi.ContentType)
	//
	// 	fmt.Printf("%s	%s	%s\n", req.Method, req.URL.Path, req.Proto)
	// 	fmt.Println(reqBody)
	//
	// 	res, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Print("\n## RESPONSE\n")
	// 	fmt.Printf("%s	%s	%d\n", res.Proto, http.StatusText(res.StatusCode), res.StatusCode)
	// 	fmt.Println("json\n   ")
	// 	io.Copy(os.Stdout, res.Body)
	// 	fmt.Println("\n")
	// }

	{
		fmt.Print("# Request Issues\n")
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/"+path.Join(issuesEndpoint), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept", jsonapi.ContentType)
		req.Header.Set("Content-Type", jsonapi.ContentType)

		// fmt.Printf("%s	%s	%s\n", req.Method, req.URL.Path, req.Proto)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(printReqRes(res, nil))
		// fmt.Printf("%s	%s	%d\n", res.Proto, http.StatusText(res.StatusCode), res.StatusCode)
		// fmt.Println("json\n   ")
		// io.Copy(os.Stdout, res.Body)
		// fmt.Println("\n")
	}

	// Output:
}

func printReqRes(res *http.Response, requestBody []byte) string {
	req := res.Request
	var buf []string

	buf = append(buf, "## Request\n")

	url := fmt.Sprintf("%v %v %v", req.Method, req.URL.Path, req.Proto)
	buf = append(buf, url)

	buf = append(buf, "Accept: "+req.Header.Get("Accept"))
	buf = append(buf, "Content-Type: "+req.Header.Get("Content-Type"))

	// If this is a POST, add post data
	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
		buf = append(buf, "")
		buf = append(buf, strings.TrimSpace(string(requestBody)))
	}

	buf = append(buf, "\n## Response\n")
	buf = append(buf, fmt.Sprintf("%s	%s	%d", res.Proto, http.StatusText(res.StatusCode), res.StatusCode))
	buf = append(buf, "Content-Type: "+req.Header.Get("Content-Type"))

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(bodyBytes) > 0 {
		buf = append(buf, "")
		buf = append(buf, string(bodyBytes))
	}

	// Return the request as a string
	return strings.Join(buf, "\n") + "\n\n"
}

//
// HTTP/1.1	OK	200
//	{
//		"data": {
//			"type": "feature",
//			"attributes": {
//				"desc": "As a teapot, I should pour tea"
//			}
//	}
//}
