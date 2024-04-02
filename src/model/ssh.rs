use super::config::TomlSshServer;

#[derive(Debug)]
pub struct SshServer {
    pub key: String,
    pub host: String,
    pub username: String,
}

impl SshServer {
    pub fn new(key: String, toml_server: TomlSshServer) -> Self {
        Self {
            key,
            host: toml_server.host,
            username: toml_server.username,
        }
    }
}
