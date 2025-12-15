use tauri::Manager;
use tauri::path::BaseDirectory;
use tauri_plugin_shell::ShellExt;

mod commands;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_os::init())
        .setup(|app| {
            // 打开 DevTools（Release 模式也可用于调试）
            let window = app.get_webview_window("main").unwrap();
            window.open_devtools();
            
            // 启动 Go 后端（仅 Release 模式，开发模式手动启动）
            if !cfg!(debug_assertions) {
                let handle = app.handle().clone();
                let resource_path = app.path().resource_dir().ok();
                
                tauri::async_runtime::spawn(async move {
                    // 获取配置文件路径
                    let config_path = resource_path
                            .map(|p| p.join("resources").join("config.develop.yaml"))
                            .and_then(|p| p.to_str().map(|s| s.to_string()))
                            .unwrap_or_default();
                    
                    println!("[go-sea] Using config: {}", config_path);
                    
                    let cmd_result = if config_path.is_empty() {
                        handle.shell().sidecar("go-sea")
                    } else {
                        handle.shell().sidecar("go-sea")
                            .map(|c| c.args(["-config", &config_path]))
                    };
                    
                    match cmd_result {
                        Ok(cmd) => {
                            match cmd.spawn() {
                                Ok((mut rx, _child)) => {
                                    println!("[go-sea] Backend started successfully");
                                    // 监听输出
                                    tauri::async_runtime::spawn(async move {
                                        while let Some(event) = rx.recv().await {
                                            match event {
                                                tauri_plugin_shell::process::CommandEvent::Stdout(line) => {
                                                    println!("[go-sea] {}", String::from_utf8_lossy(&line));
                                                }
                                                tauri_plugin_shell::process::CommandEvent::Stderr(line) => {
                                                    eprintln!("[go-sea] {}", String::from_utf8_lossy(&line));
                                                }
                                                tauri_plugin_shell::process::CommandEvent::Terminated(payload) => {
                                                    println!("[go-sea] Process terminated with code: {:?}", payload.code);
                                                    break;
                                                }
                                                _ => {}
                                            }
                                        }
                                    });
                                }
                                Err(e) => {
                                    eprintln!("[go-sea] Failed to spawn backend: {}", e);
                                }
                            }
                        }
                        Err(e) => {
                            eprintln!("[go-sea] Failed to create sidecar command: {}", e);
                        }
                    }
                });
            }
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::greet,
            commands::get_app_version,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
