package swarm

type NodeConfig struct {
    Node struct {
        ID string `yaml:"id"`
        Domain string `yaml:"domain"`
        ListenAddr string `yaml:"listen_addr"`
        Role string `yaml:"role"`
        Region string `yaml:"region"`
        PrivateKeyFile string `yaml:"private_key_file"`
    } `yaml:"node"`
}
