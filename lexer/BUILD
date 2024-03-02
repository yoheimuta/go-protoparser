load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "lexer",
    srcs = [
        "constant.go",
        "emptyStatement.go",
        "enumType.go",
        "error.go",
        "fullIdent.go",
        "lexer.go",
        "messageType.go",
    ],
    importpath = "github.com/yoheimuta/go-protoparser/v4/lexer",
    visibility = ["//visibility:public"],
    deps = [
        "//lexer/scanner",
        "//parser/meta",
    ],
)

go_test(
    name = "lexer_test",
    srcs = [
        "constant_test.go",
        "emptyStatement_test.go",
        "enumType_test.go",
        "fullIdent_test.go",
        "messageType_test.go",
    ],
    deps = [":lexer"],
)
