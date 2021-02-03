package es

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("When no version set should return error", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(make(map[string]interface{})),
		}

		_, err := NewClient(context.Background(), ds, nil)
		require.Error(t, err)
	})

	t.Run("When no time field name set should return error", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 5,
			}),
		}

		_, err := NewClient(context.Background(), ds, nil)
		require.Error(t, err)
	})

	t.Run("When unsupported version set should return error", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 6,
				"timeField": "@timestamp",
			}),
		}

		_, err := NewClient(context.Background(), ds, nil)
		require.Error(t, err)
	})

	t.Run("When version 2 should return v2 client", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 2,
				"timeField": "@timestamp",
			}),
		}

		c, err := NewClient(context.Background(), ds, nil)
		require.NoError(t, err)
		assert.Equal(t, 2, c.GetVersion())
	})

	t.Run("When version 5 should return v5 client", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 5,
				"timeField": "@timestamp",
			}),
		}

		c, err := NewClient(context.Background(), ds, nil)
		require.NoError(t, err)
		assert.Equal(t, 5, c.GetVersion())
	})

	t.Run("When version 56 should return v5.6 client", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 56,
				"timeField": "@timestamp",
			}),
		}

		c, err := NewClient(context.Background(), ds, nil)
		require.NoError(t, err)
		assert.Equal(t, 56, c.GetVersion())
	})

	t.Run("When version 60 should return v6.0 client", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 60,
				"timeField": "@timestamp",
			}),
		}

		c, err := NewClient(context.Background(), ds, nil)
		require.NoError(t, err)
		assert.Equal(t, 60, c.GetVersion())
	})

	t.Run("When version 70 should return v7.0 client", func(t *testing.T) {
		ds := &models.DataSource{
			JsonData: simplejson.NewFromAny(map[string]interface{}{
				"esVersion": 70,
				"timeField": "@timestamp",
			}),
		}

		c, err := NewClient(context.Background(), ds, nil)
		require.NoError(t, err)
		assert.Equal(t, 70, c.GetVersion())
		// Convey("Given a fake http client", func() {
		// 	var responseBuffer *bytes.Buffer
		// 	var req *http.Request
		// 	handleResponse := func(rw http.ResponseWriter) {}
		// 	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// 		req = r
		// 		buf, err := ioutil.ReadAll(r.Body)
		// 		if err != nil {
		// 			t.Fatalf("Failed to read response body, err=%v", err)
		// 		}
		// 		responseBuffer = bytes.NewBuffer(buf)
		// 		handleResponse(rw)
		// 	}))

		// 	currentNewDatasourceHttpClient := newDatasourceHttpClient

		// 	newDatasourceHttpClient = func(ds *models.DataSource) (*http.Client, error) {
		// 		return ts.Client(), nil
		// 	}

		// 	originalTimeNowUtc := timeNowUtc
		// 	timeNowUtc = func() time.Time {
		// 		return time.Date(2018, 5, 15, 17, 50, 0, 0, time.UTC)
		// 	}
		// 	from := time.Date(2018, 5, 15, 17, 50, 0, 0, time.UTC)
		// 	to := time.Date(2018, 5, 15, 17, 55, 0, 0, time.UTC)
		// 	fromStr := fmt.Sprintf("%d", from.UnixNano()/int64(time.Millisecond))
		// 	toStr := fmt.Sprintf("%d", to.UnixNano()/int64(time.Millisecond))
		// 	timeRange := tsdb.NewTimeRange(fromStr, toStr)

		// 	Convey("and a v2.x client", func() {
		// 		ds := models.DataSource{
		// 			Database: "[metrics-]YYYY.MM.DD",
		// 			Url:      ts.URL,
		// 			JsonData: simplejson.NewFromAny(map[string]interface{}{
		// 				"esVersion": 2,
		// 				"timeField": "@timestamp",
		// 				"interval":  "Daily",
		// 			}),
		// 		}

		// 		c, err := NewClient(context.Background(), &ds, timeRange)
		// 		So(err, ShouldBeNil)
		// 		So(c, ShouldNotBeNil)

		// 		Convey("When executing multi search", func() {
		// 			ms, err := createMultisearchForTest(c)
		// 			So(err, ShouldBeNil)
		// 			c.ExecuteMultisearch(ms)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodPost)
		// 				So(req.URL.Path, ShouldEqual, "/_msearch")

		// 				So(responseBuffer, ShouldNotBeNil)

		// 				headerBytes, err := responseBuffer.ReadBytes('\n')
		// 				So(err, ShouldBeNil)
		// 				bodyBytes := responseBuffer.Bytes()

		// 				jHeader, err := simplejson.NewJson(headerBytes)
		// 				So(err, ShouldBeNil)

		// 				jBody, err := simplejson.NewJson(bodyBytes)
		// 				So(err, ShouldBeNil)

		// 				So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
		// 				So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
		// 				So(jHeader.Get("search_type").MustString(), ShouldEqual, "count")
		// 				So(jHeader.Get("max_concurrent_shard_requests").MustInt(10), ShouldEqual, 10)

		// 				Convey("and replace $__interval variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
		// 				})

		// 				Convey("and replace $__interval_ms variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
		// 				})
		// 			})
		// 		})

		// 		Convey("When get index mapping (200 OK)", func() {
		// 			handleResponse = createIndexMappingFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse successful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 200)
		// 				So(res.Error, ShouldHaveLength, 0)
		// 				So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
		// 			})
		// 		})

		// 		Convey("When get index mapping (404 Not Found)", func() {
		// 			handleResponse = createIndexMappingNotFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse unsuccessful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 404)
		// 				So(res.Mappings, ShouldHaveLength, 0)
		// 				So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
		// 			})
		// 		})
		// 	})

		// 	Convey("and a v5.x client", func() {
		// 		ds := models.DataSource{
		// 			Database: "[metrics-]YYYY.MM.DD",
		// 			Url:      ts.URL,
		// 			JsonData: simplejson.NewFromAny(map[string]interface{}{
		// 				"esVersion":                  5,
		// 				"maxConcurrentShardRequests": 100,
		// 				"timeField":                  "@timestamp",
		// 				"interval":                   "Daily",
		// 			}),
		// 		}

		// 		c, err := NewClient(context.Background(), &ds, timeRange)
		// 		So(err, ShouldBeNil)
		// 		So(c, ShouldNotBeNil)

		// 		Convey("When executing multi search", func() {
		// 			ms, err := createMultisearchForTest(c)
		// 			So(err, ShouldBeNil)
		// 			c.ExecuteMultisearch(ms)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodPost)
		// 				So(req.URL.Path, ShouldEqual, "/_msearch")

		// 				So(responseBuffer, ShouldNotBeNil)

		// 				headerBytes, err := responseBuffer.ReadBytes('\n')
		// 				So(err, ShouldBeNil)
		// 				bodyBytes := responseBuffer.Bytes()

		// 				jHeader, err := simplejson.NewJson(headerBytes)
		// 				So(err, ShouldBeNil)

		// 				jBody, err := simplejson.NewJson(bodyBytes)
		// 				So(err, ShouldBeNil)

		// 				So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
		// 				So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
		// 				So(jHeader.Get("search_type").MustString(), ShouldEqual, "query_then_fetch")
		// 				So(jHeader.Get("max_concurrent_shard_requests").MustInt(10), ShouldEqual, 10)

		// 				Convey("and replace $__interval variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
		// 				})

		// 				Convey("and replace $__interval_ms variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
		// 				})
		// 			})
		// 		})

		// 		Convey("When get index mapping (200 OK)", func() {
		// 			handleResponse = createIndexMappingFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse successful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 200)
		// 				So(res.Error, ShouldHaveLength, 0)
		// 				So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
		// 			})
		// 		})

		// 		Convey("When get index mapping (404 Not Found)", func() {
		// 			handleResponse = createIndexMappingNotFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse unsuccessful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 404)
		// 				So(res.Mappings, ShouldHaveLength, 0)
		// 				So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
		// 			})
		// 		})
		// 	})

		// 	Convey("and a v5.6 client", func() {
		// 		ds := models.DataSource{
		// 			Database: "[metrics-]YYYY.MM.DD",
		// 			Url:      ts.URL,
		// 			JsonData: simplejson.NewFromAny(map[string]interface{}{
		// 				"esVersion":                  56,
		// 				"maxConcurrentShardRequests": 100,
		// 				"timeField":                  "@timestamp",
		// 				"interval":                   "Daily",
		// 			}),
		// 		}

		// 		c, err := NewClient(context.Background(), &ds, timeRange)
		// 		So(err, ShouldBeNil)
		// 		So(c, ShouldNotBeNil)

		// 		Convey("When executing multi search", func() {
		// 			ms, err := createMultisearchForTest(c)
		// 			So(err, ShouldBeNil)
		// 			c.ExecuteMultisearch(ms)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodPost)
		// 				So(req.URL.Path, ShouldEqual, "/_msearch")

		// 				So(responseBuffer, ShouldNotBeNil)

		// 				headerBytes, err := responseBuffer.ReadBytes('\n')
		// 				So(err, ShouldBeNil)
		// 				bodyBytes := responseBuffer.Bytes()

		// 				jHeader, err := simplejson.NewJson(headerBytes)
		// 				So(err, ShouldBeNil)

		// 				jBody, err := simplejson.NewJson(bodyBytes)
		// 				So(err, ShouldBeNil)

		// 				So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
		// 				So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
		// 				So(jHeader.Get("search_type").MustString(), ShouldEqual, "query_then_fetch")
		// 				So(jHeader.Get("max_concurrent_shard_requests").MustInt(), ShouldEqual, 100)

		// 				Convey("and replace $__interval variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
		// 				})

		// 				Convey("and replace $__interval_ms variable", func() {
		// 					So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
		// 				})
		// 			})
		// 		})

		// 		Convey("When get index mapping (200 OK)", func() {
		// 			handleResponse = createIndexMappingFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse successful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 200)
		// 				So(res.Error, ShouldHaveLength, 0)
		// 				So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
		// 			})
		// 		})

		// 		Convey("When get index mapping (404 Not Found)", func() {
		// 			handleResponse = createIndexMappingNotFoundResponse
		// 			res, err := c.GetIndexMapping()
		// 			So(err, ShouldBeNil)

		// 			Convey("Should send correct request and payload", func() {
		// 				So(req, ShouldNotBeNil)
		// 				So(req.Method, ShouldEqual, http.MethodGet)
		// 				So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
		// 			})

		// 			Convey("Should parse unsuccessful response", func() {
		// 				So(res, ShouldNotBeNil)
		// 				So(res.StatusCode, ShouldEqual, 404)
		// 				So(res.Mappings, ShouldHaveLength, 0)
		// 				So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
		// 			})
		// 		})
		// 	})

		// 	Reset(func() {
		// 		newDatasourceHttpClient = currentNewDatasourceHttpClient
		// 		timeNowUtc = originalTimeNowUtc
		// 	})
		// })
	})
}

