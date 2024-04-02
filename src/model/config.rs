use serde::Deserialize;
use std::collections::HashMap;

#[derive(Deserialize, Debug)]
pub struct TomlConfig {
    pub folders: HashMap<String, TomlFolder>,
    pub ssh: HashMap<String, TomlSshServer>,
}

#[derive(Deserialize, Debug)]
pub struct TomlFolder {
    pub path: String,
    pub target: TomlType,
    pub ssh_key: Option<String>,
}

#[derive(Deserialize, Debug)]
pub struct TomlSshServer {
    pub host: String,
    pub name: String,
}

#[derive(Deserialize, Debug)]
pub enum TomlType {
    #[serde(alias = "local")]
    Local,
    #[serde(alias = "ssh")]
    Ssh,
}
