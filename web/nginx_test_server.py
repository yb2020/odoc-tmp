#!/usr/bin/env python3
import http.server
import socketserver
import os
import sys
import re
import requests
import json
from urllib.parse import urlparse

# 设置服务器根目录为 dist 目录
ROOT_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'dist')
# API 代理目标
API_PROXY_URL = 'http://192.168.213.179:8081'
# Nginx 配置文件路径
NGINX_CONF = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'nginx.conf')
PORT = 8000

# 创建会话对象来维护 Cookie
session = requests.Session()

# 解析 nginx.conf 中的路由配置
def parse_nginx_conf(conf_path):
    try:
        print(f"正在解析配置文件: {conf_path}")
        with open(conf_path, 'r') as f:
            content = f.read()
            
        # 解析 location /docs/ 块 (注意可能有尾部斜杠)
        docs_location = re.search(r'location\s+/docs/?(?:\s*|\s+[^{]*){([^}]*)}', content)
        if not docs_location:
            print(f"警告: 在 {conf_path} 中没有找到 location /docs 块")
            return None
            
        # 提取 try_files 指令
        try_files = re.search(r'try_files\s+([^;]*);', docs_location.group(1))
        if not try_files:
            print(f"警告: 在 location /docs 块中没有找到 try_files 指令")
            return None
            
        # 解析 try_files 指令中的路径
        paths = try_files.group(1).strip().split()
        print(f"从 nginx.conf 中解析到的 try_files 路径: {paths}")
        return paths
    except Exception as e:
        print(f"解析 nginx.conf 时出错: {e}")
        return None