func TestClient_ExecuteMultisearch(t *testing.T) {
	httpClientScenario(t, "Given a fake http client and a v2.x client with response", &models.DataSource{
		Database: "[metrics-]YYYY.MM.DD",
		JsonData: simplejson.NewFromAny(map[string]interface{}{
			"esVersion": 2,
			"timeField": "@timestamp",
			"interval":  "Daily",
		}),
	}, func(sc *scenarioContext) {
		sc.responseBody = `{
				"responses": [
					{
						"hits": { "hits": [], "max_score": 0, "total": 4656 },
						"status": 200
					}
				]
			}`

		ms, err := createMultisearchForTest(t, sc.client)
		require.NoError(t, err)
		res, err := sc.client.ExecuteMultisearch(ms)
		require.NoError(t, err)

		require.NotNil(t, sc.request)
		assert.Equal(t, http.MethodPost, sc.request.Method)
		assert.Equal(t, "/_msearch", sc.request.URL.Path)

		require.NotNil(t, sc.requestBody)
		headerBytes, err := sc.requestBody.ReadBytes('\n')
		require.NoError(t, err)
		bodyBytes := sc.requestBody.Bytes()

		jHeader, err := simplejson.NewJson(headerBytes)
		require.NoError(t, err)

		jBody, err := simplejson.NewJson(bodyBytes)
		require.NoError(t, err)

		assert.Equal(t, "metrics-2018.05.15", jHeader.Get("index").MustString())
		assert.True(t, jHeader.Get("ignore_unavailable").MustBool(false))
		assert.Equal(t, "count", jHeader.Get("search_type").MustString())
		assert.Empty(t, jHeader.Get("max_concurrent_shard_requests"))

		assert.Equal(t, "15000*@hostname", jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString())

		assert.Equal(t, "15s", jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString())

		assert.Equal(t, 200, res.Status)
		require.Len(t, res.Responses, 1)
	})

	httpClientScenario(t, "Given a fake http client and a v5.x client with response", &models.DataSource{
		Database: "[metrics-]YYYY.MM.DD",
		JsonData: simplejson.NewFromAny(map[string]interface{}{
			"esVersion":                  5,
			"maxConcurrentShardRequests": 100,
			"timeField":                  "@timestamp",
			"interval":                   "Daily",
		}),
	}, func(sc *scenarioContext) {
		sc.responseBody = `{
				"responses": [
					{
						"hits": { "hits": [], "max_score": 0, "total": 4656 },
						"status": 200
					}
				]
			}`

		ms, err := createMultisearchForTest(t, sc.client)
		require.NoError(t, err)
		res, err := sc.client.ExecuteMultisearch(ms)
		require.NoError(t, err)

		require.NotNil(t, sc.request)
		assert.Equal(t, http.MethodPost, sc.request.Method)
		assert.Equal(t, "/_msearch", sc.request.URL.Path)

		require.NotNil(t, sc.requestBody)

		headerBytes, err := sc.requestBody.ReadBytes('\n')
		require.NoError(t, err)
		bodyBytes := sc.requestBody.Bytes()

		jHeader, err := simplejson.NewJson(headerBytes)
		require.NoError(t, err)

		jBody, err := simplejson.NewJson(bodyBytes)
		require.NoError(t, err)

		assert.Equal(t, "metrics-2018.05.15", jHeader.Get("index").MustString())
		assert.True(t, jHeader.Get("ignore_unavailable").MustBool(false))
		assert.Equal(t, "query_then_fetch", jHeader.Get("search_type").MustString())
		assert.Empty(t, jHeader.Get("max_concurrent_shard_requests"))

		assert.Equal(t, "15000*@hostname", jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString())

		assert.Equal(t, "15s", jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString())

		assert.Equal(t, 200, res.Status)
		require.Len(t, res.Responses, 1)
	})

	httpClientScenario(t, "Given a fake http client and a v5.6 client with response", &models.DataSource{
		Database: "[metrics-]YYYY.MM.DD",
		JsonData: simplejson.NewFromAny(map[string]interface{}{
			"esVersion":                  56,
			"maxConcurrentShardRequests": 100,
			"timeField":                  "@timestamp",
			"interval":                   "Daily",
		}),
	}, func(sc *scenarioContext) {
		sc.responseBody = `{
				"responses": [
					{
						"hits": { "hits": [], "max_score": 0, "total": 4656 },
						"status": 200
					}
				]
			}`

		ms, err := createMultisearchForTest(t, sc.client)
		require.NoError(t, err)
		res, err := sc.client.ExecuteMultisearch(ms)
		require.NoError(t, err)

		require.NotNil(t, sc.request)
		assert.Equal(t, http.MethodPost, sc.request.Method)
		assert.Equal(t, "/_msearch", sc.request.URL.Path)

		require.NotNil(t, sc.requestBody)

		headerBytes, err := sc.requestBody.ReadBytes('\n')
		require.NoError(t, err)
		bodyBytes := sc.requestBody.Bytes()

		jHeader, err := simplejson.NewJson(headerBytes)
		require.NoError(t, err)

		jBody, err := simplejson.NewJson(bodyBytes)
		require.NoError(t, err)

		assert.Equal(t, "metrics-2018.05.15", jHeader.Get("index").MustString())
		assert.True(t, jHeader.Get("ignore_unavailable").MustBool(false))
		assert.Equal(t, "query_then_fetch", jHeader.Get("search_type").MustString())
		assert.Equal(t, 100, jHeader.Get("max_concurrent_shard_requests").MustInt())

		assert.Equal(t, "15000*@hostname", jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString())

		assert.Equal(t, "15s", jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString())

		assert.Equal(t, 200, res.Status)
		require.Len(t, res.Responses, 1)
	})

	httpClientScenario(t, "Given a fake http client and a v7.0 client with response", &models.DataSource{
		Database: "[metrics-]YYYY.MM.DD",
		JsonData: simplejson.NewFromAny(map[string]interface{}{
			"esVersion":                  70,
			"maxConcurrentShardRequests": 6,
			"timeField":                  "@timestamp",
			"interval":                   "Daily",
		}),
	}, func(sc *scenarioContext) {
		sc.responseBody = `{
				"responses": [
					{
						"hits": {	"hits": [], "max_score": 0,	"total": { "value": 4656, "relation": "eq"}	},
						"status": 200
					}
				]
			}`

		ms, err := createMultisearchForTest(t, sc.client)
		require.NoError(t, err)
		res, err := sc.client.ExecuteMultisearch(ms)
		require.NoError(t, err)

		require.NotNil(t, sc.request)
		assert.Equal(t, http.MethodPost, sc.request.Method)
		assert.Equal(t, "/_msearch", sc.request.URL.Path)
		assert.Equal(t, "max_concurrent_shard_requests=6", sc.request.URL.RawQuery)

		require.NotNil(t, sc.requestBody)

		headerBytes, err := sc.requestBody.ReadBytes('\n')
		require.NoError(t, err)
		bodyBytes := sc.requestBody.Bytes()

		jHeader, err := simplejson.NewJson(headerBytes)
		require.NoError(t, err)

		jBody, err := simplejson.NewJson(bodyBytes)
		require.NoError(t, err)

		assert.Equal(t, "metrics-2018.05.15", jHeader.Get("index").MustString())
		assert.True(t, jHeader.Get("ignore_unavailable").MustBool(false))
		assert.Equal(t, "query_then_fetch", jHeader.Get("search_type").MustString())
		assert.Empty(t, jHeader.Get("max_concurrent_shard_requests"))

		assert.Equal(t, "15000*@hostname", jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString())

		assert.Equal(t, "15s", jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString())

		assert.Equal(t, 200, res.Status)
		require.Len(t, res.Responses, 1)
	})
}

