use clap::Parser;

#[derive(Parser, Debug)]
pub struct CmdArgs {
    pub target: String,
    pub relative_path: Option<String>,
}
#[derive(Parser, Debug)]
pub enum CliCmd {
    Ls(CmdArgs),
    Pull(CmdArgs),
    Push(CmdArgs),
}
