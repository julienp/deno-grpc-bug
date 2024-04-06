import * as grpc from "@grpc/grpc-js";
import * as resource from "@pulumi/pulumi/proto/resource_pb.js";
import * as resourceRpc from "@pulumi/pulumi/proto/resource_grpc_pb.js";

const monitor = new resourceRpc.ResourceMonitorClient("localhost:8888", grpc.ChannelCredentials.createInsecure());

const request = new resource.default.SupportsFeatureRequest();
request.setId("secrets");

monitor.supportsFeature(request, (err, response) => {
  console.log({ err });
  console.log(`has support:`, response?.getHassupport());
});
