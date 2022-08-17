// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testutil

import (
	"path/filepath"
	"testing"
)

// TestCollection is a set of integration tests.
type TestCollection struct {
	Tests []IntegrationTest
}

// CollectTests builds TestCollection from directory structure.
func CollectTests() (*TestCollection, error) {
	collection := TestCollection{}

	paths, err := filepath.Glob("./testdata/*/test.yaml")
	if err != nil {
		return nil, err
	}

	for _, manifestPath := range paths {
		manifest, err := NewTestManifest(manifestPath)
		if err != nil {
			return nil, err
		}

		path := filepath.Dir(manifestPath)

		collection.Tests = append(collection.Tests, IntegrationTest{
			Name:     filepath.Base(path),
			Path:     path,
			Manifest: manifest,
		})
	}

	return &collection, nil
}

// Each iterates over collection providing runner for each test.
func (collection *TestCollection) Each(f func(string, func(t *testing.T))) {
	for _, test := range collection.Tests {
		f(test.Name, test.Run)
	}
}
