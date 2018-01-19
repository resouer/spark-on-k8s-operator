/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsLocalFile(t *testing.T) {
	type testcase struct {
		file    string
		isLocal bool
	}

	testFn := func(test testcase, t *testing.T) {
		isLocal, err := isLocalFile(test.file)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, test.isLocal, isLocal, "%s: expected %v got %v", test.file, test.isLocal, isLocal)
	}

	testcases := []testcase {
		{file: "/path/to/file", isLocal: true},
		{file: "file:///path/to/file", isLocal: true},
		{file: "local:///path/to/file", isLocal: false},
		{file: "http://localhost/path/to/file", isLocal: false},
	}

	for _, test := range testcases {
		testFn(test, t)
	}
}

func TestFilterLocalFiles(t *testing.T) {
	files := []string{
		"/path/to/file",
		"file:///file/to/path",
		"http://localhost/path/to/file",
		"hdfs://localhost/path/to/file",
		"gs://bucket/path/to/file",
	}

	expected := []string{
		"/path/to/file",
		"file:///file/to/path",
	}

	actual, err := filterLocalFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)
}

func TestLoadFromYAML(t *testing.T) {
	app, err := loadFromYAML("testdata/test-app.yaml")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, app.Name, "example")
	assert.Equal(t, *app.Spec.MainClass, "org.examples.SparkExample")
	assert.Equal(t, app.Spec.MainApplicationFile, "local:///path/to/example.jar")
	assert.Equal(t, app.Spec.Driver.Image, "spark")
	assert.Equal(t, app.Spec.Executor.Image, "spark")
	assert.Equal(t, int(*app.Spec.Executor.Instances), 1)
}

func TestHandleHadoopConfiguration(t *testing.T) {
	configMap, err := buildHadoopConfigMap("test", "testdata/hadoop-conf")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, configMap.Name, "test-hadoop-config")
	assert.Equal(t, len(configMap.Data), 1)
	assert.True(t, strings.Contains(configMap.Data["core-site.xml"], "fs.gs.impl"))
}