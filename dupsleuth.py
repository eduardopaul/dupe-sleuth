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

    # First group by file size.
    size_dict = defaultdict(list)
    for step in p.walk():
        root = step[0]
        for file_name in step[2]:
            file_path = root/file_name
            file_size = file_path.stat().st_size
            size_dict[file_size].append(file_path)

    # Then check duplicates by hash.
    duplicate_files = defaultdict(list)
    for size, file_list in size_dict.items():
        if len(file_list) > 1:
            for file_path in file_list:
                with open(file_path, "rb") as f:
                    hash = file_digest(f, "md5").hexdigest()
                    duplicate_files[hash].append(str(file_path.relative_to(p))) 

    repeated_files = [
        file_paths
        for hash, file_paths in duplicate_files.items()
        if len(file_paths) > 1
    ]

    output_file = Path("repeated_files.txt")
    with open(output_file, "w") as f:
        f.write("Each group is a set of tentative identical files.\n\n")
        for line in repeated_files:
            for file in line:
                f.write(str(file) + "\n")
            f.write("\n")

else:
    print(f"The directory \"{p}\" does not exist.")

