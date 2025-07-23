```bash
sudo pacman -S --needed - < pkglist
yay -S --needed - < aurpkglist
```
```bash
sed -i "s/enp3s0/interface/" suckless/slstatus-1.0/config.h
```
## zsh
```bash
mkdir -p ~/.config/zsh/plugins
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.config/zsh/plugins/zsh-syntax-highlighting
git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.config/zsh/plugins/zsh-autosuggestions
```
