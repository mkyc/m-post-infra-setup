package cmd

type ZeroLogger struct{}

func (z ZeroLogger) Trace(format string, v ...interface{}) {
	logger.
		Trace().
		Msgf(format, v...)
}

func (z ZeroLogger) Debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func (z ZeroLogger) Info(format string, v ...interface{}) {
	logger.
		Info().
		Msgf(format, v...)
}

func (z ZeroLogger) Warn(format string, v ...interface{}) {
	logger.
		Warn().
		Msgf(format, v...)
}

func (z ZeroLogger) Error(format string, v ...interface{}) {
	logger.
		Error().
		Msgf(format, v...)
}

func (z ZeroLogger) Fatal(format string, v ...interface{}) {
	logger.
		Fatal().
		Msgf(format, v...)
}

func (z ZeroLogger) Panic(format string, v ...interface{}) {
	logger.
		Panic().
		Msgf(format, v...)
}
