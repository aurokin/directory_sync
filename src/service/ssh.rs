use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use std::collections::HashMap;
use std::process::{Command, Stdio};

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

            let ssh_output = Command::new("ssh")
                .arg(format!("{}@{}", ssh_server.username, ssh_server.host))
                .arg("ls")
                .arg("-l")
                .arg(folder.path)
                .stdout(Stdio::piped())
                .output()
                .expect("Failed to LS");
            let ssh_output = String::from_utf8(ssh_output.stdout).expect("Error converting Stdout");
            print!("{}", ssh_output);
        }
        FolderType::Local => {
            let ls_output = Command::new("ls")
                .arg("-l")
                .arg(folder.path)
                .stdout(Stdio::piped())
                .output()
                .expect("Failed to LS");

            let ls_output = String::from_utf8(ls_output.stdout).expect("Error converting Stdout");

            print!("{}", ls_output);
        }
    }
    // Command::new("ssh").arg(format!("{}", folder.name))
}
