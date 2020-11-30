// Copyright 2020 The Tekton Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pipelinerun

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	tb "github.com/tektoncd/experimental/commit-status-tracker/test/builder"
	"knative.dev/pkg/apis"

	pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	resource "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
)

const (
	exampleRepoURL = "https://example.com/test/repo.git"
	exampleSHA     = "aa0a35a867ed2094da60042062e8f3d6000e3952"
)

// TODO: Test for no commit at all returning an error
func Test_repoAndSHA(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(a *pipelinev1.PipelineRun)
		wantRepo string
		wantSHA  string
		wantErr  string
	}{
		{
			"with git resources",
			func(p *pipelinev1.PipelineRun) {
				p.Spec.Resources = []pipelinev1.PipelineResourceBinding{
					{
						Name: "source-repo",
						ResourceSpec: &resource.PipelineResourceSpec{
							Type: pipelinev1.PipelineResourceTypeGit,
							Params: []resource.ResourceParam{
								{Name: "revision", Value: exampleSHA},
								{Name: "url", Value: exampleRepoURL},
							},
						},
					},
				}
			},
			strings.TrimSuffix(exampleRepoURL, ".git"),
			exampleSHA,
			"",
		},
		{
			"with git annotations",
			func(p *pipelinev1.PipelineRun) {
				p.ObjectMeta.Annotations = map[string]string{
					sourceURLAnnotation: exampleRepoURL,
					sourceSHAAnnotation: exampleSHA,
				}
			},
			strings.TrimSuffix(exampleRepoURL, ".git"),
			exampleSHA,
			"",
		},
		{
			"with git annotations",
			func(p *pipelinev1.PipelineRun) {
			},
			"",
			"",
			"failed to find commit for pipelinerun",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := makePipelineRun()
			tt.modifier(pr)

			c, err := commitFromPR(pr)
			if !matchError(t, tt.wantErr, err) {
				t.Errorf("repoAndSHA() %s: got error %v, want %s", tt.name, err, tt.wantErr)
				return
			}
			if c == nil {
				return
			}
			if c.repoURL != tt.wantRepo {
				t.Errorf("got repo %q, want %q", c.repoURL, tt.wantRepo)
			}
			if c.sha != tt.wantSHA {
				t.Errorf("got sha %q, want %q", c.sha, tt.wantSHA)
			}
		})
	}
}

