```bash
sudo pacman -S --needed - < pkglist
yay -S --needed - < aurpkglist
```
```bash
sed -i "s/enp3s0/interface/" suckless/slstatus-1.0/config.h zshrc
sed -i "s/roe/user/" suckless/slock-1.5/config.h
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
