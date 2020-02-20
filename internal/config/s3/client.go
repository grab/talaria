// Copyright 2019-2020 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package s3

import (
	"context"
	"sync"
	"time"

	"github.com/grab/talaria/internal/monitor/errors"
	"github.com/kelindar/loader/s3"
)

type downloader interface {
	DownloadIf(ctx context.Context, bucket, prefix string, updatedSince time.Time) ([]byte, error)
}

type client struct {
	sync.Mutex
	updatedAt  time.Time  // The last updated time
	configData []byte     // The configuration file last downloaded
	downloader downloader // The downloader to use
}

// newClient a new S3 Client.
func newClient(dl downloader) (*client, error) {
	if dl != nil {
		return &client{
			downloader: dl,
		}, nil
	}

	c, err := s3.New("", 5)
	if err != nil {
		return nil, errors.Internal("unable to create loader client for s3", err)
	}

	return &client{
		downloader: c,
	}, nil
}

// Download a specific key from the bucket
func (s *client) Download(ctx context.Context, bucket, key string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	// Try to download newer file
	data, err := s.downloader.DownloadIf(ctx, bucket, key, s.updatedAt)
	if err != nil {
		return nil, err
	}

	// If we already have the data, return cached one
	if data == nil {
		return s.configData, nil
	}

	// Update the cache and time
	s.updatedAt = time.Now()
	s.configData = data
	return data, nil
}
