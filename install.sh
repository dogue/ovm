#!/usr/bin/env bash

# OVM install script - v0.1.5 - OVM: https://github.com/dogue/ovm
# Install script shamelessly stolen with permission from https://github.com/tristanisham

ARCH=$(uname -m)
OS=$(uname -s)


if [ $ARCH = "x86_64" ]; then
    ARCH="amd64"
fi

# echo "Installing ovm-$OS-$ARCH"

install_latest() {
    echo -e "Installing $1 in $(pwd)/ovm"
    if [ "$(uname)" = "Darwin" ]; then
     # Do something under MacOS platform

        if command -v wget >/dev/null 2>&1; then
    
            echo "wget is installed. Using wget..."
            wget -q --show-progress --max-redirect 5 -O ovm.tar "https://github.com/dogue/ovm/releases/latest/download/$1"
        else
            echo "wget is not installed. Using curl..."
            curl -L --max-redirs 5 "https://github.com/dogue/ovm/releases/latest/download/$1" -o ovm.tar
        fi
        
        mkdir -p $HOME/.ovm/self
        tar -xf ovm.tar -C $HOME/.ovm/self
        rm "ovm.tar"
        
    elif [ $OS = "Linux" ]; then
     # Do something under GNU/Linux platform
        if command -v wget >/dev/null 2>&1; then
    
            echo "wget is installed. Using wget..."
            wget -q --show-progress --max-redirect 5 -O ovm.tar "https://github.com/dogue/ovm/releases/latest/download/$1"
        else
            echo "wget is not installed. Using curl..."
            curl -L --max-redirs 5 "https://github.com/dogue/ovm/releases/latest/download/$1" -o ovm.tar
        fi
        
        mkdir -p $HOME/.ovm/self
        tar -xf ovm.tar -C $HOME/.ovm/self
        rm "ovm.tar"
    elif [ $OS = "MINGW32_NT" ]; then
    # Do something under 32 bits Windows NT platform
        curl -L --max-redirs 5 "https://github.com/dogue/ovm/releases/latest/download/$($1)" -o ovm.zip

    elif [ $OS == "MINGW64_NT" ]; then
    # Do something under 64 bits Windows NT platform
        curl -L --max-redirs 5 "https://github.com/dogue/ovm/releases/latest/download/$($1)" -o ovm.zip

    fi
}



if [ "$(uname)" = "Darwin" ]; then
    # Do something under Mac OS X platform
    install_latest "ovm-darwin-$ARCH.tar"
elif [ $OS = "Linux" ]; then
     # Do something under GNU/Linux platform
    install_latest "ovm-linux-$ARCH.tar"
elif [ $OS = "MINGW32_NT" ]; then
    # Do something under 32 bits Windows NT platform
    install_latest "ovm-windows-$ARCH.zip"
elif [ $OS == "MINGW64_NT" ]; then
    # Do something under 64 bits Windows NT platform
    install_latest "ovm-windows-$ARCH.zip"
fi

echo
echo "Run the following commands to put OVM on your path via $HOME/.profile"
echo 
# Check if TERM is set to a value that typically supports colors
if [[ "$TERM" == "xterm" || "$TERM" == "xterm-256color" || "$TERM" == "screen" || "$TERM" == "tmux" ]]; then
    # Colors
    RED='\033[0;31m'        # For strings
    GREEN='\033[0;32m'      # For commands
    BLUE='\033[0;34m'       # For variables
    NC='\033[0m'            # No Color

    echo -e "${GREEN}echo${NC} ${RED}\"# OVM\"${NC} ${GREEN}>>${NC} ${BLUE}\$HOME/.profile${NC}"
    echo -e "${GREEN}echo${NC} ${RED}'export OVM_INSTALL=\"\$HOME/.ovm/self\"'${NC} ${GREEN}>>${NC} ${BLUE}\$HOME/.profile${NC}"
    echo -e "${GREEN}echo${NC} ${RED}'export PATH=\"\$PATH:\$HOME/.ovm/bin\"'${NC} ${GREEN}>>${NC} ${BLUE}\$HOME/.profile${NC}"
    echo -e "${GREEN}echo${NC} ${RED}'export PATH=\"\$PATH:\$OVM_INSTALL/\"'${NC} ${GREEN}>>${NC} ${BLUE}\$HOME/.profile${NC}"

    echo -e "Run '${GREEN}source ~/.profile${NC}' to start using OVM in this shell!"

else
    echo 'echo "# OVM" >> $HOME/.profile'
    echo 'echo '\''export OVM_INSTALL="$HOME/.ovm/self"'\'' >> $HOME/.profile'
    echo 'echo '\''export PATH="$PATH:$HOME/.ovm/bin"'\'' >> $HOME/.profile'
    echo 'echo '\''export PATH="$PATH:$OVM_INSTALL/"'\'' >> $HOME/.profile'

    echo "Run 'source ~/.profile' to start using OVM in this shell!"

fi
    
echo
