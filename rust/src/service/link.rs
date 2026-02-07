use crate::model::link::Link;
use std::collections::HashMap;

pub fn get(name: String, links: HashMap<String, Link>) -> Option<Link> {
    let mut found = None;
    for link in links {
        if link.0 == name {
            found = Some(link.1);
            break;
        }
    }
    return found;
}
