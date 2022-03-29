// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	github "github.com/google/go-github/v32/github"

	mock "github.com/stretchr/testify/mock"
)

// GithubClient is an autogenerated mock type for the GithubClient type
type GithubClient struct {
	mock.Mock
}

// UsersGet provides a mock function with given fields: ctx, user
func (_m *GithubClient) UsersGet(ctx context.Context, user string) (*github.User, *github.Response, error) {
	ret := _m.Called(ctx, user)

	var r0 *github.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *github.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.User)
		}
	}

	var r1 *github.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *github.Response); ok {
		r1 = rf(ctx, user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, user)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
