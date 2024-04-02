// folder.rs

use crate::model::folder::Folder;
use std::collections::HashMap;

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
