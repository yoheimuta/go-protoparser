package parser_test

import (
	"strings"

	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type protoTestVisitor struct {
	buffers []string
}

func (p *protoTestVisitor) String() string {
	return strings.Join(p.buffers, "\n")
}

func (p *protoTestVisitor) VisitComment(c *parser.Comment) {
	p.buffers = append(p.buffers, "Comment: "+c.Raw)
}

func (p *protoTestVisitor) VisitEmptyStatement(*parser.EmptyStatement) bool {
	p.buffers = append(p.buffers, "EmptyStatement")
	return true
}

func (p *protoTestVisitor) VisitEnum(e *parser.Enum) bool {
	p.buffers = append(p.buffers, "Enum: "+e.EnumName)
	return true
}

func (p *protoTestVisitor) VisitEnumField(f *parser.EnumField) bool {
	p.buffers = append(p.buffers, "EnumField: "+f.Ident)
	return true
}

func (p *protoTestVisitor) VisitExtend(e *parser.Extend) bool {
	p.buffers = append(p.buffers, "Extend: "+e.MessageType)
	return true
}

func (p *protoTestVisitor) VisitExtensions(r *parser.Extensions) bool {
	p.buffers = append(p.buffers, "Extensions")
	return true
}

func (p *protoTestVisitor) VisitField(f *parser.Field) bool {
	p.buffers = append(p.buffers, "Field: "+f.FieldName)
	return true
}

func (p *protoTestVisitor) VisitGroupField(f *parser.GroupField) bool {
	p.buffers = append(p.buffers, "GroupField: "+f.GroupName)
	return true
}

func (p *protoTestVisitor) VisitImport(i *parser.Import) bool {
	p.buffers = append(p.buffers, "Import: "+i.Location)
	return true
}

func (p *protoTestVisitor) VisitMapField(m *parser.MapField) bool {
	p.buffers = append(p.buffers, "MapField: "+m.MapName)
	return true
}

func (p *protoTestVisitor) VisitMessage(m *parser.Message) bool {
	p.buffers = append(p.buffers, "Message: "+m.MessageName)
	return true
}

func (p *protoTestVisitor) VisitOneof(o *parser.Oneof) bool {
	p.buffers = append(p.buffers, "Oneof: "+o.OneofName)
	return true
}

func (p *protoTestVisitor) VisitOneofField(f *parser.OneofField) bool {
	p.buffers = append(p.buffers, "OneofField: "+f.FieldName)
	return true
}

func (p *protoTestVisitor) VisitOption(o *parser.Option) bool {
	p.buffers = append(p.buffers, "Option: "+o.OptionName)
	return true
}

func (p *protoTestVisitor) VisitPackage(pa *parser.Package) bool {
	p.buffers = append(p.buffers, "Package: "+pa.Name)
	return true
}

func (p *protoTestVisitor) VisitReserved(r *parser.Reserved) bool {
	p.buffers = append(p.buffers, "Reserved")
	return true
}

func (p *protoTestVisitor) VisitRPC(r *parser.RPC) bool {
	p.buffers = append(p.buffers, "RPC: "+r.RPCName)
	return true
}

func (p *protoTestVisitor) VisitService(s *parser.Service) bool {
	p.buffers = append(p.buffers, "Service: "+s.ServiceName)
	return true
}

func (p *protoTestVisitor) VisitSyntax(s *parser.Syntax) bool {
	p.buffers = append(p.buffers, "Syntax: "+s.ProtobufVersion)
	return true
}

func TestProto_Accept(t *testing.T) {
	tests := []struct {
		name       string
		inputProto *parser.Proto
		wantBuffer string
	}{
		{
			name:       "parsing an empty",
			inputProto: &parser.Proto{},
		},
		{
			name: "parsing a syntax",
			inputProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "3",
				},
			},
			wantBuffer: `Syntax: 3`,
		},
		{
			name: "parsing an enum",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "TestEnum",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "TestEnumField1",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "TestEnumField2",
								Number: "0",
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: "EnumComment1",
							},
							{
								Raw: "EnumComment2",
							},
						},
						InlineComment: &parser.Comment{
							Raw: "EnumInlineComment",
						},
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "EnumInlineCommentBehindLeftCurly",
						},
					},
				},
			},
			wantBuffer: `Enum: TestEnum
EnumField: TestEnumField1
EnumField: TestEnumField2
Comment: EnumComment1
Comment: EnumComment2
Comment: EnumInlineComment
Comment: EnumInlineCommentBehindLeftCurly`,
		},
		{
			name: "parsing a service",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "TestService",
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "TestRPC1",
							},
							&parser.RPC{
								RPCName: "TestRPC2",
							},
						},
					},
				},
			},
			wantBuffer: `Service: TestService
RPC: TestRPC1
RPC: TestRPC2`,
		},
		{
			name: "parsing a message",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageName: "TestMessage",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "TestMessageInner",
							},
						},
					},
				},
			},
			wantBuffer: `Message: TestMessage
Message: TestMessageInner`,
		},
		{
			name: "parsing an import and a package ",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Import{
						Location: "Test_other.proto",
					},
					&parser.Package{
						Name: "TestMainPb",
					},
				},
			},
			wantBuffer: `Import: Test_other.proto
Package: TestMainPb`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			visitor := &protoTestVisitor{}
			test.inputProto.Accept(visitor)

			got := visitor.String()
			if got != test.wantBuffer {
				t.Errorf("got %s, but want %s", got, test.wantBuffer)
			}
		})
	}

}
