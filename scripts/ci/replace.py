#!/usr/bin/env python3

from pathlib import Path
import sys
import re
             
PROD_URL = "https://jdeko.me/"
LOCL_URL = "http://localhost:8080/"
URLS_EXCL = [".git", "wiki", "log", "z_log", "puml", "bin"]
URLS_FTYP = [".js", ".html", ".css"]
PROD_CPU = "arm64"
LOCL_CPU = "amd64"

RE_URLS = r'https?://(?:localhost|jdeko(?:.me)?):?[0-9]*/?'
RE_PROD_URL = rf'^(\s*PROD_URL\s+=\s+")({RE_URLS})("\s*)$'
RE_IS_PROD = r"^(\s*IS_PROD\s*=\s*)(true|false)(\s*$)"
RE_GOARCH = rf"^(.*GOARCH=)({LOCL_CPU}|{PROD_CPU})(.*)$"

def main():
    loc = False
    if len(sys.argv) > 1:
        if sys.argv[1] == "local":
            loc = True            
            
    to_replace = [
        ReReplace(loc, "URLS", ".", LOCL_URL, PROD_URL, 0, 0, RE_URLS, URLS_FTYP, URLS_EXCL),
        ReReplace(loc, "PROD_URL", "./main/main.go", PROD_URL, PROD_URL, 3, 2, RE_PROD_URL, [], []),
        ReReplace(loc, "IS_PROD", "./main/main.go", "false", "true", 3, 2, RE_IS_PROD, [], []),
        ReReplace(loc, "GOARCH", "./Dockerfile", LOCL_CPU, PROD_CPU, 3, 2, RE_GOARCH, [], [])
    ]
    
    files = 0
    rplcmnts = 0
    found = 0
    for r in to_replace:
        cnt = r.replace()
        files += cnt[0]
        found += cnt[1]
        rplcmnts += cnt[2]
        print(f"{r.rename} SUMMARY: {cnt[0]} files changed | {cnt[1]} matches | {cnt[2]} replacements | {r.p}\n")
    print(f"COMPLETE | {files} files changed | {found} matches | {rplcmnts} replacements")
    

class ReReplace:
    def __init__(self, loc: bool, name: str,
                path, local_repl, prod_repl: str, 
                capt_grps: int, grp_pos: int,
                pattern: str, filetypes: tuple[str], exclude_dirs: tuple[str]
                ):
        self.p = Path(path)
        self.rename = name
        self.reptrn = pattern
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
    
    def replace(self) -> tuple:
        found = 0
        rplcd = 0
        fc = 0
        print(f"{self.rename}: searching for {self.reptrn} in {self.p}")
        for file in self.rfiles:
            f = Path(file)
            ff = 0
            rf = 0
            cnt = 0
            try:
                txt = f.read_text()
                new_txt, cnt = re.subn(self.reptrn, self.rplc_ptrn, txt, flags=re.MULTILINE)
            except Exception as e:
                print(f"error reading {f} {e}\ncontinuing...")
                continue
            if cnt > 0:
                ff += cnt
                if txt != new_txt:
                    fc += 1
                    rf += cnt                      
                    try:
                        f.write_text(new_txt)
                    except Exception as e: 
                        print(f"error writing to {f}: {e}\ncontinuing...")
                        continue
            rplcd += rf
            found += ff
        return [fc, found, rplcd]
                                
if __name__ == "__main__":
    main()