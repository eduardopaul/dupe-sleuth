from collections import defaultdict
from hashlib import file_digest
from pathlib import Path


p = Path('files/')

d = defaultdict(list)
for step in p.walk():
    root = step[0]
    for file_name in step[2]:
        f_path = root/file_name
        with open(f_path, 'rb') as f:
            hash = file_digest(f, 'md5').hexdigest()
            d[hash].append(str(f_path)) 

repeated_files = [f_paths for hash, f_paths in d.items() if len(f_paths) > 1]

for line in repeated_files:
    print(line)

