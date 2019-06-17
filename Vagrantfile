Vagrant.configure("2") do |config|
  (1..3).each do |i|
    config.vm.define "cassandra#{i}" do |node|
      node.vm.box = "cassandra"
      node.vm.box_url = "cassandra#{i}.example.com"
      node.vm.hostname = "cassandra#{i}.example.com"
      node.vm.network :private_network, ip: "192.168.50.4#{i}"
      node.ssh.insert_key = false
 
      # Replace the "listen_address" line in the conf file.
      node.vm.provision "shell", inline: "sed -i 's/^cluster_name: .*/cluster_name: \"My Cluster\"/' ~/apache-cassandra-3.11.4/conf/cassandra.yaml", privileged: false
      node.vm.provision "shell", inline: "sed -i 's/- seeds: .*/- seeds: \"192.168.50.41\"/' ~/apache-cassandra-3.11.4/conf/cassandra.yaml", privileged: false
      node.vm.provision "shell", inline: "sed -i 's/listen_address: .*/listen_address: \"192.168.50.4#{i}\"/' ~/apache-cassandra-3.11.4/conf/cassandra.yaml", privileged: false
    end
  end
end
