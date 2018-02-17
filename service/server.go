package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"github.com/theshadow/ushadow/audify/api"
)

type Version struct{
	Binary string
	Dependencies[] string
}

func (v Version) String() string {
	return fmt.Sprintf("binary: %s\n%s", v.Binary, strings.Join(v.Dependencies, "\n"))
}

type Server struct{
	version Version
	rpcSrv *grpc.Server
	api    *api.Client
	done   chan struct{}
}

func New(ver Version, rpc *grpc.Server, api *api.Client, done chan struct{}) *Server {
	return &Server{version: ver, rpcSrv: rpc, api: api, done: done}
}

func (s *Server) Search(req *SearchRequest, srv Audify_SearchServer) error {
	var tags []string
	for _, t := range req.Tags {
		tags = append(tags, t.Tag)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second * 6)
	apiReq := api.Request{
		Source: req.Source,
		Tags: tags,
	}

	items, err := s.api.Search(ctx, apiReq)
	if err != nil {
		return err
	}

	for _, item := range items {
		var resp SearchResponse
		Unmarshal(item, &resp)
		if err := srv.Send(&resp); err != nil {
			return err
		}
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context, in *ShutdownRequest) (*ShutdownResponse, error) {
	close(s.done)
	return &ShutdownResponse{}, nil
}

func (s *Server) Version(ctx context.Context, in *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{
		Version: s.version.Binary,
		Dependencies: s.version.Dependencies,
	}, nil
}

func Unmarshal(item api.Item, resp *SearchResponse) {
		resp.Title = item.Title
		resp.Summary = item.Summary
		resp.DateURL = item.DateURL
		resp.AudioURL = item.AudioURL
		resp.ImageURL = item.ImageURL
		resp.ArticleURL = item.ArticleURL
		resp.Duration = item.Duration
		resp.FileSizeInBytes = item.FileSizeInBytes
		resp.NumPlays = item.NumPlays
		resp.SourceID = item.SourceID
		resp.GUID = item.GUID
		resp.PublishedAt = item.PublishedAt
}



