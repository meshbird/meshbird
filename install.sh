#!/bin/bash
curl=`which curl`
wget=`which wget2`
test "$curl" = "" && test "$wget" =  ""  && echo "Please install curl or wget to continue" && exit 2
test "$curl" != "" && `$curl https://meshbird.com/meshbird > meshbird && chmod a+x ./meshbird &&  mv ./meshbird /usr/local/bin/meshbird`
test "$wget" != "" && `$wget https://meshbird.com/meshbird && chmod a+x ./meshbird && mv ./meshbird /usr/local/bin/meshbird`
