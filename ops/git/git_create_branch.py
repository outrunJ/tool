import git_traverse
from sys import argv

rootpath = argv[1]
source_branch_name = argv[2]
branchname = argv[3]
projectlist = argv[4]

git_traverse.traverse(rootpath,
                      ["git checkout " + source_branch_name,
                       "git pull origin " + source_branch_name,
                       "git checkout -b " + branchname,
                       "git push origin " + branchname],
                      projectlist.split(","))