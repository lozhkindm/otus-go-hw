package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	defer truncateTable()

	req := eventCreateRequest{
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode, response, err := createEvent(req)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	expected := eventResponse{
		ID:          response.ID,
		UserID:      100500,
		Title:       "title 100500",
		Description: stringPtr("description 100500"),
		StartAt:     time.Unix(1658584623, 0),
		EndAt:       time.Unix(1658584624, 0),
		NotifyAt:    timePtr(time.Unix(1658584625, 0)),
	}

	require.Equal(t, expected, response)
}

func TestUpdateEvent(t *testing.T) {
	defer truncateTable()

	createReq := eventCreateRequest{
		UserID:      100600,
		Title:       "title 100600",
		Description: stringPtr("description 100600"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode, createRes, err := createEvent(createReq)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	updateReq := eventUpdateRequest{
		Title:       "title 100600 new",
		Description: stringPtr("description 100600 new"),
		StartAt:     1658584624,
		EndAt:       1658584625,
		NotifyAt:    int64Ptr(1658584626),
	}

	var body bytes.Buffer
	err = json.NewEncoder(&body).Encode(updateReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/events/%d", host, createRes.ID), &body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var updateResponse eventResponse
	err = json.NewDecoder(res.Body).Decode(&updateResponse)
	require.NoError(t, err)

	expected := eventResponse{
		ID:          updateResponse.ID,
		UserID:      100600,
		Title:       "title 100600 new",
		Description: stringPtr("description 100600 new"),
		StartAt:     time.Unix(1658584624, 0),
		EndAt:       time.Unix(1658584625, 0),
		NotifyAt:    timePtr(time.Unix(1658584626, 0)),
	}

	require.Equal(t, expected, updateResponse)
}

func TestDeleteEvent(t *testing.T) {
	defer truncateTable()

	createReq := eventCreateRequest{
		UserID:      100700,
		Title:       "title 100700",
		Description: stringPtr("description 100700"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode, createRes, err := createEvent(createReq)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/events/%d", host, createRes.ID), nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)
}

func TestListEvents(t *testing.T) {
	defer truncateTable()

	createReq1 := eventCreateRequest{
		UserID:      100800,
		Title:       "title 100800",
		Description: stringPtr("description 100800"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode1, createRes1, err := createEvent(createReq1)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode1)

	createReq2 := eventCreateRequest{
		UserID:      100900,
		Title:       "title 100900",
		Description: stringPtr("description 100900"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode2, createRes2, err := createEvent(createReq2)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode2)

	res, err := http.Get(host + "/events")
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var listResponse eventListResponse
	err = json.NewDecoder(res.Body).Decode(&listResponse)
	require.NoError(t, err)

	expected := eventListResponse{Events: []eventResponse{
		{
			ID:          createRes1.ID,
			UserID:      100800,
			Title:       "title 100800",
			Description: stringPtr("description 100800"),
			StartAt:     time.Unix(1658584623, 0),
			EndAt:       time.Unix(1658584624, 0),
			NotifyAt:    timePtr(time.Unix(1658584625, 0)),
		},
		{
			ID:          createRes2.ID,
			UserID:      100900,
			Title:       "title 100900",
			Description: stringPtr("description 100900"),
			StartAt:     time.Unix(1658584623, 0),
			EndAt:       time.Unix(1658584624, 0),
			NotifyAt:    timePtr(time.Unix(1658584625, 0)),
		},
	}}

	require.Equal(t, expected, listResponse)
}

func TestGetEvent(t *testing.T) {
	defer truncateTable()

	createReq := eventCreateRequest{
		UserID:      100100,
		Title:       "title 100100",
		Description: stringPtr("description 100100"),
		StartAt:     1658584623,
		EndAt:       1658584624,
		NotifyAt:    int64Ptr(1658584625),
	}

	statusCode, createRes, err := createEvent(createReq)
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	res, err := http.Get(fmt.Sprintf("%s/events/%d", host, createRes.ID))
	require.NoError(t, err)

	require.Equal(t, 200, res.StatusCode)

	var getResponse eventResponse
	err = json.NewDecoder(res.Body).Decode(&getResponse)
	require.NoError(t, err)

	expected := eventResponse{
		ID:          createRes.ID,
		UserID:      100100,
		Title:       "title 100100",
		Description: stringPtr("description 100100"),
		StartAt:     time.Unix(1658584623, 0),
		EndAt:       time.Unix(1658584624, 0),
		NotifyAt:    timePtr(time.Unix(1658584625, 0)),
	}

	require.Equal(t, expected, getResponse)
}

func createEvent(req eventCreateRequest) (int, eventResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return 0, eventResponse{}, err
	}

	res, err := http.Post(host+"/events", "application/json", &body)
	if err != nil {
		return 0, eventResponse{}, err
	}

	var response eventResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, eventResponse{}, err
	}

	return res.StatusCode, response, nil
}

func truncateTable() {
	if _, err := db.Exec("truncate table events;"); err != nil {
		panic(err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(ts int64) *int64 {
	return &ts
}

func timePtr(t time.Time) *time.Time {
	return &t
}
