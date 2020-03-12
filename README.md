# logex

Package logex is a logging package with somewhat more control than the stdlib log package and is largely inspired by the logrus README file.

API is not yet set.

Interface to a Logger is as follows:

```
// Log defines an interface for a logger.
type Log interface {

	// Debugf will log a debug message formed from format string and args.
	Debugf(string, ...interface{})
	// Debugln will log args as a debug message.
	Debugln(...interface{})
	// Infof will log an info message formed from format string and args.
	Infof(string, ...interface{})
	// Infoln will log args as an info message.
	Infoln(...interface{})
	// Warningf will log a warning message formed from format string and args.
	Warningf(string, ...interface{})
	// Warningln will log args as a warning message.
	Warningln(...interface{})
	// Errorf will log an error and an error message formed from format string and args.
	Errorf(error, string, ...interface{})
	// Errorln will log an error and args as a warning message.
	Errorln(error, ...interface{})

	// Printf will log a message with a custom logging level formed from format string and args.
	Printf(LogLevel, string, ...interface{})
	// Println will log args as a message with custom logging level.
	Println(LogLevel, ...interface{})

	// ToOutputs will return a clone which will output only to specified output names.
	ToOutputs(names ...string) Log

	// Caller will append the caller field to the next logged line.
	WithCaller(skip int) Log
	// Stack will append the stack field to the next logged line.
	WithStack(skip int, depth int) Log
	// Fields will append the specified fields to the next logged line.
	WithFields(*Fields) Log
}
```

The following log levels are predefined:
```
const (
	// LevelNone is undefined logging level. It prints nothing.
	LevelNone LogLevel = iota
	// LevelMute is the silent logging level used to silence the logger.
	LevelMute
	// LevelError is the error logging level that prints errors only.
	LevelError
	// LevelWarning is the warning logging level that prints warnings and errors.
	LevelWarning
	// LevelInfo is the info logging level that prints information, warnings and errors.
	LevelInfo
	// LevelDebug is the debug logging level that prints debug messages, information, warnings and errors.
	LevelDebug
	// LevelCustom and levels up to LevelPrint are custom logging levels.
	// To define a custom logging level use: MyLevel := LogLevel(LevelCustom +1).
	LevelCustom
	// LevelPrint is the print logging level that prints everything that gets logged.
	LevelPrint = LogLevel(255)
)
```

Custom log levels are defineable in the range LevelCustom+1 ... LevelPrint.
## Usage

To create a new logger that outputs to `stdout` using default simple text formatter use `NewStd()`.

Example:

```
l := NewStd(nil)
l.Errorln(errors.New("actual error message"), "additional error message")
l.Debugln("debug info)
```

```
// Output:
[2020-03-03 12:56:49] Error: additional error message
        actual error message
[2020-03-03 12:59:13] Debug: debug info

```

To create a new custom logger that formats errors to json and outputs to custom oputput use `New()`.

Example:

```
l := New(nil)
l.AddOutput("stdout", os.StdOut, NewJSONFormatter(true))
l.Errorln(errors.New("actual error message"), "additional error message")
l.Debugln("debug info)
```

```
// Output:
{
        "error": {},
        "loglevel": 2,
        "message": "additional error message\n",
        "time": "2020-03-03T13:00:19.46159898+01:00"
}
{
        "loglevel": 5,
        "message": "debug info\n",
        "time": "2020-03-03T13:00:19.461790209+01:00"
}
```

You can add additional fields per log line.

Example:

```
l := New(nil)
l.SetLevel(LevelPrint)
l.AddOutput("stdout", os.Stdout, NewJSONFormatter(true))
f := NewFields()
f.Set("customfield", "customvalue")
l.Fields(f).Println(LevelInfo, "some log message")
```

```
// Output:
{
        "customfield": "customvalue",
        "loglevel": 255,
        "message": "some log message\n",
        "time": "2020-03-03T13:05:46.550892961+01:00"
}
```

You can also create custom formatters.

```
// Formatter formats Fields to a custom format.
type Formatter interface {
	// Format must return a string representation of Fields, such as JSON, CSV, Text, etc..
	Format(*Fields) string
}
```

## License

See included LICENSE file.