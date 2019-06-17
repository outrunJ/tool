from sys import argv
import os
import re

def parse(orig):


def main(orinpath, targetpath):
    origfile = open(orinpath, 'r')
    targetfile = open(targetpath, 'w')
    try:
        orig = origfile.read()
        target = parse(orig)
        targetfile.write(target)
    finally:
        origfile.close()
        targetfile.close()

if __name__ == '__main__':
    main(argv[1], argv[2])