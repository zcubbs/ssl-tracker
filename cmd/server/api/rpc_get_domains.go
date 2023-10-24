package api

import (
	"context"
	pb "github.com/zcubbs/ssl-tracker/pb"
)

func (s *Server) GetDomains(ctx context.Context, req *pb.Empty) (*pb.GetDomainsResponse, error) {
	domains, err := s.store.GetAllDomains(ctx)
	if err != nil {
		return nil, err
	}

	responseDomains := make([]*pb.Domain, len(domains))
	for i, domain := range domains {
		responseDomains[i] = convertDomain(domain)
	}

	return &pb.GetDomainsResponse{
		Domains: responseDomains,
	}, nil
}
