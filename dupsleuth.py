from argparse import ArgumentParser
from collections import defaultdict
from hashlib import file_digest
from pathlib import Path


parser = ArgumentParser()

parser.add_argument(
    "dir",
    help="Root directory the search should be performed in.",
    type=Path,
)

args = parser.parse_args()

p = args.dir

if p.exists():

    d = defaultdict(list)
    for step in p.walk():
        root = step[0]
        for file_name in step[2]:
            f_path = root/file_name
            with open(f_path, "rb") as f:
                hash = file_digest(f, "md5").hexdigest()
                d[hash].append(str(f_path)) 

    repeated_files = [f_paths for hash, f_paths in d.items() if len(f_paths) > 1]

    output_file = Path("repeated_files.txt")
    with open(output_file, "w") as f:
        f.write("Each line is a set of identical files.\n")
        for line in repeated_files:
            f.write(str(line) + "\n")

else:
    print(f"The directory \"{p}\" does not exist.")
