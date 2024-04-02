mod model;
use home::home_dir;
use model::config::TomlConfig;
use model::folder::Folder;
use model::folder::FolderType;
use model::ssh::SshServer;
use std::collections::HashMap;
use std::fs;
use std::option::Option;
use std::path::Path;

fn main() {
    let config = read_config().expect("Error reading config");
    let config: TomlConfig = toml::from_str(config.as_str()).expect("Error parsing config");

    let mut ssh_servers: HashMap<String, SshServer> = HashMap::new();
    for toml_ssh_server in config.ssh {
        let ssh_server = SshServer::new(toml_ssh_server.0, toml_ssh_server.1);
        ssh_servers.insert(ssh_server.key.clone(), ssh_server);
    }
    for ssh_server_map in ssh_servers {
        let ssh_server = ssh_server_map.1;
        println!("{:?}", ssh_server);
    }

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
