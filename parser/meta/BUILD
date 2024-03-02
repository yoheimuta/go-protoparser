load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "meta",
    srcs = [
        "error.go",
        "meta.go",
        "position.go",
    ],
    importpath = "github.com/yoheimuta/go-protoparser/v4/parser/meta",
    visibility = ["//visibility:public"],
)

go_test(
    name = "meta_test",
    srcs = ["position_test.go"],
    deps = [":meta"],
)
