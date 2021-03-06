package helix

import (
	gotau "github.com/Team-TAU/tau-client-go"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient_DeleteRequestReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteRequestReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.Error(t, err)
	require.IsType(t, gotau.AuthorizationError{}, err)
	require.False(t, deleted)
}

func TestClient_DeleteRequestReturnsRateLimitError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)

		w.Header().Set("Ratelimit-Reset", "1623961625")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.Error(t, err)
	require.IsType(t, RateLimitError{}, err)
	rlErr := err.(RateLimitError)
	require.NotNil(t, rlErr.ResetTime())
	require.Equal(t, 2021, rlErr.ResetTime().Year())
	require.Equal(t, time.Month(6), rlErr.ResetTime().Month())
	require.Equal(t, 17, rlErr.ResetTime().Day())
	require.False(t, deleted)
}

func TestClient_DeleteRequestReturnsGenericError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteRequest("channel_points/custom_rewards", nil)
	require.Error(t, err)
	genericError, ok := err.(gotau.GenericError)
	require.True(t, ok)
	require.False(t, deleted)

	require.Equal(t, 404, genericError.Code)
	require.Empty(t, genericError.Body)
}

func TestClient_DeleteCustomRewardReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/channel_points/custom_rewards/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "456", r.URL.Query().Get("id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteCustomReward("123", "456")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteCustomRewardReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteCustomReward("", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcaster can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteCustomReward("123", "")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)
}

func TestClient_DeleteEventSubSubscriptionReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/eventsub/subscriptions/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteEventSubSubscription("123")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteEventSubSubscriptionReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteEventSubSubscription("")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteEventSubSubscription("    ")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteEventSubSubscription("		")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)
}

func TestClient_DeleteUserFollowsReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/follows/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("from_id"))
		require.Equal(t, "456", r.URL.Query().Get("to_id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteUserFollows("123", "456")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteUserFollowsReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteUserFollows("", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, fromID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteUserFollows("    ", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, fromID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteUserFollows("		", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, fromID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteUserFollows("123", "")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, toID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteUserFollows("123", "    ")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, toID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteUserFollows("123", "	")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, toID can't be blank"})
	require.False(t, deleted)
}

func TestClient_UnblockUserReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/users/blocks/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("user_id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.UnblockUser("123")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_UnblockUserReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.UnblockUser("")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, userID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.UnblockUser("    ")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, userID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.UnblockUser("		")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, userID can't be blank"})
	require.False(t, deleted)
}

func TestClient_DeleteVideosReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/videos/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "123", r.URL.Query().Get("id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteVideos([]string{"123"})
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteVideosReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteVideos(nil)
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, IDs can't be empty"})
	require.False(t, deleted)

	deleted, err = client.DeleteVideos(make([]string, 0))
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, IDs can't be empty"})
	require.False(t, deleted)

	deleted, err = client.DeleteVideos(make([]string, 6))
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, maximum number of IDs is 5 but you supplied 6"})
	require.False(t, deleted)
}

func TestClient_DeleteChannelStreamScheduleSegmentReturnsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/twitch/helix/schedule/segment/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "141981764", r.URL.Query().Get("broadcaster_id"))
		require.Equal(t, "eyJzZWdtZW50SUQiOiI4Y2EwN2E2NC0xYTZkLTRjYWItYWE5Ni0xNjIyYzNjYWUzZDkiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyMX0=", r.URL.Query().Get("id"))
		require.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client, err := NewClient(host, portNum, "foo", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	deleted, err := client.DeleteChannelStreamScheduleSegment("141981764", "eyJzZWdtZW50SUQiOiI4Y2EwN2E2NC0xYTZkLTRjYWItYWE5Ni0xNjIyYzNjYWUzZDkiLCJpc29ZZWFyIjoyMDIxLCJpc29XZWVrIjoyMX0=")
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestClient_DeleteChannelStreamScheduleSegmentReturnsError(t *testing.T) {
	client := Client{}

	deleted, err := client.DeleteChannelStreamScheduleSegment("", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcaster can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteChannelStreamScheduleSegment("    ", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcaster can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteChannelStreamScheduleSegment("		", "123")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, broadcaster can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteChannelStreamScheduleSegment("123", "")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteChannelStreamScheduleSegment("123", "    ")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)

	deleted, err = client.DeleteChannelStreamScheduleSegment("123", "	")
	require.ErrorIs(t, err, gotau.BadRequestError{Err: "invalid request, ID can't be blank"})
	require.False(t, deleted)
}
