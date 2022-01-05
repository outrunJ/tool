import os
def traverse(root, cmd):
    for i in os.walk(root):
        dir = i[0]
        if os.path.isdir(dir):
            filename = dir.split("/")[-1]
            if filename == ".git":
                path = dir[0:dir.rfind("/")]
                print("enter path: " + path)
                os.system("cd " + path)
                status = os.system(cmd)
                if status == 0 :
                    print("success")
                else:
                    print("fail")
                os.system("cd -")