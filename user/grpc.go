package user

import (
	"context"
	"log"
	"time"

	android "github.com/u03013112/ss-pb/android"
	"google.golang.org/grpc"
)

const (
	androidAddress = "android:50003"
)

func grpcGetAndroidUserStatus(token string) (string, error) {
	conn, err := grpc.Dial(androidAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := android.NewAndroidClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetUserInfo(ctx, &android.GetUserInfoRequest{
		Token: token,
	})
	if err != nil {
		log.Printf("could not GetSSConfig: %v", err)
		return "unknow", err
	}
	return r.Status, nil
}
