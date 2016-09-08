// Program googlelog sends logs to Google Cloud Logging
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/genproto/googleapis/api/monitoredres"
	logging "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	flagProjectID = flag.String("projectID", "", "project id")
	flagLogName   = flag.String("log", "", "log name")
)

func main() {
	flag.Parse()
	if err := logentry(*flagProjectID, *flagLogName); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func logentry(projectID, logName string) error {
	ctx := context.Background()
	client, err := loggingClient(ctx)
	if err != nil {
		return err
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	_, err = client.WriteLogEntries(ctx, &logging.WriteLogEntriesRequest{
		Entries: []*logging.LogEntry{
			{
				Resource: &monitoredres.MonitoredResource{Type: "global"},
				LogName:  fmt.Sprintf("projects/%s/logs/%s", *flagProjectID, *flagLogName),
				Payload: &logging.LogEntry_TextPayload{
					TextPayload: string(input),
				},
				/*
					Payload: &logging.LogEntry_ProtoPayload{
						ProtoPayload: &any.Any{
							TypeUrl: "type.googleapis.com/google.cloud.audit.AuditLog",
							Value:   alBytes,
						},
					},
				*/
			},
		},
	})
	return err
}

func loggingClient(ctx context.Context) (logging.LoggingServiceV2Client, error) {
	creds, err := oauth.NewApplicationDefault(ctx, "https://www.googleapis.com/auth/logging.write")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial("logging.googleapis.com:443",
		grpc.WithPerRPCCredentials(creds),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		return nil, err
	}
	return logging.NewLoggingServiceV2Client(conn), nil
}
