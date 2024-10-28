```bash
sudo pacman -S --needed dunst feh go gvfs htop imagemagick jq \
	jre-openjdk-headless keepassxc leafpad lf lxappearance \
	mpv ncdu neovim noto-fonts-cjk noto-fonts-emoji pcmanfm \
	playerctl proxychains-ng streamlink tor viewnior yt-dlp \
	zsh cmus nodejs npm fzf dnsmasq nftables ttf-jetbrains-mono-nerd \
	qbittorrent ripgrep unzip maim xclip alsa-utils libnotify
```
```bash
yay -S librewolf-bin simplescreenrecorder vesktop-bin \
	picom-git xcursor-simp1e-gruvbox-dark nekoray \
	sing-geoip-db sing-geosite-db chatterino2-git
```
## zsh
```bash
mkdir -p ~/.config/zsh/plugins
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.config/zsh/plugins/zsh-syntax-highlighting
git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.config/zsh/plugins/zsh-autosuggestions
```
## webtunnel
```bash
git clone https://gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel
cd webtunnel/main/client && go build
sudo cp client /usr/local/bin/webtunnel
```
- [GTK theme](https://github.com/Fausto-Korpsvart/Gruvbox-GTK-Theme)
- [GTK icons](https://github.com/jmattheis/gruvbox-dark-icons-gtk)
- [firefox theme](https://addons.mozilla.org/en-US/firefox/addon/gruvboxgruvboxgruvboxgruvboxgr)
