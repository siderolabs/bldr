package v1alpha1

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
)

func (p *Pkg) Build() error {
	if err := environment(); err != nil {
		return err
	}

	println("Build:", os.Getenv("BUILD"))
	println("Host:", os.Getenv("HOST"))
	println("Target:", os.Getenv("TARGET"))

	for _, step := range p.Steps {
		if _, err := os.Stat("/tmp"); err != nil {
			if os.IsNotExist(err) {
				// TODO(andrewrynhard): Remove this directory.
				log.Println("Creating /tmp")
				// nolint: errcheck
				os.MkdirAll("/tmp", 0777)
			}
		}
		dir, err := ioutil.TempDir("", "pkg")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)

		if err := os.Chdir(dir); err != nil {
			log.Fatal(err)
		}
		if err := step.download(dir); err != nil {
			log.Fatal(err)
		}
		if err := step.prepare(p.Shell); err != nil {
			log.Fatal(err)
		}
		if err := step.build(p.Shell); err != nil {
			log.Fatal(err)
		}
		if err := step.install(p.Shell); err != nil {
			log.Fatal(err)
		}
		if err := step.test(p.Shell); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func environment() error {
	var result *multierror.Error

	basic := func() {
		os.Setenv("CFLAGS", "-g0 -Os")
		os.Setenv("CXXFLAGS", "-g0 -Os")
		os.Setenv("LDFLAGS", "-s")
		os.Setenv("VENDOR", "talos")
		os.Setenv("SYSROOT", "/talos")
		os.Setenv("TOOLCHAIN", "/toolchain")
		os.Setenv("PATH", fmt.Sprintf("/toolchain/bin:%s", os.Getenv("PATH")))
	}

	target := func() {
		arch, ok := os.LookupEnv("TARGETPLATFORM")
		if !ok {
			result = multierror.Append(result, xerrors.New("TARGETPLATFORM is required"))
			return
		}

		switch arch {
		case "linux/amd64":
			os.Setenv("ARCH", "x86_64")
			os.Setenv("TARGET", "x86_64-talos-linux-musl")
		case "linux/arm64":
			os.Setenv("ARCH", "aarch64")
			os.Setenv("TARGET", "aarch64-talos-linux-musl")
		case "linux/arm/v7":
			os.Setenv("ARCH", "armv7")
			os.Setenv("TARGET", "armv7-talos-linux-musl")
		default:
			result = multierror.Append(result, xerrors.Errorf("unsupported build platform: %q", arch))
			return
		}
	}

	build := func() {
		buf := &unix.Utsname{}
		if err := unix.Uname(buf); err != nil {
			result = multierror.Append(result, err)
			return
		}

		b := bytes.Trim(buf.Machine[:], "\x00")
		arch := string(b)
		switch arch {
		case "x86_64":
			os.Setenv("BUILD", "x86_64-linux-musl")
			os.Setenv("HOST", "x86_64-linux-musl")
		case "aarch64":
			os.Setenv("BUILD", "aarch64-linux-musl")
			os.Setenv("HOST", "aarch64-linux-musl")
		case "armv7l":
			os.Setenv("BUILD", "armv7-linux-musl")
			os.Setenv("HOST", "armv7-linux-musl")
		default:
			result = multierror.Append(result, xerrors.Errorf("unsupported build platform: %q", arch))
			return
		}
	}

	basic()
	target()
	build()

	return result.ErrorOrNil()
}

func download(url, filepath string) error {
	// Download the file.

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the file to disk.

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func verify(filepath string, s256, s512 string) error {
	errCh := make(chan error)

	// Verify SHA256
	go func(filepath string) {
		errCh <- func() error {
			f, err := os.Open(filepath)
			if err != nil {
				return xerrors.Errorf("%w", err)
			}
			defer f.Close()

			hash256 := sha256.New()
			if _, err := io.Copy(hash256, f); err != nil {
				return xerrors.Errorf("%w", err)
			}
			sum256 := hash256.Sum(nil)
			sha256 := hex.EncodeToString(sum256)
			if sha256 != s256 {
				return xerrors.Errorf("sha256 checksum mismatch, %s != %s", sha256, s256)
			}

			return nil
		}()
	}(filepath)

	// Verify SHA512
	go func(filepath string) {
		errCh <- func() error {
			f, err := os.Open(filepath)
			if err != nil {
				return xerrors.Errorf("%w", err)
			}

			hash512 := sha512.New()
			if _, err := io.Copy(hash512, f); err != nil {
				return xerrors.Errorf("%w", err)
			}
			sum512 := hash512.Sum(nil)
			sha512 := hex.EncodeToString(sum512)
			if sha512 != s512 {
				return xerrors.Errorf("sha512 checksum mismatch, %s != %s", sha512, s512)
			}

			return nil
		}()
	}(filepath)

	var result *multierror.Error

	for i := 0; i < 2; i++ {
		result = multierror.Append(result, <-errCh)
	}

	return result.ErrorOrNil()
}
