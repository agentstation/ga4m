```sh
                                     __ _        _  _              
                                    / _` |  __ _| || |   _ __ ___  
                                   | (_| | / _` | || |_ | '_ ` _ \ 
                                    \__, || (_| |__   _|| | | | | |
                                    |___/  \__,_|  |_|  |_| |_| |_|
```

This package provides a client for sending events to Google Analytics 4 via the Measurement Protocol (GA4M).

## Usage Example

Here's a simple example of how to use ga4m to track page views in a web application:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/agentstation/ga4m"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the AnalyticsClient
	gaClient := ga4m.NewClient("YOUR_MEASUREMENT_ID", "YOUR_API_SECRET")

	// Parse the GA cookies to get clientID and sessionID
	session := ga4m.ParseSessionFromRequest(r)

	// Prepare event parameters
	params := map[string]interface{}{
		"page_title": "Homepage",
		"page_path":  "/",
	}

	// Send the event
	err := gaClient.SendEvent(session.ClientID, "page_view", params, ga4m.WithSessionID(session.LastSessionID()))
	if err != nil {
		// Handle error
		fmt.Printf("Error sending event: %v\n", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

This example demonstrates:
- Creating a new GA4 client with your measurement ID and API secret
- Parsing GA cookies from incoming HTTP requests
- Sending a page view event with custom parameters
- Using the session ID option for better session tracking

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# ga4m

```go
import "github.com/agentstation/ga4m"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [type AnalyticsClient](<#AnalyticsClient>)
  - [func NewClient\(measurementID, apiSecret string\) \*AnalyticsClient](<#NewClient>)
  - [func \(c \*AnalyticsClient\) SendEvent\(clientID, eventName string, params map\[string\]interface\{\}, opts ...SendEventOption\) error](<#AnalyticsClient.SendEvent>)
  - [func \(c \*AnalyticsClient\) SendEvents\(clientID string, events \[\]EventParams, opts ...SendEventOption\) error](<#AnalyticsClient.SendEvents>)
  - [func \(c \*AnalyticsClient\) SetHTTPClient\(client HTTPClient\)](<#AnalyticsClient.SetHTTPClient>)
- [type AnalyticsEvent](<#AnalyticsEvent>)
- [type EventParams](<#EventParams>)
- [type HTTPClient](<#HTTPClient>)
- [type SendEventOption](<#SendEventOption>)
  - [func WithContext\(ctx context.Context\) SendEventOption](<#WithContext>)
  - [func WithDebug\(debug bool\) SendEventOption](<#WithDebug>)
  - [func WithSessionID\(sessionID string\) SendEventOption](<#WithSessionID>)
  - [func WithTimestamp\(timestamp time.Time\) SendEventOption](<#WithTimestamp>)
  - [func WithUserID\(userID string\) SendEventOption](<#WithUserID>)
- [type Session](<#Session>)
  - [func LatestSessions\(sessions ...Session\) Session](<#LatestSessions>)
  - [func ParseSessionFromRequest\(r \*http.Request\) Session](<#ParseSessionFromRequest>)
  - [func \(s Session\) LastSessionID\(\) string](<#Session.LastSessionID>)


## Constants

<a name="DefaultEngagementTimeMS"></a>

```go
const (
    // DefaultEngagementTimeMS is the default engagement time in milliseconds
    DefaultEngagementTimeMS = "100"
    // SessionIDParam is the parameter name for the session ID
    SessionIDParam = "session_id"
    // EngagementTimeParam is the parameter name for the engagement time in milliseconds
    EngagementTimeParam = "engagement_time_msec"
    // MaxEventsPerRequest is the maximum number of events per request
    MaxEventsPerRequest = 25
    // URLFormat is the format for the URL
    URLFormat = "%s?measurement_id=%s&api_secret=%s"
    // ContentTypeHeader is the header for the content type
    ContentTypeHeader = "Content-Type"
    // ContentTypeJSON is the content type for JSON
    ContentTypeJSON = "application/json"
)
```

<a name="ClientCookieName"></a>

```go
const (
    ClientCookieName  = "_ga"
    SessionCookieName = "_ga_"
)
```

## Variables

<a name="EmptySession"></a>EmptySession is an empty Google Analytics session

```go
var EmptySession = Session{}
```

<a name="AnalyticsClient"></a>
## type [AnalyticsClient](<https://github.com/agentstation/ga4m/blob/master/client.go#L14-L20>)

AnalyticsClient is the client for sending events to Google Analytics

```go
type AnalyticsClient struct {
    MeasurementID string
    APISecret     string
    Endpoint      string
    DebugEndpoint string
    HTTPClient    HTTPClient
}
```

<a name="NewClient"></a>
### func [NewClient](<https://github.com/agentstation/ga4m/blob/master/client.go#L23>)

```go
func NewClient(measurementID, apiSecret string) *AnalyticsClient
```

NewClient creates a new AnalyticsClient with the provided measurement ID and API secret

<a name="AnalyticsClient.SendEvent"></a>
### func \(\*AnalyticsClient\) [SendEvent](<https://github.com/agentstation/ga4m/blob/master/send_event.go#L44>)

```go
func (c *AnalyticsClient) SendEvent(clientID, eventName string, params map[string]interface{}, opts ...SendEventOption) error
```

SendEvent sends a single event to Google Analytics.

<a name="AnalyticsClient.SendEvents"></a>
### func \(\*AnalyticsClient\) [SendEvents](<https://github.com/agentstation/ga4m/blob/master/send_event.go#L107>)

```go
func (c *AnalyticsClient) SendEvents(clientID string, events []EventParams, opts ...SendEventOption) error
```

SendEvents sends multiple events in a single batch request to Google Analytics.

<a name="AnalyticsClient.SetHTTPClient"></a>
### func \(\*AnalyticsClient\) [SetHTTPClient](<https://github.com/agentstation/ga4m/blob/master/client.go#L34>)

```go
func (c *AnalyticsClient) SetHTTPClient(client HTTPClient)
```

SetHTTPClient allows setting a custom HTTP client

<a name="AnalyticsEvent"></a>
## type [AnalyticsEvent](<https://github.com/agentstation/ga4m/blob/master/send_event.go#L36-L41>)

AnalyticsEvent represents the payload structure for GA4 events.

```go
type AnalyticsEvent struct {
    ClientID        string        `json:"client_id"`
    Events          []EventParams `json:"events"`
    UserID          string        `json:"user_id,omitempty"`
    TimestampMicros int64         `json:"timestamp_micros,omitempty"`
}
```

<a name="EventParams"></a>
## type [EventParams](<https://github.com/agentstation/ga4m/blob/master/send_event.go#L29-L33>)

EventParams represents parameters for a GA4 event.

```go
type EventParams struct {
    Name            string                 `json:"name"`
    Params          map[string]interface{} `json:"params,omitempty"`
    TimestampMicros int64                  `json:"timestamp_micros,omitempty"`
}
```

<a name="HTTPClient"></a>
## type [HTTPClient](<https://github.com/agentstation/ga4m/blob/master/client.go#L9-L11>)

HTTPClient interface allows for mocking of http.Client in tests

```go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}
```

<a name="SendEventOption"></a>
## type [SendEventOption](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L9>)

SendEventOption allows for optional parameters when sending events.

```go
type SendEventOption func(*sendEventOptions)
```

<a name="WithContext"></a>
### func [WithContext](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L27>)

```go
func WithContext(ctx context.Context) SendEventOption
```

WithContext sets a custom context for the request.

<a name="WithDebug"></a>
### func [WithDebug](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L34>)

```go
func WithDebug(debug bool) SendEventOption
```

WithDebug enables or disables debug mode.

<a name="WithSessionID"></a>
### func [WithSessionID](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L55>)

```go
func WithSessionID(sessionID string) SendEventOption
```

WithSessionID sets the session ID for the event.

<a name="WithTimestamp"></a>
### func [WithTimestamp](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L48>)

```go
func WithTimestamp(timestamp time.Time) SendEventOption
```

WithTimestamp sets a custom timestamp for the event.

<a name="WithUserID"></a>
### func [WithUserID](<https://github.com/agentstation/ga4m/blob/master/send_event_options.go#L41>)

```go
func WithUserID(userID string) SendEventOption
```

WithUserID sets the user ID for the event.

<a name="Session"></a>
## type [Session](<https://github.com/agentstation/ga4m/blob/master/session.go#L20-L25>)

Session represents the Google Analytics session tracking data for a user.

```go
type Session struct {
    ClientID     string    // The client ID from _ga cookie.
    FirstVisit   time.Time // First visit timestamp.
    SessionCount int       // Number of sessions.
    LastSession  time.Time // Last session timestamp.
}
```

<a name="LatestSessions"></a>
### func [LatestSessions](<https://github.com/agentstation/ga4m/blob/master/session.go#L78>)

```go
func LatestSessions(sessions ...Session) Session
```

LatestSessions compares Google Analytics sessions and returns the latest one

<a name="ParseSessionFromRequest"></a>
### func [ParseSessionFromRequest](<https://github.com/agentstation/ga4m/blob/master/session.go#L33>)

```go
func ParseSessionFromRequest(r *http.Request) Session
```

ParseSessionFromRequest parses the Google Analytics cookies from an HTTP request and returns a Session.

<a name="Session.LastSessionID"></a>
### func \(Session\) [LastSessionID](<https://github.com/agentstation/ga4m/blob/master/session.go#L28>)

```go
func (s Session) LastSessionID() string
```

LastSessionID returns the Unix timestamp of the last session as a string, this can be used as a session ID.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
