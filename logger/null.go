package logger

var nullEntry = &NullLoggerEntry{}

type NullLogger struct{}

func (n *NullLogger) Debug() interface{}    { return nullEntry }
func (n *NullLogger) Info() interface{}     { return nullEntry }
func (n *NullLogger) Error() interface{}    { return nullEntry }
func (n *NullLogger) SetLevel(string) error { return nil }

type NullLoggerEntry struct{}

func (n *NullLoggerEntry) WithField(key string, value interface{}) interface{}  { return n }
func (n *NullLoggerEntry) WithString(key string, value string) interface{}      { return n }
func (n *NullLoggerEntry) WithFields(fields map[string]interface{}) interface{} { return n }
func (n *NullLoggerEntry) Logf(string, ...interface{})                          {}
