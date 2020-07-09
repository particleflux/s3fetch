package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
	"os"
	"strings"
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

	sess := session.Must(session.NewSession(&aws.Config{Region: region}))

	downloader := s3manager.NewDownloader(sess)
	f := os.Stdout
	if *output != "-" {
		f, err = os.Create(*output)
		if err != nil {
			exitErrorf("Failed to open file: %v", err)
		}
		defer f.Close()
	}

	numBytes, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(objectURL.Host),
		Key:    aws.String(strings.TrimLeft(objectURL.Path, "/")),
	})
	if err != nil {
		exitErrorf("Failed to download file: %v", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "File downloaded, %d bytes\n", numBytes)
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
