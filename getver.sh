#!/bin/bash
VER=`~/Downloads/vault version`
# NEWVER must be a string such as "0.11.2"
# The comparison is stdin against value in -compare
# If stdin < -compare then return -lessthan
# If stdin = -compare then return -equals
# If stdin > -compare then return -morethan
NEWVER="0.12"
echo "Version is: $VER"
RESULT=`echo "$VER" | ./vercmp -compare="$NEWVER" -noexitcode=true -lessthan=0 -equals=1 -morethan=2`
echo "result is $RESULT"
if [[ $RESULT -eq 0 ]]
then
  echo "Proceeding to upgrade"
else
  echo "Not upgrading"
fi

