Vagrant.configure(2) do |config|
    config.vm.box = "ubuntu/trusty64"
    config.vm.box_check_update = false

    config.vm.provision "shell", path: "provision.sh"

    config.vm.define "node01" do |node01|
        node01.vm.network "public_network"
    end

    config.vm.define "node02" do |node02|
        node02.vm.network "public_network"
    end
end