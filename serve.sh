#!/bin/sh

scene=$1

[ -n "$scene" -a -d "$scene" ] || { echo "usage: $0 <scenedir>" ; exit 1;  }

srcdir=$PWD
root=`mktemp -d talehttproot-XXXXXX`

if [ -d $root ]
then

trap 'rm -r $root' 1 2 3 15

echo "$srcdir/$scene -> $root"

ln -s $srcdir/3rdparty $root/import
ln -s $srcdir/shader $root/shader
ln -s $srcdir/talescene/build $root/talescene
ln -s $srcdir/$scene/* $root/

ls -l $root

webfsd -F -p 8080 -r $root -j -f index.html

else
echo "failed to create '$root'"
fi