func createMultisearchForTest(t *testing.T, c Client) (*MultiSearchRequest, error) {
	t.Helper()
	msb := c.MultiSearch()
	s := msb.Search(tsdb.Interval{Value: 15 * time.Second, Text: "15s"})
	s.Agg().DateHistogram("2", "@timestamp", func(a *DateHistogramAgg, ab AggBuilder) {
		a.Interval = "$__interval"

		ab.Metric("1", "avg", "@hostname", func(a *MetricAggregation) {
			a.Settings["script"] = "$__interval_ms*@hostname"
		})
	})
	return msb.Build()
}

func createSearchForTest(c Client) (*SearchRequest, error) {
	s := c.Search(tsdb.Interval{Value: 15 * time.Second, Text: "15s"})
	s.Agg().DateHistogram("2", "@timestamp", func(a *DateHistogramAgg, ab AggBuilder) {
		a.Interval = "$__interval"

		ab.Metric("1", "avg", "@hostname", func(a *MetricAggregation) {
			a.Settings["script"] = "$__interval_ms*@hostname"
		})
	})
	return s.Build()
}

// 				c, err := NewClient(context.Background(), &ds, timeRange)
// 				So(err, ShouldBeNil)
// 				So(c, ShouldNotBeNil)

// 				Convey("When executing search", func() {
// 					s, err := createSearchForTest(c)
// 					So(err, ShouldBeNil)
// 					c.ExecuteSearch(s)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_search?search_type=count&ignore_unavailable=true")

