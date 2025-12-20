#!/usr/bin/env python3

import pathlib
import sys
import re

patterns = [r"https?://localhost:[0-9]+/?", r"https?://jdeko.me[:\/]?[0-9]*/?"]
filetypes = [".js", ".html", ".css", ".go"]
exclude_dirs = [".git", "wiki", "log", "z_log", "puml", "bin"]

def main():
    repl_url = "https://jdeko.me/"
    num_replaced = 0    
    if len(sys.argv) > 1:
        if sys.argv[1] == "local":
            repl_url = "http://localhost:8080/"
    try: 
        for path in pathlib.Path(".").rglob('*'):
            if path.parts[0] in exclude_dirs: continue
            if path.suffix in filetypes:
                try:
                    for ptrn in patterns:
                        txt = path.read_text()
                        new_txt = re.compile(ptrn).sub(repl_url, txt)
                        if txt != new_txt:
                            path.write_text(new_txt)
                            num_replaced += 1
                            print(f"replaced match to {ptrn} with {repl_url} in {path}")
                except Exception as e:
                    raise SystemError(f"error occured reading {path}: {e}")
    except SystemError as e: 
        raise SystemExit(f"exiting with error: {e}")
    finally: print(f"{num_replaced} files with changs")
        
if __name__ == "__main__":
    main()