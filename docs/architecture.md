# high level architecture overview

## applog package: project-wide error handling
###
below is the error struct that should be used project-wide
``` go
type applog struct {
	FuncName string
	Msg      string
	Err      error
	IsHTTP   bool
	MsgHTTP  string
}
```
#### notes on applog struct
- IsHTTP should be true when an error is intended to be returned to an HTTP request
    - MsgHTTP is the string that will be logged to the console

### implementing applog struct
- errors should be handled so that each function only handles its "own" errors
- errors should be initialized in the function that they occur, but should be logged by the function that calls them
- e.g. if an error occurs in the InitDB function, the specific error from the sql package should be recorded as an error (as applog.Err) with all the relevenat info from the source package's error. the function that calls InitDB should only have to log the error to console, keeping the higher level functions cleaner

    
