// Copyright 2019-2020 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

// Code generated by mockery v1.0.0. DO NOT EDIT.

package s3sqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/mock"
)

// MockReader is an autogenerated mock type for the Reader type
type MockReader struct {
	mock.Mock
}

// DeleteMessage provides a mock function with given fields: msg
func (_m *MockReader) DeleteMessage(msg *sqs.Message) error {
	ret := _m.Called(msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sqs.Message) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartPolling provides a mock function with given fields: maxPerRead, sleepMs, attributeNames, messageAttributeNames
func (_m *MockReader) StartPolling(maxPerRead int64, sleepMs int64, attributeNames []*string, messageAttributeNames []*string) <-chan *sqs.Message {
	ret := _m.Called(maxPerRead, sleepMs, attributeNames, messageAttributeNames)

	var r0 <-chan *sqs.Message
	if rf, ok := ret.Get(0).(func(int64, int64, []*string, []*string) <-chan *sqs.Message); ok {
		r0 = rf(maxPerRead, sleepMs, attributeNames, messageAttributeNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *sqs.Message)
		}
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *MockReader) Close() error {
	_m.Called()
	return nil
}
