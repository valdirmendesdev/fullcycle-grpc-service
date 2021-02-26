package main

import (
	"context"
	"fmt"
	"github.com/valdirmendesdev/fc2-grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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

func AddUser(client pb.UserServiceClient)  {
	req := &pb.User{
		Id: "",
		Name: "Valdir",
		Email: "j@j.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient)  {
	req := &pb.User{
		Id: "",
		Name: "Valdir",
		Email: "j@j.com",
	}
	resStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := resStream.Recv()
		if err == io.EOF {
			break;
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient)  {
	reqs := []*pb.User{
		&pb.User{
			Id: "1",
			Name: "Valdir1",
			Email: "j1@j.com",
		},
		&pb.User{
			Id: "2",
			Name: "Valdir2",
			Email: "j2@j.com",
		},
		&pb.User{
			Id: "3",
			Name: "Valdir3",
			Email: "3j@j.com",
		},
		&pb.User{
			Id: "4",
			Name: "Valdir4",
			Email: "4j@j.com",
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

	req, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Println(req)
}


func AddUserStreamBoth(client pb.UserServiceClient)  {
	reqs := []*pb.User{
		&pb.User{
			Id: "1",
			Name: "Valdir1",
			Email: "j1@j.com",
		},
		&pb.User{
			Id: "2",
			Name: "Valdir2",
			Email: "j2@j.com",
		},
		&pb.User{
			Id: "3",
			Name: "Valdir3",
			Email: "3j@j.com",
		},
		&pb.User{
			Id: "4",
			Name: "Valdir4",
			Email: "4j@j.com",
		},
	}

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user:", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 3)
		}
		stream.CloseSend()
	}()

	wait := make(chan int)

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Println("Status", res.Status, " - ", res.GetUser())
		}
		close(wait)
	}()

	<-wait
}