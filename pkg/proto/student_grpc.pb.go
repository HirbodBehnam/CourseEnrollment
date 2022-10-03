// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: pkg/proto/student.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CourseEnrollmentServerServiceClient is the client API for CourseEnrollmentServerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CourseEnrollmentServerServiceClient interface {
	// This method must enroll the student
	StudentEnroll(ctx context.Context, in *StudentEnrollRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// This method must remove the student from an enrolled course
	StudentDisenroll(ctx context.Context, in *StudentDisenrollRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// This method must change the group of a registered user
	StudentChangeGroup(ctx context.Context, in *StudentChangeGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// This method must send the registered courses of a user
	GetStudentEnrolledCourses(ctx context.Context, in *GetStudentCoursesRequest, opts ...grpc.CallOption) (*StudentCourseDataArray, error)
	// This method will get all courses in a department
	GetCoursesOfDepartment(ctx context.Context, in *GetDepartmentCoursesRequest, opts ...grpc.CallOption) (*DepartmentCourses, error)
}

type courseEnrollmentServerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCourseEnrollmentServerServiceClient(cc grpc.ClientConnInterface) CourseEnrollmentServerServiceClient {
	return &courseEnrollmentServerServiceClient{cc}
}

func (c *courseEnrollmentServerServiceClient) StudentEnroll(ctx context.Context, in *StudentEnrollRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.CourseEnrollmentServerService/StudentEnroll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseEnrollmentServerServiceClient) StudentDisenroll(ctx context.Context, in *StudentDisenrollRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.CourseEnrollmentServerService/StudentDisenroll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseEnrollmentServerServiceClient) StudentChangeGroup(ctx context.Context, in *StudentChangeGroupRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.CourseEnrollmentServerService/StudentChangeGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseEnrollmentServerServiceClient) GetStudentEnrolledCourses(ctx context.Context, in *GetStudentCoursesRequest, opts ...grpc.CallOption) (*StudentCourseDataArray, error) {
	out := new(StudentCourseDataArray)
	err := c.cc.Invoke(ctx, "/proto.CourseEnrollmentServerService/GetStudentEnrolledCourses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseEnrollmentServerServiceClient) GetCoursesOfDepartment(ctx context.Context, in *GetDepartmentCoursesRequest, opts ...grpc.CallOption) (*DepartmentCourses, error) {
	out := new(DepartmentCourses)
	err := c.cc.Invoke(ctx, "/proto.CourseEnrollmentServerService/GetCoursesOfDepartment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourseEnrollmentServerServiceServer is the server API for CourseEnrollmentServerService service.
// All implementations must embed UnimplementedCourseEnrollmentServerServiceServer
// for forward compatibility
type CourseEnrollmentServerServiceServer interface {
	// This method must enroll the student
	StudentEnroll(context.Context, *StudentEnrollRequest) (*emptypb.Empty, error)
	// This method must remove the student from an enrolled course
	StudentDisenroll(context.Context, *StudentDisenrollRequest) (*emptypb.Empty, error)
	// This method must change the group of a registered user
	StudentChangeGroup(context.Context, *StudentChangeGroupRequest) (*emptypb.Empty, error)
	// This method must send the registered courses of a user
	GetStudentEnrolledCourses(context.Context, *GetStudentCoursesRequest) (*StudentCourseDataArray, error)
	// This method will get all courses in a department
	GetCoursesOfDepartment(context.Context, *GetDepartmentCoursesRequest) (*DepartmentCourses, error)
	mustEmbedUnimplementedCourseEnrollmentServerServiceServer()
}

// UnimplementedCourseEnrollmentServerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCourseEnrollmentServerServiceServer struct {
}

func (UnimplementedCourseEnrollmentServerServiceServer) StudentEnroll(context.Context, *StudentEnrollRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudentEnroll not implemented")
}
func (UnimplementedCourseEnrollmentServerServiceServer) StudentDisenroll(context.Context, *StudentDisenrollRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudentDisenroll not implemented")
}
func (UnimplementedCourseEnrollmentServerServiceServer) StudentChangeGroup(context.Context, *StudentChangeGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudentChangeGroup not implemented")
}
func (UnimplementedCourseEnrollmentServerServiceServer) GetStudentEnrolledCourses(context.Context, *GetStudentCoursesRequest) (*StudentCourseDataArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentEnrolledCourses not implemented")
}
func (UnimplementedCourseEnrollmentServerServiceServer) GetCoursesOfDepartment(context.Context, *GetDepartmentCoursesRequest) (*DepartmentCourses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoursesOfDepartment not implemented")
}
func (UnimplementedCourseEnrollmentServerServiceServer) mustEmbedUnimplementedCourseEnrollmentServerServiceServer() {
}

// UnsafeCourseEnrollmentServerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourseEnrollmentServerServiceServer will
// result in compilation errors.
type UnsafeCourseEnrollmentServerServiceServer interface {
	mustEmbedUnimplementedCourseEnrollmentServerServiceServer()
}

func RegisterCourseEnrollmentServerServiceServer(s grpc.ServiceRegistrar, srv CourseEnrollmentServerServiceServer) {
	s.RegisterService(&CourseEnrollmentServerService_ServiceDesc, srv)
}

func _CourseEnrollmentServerService_StudentEnroll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentEnrollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseEnrollmentServerServiceServer).StudentEnroll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CourseEnrollmentServerService/StudentEnroll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseEnrollmentServerServiceServer).StudentEnroll(ctx, req.(*StudentEnrollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseEnrollmentServerService_StudentDisenroll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentDisenrollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseEnrollmentServerServiceServer).StudentDisenroll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CourseEnrollmentServerService/StudentDisenroll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseEnrollmentServerServiceServer).StudentDisenroll(ctx, req.(*StudentDisenrollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseEnrollmentServerService_StudentChangeGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentChangeGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseEnrollmentServerServiceServer).StudentChangeGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CourseEnrollmentServerService/StudentChangeGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseEnrollmentServerServiceServer).StudentChangeGroup(ctx, req.(*StudentChangeGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseEnrollmentServerService_GetStudentEnrolledCourses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStudentCoursesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseEnrollmentServerServiceServer).GetStudentEnrolledCourses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CourseEnrollmentServerService/GetStudentEnrolledCourses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseEnrollmentServerServiceServer).GetStudentEnrolledCourses(ctx, req.(*GetStudentCoursesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseEnrollmentServerService_GetCoursesOfDepartment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepartmentCoursesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseEnrollmentServerServiceServer).GetCoursesOfDepartment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CourseEnrollmentServerService/GetCoursesOfDepartment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseEnrollmentServerServiceServer).GetCoursesOfDepartment(ctx, req.(*GetDepartmentCoursesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CourseEnrollmentServerService_ServiceDesc is the grpc.ServiceDesc for CourseEnrollmentServerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CourseEnrollmentServerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CourseEnrollmentServerService",
	HandlerType: (*CourseEnrollmentServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StudentEnroll",
			Handler:    _CourseEnrollmentServerService_StudentEnroll_Handler,
		},
		{
			MethodName: "StudentDisenroll",
			Handler:    _CourseEnrollmentServerService_StudentDisenroll_Handler,
		},
		{
			MethodName: "StudentChangeGroup",
			Handler:    _CourseEnrollmentServerService_StudentChangeGroup_Handler,
		},
		{
			MethodName: "GetStudentEnrolledCourses",
			Handler:    _CourseEnrollmentServerService_GetStudentEnrolledCourses_Handler,
		},
		{
			MethodName: "GetCoursesOfDepartment",
			Handler:    _CourseEnrollmentServerService_GetCoursesOfDepartment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/student.proto",
}
