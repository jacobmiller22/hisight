package cmd

import "fmt"

const hsServerUsage string = "Usage:\n\thsserver [http | grpc] [args]"

func HsServer(args []string) error {
	return fmt.Errorf(hsServerUsage)
}
