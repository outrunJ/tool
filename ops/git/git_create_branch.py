import git_traverse
from sys import argv

# 处理入参
rootpath = argv[1]
source_branch_name = argv[2]
branchname = argv[3]
project_names = []
if len(argv) > 4:
    projectlist = argv[4]
    project_names = projectlist.split(",")

# 执行
git_traverse.traverse(rootpath,
                      ["git checkout " + source_branch_name,
                       "git pull origin " + source_branch_name,
                       "git checkout -b " + branchname,
                       "git push origin " + branchname],
                      project_names
                      )