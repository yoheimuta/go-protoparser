package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/parser"
)

// Service consists of RPCs.
type Service struct {
	ServiceName string

	// ServiceBody is unordered in nature, but each slice field preserves the original order.
	Options []*parser.Option
	RPCs    []*parser.RPC

	// Comments are the optional ones placed at the beginning.
	Comments []*parser.Comment
}

// InterpretService interprets *parser.Service to *Service.
func InterpretService(src *parser.Service) (*Service, error) {
	if src == nil {
		return nil, nil
	}

	options, rpcs, err := interpretServiceBody(src.ServiceBody)
	if err != nil {
		return nil, err
	}
	return &Service{
		ServiceName: src.ServiceName,
		Options:     options,
		RPCs:        rpcs,
		Comments:    src.Comments,
	}, nil
}

func interpretServiceBody(src []interface{}) (
	options []*parser.Option,
	rpcs []*parser.RPC,
	err error,
) {
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Option:
			options = append(options, t)
		case *parser.RPC:
			rpcs = append(rpcs, t)
		default:
			err = fmt.Errorf("invalid ServiceBody type %v of %v", t, s)
		}
	}
	return
}
