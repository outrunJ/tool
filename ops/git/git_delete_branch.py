import git_traverse
from sys import argv

rootpath = argv[1]
branchname = argv[2]
projectlist = argv[3]

git_traverse.traverse(rootpath,
                      ["git checkout main",
                       "git push origin --delete " + branchname,
                       "git branch -D " + branchname],
                      projectlist.split(","))