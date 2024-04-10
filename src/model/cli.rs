use clap::Parser;

#[derive(Parser, Debug)]
pub struct CmdArgs {
    pub target: String,
    pub relative_path: Option<String>,
}
#[derive(Parser, Debug)]
pub struct PullArgs {
    pub target: String,
    pub relative_path: Option<String>,
}
#[derive(Parser, Debug)]
pub struct PushArgs {
    pub target: String,
    pub relative_path: Option<String>,
}
#[derive(Parser, Debug)]
pub enum CliCmd {
    Ls(CmdArgs),
    Pull(PullArgs),
    Push(PushArgs),
}
