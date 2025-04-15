# Go 开发环境

### root 用户登录并配置
1 使用 root 用户登录并配置 Linux 服务器
```shell
root@VM-0-13-debian:~# cat /etc/os-release
PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"
NAME="Debian GNU/Linux"
VERSION_ID="12"
VERSION="12 (bookworm)"
VERSION_CODENAME=bookworm
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"
```

2 在 root 用户下，创建普通用户
创建普通用户
```shell
useradd -m -d /home/chen -s /bin/bash chen
passwd chen
```

3 添加 sudoers
```shell
sed -i '/^[a-z]* ALL=(ALL) NOPASSWD: ALL/a\chen ALL=(ALL) NOPASSWD: ALL' /etc/sudoers
# 方法二 
visudo
chen ALL=(ALL:ALL) ALL
# 测试 sudo，显示 root 即为成功
sudo whoami
```

### 使用普通用户
1 登录普通用户 chen

2 配置 $HOME/.bashrc 文件
vi $HOME/.bashrc
输入 :set paste 进入粘贴模式，保存退出
```shell
# .bashrc

# Source global definition
if [ -f /etc/bashrc ]; then
  . /etc/bashrc
fi 

# Alias definitions
# You may want to put all your additions into a separate file like
# ~/.bash_aliases, instead of adding them here directly.
# See /usr/share/doc/bash-doc/examples in the bash-doc package.
if [ -f ~/.bash_aliases ]; then
  . ~/.bash_aliases
fi 

# Enable color support of ls and alse add handy aliases
if [ -x /usr/bin/dircolors ]; then
  test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
  alias ls='ls --color=auto'
  alias ll='ls --color=auto -l'
  alias la='ls --color=auto -A'
  alias l='ls --color=auto -lA'
  alias dir='dir --color=auto'
  alias vdir='vdir --color=auto'
  
  alias grep='grep --color=auto'
  alias fgrep='fgrep --color=auto'
  alias egrep='egrep --color=auto'
fi 

# Colored GCC warnings and errors
export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'

# Some aliases to prevent mistakes
alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Create a default workspace directory to keep all work files in one place
if [ ! -d $HOME/workspace ]; then
  mkdir -p $HOME/workspace
fi 

# User-specific environment settings
# Basic environment
# Set system language to en_US.UTF-8 to avoid Chinese chatacter display issues in the terminal
export LANG="en_US.UTF-8"
# The default PS1 settings displays the full path, to prevent it from bacoming too long,
# it now shows "username@dev last_directory_name"
export PS1='[\u@dev \W]\$'
# Set the workspace directory
export WORKSPACE="$HOME/workspace"
# Add $HOME/bin directory to the PATH variable
export PATH=$HOME/bin:$PATH
# Set the default editor to vim
export EDITOR=vim

# When logging into the system, default to the Workspace directory
cd $WORKSPACE

# User-specific aliases, configures and fuctions


```
### 依赖安装和配置

1 安装依赖

```shell
sudo apt-get update
sudo apt install -y build-essential jq tclsh gettext bc libcurl4-openssl-dev

```
2 安装 git
```shell
cd /tmp
wget --no-check-certificate https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.43.0.tar.gz
tar -xvzf git-2.43.0.tar.gz
cd git-2.43.0
./configure
make
sudo make install
git --version

# 将 git 二进制目录添加到 PATH 路径
tee -a $HOME/.bashrc <<'EOF'
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF

```
3 配置 git

```shell
git config --global user.name "half-coconut"
git config --global user.email "584947559@qq.com"
git config --global credential.helper store
git config --global core.longpaths true


```
### Go 语言开发环境安装和配置

1 编译环境安装和配置
```shell
wget -P /tmp https://go.dev/dl/go1.22.2.linux-amd64.tar.gz

mkdir -p $HOME/go
tar -xvzf /tmp/go1.22.2.linux-amd64.tar.gz -C $HOME/go
mv $HOME/go/go $HOME/go/go1.22.2
```
2 配置 $HOME/.bashrc 文件
```shell
tee -a $HOME/.bashrc <<'EOF'
#Go envs
# Go version settings
export GOVERSION=go1.22.2
# Go installation directory
export GO_INSTALL_DIR=$HOME/go
# GOROOT setting
export GOROOT=$GO_INSTALL_DIR/$GOVERSION
# GOPATH setting
export GOPATH=$HOME/workspace/golang
# Add the binaries from both the Go language and those installed via go install to the PATH
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
# Enable Go modules feature
export GO111MODULE="on"
# Proxy server setting for installing Go modules
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=
# Turn off checking the hash value of Go dependency package
export GOSUMDB=off
EOF




```
### Protobuf 编译环境安装

1 安装 protoc 命令
```shell
cd /tmp
wget https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip

unzip protoc-25.1-linux-x86_64.zip -d protoc-25.1-linux-x86_64
sudo cp protoc-25.1-linux-x86_64/bin/protoc /usr/local/bin/
sudo cp -r protoc-25.1-linux-x86_64/include/ /usr/local/include/

protoc --version

```

2 安装 protoc-gen-go 命令
```shell
go install github.com/golang/protobuf/protoc-gen-go@latest

go install -x github.com/golang/protobuf/protoc-gen-go@latest
```

### Go 开发 IDE 安装和配置
1 安装 Vim9
```shell
git clone https://github.com/vim/vim /tmp/vim
cd /tmp/vim

sudo apt install -y libncurses5-dev
CFLAGS="-I/usr/local/include -fPIC" ./configure --prefix=/usr/local --with-features=huge --enable-cscope --enable-multibyte --enable-rubyinterp --enable-perlinterp --enable-python3interp --enable-luainterp --with-tlib=ncurses --without-local-dir 

make
sudo make install
echo "alias vi=/usr/local/bin/vim" >> ~/.bash_aliases 

bash

```

2 Vim IDE 安装和配置

```shell
rm -rf $HOME/.vim; mkdir -p ~/.vim/pack/plugins/start/

git clone https://github.com/colin404/vim-go ~./vim/pack/plugins/start/vim-go
git clone https://github.com/VundleVim/Vundle.vim.git ~/.vim/bundle/Vundle.vim
git clone --depth=1 https://github.com/colin404/vimrc.git ~/.vim_runtime
sh ~/.vim_runtime/install_awesome_vimrc.sh
git clone https://github.com/onexstack/vimrc /tmp/vimrc
cp /tmp/vimrc/vimrc $HOME/.vim_runtime/my_configs.vim


vi /tmp/test.go
:GoInstallBinaries


```






