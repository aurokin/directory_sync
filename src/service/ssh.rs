use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use std::collections::HashMap;
use std::process::{Command, Stdio};

fn build_path(folder: &Folder, relative_path: &Option<String>) -> String {
    if let Some(relative_path) = relative_path {
        return format!("{}/{}", folder.path.clone(), relative_path);
    }

    return folder.path.clone();
}

fn get(name: String, ssh_servers: &HashMap<String, SshServer>) -> Option<&SshServer> {
    for ssh_server in ssh_servers {
        if name == *ssh_server.0 {
            let ssh_server = ssh_server.1;
            return Some(ssh_server);
        }
    }

    return None;
}

fn ssh_cmd(folder: &Folder, ssh_servers: &HashMap<String, SshServer>) -> String {
    let ssh_key = folder.ssh_key.clone().expect("No SSH Key Found");
    let ssh_server = get(ssh_key, ssh_servers).expect("No SSH Server Found");
    let ssh_connection_str = format!("{}@{}", ssh_server.username, ssh_server.host);
    return ssh_connection_str;
}

pub fn ls(
    folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    relative_path: &Option<String>,
) -> () {
    let path = build_path(folder, relative_path);
    println!("SSH - {:?} - {:?}", folder, path);

    let mut cmd_args: Vec<String> = Vec::new();
    match folder.target {
        FolderType::Ssh => {
            let ssh_cmd = ssh_cmd(folder, ssh_servers);
            cmd_args.push("ssh".to_string());
            cmd_args.push(ssh_cmd);
        }
        _ => {}
    }
    cmd_args.push("ls".to_string());
    cmd_args.push("-l".to_string());
    cmd_args.push(path.clone());
    let ls_cmd = Command::new(cmd_args)
        .stdout(Stdio::piped())
        .output()
        .expect("Failed to LS");

    let ls_cmd = String::from_utf8(ls_cmd.stdout).expect("Error converting Stdout");
    println!("{}", ls_cmd);
}
pub fn sync(
    from_folder: &Folder,
    to_folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    relative_path: &Option<String>,
) -> () {
    let from_path = build_path(from_folder, relative_path);
    let to_path = build_path(to_folder, relative_path);
    println!("{:?} {:?}", from_path, to_path);
}