func Test_findGitResourceWithNoRepository(t *testing.T) {
	pipelineRun := makePipelineRun()

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func Test_findGitResourceWithRepository(t *testing.T) {
	pipelineRun := makePipelineRun(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"))

	want := &pipelinev1alpha1.PipelineResourceSpec{
		Type: "git",
		Params: []pipelinev1.ResourceParam{
			{
				Name:  "url",
				Value: "https://github.com/tektoncd/triggers",
			},
			{
				Name:  "revision",
				Value: "master",
			},
		},
	}

	r, err := findGitResource(pipelineRun)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v, want %+v", r, want)
	}
}

func TestFindGitResourceWithMultipleRepositories(t *testing.T) {
	pipelineRun := makePipelineRun(
		makeGitResourceBinding("https://github.com/tektoncd/triggers", "master"),
		makeGitResourceBinding("https://github.com/tektoncd/pipeline", "master"))

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestFindGitResourceWithNonGitResource(t *testing.T) {
	pipelineRun := makePipelineRun(
		makeImageResourceBinding("example.com/project/myimage"))

	_, err := findGitResource(pipelineRun)
	if err == nil {
		t.Fatal("did not get an error with no git resource")
	}
}

func TestGetRepoAndSHA(t *testing.T) {
	resourceTests := []struct {
		name     string
		resType  pipelinev1.PipelineResourceType
		url      string
		revision string
		repo     string
		sha      string
		wantErr  string
	}{
		{"non-git resource", pipelinev1.PipelineResourceTypeImage, "", "", "", "", "non-git resource"},
		{"git resource with no url", pipelinev1.PipelineResourceTypeGit, "", "master", "", "", "failed to find param url"},
		{"git resource with no revision", pipelinev1.PipelineResourceTypeGit, exampleRepoURL, "", "", "", "failed to find param revision"},
		{"git resource", pipelinev1.PipelineResourceTypeGit, exampleRepoURL, "master", strings.TrimSuffix(exampleRepoURL, ".git"), "master", ""},
		{"git resource with .git", pipelinev1.PipelineResourceTypeGit, exampleRepoURL + ".git", "master", exampleRepoURL, "master", ""},
	}

	for _, tt := range resourceTests {
		res := makePipelineResource(tt.resType, tt.url, tt.revision)
		repo, sha, err := getRepoAndSHA(res)
		if !matchError(t, tt.wantErr, err) {
			t.Errorf("getRepoAndSHA() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("getRepoAndSHA() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}

		if tt.sha != sha {
			t.Errorf("getRepoAndSHA() %s: got SHA %s, want %s", tt.name, sha, tt.sha)
		}
	}
}

func TestGetDriverName(t *testing.T) {
	tests := []struct {
		url    string
		driver string
		errMsg string
	}{
		{"http://github.com/", "github", ""},
		{"http://github.com/foo/bar", "github", ""},
		{"https://githuB.com/foo/bar.git", "github", ""},
		{"http://gitlab.com/foo/bar.git2", "gitlab", ""},
		{"http://gitlab/foo/bar/", "", "unable to determine type of Git host from: http://gitlab/foo/bar/"},
		{"https://gitlab.a.b/foo/bar/bar", "", "unable to determine type of Git host from: https://gitlab.a.b/foo/bar/bar"},
		{"https://gitlab.org2/f.b/bar.git", "", "unable to determine type of Git host from: https://gitlab.org2/f.b/bar.git"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(rt *testing.T) {
			gotDriver, err := getDriverName(tt.url)
			if !matchError(t, tt.errMsg, err) {
				rt.Errorf("driver errMsg mismatch: got error %v, want %v", err, tt.errMsg)
			}
			if tt.driver != gotDriver {
				rt.Errorf("driver mismatch: got %v, want %v", gotDriver, tt.driver)
			}
		})
	}
}

func TestExtractRepoPath(t *testing.T) {
	repoURLTests := []struct {
		name    string
		url     string
		repo    string
		wantErr string
	}{
		{"standard URL", "https://github.com/tektoncd/triggers", "tektoncd/triggers", ""},
		{"invalid URL", "http://192.168.0.%31/test/repo", "", "failed to parse repo URL.*invalid URL escape"},
		{"url with no repo path", "https://github.com/", "", "could not determine repo from URL"},
		{"gitlab URL path", "https://gitlab.com/org/group/repo", "org/group/repo", ""},
	}

	for _, tt := range repoURLTests {
		repo, err := extractRepoPath(tt.url)
		if !matchError(t, tt.wantErr, err) {
			t.Errorf("extractRepoPath() %s: got error %v, want %s", tt.name, err, tt.wantErr)
			continue
		}

		if tt.repo != repo {
			t.Errorf("getRepoAndSHA() %s: got repo %s, want %s", tt.name, repo, tt.repo)
		}
	}
}

func makePipelineRun(opts ...tb.PipelineRunSpecOp) *pipelinev1.PipelineRun {
	return tb.PipelineRun(pipelineRunName, testNamespace, tb.PipelineRunSpec(
		"tomatoes", opts...,
	), tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
		apis.Condition{Type: apis.ConditionSucceeded}),
		tb.PipelineRunTaskRunsStatus("trname", &pipelinev1.PipelineRunTaskRunStatus{
			PipelineTaskName: "task-1",
		}),
	), tb.PipelineRunLabel("label-key", "label-value"))
}

func makeGitResourceBinding(url, rev string) tb.PipelineRunSpecOp {
	return tb.PipelineRunResourceBinding("some-resource"+randomSuffix(),
		tb.PipelineResourceBindingResourceSpec(&pipelinev1alpha1.PipelineResourceSpec{
			Type: pipelinev1.PipelineResourceTypeGit,
			Params: []pipelinev1.ResourceParam{{
				Name:  "url",
				Value: url,
			}, {
				Name:  "revision",
				Value: rev,
			}}}))
}

func makeImageResourceBinding(url string) tb.PipelineRunSpecOp {
	return tb.PipelineRunResourceBinding("some-resource"+randomSuffix(),
		tb.PipelineResourceBindingResourceSpec(&pipelinev1alpha1.PipelineResourceSpec{
			Type: pipelinev1.PipelineResourceTypeImage,
			Params: []pipelinev1.ResourceParam{{
				Name:  "url",
				Value: url,
			},
			}}))
}

func makePipelineResource(resType pipelinev1.PipelineResourceType, url, rev string) *pipelinev1alpha1.PipelineResourceSpec {
	spec := &pipelinev1alpha1.PipelineResourceSpec{
		Type: resType,
	}
	if url != "" {
		spec.Params = append(spec.Params,
			pipelinev1.ResourceParam{
				Name:  "url",
				Value: url,
			})
	}
	if rev != "" {
		spec.Params = append(spec.Params,
			pipelinev1.ResourceParam{
				Name:  "revision",
				Value: rev,
			})
	}
	return spec
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomSuffix() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
