package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/armiariyan/logger"
	"github.com/armiariyan/synapsis/internal/config"
	"github.com/spf13/cast"
)

var pkgLogger logger.Logger

func New() {
	opt := logger.Options{
		Name: config.GetString("app.name"),
		SysOptions: logger.OptionsLogger{
			Type: logger.File,
			OptionsFile: logger.OptionsFile{
				Stdout:       config.GetBool("logger.sys.stdout"),
				FileLocation: config.GetString("logger.sys.fileLocation"),
				FileMaxAge:   time.Duration(config.GetInt("logger.sys.fileMaxAge") * int(time.Hour) * 24),
				Mask:         config.GetBool("logger.sys.mask"),
			},
		},
		TdrOptions: logger.OptionsLogger{
			Type: logger.File,
			OptionsFile: logger.OptionsFile{
				Stdout:       config.GetBool("logger.tdr.stdout"),
				FileLocation: config.GetString("logger.tdr.fileLocation"),
				FileMaxAge:   time.Duration(config.GetInt("logger.tdr.fileMaxAge") * int(time.Hour) * 24),
				Mask:         config.GetBool("logger.tdr.mask"),
			},
		},
	}

	pkgLogger = logger.SetupLoggerCombine(opt)
}

func Info(ctx context.Context, title string, messages ...interface{}) {
	fields := formatLogs(messages...)
	pkgLogger.Info(ctx, title, fields...)
}

func Warn(ctx context.Context, title string, messages ...interface{}) {
	fields := formatLogs(messages...)
	pkgLogger.Warn(ctx, title, fields...)
}

func Fatal(ctx context.Context, title string, messages ...interface{}) {
	fields := formatLogs(messages...)
	pkgLogger.Fatal(ctx, title, fields...)
}

func Error(ctx context.Context, title string, messages ...interface{}) {
	fields := formatLogs(messages...)
	pkgLogger.Error(ctx, title, fields...)
}

func T2(ctx context.Context, title string, messages ...interface{}) {
	fields := formatLogs(messages...)
	pkgLogger.Info(ctx, title, fields...)
}

func T3(ctx context.Context, title string, startProcessTime time.Time, messages ...interface{}) {
	stop := time.Now()

	fields := formatLogs(messages...)
	fields = append(fields, logger.ToField("_process_time", fmt.Sprintf("%d ms", stop.Sub(startProcessTime).Nanoseconds()/1000000)))

	pkgLogger.Info(ctx, title, fields...)
}

func TDR(ctx context.Context, request, response []byte) {
	rt := time.Now().Sub(GetRequestTimeFromContext(ctx)).Nanoseconds() / 1000000
	ctxLogger := logger.ExtractCtx(ctx)

	tdrLog := logger.LogTdrModel{
		AppName:    ctxLogger.ServiceName,
		AppVersion: ctxLogger.ServiceVersion,
		Port:       ctxLogger.ServicePort,
		Method:     ctxLogger.ReqMethod,
		SrcIP:      GetRequestIPFromContext(ctx),
		RespTime:   rt,
		Path:       ctxLogger.ReqURI,
		Header:     GetRequestHeaderFromContext(ctx),
		Request:    string(request),
		Response:   string(response),
		ThreadID:   ctxLogger.ThreadID,
		JourneyID:  ctxLogger.JourneyID,
		Error:      GetErrorMessageFromContext(ctx),
	}

	responseMap := make(map[string]interface{})
	if err := json.Unmarshal(response, &responseMap); err != nil {
		pkgLogger.TDR(ctx, tdrLog)
		return
	}

	if responseMap["status"] != nil {
		tdrLog.ResponseCode = cast.ToString(responseMap["status"])
	}

	pkgLogger.TDR(ctx, tdrLog)
}

func formatLogs(messages ...interface{}) (logRecord []logger.Field) {
	for index, msg := range messages {
		logRecord = append(logRecord, logger.ToField("_message_"+cast.ToString(index), msg))
	}

	return
}
