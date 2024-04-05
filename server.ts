import * as grpc from "npm:@grpc/grpc-js@1.10.6";
import resource from "npm:@pulumi/pulumi@3.112.0/proto/resource_pb.js";
import resourceRpc from "npm:@pulumi/pulumi@3.112.0/proto/resource_grpc_pb.js";

const service = {
    supportsFeature: (call: any, callback: any) => {
        console.log({ supportsFeature: call.request })
        call.on('error', (args: Error) => {
            console.log("supportsFeature() got error:", args)
        })

        const response = new resource.SupportsFeatureResponse()
        response.setHassupport(true)

        callback(null, response)
    }
}

const server = new grpc.Server();
server.addService(resourceRpc.ResourceMonitorService, service)

server.bindAsync("0.0.0.0:8888", grpc.ServerCredentials.createInsecure(),
    (err: Error | null, port: number) => {
        if (err) {
            console.error(`Server error: ${err.message}`);
        } else {
            console.log(`Server bound on port: ${port}`);
        }
    });
