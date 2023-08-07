package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupTest(t *testing.T, reqMethod, urlPath, body string, handlerFunc http.HandlerFunc) (*http.Request, *httptest.ResponseRecorder) {
	Log.SetOutput(io.Discard) //Disable logging for testing
	DataFilePath = filepath.Join("testdata", "testdata.json")

	if err := LoadData(); err != nil {
		t.Fatal(err)
	}

	router, err := NewRouter(DataFilePath)
	if err != nil {
		t.Errorf("Failed to initialize router: %v", err)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(reqMethod, urlPath, bytes.NewBuffer([]byte(body)))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)
	return req, rr
}

func assertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expectedCode int) {
	if status := rr.Code; status != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedCode)
	}
}

func assertResponseBody(t *testing.T, rr *httptest.ResponseRecorder, expectedBody string) {
	var got, want interface{}
	if expectedBody != "" {
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("Failed to unmarshal got body: %v", err)
		}
		if err := json.Unmarshal([]byte(expectedBody), &want); err != nil {
			t.Fatalf("Failed to unmarshal expected body: %v", err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("handler returned unexpected body (-want +got):\n%s", diff)
		}
	}
}

func TestEndpoints(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		body         string
		handler      http.HandlerFunc
		expectedCode int
		expectedBody string
	}{
		{
			name:         "TestListRegions",
			method:       "GET",
			url:          "/regions",
			handler:      ListRegions,
			expectedCode: http.StatusOK,
			expectedBody: `{"amer":{"dev":{"name":"dev","apps":{"app1":{"name":"Application-1","version":"258.2.0"},"app2":{"name":"Application-2","version":"4.2.0"}}},"dr":{"name":"dr","apps":{}},"production":{"name":"production","apps":{"app3":{"name":"Application-3","version":"2.0.0"}}},"qa":{"name":"qa","apps":{}},"uat":{"name":"uat","apps":{"app4":{"name":"Application-4","version":"0.9.0"}}}},"apac":{"dev":{"name":"dev","apps":{}},"dr":{"name":"dr","apps":{}},"production":{"name":"production","apps":{}},"qa":{"name":"qa","apps":{}},"uat":{"name":"uat","apps":{}}},"emea":{"dev":{"name":"dev","apps":{}},"dr":{"name":"dr","apps":{}},"production":{"name":"production","apps":{}},"qa":{"name":"qa","apps":{}},"uat":{"name":"uat","apps":{}}}}`,
		},
		{
			name:         "TestListEnvironments",
			method:       "GET",
			url:          "/regions/amer/environments",
			handler:      ListEnvironments,
			expectedCode: http.StatusOK,
			expectedBody: `{"dev":{"name":"dev","apps":{"app1":{"name":"Application-1","version":"258.2.0"},"app2":{"name":"Application-2","version":"4.2.0"}}},"dr":{"name":"dr","apps":{}},"production":{"name":"production","apps":{"app3":{"name":"Application-3","version":"2.0.0"}}},"qa":{"name":"qa","apps":{}},"uat":{"name":"uat","apps":{"app4":{"name":"Application-4","version":"0.9.0"}}}}`,
		},
		{
			name:         "TestListApps",
			method:       "GET",
			url:          "/regions/amer/environments/dev/apps",
			handler:      ListApps,
			expectedCode: http.StatusOK,
			expectedBody: `{"app1":{"name":"Application-1","version":"258.2.0"},"app2":{"name":"Application-2","version":"4.2.0"}}`,
		},
		{
			name:         "TestCreateRegion",
			method:       "POST",
			url:          "/regions/us-west-1",
			body:         "",
			handler:      CreateRegion,
			expectedCode: http.StatusCreated,
			expectedBody: ``,
		},

		{
			name:         "TestCreateEnvironment",
			method:       "POST",
			url:          "/regions/amer/sit",
			body:         "",
			handler:      CreateEnvironment,
			expectedCode: http.StatusCreated,
			expectedBody: ``,
		},
		{
			name:         "TestCreateApp",
			method:       "POST",
			url:          "/regions/amer/dev/testapp",
			body:         "",
			handler:      CreateApp,
			expectedCode: http.StatusCreated,
			expectedBody: ``,
		},
		{
			name:         "TestGetApp",
			method:       "GET",
			url:          "/regions/amer/environments/dev/apps/app1",
			handler:      GetApp,
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"Application-1","version":"258.2.0"}`,
		},
		{
			name:         "TestAppNotFound",
			method:       "GET",
			url:          "/regions/amer/environments/dev/apps/Application-10",
			handler:      GetApp,
			expectedCode: http.StatusNotFound,
			expectedBody: ``,
		},
		// More test cases can be added following this pattern
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, rr := setupTest(t, tt.method, tt.url, tt.body, tt.handler)
			assertStatusCode(t, rr, tt.expectedCode)
			assertResponseBody(t, rr, tt.expectedBody)
		})
	}
}
