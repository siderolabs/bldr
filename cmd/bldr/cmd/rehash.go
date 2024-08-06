// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rehashCmd)
}

var rehashCmd = &cobra.Command{
	Use:   "update",
	Short: "Update pkgs",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		for _, d := range datas {
			if err := doTemplate(
				template.Must(template.New(d["url"]).Funcs(sprig.HermeticTxtFuncMap()).Parse(d["url"])),
				d,
			); err != nil {
				return err
			}
		}

		return nil
	},
}

// datas is a list of deps to update. Provided as an example.
//
//nolint:lll
var datas = []map[string]string{
	{"VERSION": "v1.30.3", "url": "https://github.com/kubernetes/cloud-provider-aws/archive/refs/tags/{{ .VERSION }}.tar.gz"},
	{"GLIB_VERSION": "2.81.1", "url": "https://download.gnome.org/sources/glib/{{ regexReplaceAll \".\\\\d+$\" .GLIB_VERSION \"${1}\" }}/glib-{{ .GLIB_VERSION }}.tar.xz"},
	{"GVISOR_VERSION": "20240729.0", "url": "https://github.com/google/gvisor/archive/3f38cb19ba373b027f5220450591daa3ab767145.tar.gz"},
	{"QEMU_VERSION": "9.0.2", "url": "https://download.qemu.org/qemu-{{ .QEMU_VERSION }}.tar.xz"},
	{"SPIN_VERSION": "v0.15.1", "url": "https://github.com/spinkube/containerd-shim-spin/releases/download/{{ .SPIN_VERSION }}/containerd-shim-spin-v2-linux-aarch64.tar.gz"},
	{"SPIN_VERSION": "v0.15.1", "url": "https://github.com/spinkube/containerd-shim-spin/releases/download/{{ .SPIN_VERSION }}/containerd-shim-spin-v2-linux-x86_64.tar.gz"},
	{"TAILSCALE_VERSION": "1.70.0", "url": "https://github.com/tailscale/tailscale/archive/refs/tags/v{{ .TAILSCALE_VERSION }}.tar.gz"},
	{"UTIL_LINUX_VERSION": "2.40.2", "url": "https://www.kernel.org/pub/linux/utils/util-linux/v{{ regexReplaceAll \".\\\\d+$\" .UTIL_LINUX_VERSION \"${1}\" }}/util-linux-{{  regexReplaceAll \"\\\\.0$\" .UTIL_LINUX_VERSION \"${1}\" }}.tar.xz"},
}

func doTemplate(tmpl *template.Template, data map[string]string) error {
	var bldr strings.Builder

	err := tmpl.Execute(&bldr, data)
	if err != nil {
		return fmt.Errorf("failed to execute the template '%s': %w", tmpl.Name(), err)
	}

	s256, s512 := sha256.New(), sha512.New()

	result, err := http.Get(bldr.String()) //nolint:noctx
	if err != nil {
		return fmt.Errorf("failed to fetch the URL '%s': %w", bldr.String(), err)
	}

	defer result.Body.Close() //nolint:errcheck

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("code %d != 200 for URL: %s", result.StatusCode, bldr.String())
	}

	_, err = io.Copy(io.MultiWriter(s256, s512), result.Body)
	if err != nil {
		return fmt.Errorf("failed to process the content of the URL '%s': %w", bldr.String(), err)
	}

	fmt.Println("URL:", bldr.String())
	fmt.Println("sha256:", fmt.Sprintf("%x", s256.Sum(nil)))
	fmt.Println("sha512:", fmt.Sprintf("%x", s512.Sum(nil)))

	return nil
}
