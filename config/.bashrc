if [ $DEBUG_BASH ] ; then
  echo In $BIN_PUBLIC_ROOT/config/.bashrc
fi

# ENVIRONMENT VARIABLES

export EDITOR=vim

#export LS_OPTIONS="-F -N --color=auto"
# For Mac OS 10.5
export LSCOLORS="gxfxcxdxbxegedabagacad"
# For linux
export LS_COLORS=$LS_COLORS"di=01;36:"

export PS1="\[\e[33m\]\h:\w>\[\e[0m\]"

# ALIASES

alias cp="cp -i"
alias mv="mv -i"
alias rm="rm -i"
alias e="$EDITOR"
alias sbrc="source ~/.bashrc"

# export VIM="/Applications/MacVim.app/Contents/MacOS/Vim"
# export VIMDIFF="/Applications/MacVim.app/Contents/MacOS/Vim -d"
# alias vi=$VIM

alias screenstart='screen -c $BIN_PUBLIC_ROOT/config/.screenrc'
alias find_untracked_files='find . -type f -print0 | xargs -0 p4 fstat >/dev/null'

bind '"\e[A"':history-search-backward
bind '"\e[B"':history-search-forward
bind Space:magic-space    

# FUNCTIONS

