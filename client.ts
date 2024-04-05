import * as grpc from "npm:@grpc/grpc-js@1.10.6";
import * as resource from "npm:@pulumi/pulumi@3.112.0/proto/resource_pb.js";
import * as resourceRpc from "npm:@pulumi/pulumi@3.112.0/proto/resource_grpc_pb.js";

const monitor = new resourceRpc.ResourceMonitorClient(
    "localhost:8888",
    grpc.ChannelCredentials.createInsecure(),
);

const request = new resource.default.SupportsFeatureRequest();
request.setId("secrets");

monitor.supportsFeature(
    request,
    (err, response) => {
        console.log({ err })
        console.log(`has support:`, response?.getHassupport())
    },
)
