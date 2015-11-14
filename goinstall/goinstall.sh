#!/bin/sh 
# Copyright 2013 Yasutaka Kawamoto. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# const 
export PS4='+ [`basename ${BASH_SOURCE[0]}`:$LINENO ${FUNCNAME[0]} \D{%F %T} $$ ] '
CURDIR=$(cd "$(dirname "$0")"; pwd);
MYNAME="${0##*/}"

g_INSTALL_DIR="/usr/local/go"
g_USER=""

#########################
_report_err() { echo "${MYNAME}: Error: $*" >&2 ; }

if [ -t 1 ]
then
    RED="$( echo -e "\e[31m" )"
    HL_RED="$( echo -e "\e[31;1m" )"
    HL_BLUE="$( echo -e "\e[34;1m" )"

    NORMAL="$( echo -e "\e[0m" )"
fi

_hl_red()    { echo "$HL_RED""$@""$NORMAL";}
_hl_blue()   { echo "$HL_BLUE""$@""$NORMAL";}

_trace() {
    echo $(_hl_blue '  ->') "$@" >&2
}

_print_fatal() {
    echo $(_hl_red '==>') "$@" >&2
}

readlinkf() { perl -MCwd -e 'print Cwd::abs_path shift' $1;}
# ABSPATH="$(readlinkf ./non-absolute/file)"

_usage() {
    cat << USAGE
Usage: bash ${MYNAME} -u username path.

Options:
    -u, --user      user   Login user install golang.
    -h, --help             Print this help infomation.

Require:
    path             Path string for install golang.

USAGE

    exit $RET_OK
}

#
# Parses command-line options.
#  usage: _parse_options "$@" || exit $?
#
_parse_options()
{
    declare -a argv

    while [[ $# -gt 0 ]]; do
        case $1 in
            -u|--user)
                g_USER=${2}
                shift 2
                ;;
            -h|--help)
                _usage
                exit
                ;;
            --)
                shift
                argv=("${argv[@]}" "${@}")
                break
                ;;
            -*)
                _print_fatal "command line: unrecognized option $1" >&2
                return 1
                ;;
            *)
                argv=("${argv[@]}" "${1}")
                shift
                ;;
        esac
    done

    case ${#argv[@]} in
        1)
            command -v greadlink >/dev/null 2>&1 && g_INSTALL_DIR=$(greadlink -f "${argv[0]}") || g_INSTALL_DIR=$(readlink -f "${argv[0]}")
            ;;
        0|*)
            _usage 1>&2
            return 1
    ;;
    esac
}

################################## main route #################################
_parse_options "${@}" || _usage

if [ x"$g_USER" == "x" ]; then
    _print_fatal "Please privoid user name install golang."
    _usage
    exit 1
fi

ret=`id -u $g_USER`
if [ $? -ne 0 ]; then
    _print_fatal "user $g_USER is not exist."
    exit 1
fi

# Check whether go is installed.
if  [ ! -f $g_INSTALL_DIR/bin/go ]; then
    echo "go is not Installed. Continue..."
else
    echo "Go is installed."
    echo "Exit."
    #exit 1
fi

#Check OS
if [ -f /etc/lsb-release ]; then
    . /etc/lsb-release
    dlcmd="apt-get -y"
    homedir="/home/$1"
elif [ -f /etc/debian_version ]; then
    dlcmd="apt-get -y"
    homedir="/home/$1"
elif [ -f /etc/redhat-release ]; then
    dlcmd="yum -y"
    homedir="/home/$1"
elif [ -f /etc/system-release ]; then
    dlcmd="yum -y"
    homedir="/home/$1"
elif [ `uname` = "Darwin" ]; then #for Mac
    homedir="/Users/$1"
else
    echo "not Linux or Mac"
    exit 1
fi


# Install gcc
if  ! type >/dev/null "gcc" 2>&1 ; then
    if [ `uname` = "Darwin" ]; then #for Mac
        echo "You need to install Xcode and gcc."
        exit 1
    else
        echo "Installing gcc ..." 
        $dlcmd install gcc
        echo "Done"
    fi
else
    echo "gcc is installed."
fi

# Install Marcurial
if  ! type >/dev/null "hg" 2>&1 ; then
    if [ `uname` = "Darwin" ]; then #for Mac
        echo "You need to install Marcurial."
        echo "http://mercurial.selenic.com/downloads/"
        exit 1
    else
        echo "Installing Mercurial ..." 
        $dlcmd install mercurial
        echo "Done"
    fi
else
    echo "Mercurial is installed."
fi

# Install Go
if [ -d $g_INSTALL_DIR ]
then 
    echo "go dir is existed."
else
    echo "downloading go ..."
    `hg clone -u release https://code.google.com/p/go $g_INSTALL_DIR`
    echo "Done"
fi

#echo "move directory to go/src."
cd $g_INSTALL_DIR/src/
echo `pwd`
if  ! [ -d $g_INSTALL_DIR/bin -a -d $g_INSTALL_DIR/pkg ]; then
    echo "executing ./all.bash ..."
    { . ./all.bash; }
    echo "Done"
else
    echo "executed ./all.bash"
fi

# Set environment variables
cd $CURDIR
{ . ./gosetting.sh $g_USER $g_INSTALL_DIR; }

echo ""
echo "ALL DONE"
echo ""

