load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "dump_lib",
    srcs = ["main.go"],
    importpath = "github.com/yoheimuta/go-protoparser/v4/_example/dump",
    visibility = ["//visibility:private"],
    deps = ["//:go-protoparser"],
)

go_binary(
    name = "dump",
    embed = [":dump_lib"],
    visibility = ["//visibility:public"],
)
