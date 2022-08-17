// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github_test

import (
	"context"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/elastic/elastic-agent-changelog-tool/internal/githubtest"
)

func TestPRHasLabels(t *testing.T) {
	r, hc := githubtest.GetHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		c      *github.Client
		owner  string
		repo   string
		pr     int
		labels []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "verify throws error on non existing PR",
			args:    args{ctx: ctx, c: c, owner: "elastic", repo: "elastic-agent-changelog-tool", pr: -1, labels: []string{}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "verify has label",
			args:    args{ctx: ctx, c: c, owner: "elastic", repo: "elastic-agent-changelog-tool", pr: 28, labels: []string{"ci"}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "verify does not have label",
			args:    args{ctx: ctx, c: c, owner: "elastic", repo: "elastic-agent-changelog-tool", pr: 28, labels: []string{"does-not-exists"}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := github.PRHasLabels(tt.args.ctx, tt.args.c, tt.args.owner, tt.args.repo, tt.args.pr, tt.args.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("PRHasLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PRHasLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}
