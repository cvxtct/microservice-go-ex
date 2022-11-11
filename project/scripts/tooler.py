from genericpath import isfile
#TODO finish
# https://marcofranssen.nl/manage-go-tools-via-go-modules

import os
import pathlib

root_path = str(pathlib.Path(__file__).parents[2])

mod_list = []

for subdir, dirs, files in os.walk(root_path):
    for file in files:
        if file == 'go.mod':
            print(os.path.join(subdir, file))
            mod_list.append(os.path.join(subdir, file))

for file in mod_list:
    with open(file=file) as f:
        for line in f:
            print(line.rstrip())

