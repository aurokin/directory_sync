use serde::Deserialize;
use std::collections::HashMap;

#[derive(Deserialize, Debug)]
pub struct TomlConfig {
    pub local_work_dir: String,
    pub folders: HashMap<String, TomlFolder>,
    pub links: HashMap<String, TomlLink>,
    pub ssh: HashMap<String, TomlSshServer>,
}

#[derive(Deserialize, Debug)]
pub struct TomlFolder {
    pub path: String,
    pub target: TomlType,
    pub ssh_key: Option<String>,
}

#[derive(Deserialize, Debug)]
pub struct TomlLink {
    pub local: String,
    pub target: String,
    pub paths: Vec<String>,
    pub partial_only: String,
}

#[derive(Deserialize, Debug)]
pub struct TomlSshServer {
    pub host: String,
    pub username: String,
    pub work_dir: String,
    pub port: Option<String>,
}

#[derive(Deserialize, Debug)]
pub enum TomlType {
    #[serde(alias = "local")]
    Local,
    #[serde(alias = "ssh")]
    Ssh,
}
