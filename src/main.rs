use home::home_dir;
use serde::Deserialize;
use std::collections::HashMap;
use std::fs;
use std::path::Path;

#[derive(Deserialize, Debug)]
struct TomlConfig {
    folders: HashMap<String, TomlFolder>,
}

#[derive(Deserialize, Debug)]
struct TomlFolder {
    path: String,
}

struct Folder {
    name: String,
    path: String,
}

fn main() {
    let config = read_config().expect("Error reading config");
    let config: TomlConfig = toml::from_str(config.as_str()).expect("Error parsing config");

    let mut folders: HashMap<String, Folder> = HashMap::new();
    for toml_folder in config.folders {
        let folder = Folder {
            name: toml_folder.0,
            path: toml_folder.1.path,
        };
        folders.insert(folder.name.clone(), folder);
    }

    for folder_map in folders {
        let folder = folder_map.1;
        println!("{}: {}", folder.name, folder.path);
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
