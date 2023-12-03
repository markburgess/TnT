
# Trust and Trustworthiness

This library is a stripped down version of the Trustability project https://github.com/markburgess/Trustability. There is no rocket science in these methods. They provide simple solution methods that are just nor provided by other projects. They can be viewed as pedagogical examples to be incorporated into other code, perhaps replacing files with embedded databases etc, where appropriate.

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
