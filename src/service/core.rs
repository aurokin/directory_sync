use crate::model::{
    folder::{Folder, FolderType},
    ssh::SshServer,
};
use crate::service::ssh::{add_ssh_cmd, scp_cmd};
use crate::service::tar::{tar_directory, untar_directory};
use std::process::{Command, Stdio};
use std::{collections::HashMap, io};

fn build_path(folder: &Folder, relative_path: &Option<String>) -> String {
    if let Some(relative_path) = relative_path {
        return format!("{}/{}", folder.path.clone(), relative_path);
    }

    return folder.path.clone();
}

pub fn ls(
    folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    relative_path: &Option<String>,
) -> () {
    let path = build_path(folder, relative_path);
    println!("SSH - {:?} - {:?}", folder, path);

    let mut cmd_args = add_ssh_cmd(folder, ssh_servers, &mut Vec::new());
    cmd_args.push("ls".to_string());
    cmd_args.push("-l".to_string());
    cmd_args.push(path.clone());
    let first_cmd = cmd_args.first().expect("First argument required");
    let mut ls_cmd = Command::new(first_cmd);
    for cmd_arg in &cmd_args[1..] {
        ls_cmd.arg(cmd_arg);
    }

    let ls_output = ls_cmd
        .stdout(Stdio::inherit())
        .output()
        .expect("Failed to LS");
    let ls_output = String::from_utf8(ls_output.stdout).expect("Error converting Stdout");
    println!("{}", ls_output);
}

pub fn sync(
    from_folder: &Folder,
    to_folder: &Folder,
    work_folder: &Folder,
    ssh_servers: &HashMap<String, SshServer>,
    relative_path: &Option<String>,
    force: bool,
) -> () {
    // preview argument, to help build prompts
    // add ssh connection checks
    let from_path = build_path(from_folder, relative_path);
    let to_path = build_path(to_folder, relative_path);
    println!("Sync: {:?} - {:?}", from_path, to_path);

    let is_from_ssh = match from_folder.target {
        FolderType::Ssh => true,
        _ => false,
    };
    let is_to_ssh = match to_folder.target {
        FolderType::Ssh => true,
        _ => false,
    };
    if is_from_ssh && is_to_ssh {
        println!("Only one folder can be remote");
        return;
    }

    let from_work_folder = if is_from_ssh {
        let ssh_server =
            crate::service::ssh::get(from_folder.ssh_key.clone().unwrap(), ssh_servers)
                .expect("SSH Server not found");
        &ssh_server.work_folder
    } else {
        work_folder
    };

    let to_work_folder = if is_to_ssh {
        let ssh_server = crate::service::ssh::get(to_folder.ssh_key.clone().unwrap(), ssh_servers);
        let ssh_server = ssh_server.expect("SSH Server not found");
        &ssh_server.work_folder
    } else {
        work_folder
    };

    let (tar_name, mut tar_source_exists_args, mut create_tar_args, mut delete_from_tar_args) =
        tar_directory(from_path.clone(), from_work_folder);
    let tar_source_exists_args = add_ssh_cmd(from_folder, ssh_servers, &mut tar_source_exists_args);
    let create_tar_args = add_ssh_cmd(from_folder, ssh_servers, &mut create_tar_args);
    let delete_from_tar_args = add_ssh_cmd(from_folder, ssh_servers, &mut delete_from_tar_args);
    let mut copy_to_folder: Vec<String> = Vec::new();

    let (
        mut verify_tar_args,
        mut make_path_to_target_folder_args,
        mut delete_target_folder_args,
        mut untar_folder_args,
        mut delete_to_tar_args,
    ) = untar_directory(to_path.clone(), to_work_folder, tar_name.clone());
    let verify_tar_args = add_ssh_cmd(to_folder, ssh_servers, &mut verify_tar_args);
    let make_path_to_target_folder_args =
        add_ssh_cmd(to_folder, ssh_servers, &mut make_path_to_target_folder_args);
    let delete_target_folder_args =
        add_ssh_cmd(to_folder, ssh_servers, &mut delete_target_folder_args);
    let untar_folder_args = add_ssh_cmd(to_folder, ssh_servers, &mut untar_folder_args);
    let delete_to_tar_args = add_ssh_cmd(to_folder, ssh_servers, &mut delete_to_tar_args);

    if is_from_ssh || is_to_ssh {
        let scp_cmd = scp_cmd(
            from_folder,
            to_folder,
            tar_name.clone(),
            from_work_folder,
            to_work_folder,
            ssh_servers,
        );
        for cmd in scp_cmd.iter() {
            copy_to_folder.push(cmd.to_string());
        }
    } else {
        copy_to_folder.push("cp".to_string());
        copy_to_folder.push("-R".to_string());
        copy_to_folder.push(from_path.clone());
        copy_to_folder.push(to_path.clone());
    };

    let check_folder_arg = tar_source_exists_args
        .first()
        .expect("First argument required");
    let mut check_folder_cmd = Command::new(check_folder_arg);
    for folder_arg in &tar_source_exists_args[1..] {
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
    println!("- {}", create_tar_args.join(" "));
    println!("- {}", copy_to_folder.join(" "));
    println!("- {}", make_path_to_target_folder_args.join(" "));
    println!("- {}", verify_tar_args.join(" "));
    println!("- {}", delete_target_folder_args.join(" "));
    println!("- {}", untar_folder_args.join(" "));
    println!("- {}", delete_to_tar_args.join(" "));
    println!("- {}", delete_from_tar_args.join(" "));

    if !force {
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

    run_cmd(create_tar_args, true, "Failed to Create Tar".to_string());
    run_cmd(copy_to_folder, true, "Failed to Copy Files".to_string());
    run_cmd(
        make_path_to_target_folder_args,
        true,
        "Make Target Directories".to_string(),
    );
    run_cmd(verify_tar_args, true, "Failed to Verify Tar".to_string());
    run_cmd(
        delete_target_folder_args,
        true,
        "Failed to Delete Target Folder".to_string(),
    );
    run_cmd(
        untar_folder_args,
        true,
        "Failed to Untar Archive".to_string(),
    );
    run_cmd(
        delete_to_tar_args,
        true,
        "Failed to Delete To Tar".to_string(),
    );
    run_cmd(
        delete_from_tar_args,
        true,
        "Failed to Delete From Tar".to_string(),
    );
}

fn run_cmd(cmd_args: Vec<String>, print: bool, failure_msg: String) -> () {
    let first_arg = cmd_args.first().expect("First argument required");
    let mut cmd = Command::new(first_arg);
    for folder_arg in &cmd_args[1..] {
        cmd.arg(folder_arg);
    }
    if print {
        cmd.stdout(Stdio::inherit());
    }
    cmd.output().expect(&failure_msg);
}
