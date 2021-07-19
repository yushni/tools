package upload

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type uploadFunc func(uploader uploader, stream io.ReadCloser, hashFromRequest string) error

func Do(ctx context.Context, stream io.ReadCloser, testCase uploadFunc) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return testCase(uploader{}, stream, "5f17f65591794efc048db9bea5132de5")
}

func Fixed(uploader uploader, stream io.ReadCloser, hashFromRequest string) error {
	defer stream.Close()

	hash := md5.New()
	pr, pw := io.Pipe()

	go func() {
		if err := uploader.uploadStream(pr); err != nil {
			_ = pr.CloseWithError(err)
		}
	}()

	writers := io.MultiWriter(hash, pw)
	if _, err := io.Copy(writers, stream); err != nil {
		return err
	}
	if err := pw.Close(); err != nil {
		return err
	}

	expHash := hex.EncodeToString(hash.Sum(nil))
	if hashFromRequest != expHash {
		return errors.New("invalid md5 hash")
	}

	return nil
}

func NotFixed(uploader uploader, stream io.ReadCloser, hashFromRequest string) error {
	defer stream.Close()

	hash := md5.New()
	buff := new(bytes.Buffer)

	writers := io.MultiWriter(hash, buff)
	if _, err := io.Copy(writers, stream); err != nil {
		return err
	}

	expHash := hex.EncodeToString(hash.Sum(nil))
	if hashFromRequest != expHash {
		fmt.Println(expHash)
		return errors.New("invalid md5 hash")
	}

	if err := uploader.uploadStream(buff); err != nil {
		return err
	}

	return nil
}

type uploader struct{}

func (u uploader) uploadStream(r io.Reader) error {
	f, err := os.Create("/Users/shni/temp/" + time.Now().String())
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return nil
}