// 						So(responseBuffer, ShouldNotBeNil)
// 						bodyBytes := responseBuffer.Bytes()
// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})
// 				})

// 				Convey("When executing multi search", func() {
// 					ms, err := createMultisearchForTest(c)
// 					So(err, ShouldBeNil)
// 					handleResponse = createMultiSearchSuccessResponse
// 					res, err := c.ExecuteMultisearch(ms)
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/_msearch")

// 						So(responseBuffer, ShouldNotBeNil)

// 						headerBytes, err := responseBuffer.ReadBytes('\n')
// 						So(err, ShouldBeNil)
// 						bodyBytes := responseBuffer.Bytes()

// 						jHeader, err := simplejson.NewJson(headerBytes)
// 						So(err, ShouldBeNil)

// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
// 						So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
// 						So(jHeader.Get("search_type").MustString(), ShouldEqual, "count")
// 						So(jHeader.Get("max_concurrent_shard_requests").MustInt(10), ShouldEqual, 10)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Responses, ShouldHaveLength, 1)
// 						So(res.Responses[0].Error, ShouldHaveLength, 0)
// 						So(res.Responses[0].StatusCode, ShouldEqual, 200)
// 						So(res.Responses[0].Hits.Total, ShouldEqual, 100)
// 						So(res.Responses[0].Aggregations["aggstest"].(string), ShouldEqual, "aggstest")
// 					})
// 				})

