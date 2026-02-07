use crate::model::folder::Folder;
use rand::distributions::Alphanumeric;
use rand::Rng;

pub fn tar_directory(
    target_path: String,
    work_folder: &Folder,
) -> (String, Vec<String>, Vec<String>, Vec<String>) {
    let random_name: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(7)
        .map(char::from)
        .collect();
    let tar_name = format!("{}.tar.gz", random_name);
    let tar_path = format!("{}/{}", work_folder.path, &tar_name);

    let mut tar_source_exists_args: Vec<String> = Vec::new();
    tar_source_exists_args.push("ls".to_string());
    tar_source_exists_args.push(target_path.clone());

    let target_path_split = split_path(target_path.clone());
    let mut create_tar_args: Vec<String> = Vec::new();
    create_tar_args.push("tar".to_string());
    create_tar_args.push("-cf".to_string());
    create_tar_args.push(tar_path.clone());
    create_tar_args.push("-C".to_string());
    create_tar_args.push(target_path_split.0);
    create_tar_args.push(target_path_split.1);

    let mut delete_tar_args: Vec<String> = Vec::new();
    delete_tar_args.push("rm".to_string());
    delete_tar_args.push(tar_path);
    (
        tar_name,
        tar_source_exists_args,
        create_tar_args,
        delete_tar_args,
    )
}

pub fn untar_directory(
    target_path: String,
    work_folder: &Folder,
    tar_name: String,
) -> (
    Vec<String>,
    Vec<String>,
    Vec<String>,
    Vec<String>,
    Vec<String>,
) {
    let mut verify_tar: Vec<String> = Vec::new();
    let mut make_path_to_target_folder: Vec<String> = Vec::new();
    let mut delete_target_folder: Vec<String> = Vec::new();
    let mut untar_folder: Vec<String> = Vec::new();
    let mut delete_tar: Vec<String> = Vec::new();

    let tar_path = format!("{}/{}", work_folder.path, tar_name);
    verify_tar.push("ls".to_string());
    verify_tar.push(tar_path.clone());

    make_path_to_target_folder.push("mkdir".to_string());
    make_path_to_target_folder.push("-p".to_string());
    make_path_to_target_folder.push(target_path.clone());

    delete_target_folder.push("rm".to_string());
    delete_target_folder.push("-rf".to_string());
    delete_target_folder.push(target_path.clone());

    let target_path_split = split_path(target_path.clone());
    untar_folder.push("tar".to_string());
    untar_folder.push("-xf".to_string());
    untar_folder.push(tar_path.clone());
    untar_folder.push("-C".to_string());
    untar_folder.push(target_path_split.0);

    delete_tar.push("rm".to_string());
    delete_tar.push(tar_path.clone());

    (
        verify_tar,
        make_path_to_target_folder,
        delete_target_folder,
        untar_folder,
        delete_tar,
    )
}

fn split_path(path: String) -> (String, String) {
    let mut split_path: Vec<&str> = path.split("/").collect();
    let target = if let Some(..) = split_path.last() {
        split_path.pop().expect("One argument required").to_string()
    } else {
        String::new()
    };

    let mut base_path = String::new();
    for arg in split_path {
        if arg.trim().len() > 0 {
            base_path = format!("{}/{}", base_path, arg);
        }
    }
    return (base_path, target);
}
