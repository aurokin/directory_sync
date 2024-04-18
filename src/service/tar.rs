// tar.rs
use crate::model::folder::Folder;
use rand::distributions::Alphanumeric;
use rand::Rng;

pub fn tar_directory(
    target_folder: &Folder,
    work_folder: &Folder,
) -> (Vec<String>, Vec<String>, Vec<String>) {
    let random_name: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(7)
        .map(char::from)
        .collect();
    let tar_path = format!("{}/{}.tar.gz", work_folder.path, random_name);

    let mut tar_source_exists_args: Vec<String> = Vec::new();
    tar_source_exists_args.push("ls".to_string());
    tar_source_exists_args.push(target_folder.path.clone());

    let mut create_tar_args: Vec<String> = Vec::new();
    create_tar_args.push("tar".to_string());
    create_tar_args.push("-c".to_string());
    create_tar_args.push("-f".to_string());
    create_tar_args.push(tar_path.clone());
    create_tar_args.push(target_folder.path.clone());

    let mut delete_tar_args: Vec<String> = Vec::new();
    delete_tar_args.push("rm".to_string());
    delete_tar_args.push(tar_path);
    (tar_source_exists_args, create_tar_args, delete_tar_args)
}

pub fn untar_directory(
    target_folder: &Folder,
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
    make_path_to_target_folder.push(target_folder.path.clone());

    delete_target_folder.push("rm".to_string());
    delete_target_folder.push("-rf".to_string());
    delete_target_folder.push(target_folder.path.clone());

    untar_folder.push("tar".to_string());
    untar_folder.push("-x".to_string());
    untar_folder.push("-f".to_string());
    untar_folder.push(tar_path.clone());
    untar_folder.push("-C".to_string());
    untar_folder.push(target_folder.path.clone());

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
