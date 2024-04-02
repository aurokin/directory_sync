use super::config::TomlType;

#[derive(Debug)]
pub enum FolderType {
    Local,
    Ssh,
}
impl FolderType {
    pub fn get_folder_type(toml_type: TomlType) -> Self {
        let folder_type = match toml_type {
            TomlType::Local => FolderType::Local,
            TomlType::Ssh => FolderType::Ssh,
        };

        return folder_type;
    }
}
#[derive(Debug)]
pub struct Folder {
    pub name: String,
    pub path: String,
    pub target: FolderType,
    pub ssh_key: Option<String>,
}
