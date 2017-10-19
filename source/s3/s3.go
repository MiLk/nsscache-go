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

type JSONArray []map[string]interface{}

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

	r := make(JSONArray, 0)
	err = json.Unmarshal([]byte(data), &r)

	if err != nil {
		return errors.Wrap(err, "json decoding")
	}

	for _, elem := range r {
		str, _ := json.Marshal(elem)

		e := createEntry()
		err = json.Unmarshal([]byte(str), e)
		if err != nil {
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
