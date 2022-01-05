import git_traverse
from sys import argv

git_traverse.traverse(argv[1], ["git pull origin master"])