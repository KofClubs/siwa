#!/usr/bin/python3

# Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import os


def replace_copyright_symbol_in_file(file):
    with open(file, 'r') as f:
        filestr = f.read()
        filestr = filestr.replace('(c)', '(c)')
    with open(file, 'w') as f:
        f.write(filestr)


def replace_copyright_symbol(dirpath):
    files = os.listdir(dirpath)
    for file in files:
        filepath = os.path.join(dirpath, file)
        if os.path.isdir(filepath):
            if filepath != '../.git' and filepath != '../vendor' and filepath != '../bin':
                replace_copyright_symbol(filepath)
        else:
            replace_copyright_symbol_in_file(filepath)


os.chdir(os.path.dirname(os.path.abspath(__file__)))
replace_copyright_symbol('..')
