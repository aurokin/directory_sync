use crate::model::link::Link;
use crate::model::{cli::CmdArgs, folder::Folder, ssh::SshServer};
use crate::service::folder;
use crate::service::link;
use std::collections::HashMap;

pub fn ls(
    cmd_args: CmdArgs,
    is_link: bool,
    folders: HashMap<String, Folder>,
    links: HashMap<String, Link>,
    ssh_servers: HashMap<String, SshServer>,
) -> () {
    let target = cmd_args.target;
    let relative_path = cmd_args.relative_path;
    if is_link {
        let link = link::get(target.clone(), links);
        if let Some(link) = link {
            if link.paths.len() > 0 {
                for path in link.paths {
                    let path = Some(path);
                    crate::service::core::ls(&link.local, &ssh_servers, &path);
                    crate::service::core::ls(&link.target, &ssh_servers, &path);
                }
            } else {
                crate::service::core::ls(&link.local, &ssh_servers, &relative_path);
                crate::service::core::ls(&link.target, &ssh_servers, &relative_path);
            }
        }
    } else {
        let folder = folder::get(target.clone(), folders);
        if let Some(folder) = folder {
            crate::service::core::ls(&folder, &ssh_servers, &relative_path);
        } else {
            println!("Error locating folder: {}", target)
        }
    }
}

pub fn pull(
    cmd_args: CmdArgs,
    is_link: bool,
    is_force: bool,
    work_folder: Folder,
    folders: HashMap<String, Folder>,
    links: HashMap<String, Link>,
    ssh_servers: HashMap<String, SshServer>,
) -> () {
    let target = cmd_args.target;
    let relative_path = cmd_args.relative_path;
    if is_link {
        let link = link::get(target.clone(), links);
        if let Some(link) = link {
            if link.paths.len() > 0 {
                for path in link.paths {
                    let path = Some(path);
                    crate::service::core::sync(
                        &link.target,
                        &link.local,
                        &work_folder,
                        &ssh_servers,
                        &path,
                        is_force,
                    );
                }
            } else if link.partial_only {
                println!("{} is partial only, ending task", link.name);
                return;
            } else {
                crate::service::core::sync(
                    &link.target,
                    &link.local,
                    &work_folder,
                    &ssh_servers,
                    &relative_path,
                    is_force,
                );
            }
        }
    } else {
        let current_folder: Option<Folder> = folder::get_current_dir();
        let current_folder = current_folder.expect("Unable to resolve cwd");
        let folder: Option<Folder> = folder::get(target.clone(), folders);

        if let Some(folder) = folder {
            crate::service::core::sync(
                &folder,
                &current_folder,
                &work_folder,
                &ssh_servers,
                &relative_path,
                is_force,
            );
        } else {
            println!("Error locating folder: {}", target)
        }
    }
}

pub fn push(
    cmd_args: CmdArgs,
    is_link: bool,
    is_force: bool,
    work_folder: Folder,
    folders: HashMap<String, Folder>,
    links: HashMap<String, Link>,
    ssh_servers: HashMap<String, SshServer>,
) -> () {
    let target = cmd_args.target;
    let relative_path = cmd_args.relative_path;
    if is_link {
        let link = link::get(target.clone(), links);
        if let Some(link) = link {
            if link.paths.len() > 0 {
                for path in link.paths {
                    let path = Some(path);
                    crate::service::core::sync(
                        &link.local,
                        &link.target,
                        &work_folder,
                        &ssh_servers,
                        &path,
                        is_force,
                    );
                }
            } else if link.partial_only {
                println!("{} is partial only, ending task", link.name);
                return;
            } else {
                crate::service::core::sync(
                    &link.local,
                    &link.target,
                    &work_folder,
                    &ssh_servers,
                    &relative_path,
                    is_force,
                );
            }
        }
    } else {
        let current_folder: Option<Folder> = folder::get_current_dir();
        let current_folder = current_folder.expect("Unable to resolve cwd");
        let folder: Option<Folder> = folder::get(target.clone(), folders);

        if let Some(folder) = folder {
            crate::service::core::sync(
                &current_folder,
                &folder,
                &work_folder,
                &ssh_servers,
                &relative_path,
                is_force,
            );
        } else {
            println!("Error locating folder: {}", target)
        }
    }
}
