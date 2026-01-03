#!/usr/bin/env python3

from pathlib import Path
import argparse
import subprocess
import sys
import re
             
PROD_URL = "https://jdeko.me/"
LOCL_URL = "https://dev.jdeko.me/"
# LOCL_URL = "http://localhost/"
URLS_EXCL = [".git", "wiki", "log", "z_log", "puml", "bin"]
URLS_FTYP = [".js", ".html", ".css", ".yaml", ".Dockerfile"]
PROD_CPU = "arm64"
LOCL_CPU = "amd64"

RE_URLS = r'https?://(?:localhost|jdeko(?:.me)?):?[0-9]*/?'
RE_PROD_URL = rf'^(\s*PROD_URL\s+=\s+")({RE_URLS})("\s*)$'
RE_IS_PROD = r"^(\s*IS_PROD\s*=\s*)(true|false)(\s*$)"
RE_GOARCH = rf"^(.*GOARCH=)({LOCL_CPU}|{PROD_CPU})(.*)$"
RE_SSL_DKR = r"(^|# )(COPY ssl.*$)"
RE_SSL_NGX = r"(^\s*listen 80;\s*$\n)([\s\S]*?)(^\s*access.*$)"

def main():
    args = parse_args()

    loc = args.local 
            
    comment_nginx_ssl("jdme-dkr/proxy/nginx.conf", RE_SSL_NGX, loc)
            
    to_replace = [
        ReReplace(loc, "URLS", ".", LOCL_URL, PROD_URL, 0, 0, RE_URLS, URLS_FTYP, URLS_EXCL),
        ReReplace(loc, "PROD_URL", "./main/main.go", PROD_URL, PROD_URL, 3, 2, RE_PROD_URL, [], []),
        ReReplace(loc, "IS_PROD", "./main/main.go", "false", "true", 3, 2, RE_IS_PROD, [], []),
        ReReplace(loc, "GOARCH", "./jdme-dkr/api.Dockerfile", LOCL_CPU, PROD_CPU, 3, 2, RE_GOARCH, [], []),
        ReReplace(loc, "SSL_DKR", "./jdme-dkr/proxy/nginx.Dockerfile", r"# ", r"", 2, 1, RE_SSL_DKR, [], []),
    ]
    
    files_changed = 0
    rplcmnts = 0
    found = 0
    for r in to_replace:
        cnt = r.replace()
        files_changed += cnt[0]
        found += cnt[1]
        rplcmnts += cnt[2]
        print(f"{r.rename} SUMMARY: {cnt[0]} file(s) changed | {cnt[1]} match(es) | {cnt[2]} replacement(s) | {r.p}\n")
    print(f"COMPLETE | {files_changed} file(s) changed | {found} match(es) | {rplcmnts} replacement(s)")
    
    if files_changed > 0:
        msg = f"replaced {rplcmnts} string(s) in {files_changed} file(s)"
        if args.no_push:
            print(msg)    
        else: 
            push_changes(msg)
            print("changes pushed")
    
def comment_nginx_ssl(path: Path, pattern: str, local: bool) -> str:
    p = Path(path)
    ptrn = re.compile(pattern, re.MULTILINE)
    txt = p.read_text()
    m = ptrn.search(txt)
    if not m: 
        return txt
    startline = m.group(1)
    to_comment = m.group(2)
    endline = m.group(3)
    
    commented = []
    repl_str = "# " if local else ""
    for l in to_comment.splitlines(keepends=True):
        new_line, _ = re.subn(r"(^\s*)(# |)", fr"{repl_str}\1", l)
        commented.append(new_line)
    
    repl_block = "".join(commented)
    replacement = f"{startline}{repl_block}{endline}"
    
    spos, epos = m.span()
    new_txt = txt[:spos] + replacement + txt[epos:]
    
    if new_txt != txt:
        p.write_text(new_txt)


# capt_groups -> number of capture groups in pattern
# grp_pos -> capture group index (1:) that is replace with the replacement string
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
                    finally:
                        if len(self.rfiles) > 1:
                            print(f"  | {ff} match(es) | {rf} replacement(s) | {f}")
            rplcd += rf
            found += ff
        return [fc, found, rplcd]
                                

def run(cmd:str, msg=None) -> subprocess.CompletedProcess:
    to_run = cmd.split()
    if msg is not None:
        to_run.append(msg)
    return subprocess.run(to_run, check=False, 
        stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    
def run_commands(cmds: dict[str, str | None]):
    for cmd, msg in cmds.items():
        res = run(cmd, msg)
        if res.returncode > 0:
            raise SystemError(f"command failed: {cmd} {msg} | code {res.returncode}")
        

def push_changes(commit_msg:str):
    diff = run("git diff --quiet")
    if diff.returncode == 0:
        print("no changes to push")
        return
    # figure out how to call
    run_commands({
        "git config user.name github-actions[bot]": None,
        "git config user.email github-actions[bot]@users.noreply.github.com": None,
        "git add .": None,
        'git commit -m': commit_msg,
        "git push": None
    })
    
def parse_args():
    parser = argparse.ArgumentParser(
        description="Toggle nginx + docker config between local and prod"
    )

    parser.add_argument(
        "--local", "-l",
        action="store_true",
        help="Apply local (non-prod) configuration"
    )

    parser.add_argument(
        "--no-push", "-np",
        action="store_true",
        help="Do not push git changes"
    )

    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would change, but do not write files"
    )

    return parser.parse_args()


if __name__ == "__main__":
    main()