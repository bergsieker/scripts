if [ $DEBUG_BASH ] ; then
  echo In $BIN_PUBLIC_ROOT/config/.bash_profile
fi

PATH=$BIN_PUBLIC_ROOT/local:$BIN_PUBLIC_ROOT/common:$PATH
export PATH
