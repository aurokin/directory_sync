use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use std::process::{Command, Stdio};
use std::{collections::HashMap, io};

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
    let first_cmd = cmd_args.first().expect("First argument required");
    let mut ls_cmd = Command::new(first_cmd);
    for cmd_arg in &cmd_args[1..] {
        ls_cmd.arg(cmd_arg);
    }

    let ls_output = ls_cmd
        .stdout(Stdio::piped())
        .output()
        .expect("Failed to LS");
    let ls_output = String::from_utf8(ls_output.stdout).expect("Error converting Stdout");
    println!("{}", ls_output);
}
pub fn sync(
    from_folder: &Folder,
    to_folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    relative_path: &Option<String>,
) -> () {
    // preview argument, to help build prompts
    // add ssh connection checks
    let from_path = build_path(from_folder, relative_path);
    let to_path = build_path(to_folder, relative_path);
    println!("Sync: {:?} - {:?}", from_path, to_path);

    let mut check_if_from_folders_exist: Vec<String> = Vec::new();
    let mut create_empty_to_folders: Vec<String> = Vec::new();
    let mut remove_to_folders: Vec<String> = Vec::new();
    let mut copy_to_folder: Vec<String> = Vec::new();

    let is_from_ssh = match from_folder.target {
        FolderType::Ssh => true,
        _ => false,
    };
    let is_to_ssh = match to_folder.target {
        FolderType::Ssh => true,
        _ => false,
    };

    if is_from_ssh {
        let ssh_cmd = ssh_cmd(from_folder, ssh_servers);
        check_if_from_folders_exist.push("ssh".to_string());
        check_if_from_folders_exist.push(ssh_cmd);
    }
    check_if_from_folders_exist.push("ls".to_string());
    check_if_from_folders_exist.push(from_path.clone());

    if is_to_ssh {
        let ssh_cmd = ssh_cmd(to_folder, ssh_servers);
        create_empty_to_folders.push("ssh".to_string());
        create_empty_to_folders.push(ssh_cmd);
    }
    create_empty_to_folders.push("mkdir".to_string());
    create_empty_to_folders.push("-p".to_string());
    create_empty_to_folders.push(to_path.clone());

    if is_to_ssh {
        let ssh_cmd = ssh_cmd(to_folder, ssh_servers);
        remove_to_folders.push("ssh".to_string());
        remove_to_folders.push(ssh_cmd);
    }
    remove_to_folders.push("rm".to_string());
    remove_to_folders.push("-rf".to_string());
    remove_to_folders.push(to_path.clone());

    if is_from_ssh || is_to_ssh {
        copy_to_folder.push("scp".to_string());
        copy_to_folder.push("-r".to_string());
    } else {
        copy_to_folder.push("cp".to_string());
        copy_to_folder.push("-R".to_string());
    }
    if is_from_ssh {
        let ssh_cmd = ssh_cmd(from_folder, ssh_servers);
        let scp_path = format!("{}:'{}'", ssh_cmd, from_path.clone());
        copy_to_folder.push(scp_path);
    } else {
        copy_to_folder.push(from_path.clone());
    }
    if is_to_ssh {
        let ssh_cmd = ssh_cmd(to_folder, ssh_servers);
        let scp_path = format!("{}:'{}'", ssh_cmd, from_path.clone());
        copy_to_folder.push(scp_path);
    } else {
        copy_to_folder.push(to_path.clone());
    }

    let check_folder_arg = check_if_from_folders_exist
        .first()
        .expect("First argument required");
    let mut check_folder_cmd = Command::new(check_folder_arg);
    for folder_arg in &check_if_from_folders_exist[1..] {
        check_folder_cmd.arg(folder_arg);
    }
    let check_folder_output = check_folder_cmd
        .stdout(Stdio::piped())
        .output()
        .expect("Failed to Check Folder");
    let check_folder_output =
        String::from_utf8(check_folder_output.stdout).expect("Error converting Stdout");
    if check_folder_output.contains(&"No such".to_string()) {
        println!("Error: From Folder Does Not Exist");
        return;
    }
    println!("Ready for transfer, would you like to continue? The following commands will run");
    println!("- {}", create_empty_to_folders.join(" "));
    println!("- {}", remove_to_folders.join(" "));
    println!("- {}", copy_to_folder.join(" "));
    println!("Enter y to continue!");

    let mut user_run_input = String::from("");
    io::stdin()
        .read_line(&mut user_run_input)
        .expect("Failed to read line");
    let user_run_input = user_run_input.trim().to_string();
    if user_run_input != "y".to_string() {
        println!("Skipping this folder because the user did not input 'y'");
        return;
    }
}
