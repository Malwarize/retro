.PHONY: build run clean

all: clean build run

build: 
	go build -o bin/goplay client/main.go
	go build -o bin/goplayer server/main.go

s:
	go run server/main.go

run:
	./build/goplay

clean:
	rm -rf bin/
.PHONY: install check_dependencies get_goplayer get_goplay install_service stop_service

install: check_dependencies stop_service get_goplayer get_goplay install_service

check_dependencies:
	@echo "Checking dependencies"
	@for dependency in yt-dlp ffmpeg; do \
		command -v $$dependency > /dev/null || make install_$$dependency; \
	done

install_yt-dlp:
	@echo "Installing yt-dlp"
	@sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -o /usr/local/bin/yt-dlp
	@sudo chmod a+rx /usr/local/bin/yt-dlp

install_ffmpeg:
	@echo "Installing ffmpeg"
	@command -v apt > /dev/null && sudo apt install ffmpeg -y || sudo pacman -S ffmpeg --noconfirm || sudo dnf install ffmpeg -y || echo "Please install ffmpeg manually"

get_goplayer: build
	@echo "Downloading goplayer"
	@cp ./bin/goplayer ~/.local/bin/goplayer
	@chmod +x ~/.local/bin/goplayer
	@sudo rm -f /usr/local/bin/goplayer
	@sudo ln -s ~/.local/bin/goplayer /usr/local/bin/goplayer

get_goplay: build
	@echo "Downloading goplay"
	@cp ./bin/goplay ~/.local/bin/goplay
	@chmod +x ~/.local/bin/goplay
	@sudo rm -f /usr/local/bin/goplay
	@sudo ln -s ~/.local/bin/goplay /usr/local/bin/goplay

install_service:
	@echo "Installing service"
	@sudo cp ./etc/goplay.service /etc/systemd/user/goplay.service
	@systemctl --user daemon-reload
	@systemctl --user start goplay

stop_service:
	@systemctl --user stop goplay
	@systemctl --user disable goplay
	@systemctl --user daemon-reload

