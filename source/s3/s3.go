package s3

import (
	"encoding/json"
	"fmt"

	"github.com/MiLk/nsscache-go/cache"
	"github.com/MiLk/nsscache-go/source"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/pkg/errors"
)

type S3Source struct {
	prefix string
	bucket string
	client s3iface.S3API
}

func CreateS3Source(client s3iface.S3API, prefix string, bucket string) source.Source {
	return &S3Source{
		client: client,
		prefix: prefix,
		bucket: bucket,
	}
}

func (s *S3Source) run(key string, e cache.Entry, c *cache.Cache) error {
	if s.prefix != "" {
		key = fmt.Sprintf("%s/%s", s.prefix, key)
	}

	data, err := DownloadS3Data(s.client, s.bucket, key)

	if err != nil {
		return errors.Wrap(err, "downloading from S3")
	}

	err = json.Unmarshal([]byte(data), e)

	if err != nil {
		return errors.Wrap(err, "json decoding")
	}

	c.Add(e)
	return nil
}

// FillPasswdCache downloads shadow file from S3, parses the JSON and writes the passwd NSS cache file to disk.
func (s *S3Source) FillPasswdCache(c *cache.Cache) error {
	return s.run("passwd", &cache.PasswdEntry{}, c)
}

// FillShadowCache downloads shadow file from S3, parses the JSON and writes the shadow NSS cache file to disk.
func (s *S3Source) FillShadowCache(c *cache.Cache) error {
	return s.run("shadow", &cache.ShadowEntry{}, c)
}

// FillGroupCache downloads shadow file from S3, parses the JSON and writes the group NSS cache file to disk.
func (s *S3Source) FillGroupCache(c *cache.Cache) error {
	return s.run("group", &cache.GroupEntry{}, c)
}
