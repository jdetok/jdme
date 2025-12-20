#!/usr/bin/env python3

import pathlib
import re

# local_ptrn = re.compile(r"http://localhost:[0-9]+")
patterns = [r"https?://localhost:[0-9]+",
            r"https?://jdeko.me[:\/]?[0-9]*"
            ]

prod_url = "https://jdeko.me"

num_replaced = 0
try: 
    for path in pathlib.Path(".").rglob('*'):
        print(f"in {path}")
        if path.suffix in [".js", ".html", ".css", ".go"]:
            try:
                for ptrn in patterns:
                    txt = path.read_text()
                    # new_txt = local_ptrn.sub(prod_url, txt)
                    new_txt = re.compile(ptrn).sub(prod_url, txt)
                    if txt != new_txt:
                        path.write_text(new_txt)
                        num_replaced += 1
                        print(f"replaced match to {ptrn} in {path}")
            except Exception as e:
                print(f"error occured reading {path}: {e}")
                raise SystemExit
except: raise SystemExit
finally: print(f"{num_replaced} files with changs")
    
            