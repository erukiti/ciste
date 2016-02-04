Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
#  config.vm.network "forwarded_port", guest: 80, host: 8080
#  config.vm.network "forwarded_port", guest: 8500, host: 8500

  # config.vm.network "private_network", ip: "192.168.33.10"

  config.vm.network "public_network", mac: "080027c84905", bridge: "en0: Ethernet"

  config.vm.synced_folder ".", "/ciste"

  config.vm.provider "virtualbox" do |vb|
    # vb.gui = true
  
    vb.memory = "8192"
  end
  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
    sudo apt-get install -y ca-certificates
    sudo apt-get install -y git
    sudo apt-get install -y wget
    sudo wget -qO- https://get.docker.com/ | sh
    sudo update-rc.d docker defaults
    sudo useradd -m -G docker git
  SHELL
end