// 				Convey("When get index mapping (200 OK)", func() {
// 					handleResponse = createIndexMappingFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Error, ShouldHaveLength, 0)
// 						So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
// 					})
// 				})

// 				Convey("When get index mapping (404 Not Found)", func() {
// 					handleResponse = createIndexMappingNotFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse unsuccessful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 404)
// 						So(res.Mappings, ShouldHaveLength, 0)
// 						So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
// 					})
// 				})
// 			})

// 			Convey("and a v5.x client", func() {
// 				ds := models.DataSource{
// 					Database: "[metrics-]YYYY.MM.DD",
// 					Url:      ts.URL,
// 					JsonData: simplejson.NewFromAny(map[string]interface{}{
// 						"esVersion":                  5,
// 						"maxConcurrentShardRequests": 100,
// 						"timeField":                  "@timestamp",
// 						"interval":                   "Daily",
// 					}),
// 				}

// 				c, err := NewClient(context.Background(), &ds, timeRange)
// 				So(err, ShouldBeNil)
// 				So(c, ShouldNotBeNil)

// 				Convey("When executing search", func() {
// 					s, err := createSearchForTest(c)
// 					So(err, ShouldBeNil)
// 					c.ExecuteSearch(s)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_search?search_type=query_then_fetch&ignore_unavailable=true")

