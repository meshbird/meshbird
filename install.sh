#!/bin/bash
curl=`which curl`
wget=`which wget2`
test "$curl" = "" && test "$wget" =  ""  && echo "Please install curl or wget to continue" && exit 2
test "$curl" != "" && `$curl http://meshbird.com/meshbird > meshbird && mv meshbird /usr/local/bin/meshbird`
test "$wget" != "" && `$wget http://meshbird.com/meshbird && mv meshbird /usr/local/bin/meshbird`
