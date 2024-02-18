load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "unordered",
    srcs = [
        "enum.go",
        "extend.go",
        "message.go",
        "proto.go",
        "service.go",
    ],
    importpath = "github.com/yoheimuta/go-protoparser/v4/interpret/unordered",
    visibility = ["//visibility:public"],
    deps = [
        "//parser",
        "//parser/meta",
    ],
)

go_test(
    name = "unordered_test",
    srcs = [
        "enum_test.go",
        "message_test.go",
        "proto_test.go",
        "service_test.go",
    ],
    deps = [
        ":unordered",
        "//parser",
        "//parser/meta",
    ],
)
