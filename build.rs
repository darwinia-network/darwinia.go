use std::process::Command;

fn main() {
    println!(r"cargo:rustc-link-search=target/debug");
    println!("cargo:rerun-if-changed=internal/ffi/mod.go");

    let os = Command::new("uname").output().unwrap();
    let ext = match String::from_utf8_lossy(os.stdout.as_slice())
        .into_owned()
        .trim_end()
        .as_ref()
    {
        "Darwin" => "dylib",
        _ => "so",
    };

    Command::new("go")
        .args(&[
            "build",
            "-o",
            &format!("target/debug/libeth.{}", ext),
            "-buildmode=c-shared",
            "internal/ffi/mod.go",
        ])
        .status()
        .unwrap();
}