// 						So(responseBuffer, ShouldNotBeNil)
// 						bodyBytes := responseBuffer.Bytes()
// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})
// 				})

// 				Convey("When executing multi search", func() {
// 					ms, err := createMultisearchForTest(c)
// 					So(err, ShouldBeNil)
// 					handleResponse = createMultiSearchSuccessResponse
// 					res, err := c.ExecuteMultisearch(ms)
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/_msearch")

// 						So(responseBuffer, ShouldNotBeNil)

// 						headerBytes, err := responseBuffer.ReadBytes('\n')
// 						So(err, ShouldBeNil)
// 						bodyBytes := responseBuffer.Bytes()

// 						jHeader, err := simplejson.NewJson(headerBytes)
// 						So(err, ShouldBeNil)

// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
// 						So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
// 						So(jHeader.Get("search_type").MustString(), ShouldEqual, "query_then_fetch")
// 						So(jHeader.Get("max_concurrent_shard_requests").MustInt(10), ShouldEqual, 10)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Responses, ShouldHaveLength, 1)
// 						So(res.Responses[0].Error, ShouldHaveLength, 0)
// 						So(res.Responses[0].StatusCode, ShouldEqual, 200)
// 						So(res.Responses[0].Hits.Total, ShouldEqual, 100)
// 						So(res.Responses[0].Aggregations["aggstest"].(string), ShouldEqual, "aggstest")
// 					})
// 				})

