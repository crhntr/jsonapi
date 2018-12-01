// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/pivotal/go/src/github.com/crhntr/jsonapi/response_writers.go

// Package mock_jsonapi is a generated GoMock package.
package mock_jsonapi

import (
	jsonapi "github.com/crhntr/jsonapi"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDataSetter is a mock of DataSetter interface
type MockDataSetter struct {
	ctrl     *gomock.Controller
	recorder *MockDataSetterMockRecorder
}

// MockDataSetterMockRecorder is the mock recorder for MockDataSetter
type MockDataSetterMockRecorder struct {
	mock *MockDataSetter
}

// NewMockDataSetter creates a new mock instance
func NewMockDataSetter(ctrl *gomock.Controller) *MockDataSetter {
	mock := &MockDataSetter{ctrl: ctrl}
	mock.recorder = &MockDataSetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataSetter) EXPECT() *MockDataSetterMockRecorder {
	return m.recorder
}

// SetData mocks base method
func (m *MockDataSetter) SetData(resourceType, id string, attributes interface{}, relationships jsonapi.Relationships, links jsonapi.Linker, meta jsonapi.Meta) error {
	ret := m.ctrl.Call(m, "SetData", resourceType, id, attributes, relationships, links, meta)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetData indicates an expected call of SetData
func (mr *MockDataSetterMockRecorder) SetData(resourceType, id, attributes, relationships, links, meta interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetData", reflect.TypeOf((*MockDataSetter)(nil).SetData), resourceType, id, attributes, relationships, links, meta)
}

// MockDataAppender is a mock of DataAppender interface
type MockDataAppender struct {
	ctrl     *gomock.Controller
	recorder *MockDataAppenderMockRecorder
}

// MockDataAppenderMockRecorder is the mock recorder for MockDataAppender
type MockDataAppenderMockRecorder struct {
	mock *MockDataAppender
}

// NewMockDataAppender creates a new mock instance
func NewMockDataAppender(ctrl *gomock.Controller) *MockDataAppender {
	mock := &MockDataAppender{ctrl: ctrl}
	mock.recorder = &MockDataAppenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataAppender) EXPECT() *MockDataAppenderMockRecorder {
	return m.recorder
}

// AppendData mocks base method
func (m *MockDataAppender) AppendData(resourceType, id string, attributes interface{}, relationships jsonapi.Relationships, links jsonapi.Linker, meta jsonapi.Meta) error {
	ret := m.ctrl.Call(m, "AppendData", resourceType, id, attributes, relationships, links, meta)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendData indicates an expected call of AppendData
func (mr *MockDataAppenderMockRecorder) AppendData(resourceType, id, attributes, relationships, links, meta interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendData", reflect.TypeOf((*MockDataAppender)(nil).AppendData), resourceType, id, attributes, relationships, links, meta)
}

// MockIdentifierSetter is a mock of IdentifierSetter interface
type MockIdentifierSetter struct {
	ctrl     *gomock.Controller
	recorder *MockIdentifierSetterMockRecorder
}

// MockIdentifierSetterMockRecorder is the mock recorder for MockIdentifierSetter
type MockIdentifierSetterMockRecorder struct {
	mock *MockIdentifierSetter
}

// NewMockIdentifierSetter creates a new mock instance
func NewMockIdentifierSetter(ctrl *gomock.Controller) *MockIdentifierSetter {
	mock := &MockIdentifierSetter{ctrl: ctrl}
	mock.recorder = &MockIdentifierSetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentifierSetter) EXPECT() *MockIdentifierSetterMockRecorder {
	return m.recorder
}

// SetIdentifier mocks base method
func (m *MockIdentifierSetter) SetIdentifier(resourceType, id string) error {
	ret := m.ctrl.Call(m, "SetIdentifier", resourceType, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetIdentifier indicates an expected call of SetIdentifier
func (mr *MockIdentifierSetterMockRecorder) SetIdentifier(resourceType, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIdentifier", reflect.TypeOf((*MockIdentifierSetter)(nil).SetIdentifier), resourceType, id)
}

// MockIdentifierAppender is a mock of IdentifierAppender interface
type MockIdentifierAppender struct {
	ctrl     *gomock.Controller
	recorder *MockIdentifierAppenderMockRecorder
}

// MockIdentifierAppenderMockRecorder is the mock recorder for MockIdentifierAppender
type MockIdentifierAppenderMockRecorder struct {
	mock *MockIdentifierAppender
}

// NewMockIdentifierAppender creates a new mock instance
func NewMockIdentifierAppender(ctrl *gomock.Controller) *MockIdentifierAppender {
	mock := &MockIdentifierAppender{ctrl: ctrl}
	mock.recorder = &MockIdentifierAppenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentifierAppender) EXPECT() *MockIdentifierAppenderMockRecorder {
	return m.recorder
}

// AppendIdentifier mocks base method
func (m *MockIdentifierAppender) AppendIdentifier(resourceType, id string) error {
	ret := m.ctrl.Call(m, "AppendIdentifier", resourceType, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendIdentifier indicates an expected call of AppendIdentifier
func (mr *MockIdentifierAppenderMockRecorder) AppendIdentifier(resourceType, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendIdentifier", reflect.TypeOf((*MockIdentifierAppender)(nil).AppendIdentifier), resourceType, id)
}

// MockErrorAppender is a mock of ErrorAppender interface
type MockErrorAppender struct {
	ctrl     *gomock.Controller
	recorder *MockErrorAppenderMockRecorder
}

// MockErrorAppenderMockRecorder is the mock recorder for MockErrorAppender
type MockErrorAppenderMockRecorder struct {
	mock *MockErrorAppender
}

// NewMockErrorAppender creates a new mock instance
func NewMockErrorAppender(ctrl *gomock.Controller) *MockErrorAppender {
	mock := &MockErrorAppender{ctrl: ctrl}
	mock.recorder = &MockErrorAppenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockErrorAppender) EXPECT() *MockErrorAppenderMockRecorder {
	return m.recorder
}

// AppendError mocks base method
func (m *MockErrorAppender) AppendError(err error) {
	m.ctrl.Call(m, "AppendError", err)
}

// AppendError indicates an expected call of AppendError
func (mr *MockErrorAppenderMockRecorder) AppendError(err interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendError", reflect.TypeOf((*MockErrorAppender)(nil).AppendError), err)
}

// MockIncluder is a mock of Includer interface
type MockIncluder struct {
	ctrl     *gomock.Controller
	recorder *MockIncluderMockRecorder
}

// MockIncluderMockRecorder is the mock recorder for MockIncluder
type MockIncluderMockRecorder struct {
	mock *MockIncluder
}

// NewMockIncluder creates a new mock instance
func NewMockIncluder(ctrl *gomock.Controller) *MockIncluder {
	mock := &MockIncluder{ctrl: ctrl}
	mock.recorder = &MockIncluderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIncluder) EXPECT() *MockIncluderMockRecorder {
	return m.recorder
}

// Include mocks base method
func (m *MockIncluder) Include(resourceType, id string, attributes interface{}, links jsonapi.Linker, meta jsonapi.Meta) error {
	ret := m.ctrl.Call(m, "Include", resourceType, id, attributes, links, meta)
	ret0, _ := ret[0].(error)
	return ret0
}

// Include indicates an expected call of Include
func (mr *MockIncluderMockRecorder) Include(resourceType, id, attributes, links, meta interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Include", reflect.TypeOf((*MockIncluder)(nil).Include), resourceType, id, attributes, links, meta)
}
