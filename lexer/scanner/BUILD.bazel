load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "scanner",
    srcs = [
        "boolLit.go",
        "comment.go",
        "eof.go",
        "error.go",
        "floatLit.go",
        "ident.go",
        "lettersdigits.go",
        "mode.go",
        "numberLit.go",
        "position.go",
        "quote.go",
        "scanner.go",
        "strLit.go",
        "token.go",
    ],
    importpath = "github.com/yoheimuta/go-protoparser/v4/lexer/scanner",
    visibility = ["//visibility:public"],
    deps = ["//parser/meta"],
)

go_test(
    name = "scanner_test",
    srcs = [
        "position_test.go",
        "scanner_test.go",
    ],
    deps = [
        ":scanner",
        "//parser/meta",
    ],
)
