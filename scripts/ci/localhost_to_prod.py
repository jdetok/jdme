#!/usr/bin/env python3

import pathlib
import re

local_ptrn = re.compile(r"http://localhost:[0-9]+")
prod_url = "https://jdeko.me"

for path in pathlib.Path(".").rglob('*'):
    if path.suffix in [".js", ".html", ".css", ".go"]:
        txt = path.read_text()
        new_txt = local_ptrn.sub(prod_url, txt)
        if txt != new_txt:
            path.write_text(new_txt)