use home::home_dir;
use serde::Deserialize;
use std::collections::HashMap;
use std::fs;
use std::option::Option;
use std::path::Path;

#[derive(Deserialize, Debug)]
struct TomlConfig {
    folders: HashMap<String, TomlFolder>,
    ssh: HashMap<String, TomlSshServer>,
}

#[derive(Deserialize, Debug)]
struct TomlFolder {
    path: String,
    target: TomlType,
    ssh_key: Option<(String)>,
}

#[derive(Deserialize, Debug)]
struct TomlSshServer {
    host: String,
    name: String,
}

#[derive(Deserialize, Debug)]
enum TomlType {
    #[serde(alias = "local")]
    Local,
    #[serde(alias = "ssh")]
    Ssh,
}

#[derive(Debug)]
enum FolderType {
    Local,
    Ssh,
}
impl FolderType {
    fn get_folder_type(toml_type: TomlType) -> Self {
        let folder_type = match toml_type {
            TomlType::Local => FolderType::Local,
            TomlType::Ssh => FolderType::Ssh,
        };

        return folder_type;
    }
}
#[derive(Debug)]
struct Folder {
    name: String,
    path: String,
    target: FolderType,
}

fn main() {
    let config = read_config().expect("Error reading config");
    let config: TomlConfig = toml::from_str(config.as_str()).expect("Error parsing config");

    let mut folders: HashMap<String, Folder> = HashMap::new();
    for toml_folder in config.folders {
        let folder = Folder {
            name: toml_folder.0,
            path: toml_folder.1.path,
            target: FolderType::get_folder_type(toml_folder.1.target),
        };
        folders.insert(folder.name.clone(), folder);
    }

    for folder_map in folders {
        let folder = folder_map.1;
        println!("{:?}: {:?}, {:?}", folder.name, folder.path, folder.target);
    }
}

fn read_config() -> Option<String> {
    let home_dir = home_dir().expect("No Home Dir");
    let path = Path::new(&home_dir.as_os_str()).join(".dirsync.toml");
    let path_exists = path.exists();
    let path_is_file = path.is_file();

    if path_exists && path_is_file {
        let file = fs::read_to_string(path).expect("Error reading file");
        return Some(file);
    }

    return None;
}
