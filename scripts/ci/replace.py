#!/usr/bin/env python3

from pathlib import Path
import sys
import re
                                
def main():
    loc = False
    if len(sys.argv) > 1:
        if sys.argv[1] == "local":
            loc = True            
    
    goarch_opts = {"local": "amd64", "prod": "arm64"}
    to_replace = [
        ReReplace(loc, "URLS", ".", "http://localhost:8080/", "https://jdeko.me/", 0, 0, 
                          [r"https?://localhost:[0-9]+/?", r"https?://jdeko.me[:\/]?[0-9]*/?"], 
                          [".js", ".html", ".css"], [".git", "wiki", "log", "z_log", "puml", "bin"]
        ),
        ReReplace(loc, "IS_PROD", "./main/main.go", "false", "true", 3, 2, 
                              [r"^(\s*IS_PROD\s*=\s*)(true|false)(\s*$)"], [], []
        ),
        ReReplace(loc, "GOARCH", "./Dockerfile", goarch_opts['local'], goarch_opts['prod'], 3, 2, 
                              [rf"^(.*GOARCH=)({goarch_opts['local']}|{goarch_opts['prod']})(.*)$"], [], []
        )
    ]
    num_replaced = 0
    for r in to_replace:
        num_replaced += r.replace()

    print(f"replacements complete: {num_replaced} files with changes")
    

class ReReplace:
    def __init__(self, loc: bool, name: str,
                path, local_repl, prod_repl: str, 
                capt_grps: int, grp_pos: int,
                patterns: tuple[str], filetypes: tuple[str], exclude_dirs: tuple[str]
                ):
        self.p = Path(path)
        self.rename = name
        self.rptrns = patterns
        self.ftypes = filetypes
        self.exdirs = exclude_dirs
        self.rcgrps = capt_grps
        self.rcgpos = grp_pos
        self.rplcmt = local_repl if loc else prod_repl
        self.rplc_ptrn = self.get_repl_ptrn()
        self.rfiles = self.get_files() if str(self.p) == "." else (self.p,)

    def get_repl_ptrn(self) -> str:
        if self.rcgrps == 0 or self.rcgrps is None:
            return self.rplcmt
        repl_ptrn = ""
        pos = 1
        while pos <= self.rcgrps:
            if pos == self.rcgpos:
                repl_ptrn += self.rplcmt
            else: repl_ptrn += fr"\{pos}"
            pos += 1
        return repl_ptrn
    
    def get_files(self) -> tuple:
        rfiles = []
        for p in self.p.rglob("*"):
            if p.parts[0] in self.exdirs: continue
            if p.suffix in self.ftypes:
                rfiles.append(p)
        return tuple(rfiles)
    
    def replace(self) -> int:
        found = 0
        rplcd = 0
        for file in self.rfiles:
            try:
                f = Path(file)
                ff = 0
                rf = 0
                txt = f.read_text()
                for ptrn in self.rptrns:
                    new_txt, cnt = re.subn(ptrn, self.rplc_ptrn, txt, flags=re.MULTILINE)
                    if cnt > 0:
                        ff += 1
                        if txt != new_txt:
                            rf += 1                        
                            f.write_text(new_txt)
                            # print(f"{self.rename} set to {self.rplcmt} in {f}")
                        print(f"found {ff} | replaced {rf} | {f}")
                rplcd += rf
                found += ff
            except Exception as e:
                raise SystemError(f"error replacing {self.rename}: {e}")
        print(f"{self.rename} REPLACEMENT SUMMARY: found {found} | replaced {rplcd} | {self.p}")
        return rplcd
                                

        
if __name__ == "__main__":
    main()