#!/bin/bash

MYNAME="${0##*/}"

#https://www.golangtc.com/download
#wget https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz go1.7.tar.gz

function info() {
    echo -e "\033[1;34m$1 \033[0m"
}

function warn() {
    echo  -e "\033[0;33m$1 \033[0m"
}

function error() {
    echo  -e "\033[0;31m$1 \033[0m"
}

function usage() {
    info "Upgrade or install golang..."
    info "USAGE:"
    info "     ./${MYNAME} tar_file gopath"
    info "          tar_file  specify where is the tar file of go binary file"
    info "          gopath    specify where is the go workspace, include src, bin, pkg folder"
}

function createGoPath() {
    if [ ! -d $1 ]; then
        mkdir -p $1
    fi
    if [ ! -d "$1/src" ]; then
        mkdir "$1/src"
    fi
    if [ ! -d "$1/bin" ]; then
        mkdir "$1/bin"
    fi
    if [ ! -d "$1/pkg" ]; then
        mkdir "$1/pkg"
    fi
}

if [ -z $1 ];
then
    usage
    exit 1
fi

file=$1
if [ ! -f $file ];
then
    error "${file} not exist..."
    exit 1
fi

unzipPath="`pwd`/tmp_unzip_path/"
info "tmp unzip path: $unzipPath"

if [ ! -d $unzipPath ];
then
    info "$unzipPath not exist"
    mkdir $unzipPath
fi

tar -zxf $file -C $unzipPath

goroot=$GOROOT
if [ -z $GOROOT ];
then
    user=`whoami`
    goroot="/home/${user}/programs/go"
    warn "Use default go root ${goroot}"
fi


gopath=$2
info "Create go workspace, include src,bin,pkg folder..."
if [ -z $2 ]; then
    user=`whoami`
    gopath="/home/${user}/programs/golib"
    warn "Use $gopath as golang workspace..."
    if [ ! -d $gopath ]; then
        mkdir -p $gopath
    fi
fi

createGoPath $gopath

info "Copy go unzip files to $goroot"
sudo cp -r "$unzipPath/go" $goroot
rm -rf $unzipPath

etcProfile="/etc/profile"
exportGoroot="export GOROOT=$goroot"
if [ ! -z $GOROOT ];
then
    cat $etcProfile | sed 's/^export.GOROOT.*//' | sudo tee $etcProfile > /dev/null
fi
echo $exportGoroot | sudo tee -a $etcProfile

exportGopath="export GOROOT=$gopath"
if [ ! -z $GOPATH ];
then
    cat $etcProfile | sed 's/^export.GOPATH.*//' | sudo tee $etcProfile > /dev/null
fi
echo "export GOPATH=$gopath" | sudo tee -a  $etcProfile

echo 'export PATH=$GOROOT/bin:$GOPATH/bin:$PATH' | sudo tee -a  $etcProfile

# ## Replace multiple empty lines with one empty line
cat $etcProfile -s | sudo tee $etcProfile > /dev/null

info "To make configuration take effect, will reboot, pls enter[y/n]"
read -p "[y/n]" isReboot
if [ $isReboot = "y" ];
then
    sudo reboot
fi
