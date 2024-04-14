use crate::model::folder::{Folder, FolderType};
use std::{collections::HashMap, env};

pub fn get(name: String, folders: HashMap<String, Folder>) -> Option<Folder> {
    let mut found = None;
    for folder in folders {
        if folder.0 == name {
            found = Some(folder.1);
            break;
        }
    }
    return found;
}

pub fn get_current_dir() -> Option<Folder> {
    let current_dir = env::current_dir();
    match current_dir {
        Ok(dir) => {
            return Some(Folder {
                name: "current_working_directory".to_string(),
                path: dir
                    .into_os_string()
                    .into_string()
                    .expect("Error converting current working dir"),
                target: FolderType::Local,
                ssh_key: None,
            });
        }
        Err(_) => return None,
    };
}
