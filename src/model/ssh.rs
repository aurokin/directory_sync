use super::config::TomlSshServer;

#[derive(Debug)]
pub struct SshServer {
    pub key: String,
    pub host: String,
    pub username: String,
    pub port: u32,
}

impl SshServer {
    pub fn new(key: String, toml_server: TomlSshServer) -> Self {
        let port = 22;
        Self {
            key,
            host: toml_server.host,
            username: toml_server.username,
            port,
        }
    }
}