// 				Convey("When get index mapping (200 OK)", func() {
// 					handleResponse = createIndexMappingFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Error, ShouldHaveLength, 0)
// 						So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
// 					})
// 				})

// 				Convey("When get index mapping (404 Not Found)", func() {
// 					handleResponse = createIndexMappingNotFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse unsuccessful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 404)
// 						So(res.Mappings, ShouldHaveLength, 0)
// 						So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
// 					})
// 				})
// 			})

// 			Convey("and a v5.6 client", func() {
// 				ds := models.DataSource{
// 					Database: "[metrics-]YYYY.MM.DD",
// 					Url:      ts.URL,
// 					JsonData: simplejson.NewFromAny(map[string]interface{}{
// 						"esVersion":                  56,
// 						"maxConcurrentShardRequests": 100,
// 						"timeField":                  "@timestamp",
// 						"interval":                   "Daily",
// 					}),
// 				}

// 				c, err := NewClient(context.Background(), &ds, timeRange)
// 				So(err, ShouldBeNil)
// 				So(c, ShouldNotBeNil)

// 				Convey("When executing search", func() {
// 					s, err := createSearchForTest(c)
// 					So(err, ShouldBeNil)
// 					c.ExecuteSearch(s)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_search?search_type=query_then_fetch&max_concurrent_shard_requests=100&ignore_unavailable=true")

// 						So(responseBuffer, ShouldNotBeNil)
// 						bodyBytes := responseBuffer.Bytes()
// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})
// 				})

// 				Convey("When executing multi search", func() {
// 					ms, err := createMultisearchForTest(c)
// 					So(err, ShouldBeNil)
// 					handleResponse = createMultiSearchSuccessResponse
// 					res, err := c.ExecuteMultisearch(ms)
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodPost)
// 						So(req.URL.Path, ShouldEqual, "/_msearch")

// 						So(responseBuffer, ShouldNotBeNil)

// 						headerBytes, err := responseBuffer.ReadBytes('\n')
// 						So(err, ShouldBeNil)
// 						bodyBytes := responseBuffer.Bytes()

// 						jHeader, err := simplejson.NewJson(headerBytes)
// 						So(err, ShouldBeNil)

// 						jBody, err := simplejson.NewJson(bodyBytes)
// 						So(err, ShouldBeNil)

// 						So(jHeader.Get("index").MustString(), ShouldEqual, "metrics-2018.05.15")
// 						So(jHeader.Get("ignore_unavailable").MustBool(false), ShouldEqual, true)
// 						So(jHeader.Get("search_type").MustString(), ShouldEqual, "query_then_fetch")
// 						So(jHeader.Get("max_concurrent_shard_requests").MustInt(), ShouldEqual, 100)

// 						Convey("and replace $__interval variable", func() {
// 							So(jBody.GetPath("aggs", "2", "aggs", "1", "avg", "script").MustString(), ShouldEqual, "15000*@hostname")
// 						})

// 						Convey("and replace $__interval_ms variable", func() {
// 							So(jBody.GetPath("aggs", "2", "date_histogram", "interval").MustString(), ShouldEqual, "15s")
// 						})
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Responses, ShouldHaveLength, 1)
// 						So(res.Responses[0].Error, ShouldHaveLength, 0)
// 						So(res.Responses[0].StatusCode, ShouldEqual, 200)
// 						So(res.Responses[0].Hits.Total, ShouldEqual, 100)
// 						So(res.Responses[0].Aggregations["aggstest"].(string), ShouldEqual, "aggstest")
// 					})
// 				})

