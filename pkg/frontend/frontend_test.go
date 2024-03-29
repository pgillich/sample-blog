package frontend

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/pgillich/sample-blog/api"
	"github.com/pgillich/sample-blog/internal/dao"
	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/internal/test"
)

func TestMain(m *testing.M) {
	logger.Init(test.GetLogLevel())
	gin.SetMode(gin.TestMode)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestGetUserPostCommentStatsFailed(t *testing.T) {
	testGetUserPostCommentStats(t, "négy", "2019-12-05 12:00:00+00:00", http.StatusBadRequest, 0, `{}`)
}

func TestGetUserPostCommentStats(t *testing.T) {
	testGetUserPostCommentStats(t, "4", "2019-12-05 12:00:00+00:00", http.StatusOK, 3, `{
  "1": {
    "userName": "kovacsj",
    "entries": 1,
    "comments": 8
  },
  "2": {
    "userName": "szabop",
    "entries": 4,
    "comments": 4
  },
  "3": {
    "userName": "kocsisi",
    "entries": 1,
    "comments": 0
  }
}`)
}

func testGetUserPostCommentStats(t *testing.T,
	days string, now string, expectedStatus int, expectedUserNum int, expectedBody string,
) {
	dbHandler, err := dao.NewHandler(
		"sqlite3", ":memory:", dao.GetDefaultSampleFill())
	if err != nil {
		logger.Get().Panic(err)
	}
	defer dbHandler.Close()

	timeNow, _ := time.Parse(time.RFC3339, now) //nolint:errcheck
	timeNowFunc := func() time.Time {
		return timeNow
	}
	dbHandler.TimeNow = timeNowFunc

	router := SetupGin(gin.New(), dbHandler, false)

	httpResponse := performRequest("GET", "/api/v1/stat/user-post-comment?days="+days, nil, nil, router)
	assert.Equal(t, expectedStatus, httpResponse.Code, "GET stat/user-post-comment")

	if httpResponse.Code != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	assert.NoError(t, err, "Body stat/user-post-comment")

	response := api.UserPostCommentStats{}

	assert.NoError(t, json.Unmarshal(body, &response), "Body stat/user-post-comment")
	assert.Equal(t, expectedUserNum, len(response), "Users stat/user-post-comment")
	assert.Equal(t, expectedBody, test.JSONMarshal(&response), "Stat stat/user-post-comment")
}

func performRequest(method, target string, header http.Header, body io.Reader, router http.Handler) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)
	if header != nil {
		r.Header = header
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	return w
}

func TestPostComment(t *testing.T) {
	testPostComment(t, "kovacsj", "kovacsj", "kovacs12", http.StatusOK, http.StatusOK)
}

func TestPostCommentFailed(t *testing.T) {
	testPostComment(t, "szabop", "kovacsj", "kovacs12", http.StatusOK, http.StatusBadRequest)
}

func testPostComment(t *testing.T,
	entryUser string, committerUser string, committerPassword string,
	expectedStatusLogin int, expectedStatus int,
) {
	dbHandler, err := dao.NewHandler(
		"sqlite3", ":memory:", dao.GetDefaultSampleFill())
	if err != nil {
		logger.Get().Panic(err)
	}
	defer dbHandler.Close()

	router := SetupGin(gin.New(), dbHandler, false)

	loginBody := login{Username: committerUser, Password: committerPassword}
	httpResponse := performRequest("POST", "/api/v1/login", test.GetHTTPHeaderJSON(),
		strings.NewReader(test.JSONMarshal(loginBody)), router)
	assert.Equal(t, expectedStatusLogin, httpResponse.Code, "POST login")

	if httpResponse.Code != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	assert.NoError(t, err, "Body login")

	loginResponse := gin.H{}
	assert.NoError(t, json.Unmarshal(body, &loginResponse), "Body login")

	token := fmt.Sprintf("%s", loginResponse["token"])
	assert.NotEmpty(t, token, "Token login")

	entries, err := dbHandler.GetUserEntriesByName(entryUser)
	assert.NoError(t, err, "Get entries")
	assert.NotEmpty(t, entries, "Get entries")

	entry := entries[0]
	text := api.Text{Text: "HELLO"}

	httpResponse = performRequest("POST", fmt.Sprintf("/api/v1/entry/%d/comment", entry.ID),
		test.GetHTTPHeaderJSONToken(token), strings.NewReader(test.JSONMarshal(text)), router)
	assert.Equal(t, expectedStatus, httpResponse.Code, "POST /api/v1/entry/%d/comment")

	if httpResponse.Code != http.StatusOK {
		return
	}

	body, err = ioutil.ReadAll(httpResponse.Body)
	assert.NoError(t, err, "Body /api/v1/entry/%d/comment")

	response := api.Comment{}
	assert.NoError(t, json.Unmarshal(body, &response), "Body /api/v1/entry/%d/comment")
}
