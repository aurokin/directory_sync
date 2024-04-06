use super::folder::Folder;

#[derive(Debug)]
pub struct Link {
    pub name: String,
    pub local: Folder,
    pub target: Folder,
    pub paths: Vec<String>,
}
