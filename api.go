package gotau

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// GetStreamers can be used to get a list of all the streamers that TAU is listening for going live alerts.
func (c *Client) GetStreamers() ([]*TAUStreamer, error) {
	body, err := c.apiRequest("streamers", nil, nil, "GET")
	if err != nil {
		return nil, err
	}

	var streamers []*TAUStreamer
	err = json.Unmarshal(body, &streamers)
	if err != nil {
		return nil, err
	}

	return streamers, nil
}

// GetLatestStreamForStreamer gets the latest stream for a given streamer
func (c *Client) GetLatestStreamForStreamer(ID string) (*TAUStream, error) {
	ID = strings.TrimSpace(ID)
	if ID == "" {
		return nil, BadRequestError{
			"invalid request, ID can't be blank",
		}
	}

	body, err := c.apiRequest(fmt.Sprintf("streamers/%s/streams/latest", ID), nil, nil, "GET")
	if err != nil {
		return nil, err
	}

	stream := new(TAUStream)
	err = json.Unmarshal(body, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// FollowStreamerOnTau follows the users and subscribes for notifications when they go live
func (c *Client) FollowStreamerOnTau(username string) (*TAUStreamer, error) {
	type tmp struct {
		Username  string `json:"twitch_username"`
		Streaming bool   `json:"streaming"`
		Disabled  bool   `json:"disabled"`
	}

	username = strings.TrimSpace(username)
	if username == "" {
		return nil, BadRequestError{
			Err: "invalid request, username can't be blank",
		}
	}

	data := tmp{
		Username: username,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.apiRequest("streamers", nil, body, "POST")
	if err != nil {
		return nil, err
	}
	streamer := new(TAUStreamer)
	err = json.Unmarshal(responseBody, streamer)
	if err != nil {
		return nil, err
	}

	return streamer, nil
}

// GetStreamsForStreamer will get n streams for a streamer.  If maximumStreams is set to -1 then all
// streams will be gathered.  This may take some time due to pagination.  In the case where there are fewer
// results than the maximumStreams, those results will be returned.  The number of results you get may be slightly
// more than maximumResults, based on the pagination of the results.  If you request 0 results, you will get 0 results.
func (c *Client) GetStreamsForStreamer(streamerID string, maximumStreams int) ([]TAUStream, error) {
	type tmp struct {
		Streams  []TAUStream `json:"results"`
		Previous *string     `json:"previous"`
		Next     *string     `json:"next"`
		Count    int         `json:"count"`
	}
	streamerID = strings.TrimSpace(streamerID)
	if streamerID == "" {
		return nil, BadRequestError{
			Err: "invalid request, streamer id can't be blank",
		}
	}

	results := make([]TAUStream, 0)
	url := fmt.Sprintf("streamers/%s/streams", streamerID)
	getAll := maximumStreams < 0
	// if there are more results than that, we've fucked up
	if getAll {
		maximumStreams = math.MaxInt64
	}

	count := 0

	for count <= maximumStreams {
		body, err := c.apiRequest(url, nil, nil, "GET")
		if err != nil {
			return nil, err
		}

		tmpData := new(tmp)
		err = json.Unmarshal(body, tmpData)
		if err != nil {
			return nil, err
		}
		results = append(results, tmpData.Streams...)

		// no data left so abort early whether or not we've got enough data
		if tmpData.Next == nil {
			break
		} else {
			url = *tmpData.Next
		}

		count += len(tmpData.Streams)
	}

	return results, nil
}
