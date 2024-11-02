import grpc
import greeting_pb2
import greeting_pb2_grpc

def run():
    # Connect to the gRPC server
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = greeting_pb2_grpc.GreeterStub(channel)
        
        # Call the SayHello method on the server
        response = stub.SayHello(greeting_pb2.HelloRequest(name='World'))
        print(f"Greeting: {response.message}")

if __name__ == '__main__':
    run()