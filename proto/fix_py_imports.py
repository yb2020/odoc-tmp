#!/usr/bin/env python3
"""
修复 protoc 生成的 Python 文件中的导入路径问题。
将 'from definitions.xxx import yyy_pb2' 转换为绝对导入 'from xxx import yyy_pb2'。
这样生成的包可以直接安装使用，只需将 gen/py 加入 PYTHONPATH 或打包安装。
"""

import os
import re
import sys


def fix_imports_in_file(filepath, gen_dir):
    """修复单个文件中的导入语句"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    # 获取当前文件所在的模块路径 (相对于 gen_dir)
    rel_path = os.path.relpath(filepath, gen_dir)
    current_module = os.path.dirname(rel_path).replace(os.sep, '.')
    
    # 匹配 from definitions.xxx.yyy import zzz_pb2 as zzz__pb2
    # 或 from definitions.xxx import yyy_pb2 as yyy__pb2
    pattern = r'from definitions\.([a-zA-Z0-9_.]+) import ([a-zA-Z0-9_]+_pb2) as ([a-zA-Z0-9_]+__pb2)'
    
    def replace_import(match):
        import_path = match.group(1)  # e.g., "translate" or "validate"
        module_name = match.group(2)  # e.g., "TranslateEnum_pb2"
        alias = match.group(3)        # e.g., "TranslateEnum__pb2"
        
        # 将 definitions.xxx 转换为 xxx（绝对导入）
        # import_path 可能是 "translate" 或 "common.sub"
        target_module = import_path.replace('.', '.')  # 保持点分隔
        
        if current_module == import_path:
            # 同一模块内，使用相对导入
            return f'from . import {module_name} as {alias}'
        else:
            # 跨模块导入，使用绝对导入（基于 gen/py 根目录）
            return f'from {import_path} import {module_name} as {alias}'
    
    content = re.sub(pattern, replace_import, content)
    
    # 同时处理没有 as 别名的情况
    pattern_no_alias = r'from definitions\.([a-zA-Z0-9_.]+) import ([a-zA-Z0-9_]+_pb2)(?!\s+as)'
    
    def replace_import_no_alias(match):
        import_path = match.group(1)
        module_name = match.group(2)
        
        if current_module == import_path:
            # 同一模块内，使用相对导入
            return f'from . import {module_name}'
        else:
            # 跨模块导入，使用绝对导入
            return f'from {import_path} import {module_name}'
    
    content = re.sub(pattern_no_alias, replace_import_no_alias, content)
    
    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    return False


def fix_all_imports(gen_dir):
    """修复目录下所有 Python 文件的导入"""
    fixed_count = 0
    for root, dirs, files in os.walk(gen_dir):
        for filename in files:
            if filename.endswith('_pb2.py'):
                filepath = os.path.join(root, filename)
                if fix_imports_in_file(filepath, gen_dir):
                    print(f"Fixed imports in: {filepath}")
                    fixed_count += 1
    return fixed_count


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python fix_py_imports.py <gen_dir>")
        print("Example: python fix_py_imports.py ./gen/py")
        sys.exit(1)
    
    gen_dir = sys.argv[1]
    if not os.path.isdir(gen_dir):
        print(f"Error: Directory not found: {gen_dir}")
        sys.exit(1)
    
    fixed = fix_all_imports(gen_dir)
    print(f"Fixed {fixed} files")
