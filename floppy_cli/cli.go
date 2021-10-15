package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc/resolver"

	pb "github.com/brotherlogic/floppy/proto"
)

func init() {
	resolver.Register(&utils.DiscoveryClientResolverBuilder{})
}

func main() {
	ctx, cancel := utils.ManualContext("floppy-cli", time.Second*10)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "floppy")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewFloppyServiceClient(conn)

	switch os.Args[1] {
	case "set":
		setFlags := flag.NewFlagSet("SetConfig", flag.ContinueOnError)
		var instanceId = setFlags.Int("id", -1, "Id of the record to add")

		if err := setFlags.Parse(os.Args[2:]); err == nil {
			if *instanceId > 0 {
				_, err := client.Register(ctx, &pb.RegisterRequest{InstanceId: *instanceId})
				if err != nil {
					log.Fatalf("Bad request: %v", err)
				}
			}
		}
	default:
		log.Fatalf("Unknown command: %v", os.Args[1])
	}
}
