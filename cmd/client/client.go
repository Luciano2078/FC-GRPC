package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Luciano2078/FC-GRCP/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Dudu",
		Email: "dudu@dudunet.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Dudu",
		Email: "dudu@dudunet.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "d1",
			Name:  "Dudu",
			Email: "dudu@dudunet.au.au",
		},
		&pb.User{
			Id:    "d2",
			Name:  "Duda",
			Email: "duda@dudunet.au.au",
		},
		&pb.User{
			Id:    "c3",
			Name:  "Chico",
			Email: "chico@dudunet.au.au",
		},
		&pb.User{
			Id:    "d4",
			Name:  "Duque",
			Email: "duque@duquenet.miau",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "d1",
			Name:  "Dudu",
			Email: "dudu@dudunet.au.au",
		},
		&pb.User{
			Id:    "d2",
			Name:  "Duda",
			Email: "duda@dudunet.au.au",
		},
		&pb.User{
			Id:    "c3",
			Name:  "Chico",
			Email: "chico@dudunet.au.au",
		},
		&pb.User{
			Id:    "d4",
			Name:  "Duque",
			Email: "duque@duquenet.miau",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data %v", err)
				break
			}
			fmt.Printf("recebendo user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
