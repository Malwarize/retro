#!/bin/bash

# Parameters
retro_binary=./bin/retro
retroPlayer_binary=./bin/retroPlayer
service_file=./etc/retro_player.service
install_path="~/.local/bin"
systemd_user_path=/etc/systemd/user
installer_path=./bin/install.sh

# Function to build binaries
function build_binary {
	local binary_name=$1
	local binary_source=$2
	go build -o $binary_name $binary_source

	if [ $? -eq 0 ]; then
		echo "Built $binary_name successfully"
	else
		echo "Failed to build $binary_name"
		exit 1
	fi
}

# Function to generate installer
function generate_installer {
	cat <<EOF >$installer_path
#!/bin/bash

retro_binary_data="$(base64 $retro_binary)"
retroPlayer_binary_data="$(base64 $retroPlayer_binary)"
service_data="$(base64 $service_file)"

install_yt-dlp() {
    installed=0
    echo "Installing yt-dlp"
    sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -o /usr/local/bin/yt-dlp && installed=1
    sudo chmod a+rx /usr/local/bin/yt-dlp && installed=1

    if [ \$installed -eq 1 ]; then
      echo "yt-dlp installed successfully"
    else
      echo "Failed to install yt-dlp, please install it manually"
      kill -2 $$
    fi
}

install_ffmpeg() {
    echo "Installing ffmpeg"
    installed=0
    if command -v apt > /dev/null; then
        sudo apt install -y ffmpeg && installed=1
    elif command -v dnf > /dev/null; then
        sudo dnf install -y ffmpeg && installed=1
    elif command -v pacman > /dev/null 
    then
        sudo pacman -S ffmpeg && installed=1
    else
        echo "Could not install ffmpeg. Please install it manually."
        kill -2 $$
    fi
    if [ \$installed -eq 1 ]; then
      echo "ffmpeg installed successfully"
    else
      echo "Failed to install ffmpeg, please install it manually"
      kill -2 $$
    fi
}

install_bash-completion() {
    # check if bash-completion is installed

    if [ -f /usr/share/bash-completion/bash_completion ]; then
        echo "bash-completion already installed"
        return
    fi

    echo "Installing bash-completion"
    if command -v apt > /dev/null; then
        sudo apt install -y bash-completion
    elif command -v dnf > /dev/null; then
        sudo dnf install -y bash-completion
    elif command -v pacman > /dev/null
    then
        sudo pacman -S bash-completion
    else
        echo "Could not install bash-completion. Please install it manually."
        kill -2 $$
    fi
}

check_dependencies() {
    echo "Checking dependencies"
    for dependency in yt-dlp ffmpeg bash-completion;do
        command -v \$dependency > /dev/null || install_\$dependency
    done
}

# stopping services
function cleanup {
    echo "Cleaning up"
    echo "Disabling and stopping retro service..."
    if systemctl --user is-active --quiet retro; then
        systemctl --user stop retro
    fi
    # check if service is enabled
    if systemctl --user is-enabled --quiet retro; then
        systemctl --user disable retro
    fi
    echo "Removing files..."
    if [ -f $systemd_user_path/retro.service ]; then
        sudo rm -rf $systemd_user_path/retro.service  # Remove the old service file
    fi
    if [ -f $install_path/retro ]; then
        sudo rm -rf $install_path/retro > /dev/null
    fi
    if [ -f $install_path/retroPlayer ]; then
        sudo rm -rf $install_path/retroPlayer > /dev/null
    fi
    # remove links 
    if [ -L /usr/local/bin/retro ]; then
        sudo rm -rf /usr/local/bin/retro
    fi
    if [ -L /usr/local/bin/retroPlayer ]; then
        sudo rm -rf /usr/local/bin/retroPlayer
    fi
    systemctl --user daemon-reload 
    echo "Cleanup done"
}


function install_retro {
    echo "Installing retro to $install_path/retro"
    mkdir -p $install_path
    echo "\$retro_binary_data" | base64 -d > $install_path/retro
    chmod +x $install_path/retro
    sudo ln -s $install_path/retro /usr/local/bin/retro
}

function install_retroPlayer {
    echo "Installing retroPlayer to $install_path/retroPlayer"
    mkdir -p $install_path
    echo "\$retroPlayer_binary_data" | base64 -d > $install_path/retroPlayer
    chmod +x $install_path/retroPlayer
    sudo ln -s $install_path/retroPlayer /usr/local/bin/retroPlayer
}

function start_services {
    echo "Starting retro service"
    systemctl --user daemon-reload
    systemctl --user is-enabled --quiet retro || systemctl --user enable retro > /dev/null
    systemctl --user is-active --quiet retro || systemctl --user start retro > /dev/null
}

function install_service {
    echo "Installing retro service"
    mkdir -p $systemd_user_path
    echo "\$service_data" | base64 -d | sudo tee $systemd_user_path/retro.service > /dev/null
    systemctl --user daemon-reload
}

function generate_completion {
    echo "Generating completion script"

    # check if zsh is installed
    if command -v zsh > /dev/null; then
        echo "Generating zsh completion"
        completion_path=~/.zsh_completion.d
        mkdir -p \$completion_path
        $install_path/retro completion zsh > \$completion_path/_retro

        # check if zshrc exists
        if [ -f ~/.zshrc ]; then
            grep -q "source \$HOME/.zsh_completion.d/_retro" ~/.zshrc || echo "source \$HOME/.zsh_completion.d/_retro" >> ~/.zshrc
        fi
    fi

    # check if bash is installed
    if command -v bash > /dev/null; then
        echo "Generating bash completion"
        completion_path=~/.bash_completion.d
        mkdir -p \$completion_path
        $install_path/retro completion bash > \$completion_path/_retro
        if [ -f ~/.bashrc ]; then
            grep -q "source \$HOME/.bash_completion.d/_retro" ~/.bashrc || echo "source \$HOME/.bash_completion.d/_retro" >> ~/.bashrc
        fi
    fi
    echo "Completion script created, restart your shell to use it"
}
function generate_uninstall {
    echo "Generating uninstall script"
      echo "#!/bin/bash" > $install_path/uninstall_retro.sh
      echo "echo \"Uninstalling retro...\"" >> $install_path/uninstall_retro.sh

      echo "echo \"Disabling and stopping retro service...\"" >> $install_path/uninstall_retro.sh
      echo "systemctl --user is-active --quiet retro && systemctl --user stop retro" >> $install_path/uninstall_retro.sh
      echo "systemctl --user is-enabled --quiet retro && systemctl --user disable retro" >> $install_path/uninstall_retro.sh

      echo "echo \"Removing files...\"" >> $install_path/uninstall_retro.sh
      echo "rm -rf $install_path/retro" >> $install_path/uninstall_retro.sh
      echo "rm -rf $install_path/retroPlayer" >> $install_path/uninstall_retro.sh
      echo "sudo rm -rf $systemd_user_path/retro.service" >> $install_path/uninstall_retro.sh
      echo "rm -rf ~/.zsh_completion.d/_retro" >> $install_path/uninstall_retro.sh
      echo "rm -rf ~/.bash_completion.d/_retro" >> $install_path/uninstall_retro.sh
      echo "rm -rf ~/.retro" >> $install_path/uninstall_retro.sh
      echo "sudo rm -rf /usr/local/bin/retro" >> $install_path/uninstall_retro.sh
      echo "sudo rm -rf /usr/local/bin/retroPlayer" >> $install_path/uninstall_retro.sh

      echo "echo \"Reloading daemon...\"" >> $install_path/uninstall_retro.sh
      echo "systemctl --user daemon-reload" >> $install_path/uninstall_retro.sh
      echo "echo \"Uninstall complete\"" >> $install_path/uninstall_retro.sh
      chmod +x $install_path/uninstall_retro.sh

      echo "Removing retro from RC files..."
      echo "[ -f ~/.config/fish/config.fish ] && sed -i '/retro/d' ~/.config/fish/config.fish" >> $install_path/uninstall_retro.sh
      echo "[ -f ~/.bashrc ] && sed -i '/retro/d' ~/.bashrc" >> $install_path/uninstall_retro.sh
      echo "[ -f ~/.zshrc ] && sed -i '/retro/d' ~/.zshrc" >> $install_path/uninstall_retro.sh

      echo "echo \"Removing uninstall script...\"" >> $install_path/uninstall_retro.sh
      echo "rm -rf $install_path/uninstall_retro.sh" >> $install_path/uninstall_retro.sh

    echo "Uninstall script created: $install_path/uninstall_retro.sh"
}
function main {
    trap cleanup SIGINT
    cleanup
    check_dependencies
    install_retro
    install_retroPlayer
    install_service
    start_services
    generate_completion
    generate_uninstall
    echo "Installation complete, use retro --help to get started"
}

main
EOF

	chmod +x $installer_path
	echo "Installer script created: $installer_path"
}

function main {
	echo "Building retro and retroPlayer"
	build_binary "$retro_binary" "client/main.go"
	build_binary "$retroPlayer_binary" "server/main.go"

	echo "Generating installer"
	generate_installer
}

main
