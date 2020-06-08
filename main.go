package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/civicinfo/v2"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func main() {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}
	civicKey := os.Getenv("CIVIC_KEY")
	ctx := context.Background()
	if civicKey == "" {
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to create secretmanager client: %v", err)
		}
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: "projects/end-qualified-immunity/secrets/civic-key/versions/latest",
		}
		result, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			log.Fatalf("failed to access secret version: %v", err)
		}
		civicKey = string(result.Payload.Data)
	}
	civicinfoService, err := civicinfo.NewService(ctx, option.WithAPIKey(civicKey))
	if err != nil {
		log.Fatal(err)
	}
	data, err := civicinfoService.Representatives.RepresentativeInfoByDivision("ocd-division/country:us", nil).Do()
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, data.Divisions)
	}))
}
