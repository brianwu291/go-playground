package main

import (
	"net/http"

	pb "github.com/brianwu291/go-playground/pb"
	utils "github.com/brianwu291/go-playground/utils"
)

// type greaterServer struct {
// 	pb.UnimplementedGreeterServer
// }

// func (s *greaterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
// 	log.Printf("Received: %v", req.GetName())
// 	return &pb.HelloResponse{Message: fmt.Sprintf("Hello, %s!", req.GetName())}, nil
// }

// func (s *greaterServer) SayHelloStream(req *pb.HelloRequest, stream pb.Greeter_SayHelloStreamServer) error {
// 	for i := 0; i < 5; i++ {
// 		message := fmt.Sprintf("Hello, %s! Message %d", req.GetName(), i+1)
// 		if err := stream.Send(&pb.HelloResponse{Message: message}); err != nil {
// 			return err
// 		}
// 		time.Sleep(time.Second)
// 	}
// 	return nil
// }

// func runServer() error {
// 	// Create listener
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		return fmt.Errorf("failed to listen: %v", err)
// 	}

// 	// Create gRPC greaterServer
// 	s := grpc.NewServer()
// 	pb.RegisterGreeterServer(s, &greaterServer{})

// 	log.Printf("Server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		return fmt.Errorf("failed to serve: %v", err)
// 	}
// 	return nil
// }

// func runClient() error {
// 	// set up connection to server
// 	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return fmt.Errorf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	// create client
// 	c := pb.NewGreeterClient(conn)

// 	// contact server
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()

// 	// call unary RPC
// 	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
// 	if err != nil {
// 		return fmt.Errorf("could not greet: %v", err)
// 	}
// 	log.Printf("Unary Response: %s", r.GetMessage())

// 	// call streaming RPC
// 	stream, err := c.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "streaming world"})
// 	if err != nil {
// 		return fmt.Errorf("could not greet: %v", err)
// 	}

// 	for {
// 		resp, err := stream.Recv()
// 		if err != nil {
// 			break
// 		}
// 		log.Printf("Stream Response: %s", resp.GetMessage())
// 	}

// 	return nil
// }

// type ChatTemplateEnvs struct {
// 	WebSocketUrl string
// }

// func serveChat(w http.ResponseWriter, r *http.Request) {
// 	data := ChatTemplateEnvs{
// 		WebSocketUrl: utils.GetEnv("WEBSOCKETURL", "ws://localhost:8080"),
// 	}

// 	tmpl, err := template.ParseFiles("templates/chat.html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	tmpl.Execute(w, data)
// }

// func manageRoomLifecycle(rt *realtimechat.RealTimeChat) {
// 	for {
// 		// sleep for 1.5 hours between cleanup cycles
// 		time.Sleep(90 * time.Minute)
// 		rooms := rt.ListRooms()
// 		for _, roomName := range rooms {
// 			if room, err := rt.GetRoom(roomName); err == nil {
// 				room.Stop()
// 			}
// 		}
// 	}
// }

// func main() {
// 	err := dotEnv.Load()
// 	if err != nil {
// 		fmt.Printf("error loading .env file: %+v", err.Error())
// 		return
// 	}
// 	// init without max clients as it's per room now
// 	chat := realtimechat.NewRealTimeChat()

// 	// start room lifecycle management
// 	go manageRoomLifecycle(chat)

// 	wsh := websockethandler.NewWebSocketHandler(chat)

// 	http.HandleFunc("/ws", wsh.HandleRealTimeChat)
// 	http.HandleFunc("/", serveChat)

// 	portStr := utils.GetEnv("PORT", "8080")
// 	fmt.Printf("listening port %+v\n", portStr)

// 	http.ListenAndServe(":"+portStr, nil)
// }

func main() {
	pb.Demo()

	portStr := utils.GetEnv("PORT", "8080")
	http.ListenAndServe(":"+portStr, nil)
}
