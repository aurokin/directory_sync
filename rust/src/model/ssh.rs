use super::{config::TomlSshServer, folder::Folder};

#[derive(Debug)]
pub struct SshServer {
    pub key: String,
    pub host: String,
    pub username: String,
    pub port: u32,
    pub work_folder: Folder,
}

impl SshServer {
    pub fn new(key: String, toml_server: TomlSshServer, work_folder: Folder) -> Self {
        let port = match toml_server.port {
            None => 22,
            Some(port) => port.parse::<u32>().unwrap(),
        };
        Self {
            key,
            host: toml_server.host,
            username: toml_server.username,
            port,
            work_folder,
        }
    }
}
