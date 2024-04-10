use crate::model::link::Link;
use crate::model::{cli::CmdArgs, folder::Folder, ssh::SshServer};
use crate::service::folder;
use crate::service::link;
use crate::service::ssh;
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
                    ssh::ls(&link.local, &ssh_servers, &path);
                    ssh::ls(&link.target, &ssh_servers, &path);
                }
            } else {
                ssh::ls(&link.local, &ssh_servers, &relative_path);
                ssh::ls(&link.target, &ssh_servers, &relative_path);
            }
        }
    } else {
        let folder = folder::get(target.clone(), folders);
        if let Some(folder) = folder {
            ssh::ls(&folder, &ssh_servers, &relative_path);
        } else {
            println!("Error locating folder: {}", target)
        }
    }
}
