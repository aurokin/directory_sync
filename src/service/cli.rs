use crate::model::link::Link;
use crate::model::{cli::CmdArgs, folder::Folder, ssh::SshServer};
use crate::service::folder;
use crate::service::link;
use crate::service::ssh;
use std::collections::HashMap;

pub fn ls(
    cmd_args: CmdArgs,
    folders: HashMap<String, Folder>,
    ssh_servers: HashMap<String, SshServer>,
    links: HashMap<String, Link>,
    is_link: bool,
) -> () {
    let target = cmd_args.target;
    if is_link {
        let link = link::get(target.clone(), links);
        if let Some(link) = link {
            ssh::ls(link.local, &ssh_servers);
            ssh::ls(link.target, &ssh_servers);
        }
    } else {
        let folder = folder::get(target.clone(), folders);
        if let Some(folder) = folder {
            ssh::ls(folder, &ssh_servers);
        } else {
            println!("Error locating folder: {}", target)
        }
    }
}
