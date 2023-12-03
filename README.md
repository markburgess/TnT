
# Trust and Trustworthiness

This library is a stripped down version of the Trustability project https://github.com/markburgess/Trustability. There is no rocket science in these methods. They provide simple solution methods that are just nor provided by other projects. They can be viewed as pedagogical examples to be incorporated into other code, perhaps replacing files with embedded databases etc, where appropriate.

## Test program examples

`test_context_eval.go` - evaluate some test policy expressions

`test_context.go` - example of using context in dynamic changes

`test_promise_wrapper.go` - example of using the promise locking wrapper

## Promise instrumentation methods


Two sets of functions for wrapping transactional events or critical sections parenthetically (with begin-end semantics).

* Functions that timestamp transactions at the moment of measurement according to the local system clock

```
 PromiseContext_Begin(g Analytics, name string) PromiseContext 
 PromiseContext_End(g Analytics, ctx PromiseContext) PromiseHistory 
```

* Functions with Golang time stamp supplied from outside, e.g. for offline analysis with epoch timestamps.

```
 StampedPromiseContext_Begin(g Analytics, name string, before time.Time) PromiseContext 
 StampedPromiseContext_End(g Analytics, ctx PromiseContext, after time.Time) PromiseHistory
```

The wrappers refuse to acquire a lock unless a minmum time has elapsed.
If a maximum holding time has expired, the lock will be forced.
The configured values are currently fixed at 30 and sixty seconds.
```
   ifelapsed := int64(30)
   expireafter := int64(60)
```

## Context methods

The method is to assign real numbers between 0and 1 to flag/signal
variables that are set as a result of conditions within a program, and
to evaluate quasi-logical probabilistic expressions to describe policy
conditions. An undefined variable or a variable with value 0 is
effectively false and any positive value is somewhat true.

-`InitializeContext()` - reset all symbols in context to undefined / false

-`ContextActive(string)` - active the symbol, set to true

-`ContextSet()` - return the set of defined symbols

-`SetContext(s string,c float64)` - set the real value of a context variable to an explicit real confidence value

-`Confidence(s string) float64` - return the real value named symbol

-`ContextEval(s string) (string,float64)` - return the real value of the expression evaluated according to AND/OR algebra rules

-`IsDefinedContext(s string) bool` - if the expression evaluates to a result greater than zero according to AND/OR algebra rules this returns true


## Running the code:

My working environment is GNU/Linux, where everything is simple. Setting up the working environment for all the parts is a little bit of work (more steps than are desirable), but it should be smooth.

1. Install git client packages on your computer.
2. Go to: https://golang.org/dl/ to download a package.
3. Some file management to create a working directory and link it to environment variables:

```
  $ mkdir -p ~/go/bin
  $ mkdir -p ~/go/src
  $ cd ~/somedirectory
  $ git clone https://github.com/markburgess/TnT.git
  $ ln -s ~/somedirectory/Trustability/pkg/TnT ~/go/src/TnT
```

4. It's useful to put this in your ~/.bashrc file
```
export PATH=$PATH:/usr/local/go/binexport GOPATH=~/go
```
Don’t forget to restart your shell or command window after editing this or do a
```
  $ source ~/.bashrc
```
5. You can get fetch the drivers for using the graph database code and the 
```
  $ go get github.com/arangodb/go-driver
```
