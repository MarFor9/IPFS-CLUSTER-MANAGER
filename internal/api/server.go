package api

import (
	"IPFS-CLUSTER-MANAGER/internal/core/ports"
	"IPFS-CLUSTER-MANAGER/internal/log"
	"bytes"
	"context"
	"errors"
)

// Server implements StrictServerInterface and holds the implementation of all API controllers
// This is the glue to the API autogenerated code
type Server struct {
	ipfsService ports.IpfsService
}

func NewServer(ipfsService ports.IpfsService) *Server {
	return &Server{
		ipfsService: ipfsService,
	}
}

func (s *Server) Check(_ context.Context, _ CheckRequestObject) (CheckResponseObject, error) {
	return Check200Response{}, nil
}

func (s *Server) GetPinnedFiles(ctx context.Context, _ GetPinnedFilesRequestObject) (GetPinnedFilesResponseObject, error) {
	log.Info(ctx, "[GetPinnedFiles] <- Enter")
	pins, err := s.ipfsService.GetPins(ctx)
	if err != nil {
		if err.Error() == "no content" {
			return GetPinnedFiles200JSONResponse([]Pin{}), nil
		}
		return GetPinnedFiles500JSONResponse{N500JSONResponse{Message: err.Error()}}, nil
	}

	log.Info(ctx, "[GetPinnedFiles] <- Leave")
	return GetPinnedFiles200JSONResponse(toApiPins(pins)), nil
}

func (s *Server) AddClusterNodePair(ctx context.Context, request AddClusterNodePairRequestObject) (AddClusterNodePairResponseObject, error) {
	log.Info(ctx, "[AddClusterNodePair] <- Enter")
	if request.Body.NodeUrl == "" {
		return AddClusterNodePair400JSONResponse{N400JSONResponse{Message: errors.New("node url is empty").Error()}}, nil
	}
	if request.Body.ClusterUrl == "" {
		return AddClusterNodePair400JSONResponse{N400JSONResponse{Message: errors.New("cluster url is empty").Error()}}, nil
	}

	err := s.ipfsService.AddClusterNodePair(ctx, request.Body.NodeUrl, request.Body.ClusterUrl)
	if err != nil {
		return AddClusterNodePair500JSONResponse{N500JSONResponse{Message: err.Error()}}, nil
	}

	log.Info(ctx, "[AddClusterNodePair] <- Leave")
	return AddClusterNodePair200Response{}, nil
}

func (s *Server) GetFile(ctx context.Context, request GetFileRequestObject) (GetFileResponseObject, error) {
	log.Info(ctx, "[GetFile] <- Enter")
	if request.Params.Arg == "" {
		return nil, errors.New("cid is empty")
	}

	fileReader, err := s.ipfsService.GetFile(ctx, request.Params.Arg)
	if err != nil {
		return GetFile500JSONResponse{N500JSONResponse{Message: err.Error()}}, nil
	}

	log.Info(ctx, "[GetFile] <- Leave")
	return GetFile200ApplicationoctetStreamResponse{
		Body:          bytes.NewReader(fileReader),
		ContentLength: int64(len(fileReader)),
	}, nil
}

func (s *Server) AddFile(ctx context.Context, request AddFileRequestObject) (AddFileResponseObject, error) {
	log.Info(ctx, "[AddFile] <- Enter")

	cid, err := s.ipfsService.AddFile(ctx, request.Body)
	if err != nil {
		return AddFile500JSONResponse{N500JSONResponse{Message: err.Error()}}, nil
	}

	log.Info(ctx, "[AddFile] <- Leave")
	return AddFile200JSONResponse{Hash: cid}, nil
}

func (s *Server) GetStatus(ctx context.Context, _ GetStatusRequestObject) (GetStatusResponseObject, error) {
	log.Info(ctx, "[GetStatus] <- Enter")

	response := s.ipfsService.GetStatus(ctx)

	log.Info(ctx, "[GetStatus] <- Leave")
	return GetStatus200JSONResponse(toApiIpfsHealthCheckResponse(response)), nil
}
