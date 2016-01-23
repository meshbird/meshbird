Vagrant.configure(2) do |config|
    config.vm.box = "ubuntu/trusty64"
    config.vm.box_check_update = false

    config.vm.provision "shell", inline: <<-SHELL
        sudo apt-get update
        sudo apt-get install git hg
        curl -O https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
        tar -xvf go1.5.3.linux-amd64.tar.gz
        sudo mv go /usr/local
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
        echo "export GOROOT=\$HOME/go" >> ~/.profile
        echo "export PATH=\$PATH:$GOROOT/bin" >> ~/.profile
        mkdir -p $HOME/go/src
        mkdir -p $HOME/go/pkg
        mkdir -p $HOME/go/bin
    SHELL

    config.vm.define "node01" do |node01|
        node01.vm.network "private_network", ip: "192.168.33.10"
    end

    config.vm.define "node02" do |node02|
        node02.vm.network "private_network", ip: "192.168.33.11"
    end
end