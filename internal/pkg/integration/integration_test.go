/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package integration_test

import (
	"testing"

	"github.com/talos-systems/bldr/internal/pkg/util/testutil"
)

func TestIntegration(t *testing.T) {
	collection, err := testutil.CollectTests()
	if err != nil {
		t.Fatalf("error collecting tests: %v", err)
	}

	collection.Each(func(name string, f func(*testing.T)) {
		t.Run(name, f)
	})
}
