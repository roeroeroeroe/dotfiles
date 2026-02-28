```bash
sudo pacman -S --needed - < pkglist
yay -S --needed - < aurpkglist
./buildcmus
```
```bash
iface="<INTERFACE>"
block_dev="<BLOCK_DEVICE>"
sed -i -e "s/\"enp3s0\"/\"$iface\"/" -e "s/\"sda\"/\"$block_dev\"/" \
    sb/config/config.go
cd sb && make TAGS=with_pulse && mv -vf sb ~/.local/bin/sb
```
## zsh
```bash
mkdir -p ~/.config/zsh/plugins
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.config/zsh/plugins/zsh-syntax-highlighting
git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.config/zsh/plugins/zsh-autosuggestions
```
