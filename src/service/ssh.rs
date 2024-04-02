use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use std::collections::HashMap;
use std::process::Command;

fn get(name: String, ssh_servers: HashMap<String, SshServer>) -> Option<SshServer> {
    for ssh_server in ssh_servers {
        if name == ssh_server.0 {
            let ssh_server = ssh_server.1;
            return Some(ssh_server);
        }
    }

    return None;
}

pub fn ls(folder: Folder, ssh_servers: HashMap<String, SshServer>) -> () {
    println!("SSH - {:?}", folder);

    match folder.target {
        FolderType::Ssh => {
            let ssh_key = folder.ssh_key.expect("No SSH Key Found");
            let ssh_server = get(ssh_key, ssh_servers).expect("No SSH Server Found");
            println!("Found SSH Server - {:?}", ssh_server);
        }
        FolderType::Local => {
            Command::new("ls -l")
                .arg(folder.path)
                .output()
                .expect("Failed to LS");
        }
    }
    // Command::new("ssh").arg(format!("{}", folder.name))
}
