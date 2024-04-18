use crate::model::config::TomlConfig;
use crate::model::folder::Folder;
use crate::model::folder::FolderType;
use crate::model::link::Link;
use crate::model::ssh::SshServer;
use home::home_dir;
use std::collections::HashMap;
use std::fs;
use std::option::Option;
use std::path::Path;

pub fn read_config() -> Option<String> {
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

pub fn parse_config(
    config: String,
) -> (
    HashMap<String, SshServer>,
    HashMap<String, Folder>,
    HashMap<String, Link>,
    Folder,
) {
    let config: TomlConfig = toml::from_str(config.as_str()).expect("Error parsing config");

    let mut ssh_servers: HashMap<String, SshServer> = HashMap::new();
    for toml_ssh_server in config.ssh {
        let ssh_server = SshServer::new(toml_ssh_server.0, toml_ssh_server.1);
        ssh_servers.insert(ssh_server.key.clone(), ssh_server);
    }

    let mut folders: HashMap<String, Folder> = HashMap::new();
    for toml_folder in config.folders {
        let folder = Folder {
            name: toml_folder.0,
            path: toml_folder.1.path,
            target: FolderType::get_folder_type(toml_folder.1.target),
            ssh_key: toml_folder.1.ssh_key,
        };
        folders.insert(folder.name.clone(), folder);
    }
    let work_folder = folders
        .get(&config.local_work_dir)
        .expect("Local work folder required")
        .clone();

    let mut links: HashMap<String, Link> = HashMap::new();
    for toml_link in config.links {
        let local_folder = folders.get(&toml_link.1.local);
        let target_folder = folders.get(&toml_link.1.target);

        if local_folder.is_some() && target_folder.is_some() {
            let link = Link {
                name: toml_link.0,
                local: local_folder.unwrap().clone(),
                target: target_folder.unwrap().clone(),
                paths: toml_link.1.paths,
            };
            links.insert(link.name.clone(), link);
        } else {
            println!("Unable to parse link in configuration");
            continue;
        }
    }
    return (ssh_servers, folders, links, work_folder);
}