class NginxCompatibleHandler(http.server.SimpleHTTPRequestHandler):
    nginx_paths = None
    
    # 添加正确的 MIME 类型映射
    extensions_map = {
        **http.server.SimpleHTTPRequestHandler.extensions_map,
        '.js': 'application/javascript',
        '.mjs': 'application/javascript',
        '.ts': 'application/javascript',  # TypeScript 文件在浏览器中也应该作为 JavaScript 处理
        '.json': 'application/json',
        '.wasm': 'application/wasm',
    }
    
    def setup(self):
        if NginxCompatibleHandler.nginx_paths is None:
            NginxCompatibleHandler.nginx_paths = parse_nginx_conf(os.path.join(os.path.dirname(os.path.abspath(__file__)), "nginx.conf"))
        super().setup()
    
    def do_GET(self):
        # 解析请求路径
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        print(f"请求路径: {path}")
        
        # 如果是 API 请求，转发到目标服务器
        if path.startswith('/api/') or path == '/report/collection_tracking0':
            return self.do_proxy()
            
        # 如果是 /docs 开头的请求，应用从 nginx.conf 解析的规则
        if path.startswith('/docs/'):
            relative_path = path[6:]  # 去掉 '/docs/'
            
            # 特殊处理：如果访问的是 index.html/ 结尾的路径，重定向到对应的目录
            if relative_path.endswith('index.html/'):
                # 去掉 'index.html/' 后缀，重定向到目录
                dir_path = relative_path[:-11]  # 去掉 'index.html/'
                redirect_path = f'/docs/{dir_path}' if dir_path else '/docs/'
                print(f"重定向 {path} 到 {redirect_path}")
                self.send_response(301)
                self.send_header('Location', redirect_path)
                self.end_headers()
                return
            
            # 创建可能的路径列表
            possible_paths = []
            
            # 如果成功解析了 nginx.conf
            if NginxCompatibleHandler.nginx_paths:
                for nginx_path in NginxCompatibleHandler.nginx_paths:
                    # 处理特殊路径标记
                    if nginx_path == '$uri':
                        possible_paths.append(os.path.join(ROOT_DIR, 'docs', relative_path))
                    elif nginx_path == '$uri.html':
                        possible_paths.append(os.path.join(ROOT_DIR, 'docs', relative_path + '.html'))
                    elif nginx_path == '$uri/index.html':
                        possible_paths.append(os.path.join(ROOT_DIR, 'docs', relative_path, 'index.html'))
                    elif nginx_path.startswith('/docs/'):
                        # 处理绝对路径
                        fallback_path = nginx_path[6:]  # 去掉 '/docs/'
                        possible_paths.append(os.path.join(ROOT_DIR, 'docs', fallback_path))
            else:
                # 如果没有解析到 nginx.conf，使用默认路径
                possible_paths = [
                    os.path.join(ROOT_DIR, 'docs', relative_path),  # 原始路径
                    os.path.join(ROOT_DIR, 'docs', relative_path + '.html'),  # 添加 .html 扩展名
                    os.path.join(ROOT_DIR, 'docs', relative_path, 'index.html'),  # 作为目录尝试 index.html
                    os.path.join(ROOT_DIR, 'docs', 'index.html')  # 最后回退到 index.html
                ]
            
            # 尝试每一个可能的路径
            for test_path in possible_paths:
                print(f"尝试路径: {test_path}")
                if os.path.exists(test_path) and not os.path.isdir(test_path):
                    # 找到文件，直接提供文件内容
                    print(f"找到文件: {test_path}")
                    return self.serve_file(test_path)
            
            # 如果所有尝试都失败，回退到 /docs/index.html
            fallback_path = os.path.join(ROOT_DIR, 'docs', 'index.html')
            if os.path.exists(fallback_path):
                print(f"回退到: {fallback_path}")
                return self.serve_file(fallback_path)
            else:
                # 如果连 index.html 都不存在，返回 404
                self.send_error(404, "File not found")
                return
        
        # 对于根路径请求，优先提供 index.html，如果不存在才重定向到 /docs/
        elif path == '/':
            index_path = os.path.join(ROOT_DIR, 'index.html')
            if os.path.exists(index_path):
                print(f"根路径请求，提供: {index_path}")
                return self.serve_file(index_path)
            else:
                # 如果 index.html 不存在，重定向到 /docs/
                print("根路径请求，index.html 不存在，重定向到 /docs/")
                self.send_response(302)
                self.send_header('Location', '/docs/')
                self.end_headers()
                return
        
        # 处理其他路径，实现与 nginx 配置中的 location / 相同的逻辑
        else:
            # 创建可能的路径列表，模拟 nginx 的 try_files $uri $uri/ /index.html
            possible_paths = [
                os.path.join(ROOT_DIR, path.lstrip('/')),  # $uri
                os.path.join(ROOT_DIR, path.lstrip('/'), 'index.html'),  # $uri/
                os.path.join(ROOT_DIR, 'index.html')  # /index.html
            ]
            
            # 尝试每一个可能的路径
            for test_path in possible_paths:
                print(f"尝试路径: {test_path}")
                if os.path.exists(test_path) and not os.path.isdir(test_path):
                    # 找到文件，直接提供文件内容
                    print(f"找到文件: {test_path}")
                    return self.serve_file(test_path)
            
            # 如果所有尝试都失败，回退到 /index.html
            fallback_path = os.path.join(ROOT_DIR, 'index.html')
            if os.path.exists(fallback_path):
                print(f"回退到: {fallback_path}")
                return self.serve_file(fallback_path)
            else:
                # 如果连 index.html 都不存在，返回 404
                self.send_error(404, "File not found")
                return
    
    def serve_file(self, file_path):
        """直接提供文件内容，避免重新调用do_GET导致循环"""
        try:
            # 规范化路径，确保使用正斜杠
            normalized_path = file_path.replace('\\', '/')
            print(f"提供文件: {normalized_path}")
            
            # 检查文件是否存在
            if not os.path.exists(file_path):
                self.send_error(404, "File not found")
                return
            
            # 发送响应头
            self.send_response(200)
            
            # 根据文件扩展名设置Content-Type
            content_type = self.guess_type(file_path)
            self.send_header('Content-Type', content_type)
            
            # 获取文件大小
            file_size = os.path.getsize(file_path)
            self.send_header('Content-Length', str(file_size))
            
            self.end_headers()
            
            # 读取并发送文件内容
            with open(file_path, 'rb') as f:
                self.wfile.write(f.read())
                
        except Exception as e:
            print(f"提供文件时出错: {e}")
            self.send_error(500, f"Internal server error: {e}")

    def do_proxy(self):
        """将请求转发到目标服务器"""
        # 构建目标 URL
        target_url = API_PROXY_URL + self.path
        print(f"转发请求到: {target_url}")
        print(f"请求方法: {self.command}")
        
        # 获取请求头部
        headers = {}
        for header in self.headers:
            if header.lower() not in ('host', 'content-length'):
                headers[header] = self.headers[header]
        
        # 打印请求头部中的 Cookie
        if 'Cookie' in self.headers:
            print(f"请求中的 Cookie: {self.headers['Cookie']}")
        
        # 添加代理相关头部
        headers['Host'] = urlparse(API_PROXY_URL).netloc
        headers['X-Real-IP'] = self.client_address[0]
        headers['X-Forwarded-For'] = self.client_address[0]
        headers['X-Forwarded-Host'] = self.headers.get('Host', '')
        headers['X-Forwarded-Proto'] = 'http'
        headers['X-Original-URI'] = self.path
        
        # 如果是登录请求，打印更详细的信息
        if self.path == '/api/oauth2/sign_in':
            print("处理登录请求...")
            print(f"请求头部: {headers}")
            if self.command == 'POST':
                content_length = int(self.headers.get('Content-Length', 0))
                if content_length > 0:
                    print(f"请求体长度: {content_length} 字节")
        
        try:
            # 准备请求体
            data = None
            if self.command in ('POST', 'PUT'):
                content_length = int(self.headers.get('Content-Length', 0))
                if content_length > 0:
                    data = self.rfile.read(content_length)
            
            # 使用 requests 发送请求
            if self.command == 'GET':
                resp = session.get(target_url, headers=headers, allow_redirects=False)
            elif self.command == 'POST':
                resp = session.post(target_url, headers=headers, data=data, allow_redirects=False)
            elif self.command == 'PUT':
                resp = session.put(target_url, headers=headers, data=data, allow_redirects=False)
            else:
                # 其他方法
                resp = session.request(self.command, target_url, headers=headers, data=data, allow_redirects=False)
            
            # 打印响应信息
            print("\n===== 收到响应 =====")
            print(f"状态码: {resp.status_code}\n")
            print("----- 响应头部 -----")
            for name, value in resp.headers.items():
                print(f"{name}: {value}")
            
            # 打印响应体
            try:
                print("\n----- 响应体 (前1000字节) -----")
                print(resp.text[:1000])
                if len(resp.text) > 1000:
                    print("... [已省略部分内容] ...")
            except UnicodeDecodeError:
                print("\n----- 响应体 (二进制数据) -----")
                print(f"长度: {len(resp.content)} 字节")
            print("\n===== 响应结束 =====\n")
            
            # 如果是登录响应，打印更详细的信息
            if self.path == '/api/oauth2/sign_in' and resp.status_code == 200:
                print("登录成功，检查 Cookie 设置情况:")
                for cookie in session.cookies:
                    print(f"  Cookie: {cookie.name}={cookie.value}; domain={cookie.domain}; path={cookie.path}")
            
            # 设置响应状态码
            self.send_response(resp.status_code)
            
            # 设置响应头部，特别处理 Set-Cookie 头
            # 打印原始响应头部的结构
            print("\n原始响应头部结构:")
            for name, value in resp.headers.items():
                print(f"{name}: {value}")
                
            # 特别处理 Set-Cookie 头
            # 在 requests 库中，多个 Set-Cookie 头可能会被合并成一个字符串
            # 我们需要解析原始响应中的 Set-Cookie 头
            raw_cookies = []
            
            # 如果有 Set-Cookie 头，尝试解析
            if 'Set-Cookie' in resp.headers:
                # 尝试从原始响应中获取所有 Set-Cookie 头
                cookie_str = resp.headers['Set-Cookie']
                print(f"\n原始 Set-Cookie 字符串: {cookie_str}")
                
                # 如果包含多个 Cookie，尝试分割
                if ';' in cookie_str and '=' in cookie_str:
                    # 先检查是否有多个完整的 cookie 定义
                    if cookie_str.count('=') > 1 and cookie_str.count(';') > 1:
                        # 尝试按照常见的分隔符分割
                        parts = cookie_str.split(', ')
                        for part in parts:
                            if '=' in part:
                                raw_cookies.append(part)
                    else:
                        raw_cookies.append(cookie_str)
                else:
                    raw_cookies.append(cookie_str)
                
            print(f"\n检测到 {len(raw_cookies)} 个 Set-Cookie 头:")
            for cookie in raw_cookies:
                print(f"Set-Cookie: {cookie}")
                # 直接转发每个 Set-Cookie 头
                self.send_header('Set-Cookie', cookie)
                
            # 处理其他头部
            for name, value in resp.headers.items():
                if name.lower() != 'set-cookie' and name.lower() not in ('server', 'date', 'transfer-encoding'):
                    self.send_header(name, value)
                    
            # 打印当前保存的 Cookie
            print("当前保存的 Cookies:")
            for cookie in session.cookies:
                print(f"  {cookie.name}: {cookie.value}")
            self.end_headers()
            
            # 转发响应体
            self.wfile.write(resp.content)
                
        except requests.exceptions.RequestException as e:
            print(f"HTTP 错误: {e}")
            self.send_response(500)
            self.send_header('Content-Type', 'text/plain')
            self.end_headers()
            self.wfile.write(str(e).encode())
        except Exception as e:
            print(f"代理错误: {e}")
            self.send_response(500)
            self.end_headers()
            self.wfile.write(str(e).encode())
            
        return
    
    def do_POST(self):
        """处理 POST 请求"""
        # 如果是 API 请求，转发到目标服务器
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path.startswith('/api/') or path == '/report/collection_tracking0':
            return self.do_proxy()
        
        # 其他 POST 请求使用默认处理
        return super().do_POST()
        
    def do_PUT(self):
        """处理 PUT 请求"""
        # 如果是 API 请求，转发到目标服务器
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path.startswith('/api/') or path == '/report/collection_tracking0':
            return self.do_proxy()
        
        # 其他 PUT 请求使用默认处理
        return super().do_PUT()

    def translate_path(self, path):
        # 设置根目录为 dist
        path = super().translate_path(path)
        rel_path = os.path.relpath(path, os.getcwd())
        return os.path.join(ROOT_DIR, rel_path)
        
# 主函数
def main():
    # 检查 dist 目录是否存在
    if not os.path.exists(ROOT_DIR):
        print(f"错误: {ROOT_DIR} 目录不存在")
        sys.exit(1)
        
    # 检查 dist/docs 目录是否存在
    if not os.path.exists(os.path.join(ROOT_DIR, 'docs')):
        print(f"错误: {os.path.join(ROOT_DIR, 'docs')} 目录不存在")
        sys.exit(1)
    
    # 检查 nginx.conf 文件是否存在
    if not os.path.exists(NGINX_CONF):
        print(f"警告: {NGINX_CONF} 文件不存在，将使用默认配置")
    
    print(f"启动服务器，根目录: {ROOT_DIR}")
    print(f"使用配置文件: {NGINX_CONF}")
    print(f"访问 http://localhost:{PORT}/docs/ 查看文档")
    
    try:
        with socketserver.TCPServer(("", PORT), NginxCompatibleHandler) as httpd:
            print(f"服务器运行在端口 {PORT}")
            httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n服务器已停止")
    except Exception as e:
        print(f"启动服务器时出错: {e}")

if __name__ == "__main__":
    main()
