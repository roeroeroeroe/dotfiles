#!/usr/bin/env bash

[ $EUID -ne 0 ] && SUDO="sudo" || SUDO=""

DNSTT_USER="dnstt"
DNSTT_USER_HOME="/home/$DNSTT_USER"
SERVICE_FILE="/etc/systemd/system/dnstt.service"

if ! id "$DNSTT_USER" &> /dev/null; then
	echo "creating user $DNSTT_USER"
	$SUDO useradd -ms /bin/bash "$DNSTT_USER"
fi

$SUDO apt update
$SUDO apt install -y git golang iptables netfilter-persistent

sudo -u $DNSTT_USER git clone https://www.bamsoftware.com/git/dnstt.git "$DNSTT_USER_HOME/dnstt"
cd "$DNSTT_USER_HOME/dnstt/dnstt-server" || exit
sudo -u $DNSTT_USER go build

sudo -u $DNSTT_USER ./dnstt-server -gen-key \
	-privkey-file "$DNSTT_USER_HOME/dnstt/dnstt-server/server.key" \
	-pubkey-file "$DNSTT_USER_HOME/dnstt/dnstt-server/server.pub"

INTERFACE=$(ip r | awk '/default/ {print $5}' | head -n 1)
read -rp "detected primary network interface: $INTERFACE, is this correct? (Y/n): " CONFIRM
[[ "$CONFIRM" =~ ^(n|N)$ ]] && read -rp "interface: " INTERFACE
unset CONFIRM

$SUDO iptables -I INPUT -p udp --dport 53 -j ACCEPT
$SUDO iptables -I INPUT -p udp --dport 5300 -j ACCEPT
$SUDO iptables -t nat -I PREROUTING -i "$INTERFACE" -p udp --dport 53 -j REDIRECT --to-port 5300
$SUDO ip6tables -I INPUT -p udp --dport 53 -j ACCEPT
$SUDO ip6tables -I INPUT -p udp --dport 5300 -j ACCEPT
$SUDO ip6tables -t nat -I PREROUTING -i "$INTERFACE" -p udp --dport 53 -j REDIRECT --to-port 5300
$SUDO netfilter-persistent save

read -rp "domain (e.g., X.EXAMPLE.COM): " DOMAIN

SSH_PORT=$(awk '/^Port / {print $2}' /etc/ssh/sshd_config | head -n 1)
SSH_PORT=${SSH_PORT:-22}
read -rp "detected SSH port: $SSH_PORT, is this correct? (Y/n): " CONFIRM
[[ "$CONFIRM" =~ ^(n|N)$ ]] && read -rp "SSH port: " SSH_PORT

$SUDO cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=dnstt-server
After=network.target

[Service]
User=$DNSTT_USER
WorkingDirectory=$DNSTT_USER_HOME/dnstt/dnstt-server
ExecStart=$DNSTT_USER_HOME/dnstt/dnstt-server/dnstt-server -udp :5300 -privkey-file $DNSTT_USER_HOME/dnstt/dnstt-server/server.key $DOMAIN 127.0.0.1:$SSH_PORT
Restart=always

[Install]
WantedBy=multi-user.target
EOF

$SUDO systemctl daemon-reload
$SUDO systemctl enable --now dnstt

echo "
copy $DNSTT_USER_HOME/dnstt/dnstt-server/server.pub to client
add your ssh key to $DNSTT_USER_HOME/.ssh/authorized_keys and run:
(client)
./dnstt-client -udp 1.1.1.1:53 -pubkey-file server.pub $DOMAIN 127.0.0.1:8000
ssh -ND 127.0.0.1:1080 -o HostKeyAlias=SERVER_IP -p 8000 $DNSTT_USER@127.0.0.1

test connectivity:
curl -x socks5://127.0.0.1:1080 https://api.ipify.org"
