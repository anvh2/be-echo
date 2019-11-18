package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	errorCodes = []string{
		"-12001", "-12002", "-12003", "-12004", "-12005", "-12006",
		"-12007", "-12008", "-12009", "-12010", "-12011", "-12012",
		"-12013", "-12014", "-12015", "-12016", "-12017", "-12018",
		"-12019", "-12020", "-12023", "-12024",
	}
	// vector of counter, with one for each error code
	counterError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zpi_be",
		Subsystem: "lixi",
		Name:      "error_code",
		Help:      "The total number of error code",
	}, []string{"kind", "code", "method"})
)

// UnaryServerInterceptor -
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			errorInterceptor(resp, info.FullMethod)
		}
		return resp, err
	}
}

func errorInterceptor(resp interface{}, fullMethod string) error {
	item, err := marshal(resp)
	if err != nil {
		return err
	}

	// get error
	if err, ok := item["error"]; ok {
		defineErrorCode(err, fullMethod)
		return nil
	}

	if err, ok := item["Error"]; ok {
		defineErrorCode(err, fullMethod)
		return nil
	}

	return errors.New("no field name Error")
}

func defineErrorCode(groupErr interface{}, fullMethod string) error {
	item, err := marshal(groupErr)
	if err != nil {
		return err
	}
	if code, ok := item["code"]; ok {
		code := fmt.Sprintf("%v", code)
		if isInErrorCodes(code) {
			counterError.WithLabelValues("grpc", code, fullMethod).Inc()
			return nil
		}
	}
	return errors.New("can't parse interface to struct group Error")
}

// convert interface to map
func marshal(item interface{}) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	itemByte, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(string(itemByte)), &fields)
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func isInErrorCodes(code string) bool {
	for _, errorCode := range errorCodes {
		if code == errorCode {
			return true
		}
	}
	return false
}
