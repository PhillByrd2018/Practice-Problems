package main

import (
	"fmt"

	"Practice-Problems/services"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fmt.Println("Starting Lambda")
	lambda.Start(services.RunThrough)

}
