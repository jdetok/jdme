#!/usr/bin/env python3

import pathlib
import re

local_ptrn = re.compile(r"http://localhost:[0-9]+")
prod_url = "https://jdeko.me"

num_replaced = 0
try: 
    for path in pathlib.Path(".").rglob('*'):
        if path.suffix in [".js", ".html", ".css", ".go"]:
            try:
                txt = path.read_text()
                new_txt = local_ptrn.sub(prod_url, txt)
                if txt != new_txt:
                    path.write_text(new_txt)
                    num_replaced += 1
                    print(f"replaced {txt} with {new_txt} in {path}")
            except: raise SystemExit
            finally: print(f"error occured reading {path}")
except: raise SystemExit
finally: print(f"{num_replaced} strings replaced")
    
            