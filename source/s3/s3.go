package s3

import (
	"encoding/json"
	"fmt"

	"github.com/MiLk/nsscache-go/cache"
	"github.com/MiLk/nsscache-go/source"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/pkg/errors"
)

/*
S3Source describes a source.Source for S3 backends:
  - prefix: the path within the S3 bucket to the passwd, shadow and group files
  - bucket: the name of the S3 bucket
  - client: the S3 client
*/
type S3Source struct {
	prefix string
	bucket string
	client s3iface.S3API
}

// CreateS3Source returns a new Source for fetching data from S3 backends
func CreateS3Source(client s3iface.S3API, prefix string, bucket string) source.Source {
	return &S3Source{
		client: client,
		prefix: prefix,
		bucket: bucket,
	}
}

func (s *S3Source) run(key string, c *cache.Cache, createEntry func() cache.Entry) error {
	if s.prefix != "" {
		key = fmt.Sprintf("%s/%s", s.prefix, key)
	}

	data, err := DownloadS3Data(s.client, s.bucket, key)

	if err != nil {
		return errors.Wrap(err, "downloading from S3")
	}

	r := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return errors.Wrap(err, "json decoding")
	}

	for _, elem := range r {
		str, _ := json.Marshal(elem)

		e := createEntry()
		if err := json.Unmarshal([]byte(str), e); err != nil {
			return errors.Wrap(err, "json does not match entry format")
		}

		c.Add(e)
	}

	return nil
}

// FillPasswdCache downloads shadow file from S3, parses the JSON and writes the passwd NSS cache file to disk.
func (s *S3Source) FillPasswdCache(c *cache.Cache) error {
	return s.run("passwd", c, func() cache.Entry {
		return &cache.PasswdEntry{}
	})
}

// FillShadowCache downloads shadow file from S3, parses the JSON and writes the shadow NSS cache file to disk.
func (s *S3Source) FillShadowCache(c *cache.Cache) error {
	return s.run("shadow", c, func() cache.Entry {
		return &cache.ShadowEntry{}
	})
}

// FillGroupCache downloads shadow file from S3, parses the JSON and writes the group NSS cache file to disk.
func (s *S3Source) FillGroupCache(c *cache.Cache) error {
	return s.run("group", c, func() cache.Entry {
		return &cache.GroupEntry{}
	})
}