// 				Convey("When get index mapping (200 OK)", func() {
// 					handleResponse = createIndexMappingFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse successful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 200)
// 						So(res.Error, ShouldHaveLength, 0)
// 						So(res.Mappings["metrics-2018.05.15"].(string), ShouldEqual, "test")
// 					})
// 				})

// 				Convey("When get index mapping (404 Not Found)", func() {
// 					handleResponse = createIndexMappingNotFoundResponse
// 					res, err := c.GetIndexMapping()
// 					So(err, ShouldBeNil)

// 					Convey("Should send correct request and payload", func() {
// 						So(req, ShouldNotBeNil)
// 						So(req.Method, ShouldEqual, http.MethodGet)
// 						So(req.URL.Path, ShouldEqual, "/metrics-2018.05.15/_mapping")
// 					})

// 					Convey("Should parse unsuccessful response", func() {
// 						So(res, ShouldNotBeNil)
// 						So(res.StatusCode, ShouldEqual, 404)
// 						So(res.Mappings, ShouldHaveLength, 0)
// 						So(res.Error["type"].(string), ShouldEqual, "index_not_found_exception")
// 					})
// 				})
// 			})

// 			Reset(func() {
// 				newDatasourceHttpClient = currentNewDatasourceHttpClient
// 				timeNowUtc = originalTimeNowUtc
// 			})
// 		})
// 	})
// }

type scenarioContext struct {
	client         Client
	request        *http.Request
	requestBody    *bytes.Buffer
	responseStatus int
	responseBody   string
}

type scenarioFunc func(*scenarioContext)

func httpClientScenario(t *testing.T, desc string, ds *models.DataSource, fn scenarioFunc) {
	t.Helper()

	t.Run(desc, func(t *testing.T) {
		sc := &scenarioContext{
			responseStatus: 200,
			responseBody:   `{ "responses": [] }`,
		}
		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			sc.request = r
			buf, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)

			sc.requestBody = bytes.NewBuffer(buf)

			rw.Header().Set("Content-Type", "application/json")
			_, err = rw.Write([]byte(sc.responseBody))
			require.NoError(t, err)
			rw.WriteHeader(sc.responseStatus)
		}))
		ds.Url = ts.URL

		from := time.Date(2018, 5, 15, 17, 50, 0, 0, time.UTC)
		to := time.Date(2018, 5, 15, 17, 55, 0, 0, time.UTC)
		fromStr := fmt.Sprintf("%d", from.UnixNano()/int64(time.Millisecond))
		toStr := fmt.Sprintf("%d", to.UnixNano()/int64(time.Millisecond))
		timeRange := tsdb.NewTimeRange(fromStr, toStr)

		c, err := NewClient(context.Background(), ds, timeRange)
		require.NoError(t, err)
		require.NotNil(t, c)
		sc.client = c

		currentNewDatasourceHTTPClient := newDatasourceHttpClient

		newDatasourceHttpClient = func(ds *models.DataSource) (*http.Client, error) {
			return ts.Client(), nil
		}

		t.Cleanup(func() {
			ts.Close()
			newDatasourceHttpClient = currentNewDatasourceHTTPClient
		})

		fn(sc)
	})
}

func createMultiSearchSuccessResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, `{
		"responses": [
			{
				"hits": {
					"total": 100
				},
				"aggregations": {
					"aggstest": "aggstest"
				}
			}
		]
	}`)
}

func createIndexMappingFoundResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, `{
		"metrics-2018.05.15": "test"
	}`)
}

func createIndexMappingNotFoundResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	io.WriteString(rw, `{
		"error": {
			"root_cause": [
				{
					"type": "index_not_found_exception",
					"reason": "no such index",
					"resource.type": "index_or_alias",
					"resource.id": "metrics-2018.05.15",
					"index_uuid": "_na_",
					"index": "metrics-2018.05.15"
				}
			],
			"type": "index_not_found_exception",
			"reason": "no such index",
			"resource.type": "index_or_alias",
			"resource.id": "metrics-2018.05.15",
			"index_uuid": "_na_",
			"index": "metrics-2018.05.15"
		},
		"status": 404
	}`)
}
