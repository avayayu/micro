package rpc

// // FromError convert error for service reply and try to convert it to grpc.Status.
// func FromError(svrErr error) (gst *status.Status) {
// 	var err error
// 	svrErr = errors.Cause(svrErr)
// 	if code, ok := svrErr.(ecode.Codes); ok {
// 		// TODO: deal with err
// 		if gst, err = gRPCStatusFromEcode(code); err == nil {
// 			return
// 		}
// 	}
// 	// for some special error convert context.Canceled to ecode.Canceled,
// 	// context.DeadlineExceeded to ecode.DeadlineExceeded only for raw error
// 	// if err be wrapped will not effect.
// 	switch svrErr {
// 	case context.Canceled:
// 		gst, _ = gRPCStatusFromEcode(ecode.Canceled)
// 	case context.DeadlineExceeded:
// 		gst, _ = gRPCStatusFromEcode(ecode.Deadline)
// 	default:
// 		gst, _ = status.FromError(svrErr)
// 	}
// 	return
// }
