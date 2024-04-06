use clap::Parser;
#[derive(Parser, Debug)]
pub struct CmdArgs {
    pub target: String,
}
#[derive(Parser, Debug)]
pub enum CliCmd {
    Ls(CmdArgs),
}
