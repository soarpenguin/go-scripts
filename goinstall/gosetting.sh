#!/bin/sh 
# Copyright 2013 Yasutaka Kawamoto. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo ""
echo "environmet variables for go."

if [ $# -lt 2 ];then
   echo "Please provide user and install path."
   exit 1
fi

installdir="$2"

os=`uname`
if [ $os = "Linux" ]; then
#    echo "Linux"
    userdir="/home/$1"
elif [ $os = "Darwin" ]; then
#    echo "Mac"
    userdir="/Users/$1"
else
    echo "not Linux or Mac"
    echo "Exit."
    exit 1
fi


set_bash(){
	#export GOOS=darwin
	#export GOARCH=amd64
	echo "export GOROOT=$installdir" >> $userdir/.bashrc
	echo "export GOBIN=$installdir/bin" >> $userdir/.bashrc
	echo "export PATH=$installdir/bin:$PATH" >> $userdir/.bashrc
}

set_tcsh(){
	echo "setenv GOROOT $installdir" >>  $userdir/$shrc
	echo "setenv GOBIN $installdir/bin" >>  $userdir/$shrc
	echo "setenv PATH $installdir/bin:$PATH" >>  $userdir/$shrc
}


sh=`echo $SHELL`
if [ $sh = "/bin/bash" ]; then
	echo "bash"
	set_bash $1
	. /home/$1/.bashrc
elif [ $sh = "/bin/tcsh" -o $sh = "/bin/csh" ]; then
	echo "tcsh"
	if [ $sh = "/bin/tcsh" ];then
		shrc=".tcshrc"
	elif [ $sh = "/bin/csh" ]; then
		shrc=".cshrc"
	fi
	set_tcsh $1
else
	echo "other"
fi


echo ""
echo "You may type \"source ~/$shrc\", if go doesn't work."
echo ""

