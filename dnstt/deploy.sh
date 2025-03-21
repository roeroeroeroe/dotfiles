#!/usr/bin/env bash

[ $EUID -ne 0 ] && SUDO="sudo" || SUDO=""

DNSTT_USER="dnstt"
DNSTT_USER_HOME="/home/$DNSTT_USER"
SERVICE_FILE="/etc/systemd/system/dnstt.service"

command -v sudo &> /dev/null || {
	echo "sudo is not installed" >&2
	exit 1
}

if [ -f /etc/debian_version ]; then
	DISTRIBUTION="debian"
	$SUDO apt update
	$SUDO apt install -y git golang iptables netfilter-persistent
elif [ -f /etc/redhat-release ]; then
	DISTRIBUTION="redhat"
	$SUDO dnf install -y git golang iptables iptables-services
elif [ -f /etc/arch-release ]; then
	DISTRIBUTION="archlinux"
	$SUDO pacman -Sy --needed --noconfirm git go iptables
else
	echo "unsupported distribution" >&2
	exit 1
fi

if ! id "$DNSTT_USER" &> /dev/null; then
	echo "creating user $DNSTT_USER"
	$SUDO useradd -ms /bin/bash "$DNSTT_USER"
fi

[ ! -d "$DNSTT_USER_HOME/dnstt" ] && \
	sudo -u $DNSTT_USER git clone https://www.bamsoftware.com/git/dnstt.git "$DNSTT_USER_HOME/dnstt"

sudo -u $DNSTT_USER bash -c "
	cd '$DNSTT_USER_HOME/dnstt/dnstt-server' && \
		go build && \
		./dnstt-server -gen-key \
			-privkey-file '$DNSTT_USER_HOME/dnstt/dnstt-server/server.key' \
			-pubkey-file '$DNSTT_USER_HOME/dnstt/dnstt-server/server.pub'
"

INTERFACE=$(ip r | awk '/default/ {print $5}' | head -n 1)
read -rp "detected primary network interface: $INTERFACE, is this correct? (Y/n): " CONFIRM
[[ "$CONFIRM" =~ ^(n|N)$ ]] && read -rp "interface: " INTERFACE
unset CONFIRM

for TABLE in "iptables" "ip6tables"; do
	for PORT in 53 5300; do
		$SUDO $TABLE -C INPUT -p udp --dport $PORT -j ACCEPT 2> /dev/null || \
			$SUDO $TABLE -I INPUT -p udp --dport $PORT -j ACCEPT
	done

	$SUDO $TABLE -t nat -C PREROUTING -i "$INTERFACE" -p udp --dport 53 -j REDIRECT --to-port 5300 2> /dev/null || \
		$SUDO $TABLE -t nat -I PREROUTING -i "$INTERFACE" -p udp --dport 53 -j REDIRECT --to-port 5300
done

case "$DISTRIBUTION" in
	debian)
		$SUDO netfilter-persistent save
		;;
	redhat)
		$SUDO iptables-save | $SUDO tee /etc/sysconfig/iptables > /dev/null
		$SUDO ip6tables-save | $SUDO tee /etc/sysconfig/ip6tables > /dev/null
		$SUDO systemctl disable --now firewalld
		$SUDO systemctl enable --now iptables ip6tables
		;;
	archlinux)
		$SUDO iptables-save | $SUDO tee /etc/iptables/iptables.rules > /dev/null
		$SUDO ip6tables-save | $SUDO tee /etc/iptables/ip6tables.rules > /dev/null
		$SUDO systemctl enable --now iptables ip6tables
		;;
esac

read -rp "nameserver domain (e.g., x.example.com): " DOMAIN

SSH_PORT=$($SUDO awk '/^Port / {print $2}' /etc/ssh/sshd_config | head -n 1)
SSH_PORT=${SSH_PORT:-22}
read -rp "detected SSH port: $SSH_PORT, is this correct? (Y/n): " CONFIRM
[[ "$CONFIRM" =~ ^(n|N)$ ]] && read -rp "SSH port: " SSH_PORT

$SUDO tee "$SERVICE_FILE" > /dev/null <<EOF
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

if [ "$DISTRIBUTION" = "redhat" ]; then
	$SUDO semanage fcontext -a -t bin_t "$DNSTT_USER_HOME/dnstt/dnstt-server/dnstt-server"
	$SUDO restorecon -v "$DNSTT_USER_HOME/dnstt/dnstt-server/dnstt-server"
fi

$SUDO systemctl daemon-reload
$SUDO systemctl enable --now dnstt

echo "
copy $DNSTT_USER_HOME/dnstt/dnstt-server/server.pub to client;
add your SSH pubkey to $DNSTT_USER_HOME/.ssh/authorized_keys and run:

(client)
# udp
./dnstt-client -udp 8.8.8.8:53 -pubkey-file server.pub $DOMAIN 127.0.0.1:8000
# doh
./dnstt-client -doh https://cloudflare-dns.com/dns-query -pubkey-file server.pub $DOMAIN 127.0.0.1:8000
ssh -ND 127.0.0.1:1080 -o HostKeyAlias=SERVER_IP -p 8000 $DNSTT_USER@127.0.0.1

test connectivity:
curl -x socks5://127.0.0.1:1080 https://api.ipify.org"
