// Copyright 2012 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/google/gousb"

	"golang.org/x/sync/errgroup"

	"github.com/puhitaku/mtplvcap/mtp"
	"github.com/puhitaku/mtplvcap/public"
	"github.com/sirupsen/logrus"
)

func main() {
	host := flag.String("host", "localhost", "hostname: default = localhost, specify 0.0.0.0 for public access")
	port := flag.Int("port", 42839, "port: default = 42839")
	backendGo := flag.Bool("backend-go", false, "use gousb as a libusb wrapper (not recommended)")
	debug := flag.String("debug", "", "comma-separated list of debugging options: usb, data, mtp, server")
	serverOnly := flag.Bool("server-only", false, "serve frontend without opening a DSLR (for devevelopment)")
	vendorID := flag.String("vendor-id", "0x0", "VID of the camera to search (in hex), default=0x0 (all)")
	productID := flag.String("product-id", "0x0", "PID of the camera to search (in hex), default=0x0 (all)")

	flag.Parse()

	debugs := map[string]bool{}
	for _, s := range strings.Split(*debug, ",") {
		debugs[s] = true
	}

	if _, ok := debugs["server"]; ok {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}

	vid, err := strconv.ParseInt(strings.ReplaceAll(*vendorID, "0x", ""), 16, 64)
	if err != nil {
		log.WithField("prefix", "main").Fatalf("Failed to parse VID: %s", err)
	}
	pid, err := strconv.ParseInt(strings.ReplaceAll(*productID, "0x", ""), 16, 64)
	if err != nil {
		log.WithField("prefix", "main").Fatalf("Failed to parse PID: %s", err)
	}

	mtp.SetLogger(log)
	var dev mtp.Device

	if *serverOnly {
		log.WithField("prefix", "mtp").Info("Server-only mode is activated, skipping USB initialization")
	} else {
		if *backendGo {
			ctx := gousb.NewContext()
			defer ctx.Close()

			devGo, err := mtp.SelectDeviceGoUSB(ctx, uint16(vid), uint16(pid))
			if err != nil {
				log.WithField("prefix", "mtp").Fatalf("Failed to detect MTP device: %s", err)
			}
			defer devGo.Close()
			dev = devGo
		} else {
			devDirect, err := mtp.SelectDeviceDirect(uint16(vid), uint16(pid))
			if err != nil {
				log.WithField("prefix", "mtp").Fatalf("Failed to detect MTP devices: %v", err)
			}
			defer devDirect.Close()
			dev = devDirect
		}

		dev.SetDebug(mtp.DebugFlags{
			MTP:  debugs["mtp"],
			USB:  debugs["data"],
			Data: debugs["usb"],
		})
		if err = dev.Configure(); err != nil {
			log.WithField("prefix", "mtp").Fatalf("Configure failed: %v", err)
		}
	}

	eg, ctx := errgroup.WithContext(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-sigChan:
				log.WithField("prefix", "signal").Info("Caught signal: %s", s)
				return errors.New(s.String())
			}
		}
	})

	lvs := mtp.NewLVServer(dev, log, ctx)
	eg.Go(lvs.Run)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, _ := public.Root.Open("/controller.html")
		_, _ = io.Copy(w, f)
	})
	router.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		f, _ := public.Root.Open("/index.html")
		_, _ = io.Copy(w, f)
	})
	router.HandleFunc("/stream", lvs.HandleStream)
	router.HandleFunc("/control", lvs.HandleControl)
	router.Handle("/assets/", http.FileServer(public.Root))

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: logging(router),
	}

	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.WithField("prefix", "http").Error(err)
			return err
		}
		return nil
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
		}
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		return srv.Shutdown(ctx)
	})

	err = eg.Wait()
	if err != nil {
		os.Exit(1)
	}
}
