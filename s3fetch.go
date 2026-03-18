package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Set on build time
var Version = ""
var Revision = ""

func main() {
	versionFlag := flag.Bool("version", false, "Version")
	region := flag.String("region", "", "AWS Region")
	output := flag.String("output", "-", "Output file")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options] s3://your-bucket/path/to/object\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if *versionFlag {
		fmt.Printf("s3fetch: %s (%s)\n%s: %s\n", Version, Revision, aws.SDKName, aws.SDKVersion)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	objectURL, err := url.Parse(flag.Args()[0])
	if err != nil {
		exitErrorf("%v", err)
		os.Exit(1)
	}
	if objectURL.Scheme != "s3" {
		exitErrorf("Only s3 URLs supported, example: s3://your-bucket/path/to/object")
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(*region),
	)
	if err != nil {
		exitErrorf("Failed to init aws config: %v", err)
	}

	svc := s3.NewFromConfig(cfg)

	downloader := transfermanager.New(svc)
	f := os.Stdout
	if *output != "-" {
		f, err = os.Create(*output)
		if err != nil {
			exitErrorf("Failed to open file: %v", err)
		}
		defer f.Close()
	}

	key := strings.TrimLeft(objectURL.Path, "/")
	dlOutput, err := downloader.DownloadObject(
		context.Background(),
		&transfermanager.DownloadObjectInput{
			Bucket:   aws.String(objectURL.Host),
			Key:      aws.String(key),
			WriterAt: f,
		},
	)
	if err != nil {
		exitErrorf("Failed to download file: %v", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "File downloaded, %d bytes\n", *dlOutput.ContentLength)
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
