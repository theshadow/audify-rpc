package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	nurl "net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

var defaultDuration = "1800"

type Client struct {
	url *nurl.URL
	httpClient *http.Client
	doer Doer
	l *log.Logger
}

func New(url string, l *log.Logger) (*Client, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		u,
		&http.Client{},
		ctxhttp.Do,
		l,
    }, nil
}

func NewWithDoer(url string, l *log.Logger, doer Doer) (*Client, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		u,
		&http.Client{},
		doer,
		l,
	}, nil
}

func (c *Client) Search(ctx context.Context, req Request) ([]Item, error) {
	var url nurl.URL
	url = *c.url

	q := url.Query()

	// Don't know what this does but the web client uses it.
	q.Add("duration", defaultDuration)

	if len(req.Tags) > 0 {
		// add any tags
		log.Debugf("Tags: %d", len(req.Tags))
		q.Add("tag", strings.Join(req.Tags, ","))
	}

	if len(req.Source) > 0 {
		// define the source
		q.Add("source", req.Source)
	}

	url.RawQuery = q.Encode()

	r, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}
	resp, err := c.doer(ctx, c.httpClient, r)
	if err != nil {
		return nil, err
	}

	// make sure that we clean up resources
	defer resp.Body.Close()

	apiResp := &Response{}
	if err := apiResp.FromJson(resp.Body); err != nil {
		return nil, err
	}

	return apiResp.Items, nil
}

// News item
type Item struct {
	Title           string  `protobuf:"bytes,1,opt,name=Title" json:"title,omitempty"`
	Summary         string  `protobuf:"bytes,2,opt,name=Summary" json:"summary,omitempty"`
	DateURL         string  `protobuf:"bytes,3,opt,name=DateURL" json:"date_url,omitempty"`
	AudioURL        string  `protobuf:"bytes,4,opt,name=AudioURL" json:"audio_url,omitempty"`
	ImageURL        string  `protobuf:"bytes,5,opt,name=ImageURL" json:"image_url,omitempty"`
	ArticleURL      string  `protobuf:"bytes,6,opt,name=ArticleURL" json:"article_url,omitempty"`
	Duration        float32 `protobuf:"fixed32,7,opt,name=Duration" json:"duration,omitempty"`
	FileSizeInBytes uint64  `protobuf:"varint,8,opt,name=FileSizeInBytes" json:"filesize_in_bytes,omitempty"`
	NumPlays        uint32  `protobuf:"varint,9,opt,name=NumPlays" json:"num_plays,omitempty"`
	SourceID        string  `protobuf:"bytes,10,opt,name=SourceID" json:"source_id,omitempty"`
	GUID            string  `protobuf:"bytes,11,opt,name=GUID" json:"guid,omitempty"`
	PublishedAt     string  `protobuf:"bytes,12,opt,name=PublishedAt" json:"published_at,omitempty"`
}

// API Response
type Response struct {
	Status uint16 `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Items []Item `json:"items,omitempty"`
	Identifiers map[string]string `json:"identifiers,omitempty"`
}

func (r *Response) FromJson(rdr io.Reader) error {
	if err := json.NewDecoder(rdr).Decode(r); err != nil {
		return err
	}
	return nil
}

type Request struct {
	Tags []string
	Source string
}
