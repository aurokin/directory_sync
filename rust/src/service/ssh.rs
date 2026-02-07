use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use std::collections::HashMap;

pub fn get(name: String, ssh_servers: &HashMap<String, SshServer>) -> Option<&SshServer> {
    for ssh_server in ssh_servers {
        if name == *ssh_server.0 {
            let ssh_server = ssh_server.1;
            return Some(ssh_server);
        }
    }

    return None;
}

pub fn add_ssh_cmd(
    folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    cmd_args: &mut Vec<String>,
) -> Vec<String> {
    match folder.target {
        FolderType::Ssh => {
            let ssh_cmd = ssh_cmd(folder, ssh_servers);
            let mut full_cmd_args: Vec<String> = Vec::new();
            for cmd in ssh_cmd.iter() {
                full_cmd_args.push(cmd.to_string());
            }
            for cmd in cmd_args.iter() {
                full_cmd_args.push(cmd.to_string());
            }
            full_cmd_args.to_vec()
        }
        FolderType::Local => cmd_args.to_vec(),
    }
}

pub fn ssh_cmd(folder: &Folder, ssh_servers: &HashMap<String, SshServer>) -> Vec<String> {
    let mut ssh_args: Vec<String> = Vec::new();
    let ssh_key = folder.ssh_key.clone().expect("No SSH Key Found");
    let ssh_server = get(ssh_key, ssh_servers).expect("No SSH Server Found");
    let ssh_connection_str = format!("{}@{}", ssh_server.username, ssh_server.host);

    ssh_args.push("ssh".to_string());
    ssh_args.push("-p".to_string());
    ssh_args.push(ssh_server.port.to_string());
    ssh_args.push(ssh_connection_str);
    return ssh_args;
}

pub fn scp_cmd(
    from_folder: &Folder,
    to_folder: &Folder,
    tar_name: String,
    from_work_folder: &Folder,
    to_work_folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
) -> Vec<String> {
    let mut scp_args: Vec<String> = Vec::new();
    let mut port: u32 = 22;

    let from_path = match from_folder.target {
        FolderType::Ssh => {
            let ssh_key = from_folder.ssh_key.clone().expect("No SSH Key Found");
            let ssh_server = get(ssh_key, ssh_servers).expect("No SSH Server Found");
            port = ssh_server.port;
            format!(
                "{}@{}:{}/{}",
                ssh_server.username,
                ssh_server.host,
                from_work_folder.path,
                tar_name.clone(),
            )
        }
        FolderType::Local => format!("{}/{}", from_work_folder.path, tar_name.clone()),
    };

    let to_path = match to_folder.target {
        FolderType::Ssh => {
            let ssh_key = to_folder.ssh_key.clone().expect("No SSH Key Found");
            let ssh_server = get(ssh_key, ssh_servers).expect("No SSH Server Found");
            port = ssh_server.port;
            format!(
                "{}@{}:{}/{}",
                ssh_server.username,
                ssh_server.host,
                to_work_folder.path,
                tar_name.clone(),
            )
        }
        FolderType::Local => format!("{}/{}", to_work_folder.path, tar_name.clone()),
    };

    scp_args.push("scp".to_string());
    scp_args.push("-P".to_string());
    scp_args.push(port.to_string());
    scp_args.push("-r".to_string());
    scp_args.push(from_path);
    scp_args.push(to_path);
    return scp_args;
}
