#!/bin/sh
# Copyright 2015 Yasutaka Kawamoto / Soarpenguin. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo ""
echo "Environmet variables for go."

if [ $# -lt 2 ];then
   echo "Please provide user and install path."
   echo "Usage: ./gosetting.sh user installdir"
   exit 1
fi

user="$1"
installdir="$2"

os=`uname`
if [ $os = "Linux" ]; then
#    echo "Linux"
    userdir="/home/$user"
elif [ $os = "Darwin" ]; then
#    echo "Mac"
    userdir="/Users/$user"
else
    echo "Not Linux or Mac"
    echo "Exit."
    exit 1
fi

check_gopath() {
    local shrc=$1

    if [ "x$shrc" = "x" ]; then
        echo "Please provide sh rc file."
        exit 1
    fi

    grep -i GOPATH "$shrc" &>/dev/null
    return $?
}

set_shrc() {
    local shrc=$1

    if [ "x$shrc" = "x" ]; then
        echo "Please provide sh rc file."
        exit 1
    fi

    check_gopath $shrc && return 0

    echo "export GOROOT=$installdir" >>  $shrc
    echo "export GOBIN=$installdir/bin" >>  $shrc
    echo "export PATH=\$GOBIN:\$PATH" >>  $shrc
}

set_tcsh(){
    local shrc=$1

    check_gopath $shrc && return 0

    echo "setenv GOROOT $installdir" >>  $shrc
    echo "setenv GOBIN $installdir/bin" >>  $shrc
    echo "setenv PATH \$GOBIN:\$PATH" >>  $shrc
}

sh=`echo $SHELL`
if [ $sh = "/bin/bash" ]; then
    echo "Setting bash ..."
    shrc=".bashrc"
    set_shrc "${userdir}/.bashrc"
    . "${userdir}/.bashrc"
elif [ $sh = "/bin/tcsh" -o $sh = "/bin/csh" ]; then
    echo "Setting tcsh ..."
    if [ $sh = "/bin/tcsh" ];then
        shrc=".tcshrc"
    elif [ $sh = "/bin/csh" ]; then
        shrc=".cshrc"
    fi
    set_tcsh "${userdir}/$shrc"
elif [ $sh = "/bin/zsh" ];then
    echo "Setting zsh ..."
    shrc=".zshrc"
    set_shrc "${userdir}/.zshrc"
    . "${userdir}/.zshrc"
else
    echo "other"
fi


echo ""
echo "You may type \"source ~/$shrc\", if go doesn't work."
echo ""

