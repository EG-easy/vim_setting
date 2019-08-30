# Lines configured by zsh-newuser-install
HISTFILE=~/.histfile
HISTSIZE=1000
SAVEHIST=1000
bindkey -e
# End of lines configured by zsh-newuser-install
# The following lines were added by compinstall
zstyle :compinstall filename '/Users/shizuka/.zshrc'

autoload -Uz compinit
compinit
# End of lines added by compinstall
source "${ZDOTDIR:-$HOME}/.zprezto/init.zsh"

autoload -Uz colors
colors

#PROMPT="%{$fg[green]%}%D{%T}%  %~$ %{${reset_color}%}" 
PROMPT="%{$fg[green]%}%D{%T} % $ ""% " 
alias ls='ls -G'
# export LSCOLORS=gxfxcxdxbxegedabagacad
alias molpath='cd ~/go/src/github.com/moldcoin/moldex/'
alias votum='cd ~/go/src/github.com/EG-easy/votumchain/'

# ctag
alias tag='ctags -R -f .tags'


#c-lang
export PATH=$PATH:/usr/local/bin

#golang
export PATH="$PATH:/usr/local/go/bin"
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
export GOBIN=$GOPATH/bin
export GO111MODULE=on

#node
export PATH=$PATH:/Users/shizuka/.nodebrew/current/bin

#node_modules
export NODE_PATH=/usr/local/lib/node_modules

#pyenv
export PYENV_ROOT="$HOME/.pyenv"
export PATH="$PYENV_ROOT/bin:$PATH"
if command -v pyenv 1>/dev/null 2>&1; then
  eval "$(pyenv init -)"
fi

# hyper shell icon setting
# Override auto-title when static titles are desired ($ title My new title)
title() { export TITLE_OVERRIDDEN=1; echo -en "\e]0;$*\a"}
# Turn off static titles ($ autotitle)
autotitle() { export TITLE_OVERRIDDEN=0 }; autotitle
# Condition checking if title is overridden
overridden() { [[ $TITLE_OVERRIDDEN == 1 ]]; }
# Echo asterisk if git state is dirty
gitDirty() { [[ $(git status 2> /dev/null | grep -o '\w\+' | tail -n1) != ("clean"|"") ]] && echo "*" }

# Show cwd when shell prompts for input.
tabtitle_precmd() {
   if overridden; then return; fi
   pwd=$(pwd) # Store full path as variable
   cwd=${pwd##*/} # Extract current working dir only
   print -Pn "\e]0;$cwd$(gitDirty)\a" # Replace with $pwd to show full path
}
[[ -z $precmd_functions ]] && precmd_functions=()
precmd_functions=($precmd_functions tabtitle_precmd)

# Prepend command (w/o arguments) to cwd while waiting for command to complete.
tabtitle_preexec() {
   if overridden; then return; fi
   printf "\033]0;%s\a" "${1%% *} | $cwd$(gitDirty)" # Omit construct from $1 to show args
}
[[ -z $preexec_functions ]] && preexec_functions=()
preexec_functions=($preexec_functions tabtitle_preexec)

#powerline-go
# function powerline_precmd() {
#     eval "$($GOPATH/bin/powerline-go -error $? -shell zsh -eval -modules "venv,ssh,cwd,perms,gitlite,hg,jobs,exit,root,vgo" -modules-right "time")"
# }
#
# function install_powerline_precmd() {
#   for s in "${precmd_functions[@]}"; do
#     if [ "$s" = "powerline_precmd" ]; then
#       return
#     fi
#   done
#   precmd_functions+=(powerline_precmd)
# }
#
# if [ "$TERM" != "linux" ]; then
#     install_powerline_precmd
# fi
#
