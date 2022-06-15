// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"reflect"
	"testing"
)

func Test_collectKinds(t *testing.T) {
	type args struct {
		items []Entry
	}
	tests := []struct {
		name string
		args args
		want map[Kind]bool
	}{
		{
			"no kind",
			args{[]Entry{}},
			map[Kind]bool{},
		},
		{
			"one kind",
			args{[]Entry{
				{Kind: Feature},
			}},
			map[Kind]bool{
				Feature: true,
			},
		},
		{
			"more kinds",
			args{[]Entry{
				{Kind: Feature},
				{Kind: BugFix},
			}},
			map[Kind]bool{
				Feature: true,
				BugFix:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectKinds(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectKinds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collectByKind(t *testing.T) {
	type args struct {
		items []Entry
		k     Kind
	}
	tests := []struct {
		name string
		args args
		want []Entry
	}{
		{
			"no entries",
			args{
				items: []Entry{
					{Summary: "foobar", Kind: Feature},
				},
				k: BugFix,
			},
			[]Entry{},
		},
		{
			"one entry",
			args{
				items: []Entry{
					{Summary: "foobar", Kind: Feature},
				},
				k: Feature,
			},
			[]Entry{{Summary: "foobar", Kind: Feature}},
		},
		{
			"multiple entry same kind",
			args{
				items: []Entry{
					{Summary: "foobar", Kind: Feature},
					{Summary: "foo", Kind: Feature},
				},
				k: Feature,
			},
			[]Entry{{Summary: "foobar", Kind: Feature}, {Summary: "foo", Kind: Feature}},
		},
		{
			"multiple entry different kind",
			args{
				items: []Entry{
					{Summary: "foobar", Kind: Feature},
					{Summary: "foo", Kind: BugFix},
				},
				k: Feature,
			},
			[]Entry{{Summary: "foobar", Kind: Feature}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectByKind(tt.args.items, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectByKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
