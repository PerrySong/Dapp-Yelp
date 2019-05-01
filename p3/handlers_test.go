package p3

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//func TestAskForBlock(t *testing.T) {
//
//	router := NewRouter()
//	go log.Fatal(http.ListenAndServe(":6688", router))
//	AskForBlock(1, "genesis")
//
//}

func TestStartTryingNonces(t *testing.T) {
	StartTryingNonces()
}

//
//func TestDownload(t *testing.T) {
//	Download()
//}
//
//func TestForwardHeartBeat(t *testing.T) {
//
//}
//
//func TestHeartBeatReceive(t *testing.T) {
//
//}
//
//func TestRegister(t *testing.T) {
//
//}
//
//func TestShow(t *testing.T) {
//
//}
//
//func TestStart(t *testing.T) {
//
//}

func TestUploadBlock(t *testing.T) {
	Download()
	req, err := http.NewRequest("GET", "/block/1/3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48", nil)
	resp := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(UploadBlock)
	handler.ServeHTTP(resp, req)

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		if strings.Contains(string(p), "Error") {
			t.Errorf("header response shouldn't return error: %s", p)
		}
	}
	fmt.Println("body = ", resp.Body)
}

//func TestStartHeartBeat(t *testing.T) {
//
//}
//
//func TestUpload(t *testing.T) {
//
//}
