package interceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/sys"
	"github.com/Muvi7z/chat-auth-s/internal/sys/codes"
	"github.com/Muvi7z/chat-auth-s/internal/sys/validate"
	"google.golang.org/grpc"
	codes2 "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCStatusInterface interface {
	Status() status.Status
}

func ErrorCodesInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	res, err := handler(ctx, req)
	if nil == err {
		return res, nil
	}

	fmt.Printf("err: %s\n", err.Error())

	switch {
	case sys.IsCommonError(err):
		commErr := sys.GetCommonError(err)
		code := toGRPCCode(commErr.Code())

		err = status.Error(code, commErr.Error())
	case validate.IsValidationError(err):
		err = status.Error(codes2.InvalidArgument, err.Error())

	default:
		//var se GRPCStatusInterface
		//if errors.As(err, &se) {
		//	return nil, se.Status().Err()
		//} else {
		if errors.Is(err, context.DeadlineExceeded) {
			err = status.Error(codes2.DeadlineExceeded, err.Error())
		} else if errors.Is(err, context.Canceled) {
			err = status.Error(codes2.Canceled, err.Error())
		} else {
			err = status.Error(codes2.Internal, "internal error")
		}
	}

	return nil, err

	//}
}

func toGRPCCode(code codes.Code) codes2.Code {
	var res codes2.Code

	switch code {
	case codes.Ok:
		res = codes2.OK
	case codes.Canceled:
		res = codes2.Canceled
	case codes.InvalidArgument:
		res = codes2.InvalidArgument
	case codes.DeadlineExceeded:
		res = codes2.DeadlineExceeded
	case codes.NotFound:
		res = codes2.NotFound
	case codes.AlreadyExists:
		res = codes2.AlreadyExists
	case codes.PermissionDenied:
		res = codes2.PermissionDenied
	case codes.ResourceExhausted:
		res = codes2.ResourceExhausted
	case codes.FailedPrecondition:
		res = codes2.FailedPrecondition
	case codes.Aborted:
		res = codes2.Aborted
	case codes.OutOfRange:
		res = codes2.OutOfRange
	case codes.Unimplemented:
		res = codes2.Unimplemented
	case codes.Internal:
		res = codes2.Internal
	case codes.Unavailable:
		res = codes2.Unavailable
	case codes.DataLoss:
		res = codes2.DataLoss
	case codes.Unauthenticated:
		res = codes2.Unauthenticated
	default:
		res = codes2.Unknown
	}

	return res
}
