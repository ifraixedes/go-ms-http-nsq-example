// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ifraixedes/go-ms-http-nsq-example/driver-location (interfaces: Service)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	driver_location "github.com/ifraixedes/go-ms-http-nsq-example/driver-location"
	reflect "reflect"
)

// Service is a mock of Service interface
type Service struct {
	ctrl     *gomock.Controller
	recorder *ServiceMockRecorder
}

// ServiceMockRecorder is the mock recorder for Service
type ServiceMockRecorder struct {
	mock *Service
}

// NewService creates a new mock instance
func NewService(ctrl *gomock.Controller) *Service {
	mock := &Service{ctrl: ctrl}
	mock.recorder = &ServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Service) EXPECT() *ServiceMockRecorder {
	return m.recorder
}

// LocationsForLastMinutes mocks base method
func (m *Service) LocationsForLastMinutes(arg0 context.Context, arg1 uint64, arg2 uint16) ([]driver_location.Location, error) {
	ret := m.ctrl.Call(m, "LocationsForLastMinutes", arg0, arg1, arg2)
	ret0, _ := ret[0].([]driver_location.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LocationsForLastMinutes indicates an expected call of LocationsForLastMinutes
func (mr *ServiceMockRecorder) LocationsForLastMinutes(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocationsForLastMinutes", reflect.TypeOf((*Service)(nil).LocationsForLastMinutes), arg0, arg1, arg2)
}

// SetLocation mocks base method
func (m *Service) SetLocation(arg0 context.Context, arg1 uint64, arg2 driver_location.Location) error {
	ret := m.ctrl.Call(m, "SetLocation", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLocation indicates an expected call of SetLocation
func (mr *ServiceMockRecorder) SetLocation(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLocation", reflect.TypeOf((*Service)(nil).SetLocation), arg0, arg1, arg2)
}
