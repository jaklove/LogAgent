tar -xf etcd-v3.3.5-linux-amd64.tar.gz
cd etcd-v3.3.5-linux-amd64
./etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379