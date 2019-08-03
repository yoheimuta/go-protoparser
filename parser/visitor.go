package parser

// Visitor is for dispatching Protocol Buffer elements.
type Visitor interface {
	VisitComment(*Comment)
	VisitEmptyStatement(*EmptyStatement) (next bool)
	VisitEnum(*Enum) (next bool)
	VisitEnumField(*EnumField) (next bool)
	VisitExtend(*Extend) (next bool)
	VisitExtensions(*Extensions) (next bool)
	VisitField(*Field) (next bool)
	VisitGroupField(*GroupField) (next bool)
	VisitImport(*Import) (next bool)
	VisitMapField(*MapField) (next bool)
	VisitMessage(*Message) (next bool)
	VisitOneof(*Oneof) (next bool)
	VisitOneofField(*OneofField) (next bool)
	VisitOption(*Option) (next bool)
	VisitPackage(*Package) (next bool)
	VisitReserved(*Reserved) (next bool)
	VisitRPC(*RPC) (next bool)
	VisitService(*Service) (next bool)
	VisitSyntax(*Syntax) (next bool)
}
