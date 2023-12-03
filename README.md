
# Trust and Trustworthiness

This library is a stripped down version of the Trustability project https://github.com/markburgess/Trustability. There is no rocket science in these methods. They provide simple solution methods that are just nor provided by other projects. They can be viewed as pedagogical examples to be incorporated into other code, perhaps replacing files with embedded databases etc, where appropriate.

## Method

The method is to assign real numbers between 0and 1 to flag/signal variables that are set as a result
of conditions within a program, and to evaluate quasi-logical probabilistic expressions to
describe policy conditions. An undefined variable or a variable with value 0 is effectively false
and any positive value is somewhat true.

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
