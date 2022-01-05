import os
def traverse(root, cmds, matches=None):
    for i in os.walk(root):
        dir = i[0]
        if os.path.isdir(dir):
            filename = dir.split("/")[-1]
            if filename == ".git":
                path = dir[0:dir.rfind("/")]
                projectname = path.split("/")[-1]

                # match
                matched = True
                if matches and projectname:
                    matched = projectname in matches
                if not matched:
                    continue

                # cmd
                print("path: " + path)
                print("project name: " + projectname)
                for cmd in cmds:
                    status = os.system("cd " + path + " && " + cmd)
                    if status == 0:
                        print("success")
                    else:
                        print("fail")
                        break
