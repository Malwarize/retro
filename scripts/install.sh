#!/bin/bash

install_yt-dlp() {
    echo "Installing yt-dlp"
    sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -o /usr/local/bin/yt-dlp
    sudo chmod a+rx /usr/local/bin/yt-dlp
}
install_ffmpeg() { echo "Installing ffmpeg"
    if command -v apt > /dev/null; then
        sudo apt install -y ffmpeg
    elif command -v dnf > /dev/null; then
        sudo dnf install -y ffmpeg
    elif command -v pacman > /dev/null
    then
        sudo pacman -S ffmpeg
    else
        echo "Could not install ffmpeg. Please install it manually."
        exit 1
    fi
}

get_goplayer() {
    echo "Downloading goplayer"
    cp ./bin/goplayer ~/.local/bin/goplayer
    chmod +x ~/.local/bin/goplayer
    sudo rm -f /usr/local/bin/goplayer
    sudo ln -s ~/.local/bin/goplayer /usr/local/bin/goplayer
}

get_goplay() {
    echo "Downloading goplay"
    cp ./bin/goplay ~/.local/bin/goplay
    chmod +x ~/.local/bin/goplay
    sudo rm -f /usr/local/bin/goplay
    sudo ln -s ~/.local/bin/goplay /usr/local/bin/goplay
}

install_service() {
    echo "Installing service"
    sudo cp ./etc/goplay.service /etc/systemd/user/goplay.service
    systemctl --user daemon-reload
    systemctl --user start goplay
}

# Main script starts here

check_dependencies() {
    echo "Checking dependencies"
    for dependency in yt-dlp ffmpeg; do
        command -v $dependency > /dev/null || install_$dependency
    done
}

stop_service() {
    echo "Stopping service"
    systemctl --user stop goplay
    systemctl --user disable goplay
    sudo rm -f /etc/systemd/user/goplay.service
    systemctl --user daemon-reload
}

install() {
    check_dependencies
    stop_service
    get_goplayer
    get_goplay
    install_service
}

# Run the installation
install

exit 0

